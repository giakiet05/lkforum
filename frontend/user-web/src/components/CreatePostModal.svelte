<script lang="ts">
  import DraftsModal from "./DraftsModal.svelte";
  import { createPost, uploadPostImages } from "../services/post-service";
  import { getCommunitiesByUserId } from "../services/community-service";
  import { toastStore } from "../stores/toast-store";
  import {
    createDraft,
    updateDraft,
    getDraftById,
    getDrafts,
    deleteDraft,
  } from "../services/draft-service";
  import { authStore } from "../stores/auth-store";
  import type { CommunityResponse } from "../dtos/community-dto";

  interface Props {
    show: boolean;
    onClose: () => void;
    communityName?: string; // Nếu có thì auto-fill community
    onPostCreated?: () => void; // Callback khi đăng bài thành công
  }

  let { show, onClose, communityName, onPostCreated }: Props = $props();

  let activeTab = $state<"text" | "images" | "link" | "poll">("text");
  let selectedCommunity = $state(communityName || "");
  let title = $state("");
  let tags = $state<string[]>([]);
  let bodyText = $state("");
  let linkUrl = $state("");
  let mediaFiles = $state<File[]>([]);

  // Poll state
  let pollQuestion = $state("");
  let pollOptions = $state<string[]>(["", ""]);
  let pollDuration = $state("3"); // days
  let allowMultiple = $state(false);

  let isDragging = $state(false);
  let showCommunitySearch = $state(false);
  let communitySearchQuery = $state("");
  let showDraftsModal = $state(false);
  let isSubmitting = $state(false);
  let draftCount = $state(0);
  let errorMessage = $state<string | null>(null);
  let allCommunities = $state<CommunityResponse[]>([]);
  let isLoadingCommunities = $state(false);

  // Track current draft being edited
  let currentDraftId = $state<string | null>(null);

  // Load communities and draft count from API
  $effect(() => {
    if (show) {
      if (allCommunities.length === 0 && !isLoadingCommunities) {
        loadCommunities();
      }
      loadDraftCount();
    }
  });

  async function loadCommunities() {
    try {
      isLoadingCommunities = true;
      const currentUser = $authStore.user;
      if (currentUser) {
        // Load communities user has joined
        allCommunities = (await getCommunitiesByUserId(currentUser.id)) || [];
      } else {
        allCommunities = [];
      }
    } catch (error) {
      console.error("Failed to load communities:", error);
      allCommunities = [];
    } finally {
      isLoadingCommunities = false;
    }
  }

  async function loadDraftCount() {
    try {
      const response = await getDrafts(1, 1); // Just get count, not all drafts
      draftCount = response.pagination?.total || 0;
    } catch (error) {
      console.error("Failed to load draft count:", error);
      draftCount = 0;
    }
  }

  $effect(() => {
    if (communityName) {
      selectedCommunity = communityName;
    }
  });

  const filteredCommunities = $derived(
    allCommunities.filter((c) =>
      c.name.toLowerCase().includes(communitySearchQuery.toLowerCase()),
    ),
  );

  // Get selected community details
  const selectedCommunityData = $derived(
    allCommunities.find((c) => c.name === selectedCommunity),
  );

  function handleClose() {
    // Reset form
    activeTab = "text";
    selectedCommunity = communityName || "";
    title = "";
    tags = [];
    bodyText = "";
    linkUrl = "";
    mediaFiles = [];
    pollQuestion = "";
    pollOptions = ["", ""];
    pollDuration = "3";
    allowMultiple = false;
    showCommunitySearch = false;
    communitySearchQuery = "";
    errorMessage = null;
    isSubmitting = false;
    currentDraftId = null; // Reset draft tracking
    onClose();
  }

  function handleCommunitySelect(communityName: string) {
    selectedCommunity = communityName;
    showCommunitySearch = false;
    communitySearchQuery = "";
  }

  function toggleCommunitySearch() {
    showCommunitySearch = !showCommunitySearch;
    if (showCommunitySearch) {
      // Focus on search input after a small delay
      setTimeout(() => {
        document.getElementById("community-search-input")?.focus();
      }, 100);
    }
  }

  function handleOpenDrafts() {
    showDraftsModal = true;
  }

  function handleCloseDrafts() {
    showDraftsModal = false;
  }

  async function handleEditDraft(draftId: string) {
    try {
      const draft = await getDraftById(draftId);

      // Load draft data
      title = draft.title || "";
      tags = draft.tags || [];
      currentDraftId = draftId; // Track current draft

      // Set community if available
      if (draft.community_id) {
        const community = allCommunities.find(
          (c) => c.id === draft.community_id,
        );
        if (community) {
          selectedCommunity = community.name;
        }
      }

      // Set type and content
      if (draft.type) {
        if (draft.type === "text") {
          activeTab = "text";
          bodyText = draft.content?.text || "";
        } else if (draft.type === "poll") {
          activeTab = "poll";
          bodyText = draft.content?.text || "";
          if (draft.content?.poll) {
            pollQuestion = draft.content.poll.question || "";
            // Options từ backend là {id, text, votes}, cần convert sang string array
            if (Array.isArray(draft.content.poll.options)) {
              pollOptions = draft.content.poll.options.map((opt: any) =>
                typeof opt === "string" ? opt : opt.text || "",
              );
              // Đảm bảo ít nhất 2 options
              while (pollOptions.length < 2) {
                pollOptions.push("");
              }
            }
            allowMultiple = draft.content.poll.allow_multiple || false;
          }
        } else if (draft.type === "image") {
          activeTab = "images";
          bodyText = draft.content?.text || "";
          // Note: Images saved in draft are base64 URLs, but we can't convert them back to File objects
          // User will need to re-upload images when editing from draft
        }
      }

      showDraftsModal = false;
    } catch (error) {
      console.error("Failed to load draft:", error);
      toastStore.error("Failed to load draft. Please try again.");
    }
  }

  // Helper to convert File to base64
  function fileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }

  async function handleSaveDraft() {
    try {
      isSubmitting = true;
      errorMessage = null;

      // Get community ID if selected
      const communityId = selectedCommunityData?.id;

      // Prepare draft data based on active tab
      const draftData: any = {
        community_id: communityId,
        title: title || undefined,
        tags: tags.length > 0 ? tags : undefined,
      };

      // Add type-specific content
      if (activeTab === "text") {
        draftData.type = "text";
        draftData.text = bodyText || undefined;
      } else if (activeTab === "poll") {
        draftData.type = "poll";
        draftData.text = bodyText || undefined;
        if (pollQuestion) {
          // Convert options strings to PollOption objects for backend
          const filteredOptions = pollOptions
            .filter((opt) => opt.trim() !== "")
            .map((text, index) => ({
              id: `option_${index}`,
              text: text.trim(),
              votes: 0,
            }));

          draftData.poll = {
            question: pollQuestion,
            options: filteredOptions,
            allow_multiple: allowMultiple,
            expires_at: pollDuration
              ? new Date(
                  Date.now() + parseInt(pollDuration) * 24 * 60 * 60 * 1000,
                ).toISOString()
              : undefined,
          };
        }
      } else if (activeTab === "images") {
        draftData.type = "image";
        draftData.text = bodyText || undefined;

        // Convert images to base64 for draft storage
        if (mediaFiles.length > 0) {
          const imagePromises = mediaFiles
            .filter((f) => f.type.startsWith("image/"))
            .map(async (file) => {
              const base64 = await fileToBase64(file);
              return { url: base64 };
            });
          const videoPromises = mediaFiles
            .filter((f) => f.type.startsWith("video/"))
            .map(async (file) => {
              const base64 = await fileToBase64(file);
              return { url: base64 };
            });

          const images = await Promise.all(imagePromises);
          const videos = await Promise.all(videoPromises);

          if (images.length > 0) draftData.images = images;
          if (videos.length > 0) draftData.videos = videos;

          toastStore.info("Hình ảnh/video đã được lưu vào bản nháp");
        }
      } else if (activeTab === "link") {
        draftData.type = "text";
        draftData.text = linkUrl
          ? `${bodyText}\n\n${linkUrl}`
          : bodyText || undefined;
      }

      // Update existing draft or create new one
      if (currentDraftId) {
        await updateDraft(currentDraftId, draftData);
      } else {
        await createDraft(draftData);
      }

      toastStore.success("Đã lưu bản nháp!");
      handleClose();
    } catch (error) {
      console.error("Failed to save draft:", error);
      errorMessage =
        error instanceof Error ? error.message : "Không thể lưu bản nháp";
    } finally {
      isSubmitting = false;
    }
  }

  function addPollOption() {
    if (pollOptions.length < 6) {
      pollOptions = [...pollOptions, ""];
    }
  }

  function removePollOption(index: number) {
    if (pollOptions.length > 2) {
      pollOptions = pollOptions.filter((_, i) => i !== index);
    }
  }

  function updatePollOption(index: number, value: string) {
    pollOptions[index] = value;
    pollOptions = [...pollOptions];
  }

  async function handlePost() {
    // Validation
    if (!title.trim()) {
      errorMessage = "Tiêu đề là bắt buộc!";
      return;
    }

    if (title.trim().length < 3) {
      errorMessage = "Tiêu đề phải có ít nhất 3 ký tự!";
      return;
    }

    const targetCommunity = selectedCommunity || communityName;
    if (!targetCommunity) {
      errorMessage = "Vui lòng chọn cộng đồng!";
      return;
    }

    // Poll validation
    if (activeTab === "poll") {
      if (!pollQuestion.trim()) {
        errorMessage = "Câu hỏi khảo sát là bắt buộc!";
        return;
      }
      const filledOptions = pollOptions.filter((opt) => opt.trim() !== "");
      if (filledOptions.length < 2) {
        errorMessage = "Khảo sát cần ít nhất 2 lựa chọn!";
        return;
      }
    }

    // Backend doesn't support link posts yet
    if (activeTab === "link") {
      errorMessage =
        "Bài viết liên kết chưa được hỗ trợ. Vui lòng dùng tab Văn bản hoặc Hình ảnh.";
      return;
    }

    try {
      isSubmitting = true;
      errorMessage = null;

      // Get community ID from name
      let community: CommunityResponse | undefined;

      if (allCommunities.length === 0) {
        const currentUser = $authStore.user;
        if (currentUser) {
          allCommunities = (await getCommunitiesByUserId(currentUser.id)) || [];
        }
      }

      community = allCommunities.find((c) => c.name === targetCommunity);

      if (!community) {
        errorMessage = "Community not found";
        isSubmitting = false;
        return;
      }

      // Create post
      let postData: any = {
        community_id: community.id,
        title: title.trim(),
        type: activeTab === "poll" ? "poll" : "text",
        text: bodyText.trim() || "",
      };

      // Add poll data if poll type
      if (activeTab === "poll") {
        const filledOptions = pollOptions.filter((opt) => opt.trim() !== "");

        // Calculate expires_at based on duration
        const expiresAt = new Date();
        expiresAt.setDate(expiresAt.getDate() + parseInt(pollDuration));

        postData.poll = {
          question: pollQuestion.trim(),
          options: filledOptions.map((text) => text.trim()), // Array of strings, not objects
          expires_at: expiresAt.toISOString(),
          allow_multiple: allowMultiple,
        };
      }

      console.log("📤 Creating post with data:", postData);
      const post = await createPost(postData);

      // Upload images if any (2-step flow)
      if (activeTab === "images" && mediaFiles.length > 0) {
        await uploadPostImages(post.id, mediaFiles);
      }

      console.log("✅ Post created successfully:", post);

      // Delete draft if this was from a draft
      if (currentDraftId) {
        try {
          await deleteDraft(currentDraftId);
          console.log("🗑️ Draft deleted after successful post");
        } catch (err) {
          console.warn("Failed to delete draft:", err);
          // Don't show error to user - post was successful
        }
      }

      // Success! Notify parent to reload posts
      if (onPostCreated) {
        console.log("🔄 Calling onPostCreated callback to reload posts");
        onPostCreated();
      } else {
        console.warn("⚠️ onPostCreated callback not provided!");
      }

      handleClose();
    } catch (error) {
      console.error("❌ Failed to create post:", error);
      errorMessage =
        error instanceof Error ? error.message : "Failed to create post";
    } finally {
      isSubmitting = false;
    }
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    isDragging = true;
  }

  function handleDragLeave(e: DragEvent) {
    e.preventDefault();
    isDragging = false;
  }

  const MAX_FILE_SIZE = 100 * 1024 * 1024; // 100MB

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragging = false;

    const files = Array.from(e.dataTransfer?.files || []);
    const validFiles = files.filter((f) => {
      if (!f.type.startsWith("image/") && !f.type.startsWith("video/")) {
        return false;
      }
      if (f.size > MAX_FILE_SIZE) {
        toastStore.warning(`File "${f.name}" vượt quá giới hạn 100MB`);
        return false;
      }
      return true;
    });
    mediaFiles = [...mediaFiles, ...validFiles];
  }

  function handleFileSelect(e: Event) {
    const input = e.target as HTMLInputElement;
    const files = Array.from(input.files || []);
    const validFiles = files.filter((f) => {
      if (f.size > MAX_FILE_SIZE) {
        toastStore.warning(`File "${f.name}" vượt quá giới hạn 100MB`);
        return false;
      }
      return true;
    });
    mediaFiles = [...mediaFiles, ...validFiles];
    // Reset input to allow re-selecting same file
    input.value = "";
  }

  function removeMediaFile(index: number) {
    mediaFiles = mediaFiles.filter((_, i) => i !== index);
  }

  function getMediaPreviewUrl(file: File): string {
    return URL.createObjectURL(file);
  }
</script>

{#if show}
  <div class="modal-overlay" onclick={handleClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h2>Tạo bài viết</h2>
        <span
          class="drafts-indicator"
          onclick={handleOpenDrafts}
          title="Xem các bản nháp đã lưu"
        >
          Bản nháp {#if draftCount > 0}<span class="draft-count"
              >{draftCount}</span
            >{/if}
        </span>
      </div>

      <!-- Community Selector -->
      <div class="community-selector">
        {#if !showCommunitySearch}
          <!-- Button state: Show community name or "Select a community" -->
          <button class="community-display-btn" onclick={toggleCommunitySearch}>
            <div class="community-icon">
              <img
                src={selectedCommunityData?.avatar || "/LKlogo.jpg"}
                alt={selectedCommunityData?.name || "Community"}
              />
            </div>
            <span
              >{selectedCommunity
                ? `lk/${selectedCommunity}`
                : "Chọn cộng đồng"}</span
            >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path
                d="M4 6L8 10L12 6"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
        {:else}
          <!-- Search state: Input with dropdown -->
          <div class="search-input-container">
            <svg
              class="search-icon"
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
            >
              <path
                d="M9 17A8 8 0 1 0 9 1a8 8 0 0 0 0 16zM19 19l-4.35-4.35"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <input
              id="community-search-input"
              type="text"
              bind:value={communitySearchQuery}
              placeholder="Chọn cộng đồng"
              class="community-search-input"
            />
          </div>
        {/if}

        <!-- Dropdown list of communities -->
        {#if showCommunitySearch}
          <div class="community-dropdown">
            {#if isLoadingCommunities}
              <div class="loading-communities">Đang tải cộng đồng...</div>
            {:else if allCommunities.length === 0}
              <div class="no-communities-joined">
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                </svg>
                <span>Hãy tham gia cộng đồng để chia sẻ ý kiến</span>
                <a
                  href="/communities"
                  class="browse-communities-link"
                  onclick={onClose}>Khám phá cộng đồng</a
                >
              </div>
            {:else if filteredCommunities.length > 0}
              {#each filteredCommunities as community}
                <button
                  class="community-item"
                  onclick={() => handleCommunitySelect(community.name)}
                >
                  <div class="community-item-icon">
                    <img
                      src={community.avatar || "/default-community.png"}
                      alt={community.name}
                    />
                  </div>
                  <div class="community-item-info">
                    <div class="community-item-name">lk/{community.name}</div>
                    <div class="community-item-meta">
                      {community.member_count} thành viên
                    </div>
                  </div>
                </button>
              {/each}
            {:else}
              <div class="no-results">Không tìm thấy cộng đồng phù hợp</div>
            {/if}
          </div>
        {/if}
      </div>

      <!-- Tabs -->
      <div class="tabs">
        <button
          class="tab-btn"
          class:active={activeTab === "text"}
          onclick={() => (activeTab = "text")}
        >
          Văn bản
        </button>
        <button
          class="tab-btn"
          class:active={activeTab === "images"}
          onclick={() => (activeTab = "images")}
        >
          Hình ảnh & Video
        </button>
        <button
          class="tab-btn"
          class:active={activeTab === "link"}
          onclick={() => (activeTab = "link")}
        >
          Liên kết
        </button>
        <button
          class="tab-btn"
          class:active={activeTab === "poll"}
          onclick={() => (activeTab = "poll")}
        >
          Khảo sát
        </button>
      </div>

      <!-- Title Input -->
      <div class="input-group input-with-required">
        <input
          type="text"
          placeholder="Tiêu đề"
          bind:value={title}
          class="title-input"
        />
        {#if !title}
          <span class="required-mark">*</span>
        {/if}
      </div>

      <!-- Add Tags Button - tạm ẩn vì backend chưa hỗ trợ
      <button class="add-tags-btn" disabled title="Thẻ chưa được hỗ trợ"
        >Thêm thẻ</button
      >
      -->

      <!-- Content Area based on active tab -->
      {#if activeTab === "text"}
        <div class="body-text-container">
          <textarea
            placeholder="Nội dung (tùy chọn)"
            bind:value={bodyText}
            class="body-textarea"
          ></textarea>
          <div class="editor-tools">
            <button class="tool-btn" title="Thêm ảnh">
              <img
                src="/picture_icon.svg"
                alt="Picture"
                width="20"
                height="20"
              />
            </button>
            <button class="tool-btn" title="Thêm liên kết">
              <img src="/link_icon.svg" alt="Link" width="20" height="20" />
            </button>
            <button class="tool-btn" title="Thêm video">
              <img src="/video_icon.svg" alt="Video" width="20" height="20" />
            </button>
          </div>
        </div>
      {:else if activeTab === "images"}
        <div
          class="media-upload-area"
          class:dragging={isDragging}
          ondragover={handleDragOver}
          ondragleave={handleDragLeave}
          ondrop={handleDrop}
        >
          <input
            type="file"
            id="media-upload"
            accept="image/*,video/*"
            multiple
            onchange={handleFileSelect}
            style="display: none;"
          />
          <label for="media-upload" class="upload-label">
            <img
              src="/uploadmedia_icon.svg"
              alt="Upload"
              width="48"
              height="48"
            />
            <p>Kéo thả hoặc tải lên</p>
            <span class="upload-hint">Tối đa 100MB mỗi file</span>
          </label>
          {#if mediaFiles.length > 0}
            <div class="media-previews">
              {#each mediaFiles as file, index}
                <div class="media-preview-item">
                  {#if file.type.startsWith("image/")}
                    <img
                      src={getMediaPreviewUrl(file)}
                      alt={file.name}
                      class="preview-image"
                    />
                  {:else if file.type.startsWith("video/")}
                    <video
                      src={getMediaPreviewUrl(file)}
                      class="preview-video"
                      muted
                    ></video>
                    <div class="video-indicator">▶</div>
                  {/if}
                  <button
                    class="remove-media-btn"
                    onclick={() => removeMediaFile(index)}
                    title="Xóa file này"
                  >
                    ✕
                  </button>
                  <span class="file-name">{file.name}</span>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "link"}
        <div class="warning-message">
          ⚠️ Bài viết liên kết chưa được hỗ trợ. Vui lòng sử dụng tab Văn bản
          hoặc Hình ảnh.
        </div>
        <div class="input-group input-with-required">
          <input
            type="url"
            placeholder="URL liên kết"
            bind:value={linkUrl}
            class="link-input"
            disabled
          />
          {#if !linkUrl}
            <span class="required-mark">*</span>
          {/if}
        </div>
      {:else if activeTab === "poll"}
        <div class="poll-container">
          <!-- Poll Question -->
          <div class="input-group input-with-required">
            <input
              type="text"
              placeholder="Đặt câu hỏi..."
              bind:value={pollQuestion}
              class="poll-question-input"
            />
            {#if !pollQuestion}
              <span class="required-mark">*</span>
            {/if}
          </div>

          <!-- Poll Options -->
          <div class="poll-options">
            {#each pollOptions as option, index}
              <div class="poll-option-row">
                <input
                  type="text"
                  placeholder={`Lựa chọn ${index + 1}`}
                  value={option}
                  oninput={(e) =>
                    updatePollOption(index, e.currentTarget.value)}
                  class="poll-option-input"
                />
                {#if pollOptions.length > 2}
                  <button
                    class="remove-option-btn"
                    onclick={() => removePollOption(index)}
                    title="Xóa lựa chọn"
                  >
                    ✕
                  </button>
                {/if}
              </div>
            {/each}

            <!-- Add Option Button -->
            {#if pollOptions.length < 6}
              <button class="add-option-btn" onclick={addPollOption}>
                + Thêm lựa chọn
              </button>
            {/if}
          </div>

          <!-- Poll Settings -->
          <div class="poll-settings">
            <div class="setting-row">
              <label for="poll-duration">Thời gian khảo sát:</label>
              <select
                id="poll-duration"
                bind:value={pollDuration}
                class="duration-select"
              >
                <option value="1">1 ngày</option>
                <option value="3">3 ngày</option>
                <option value="7">7 ngày</option>
                <option value="14">14 ngày</option>
              </select>
            </div>

            <div class="setting-row">
              <label>
                <input
                  type="checkbox"
                  bind:checked={allowMultiple}
                  class="checkbox"
                />
                Cho phép chọn nhiều
              </label>
            </div>
          </div>
        </div>
      {/if}

      <!-- Error Message -->
      {#if errorMessage}
        <div class="error-message">
          {errorMessage}
        </div>
      {/if}

      <!-- Action Buttons -->
      <div class="modal-actions">
        <button
          class="save-draft-btn"
          onclick={handleSaveDraft}
          disabled={isSubmitting}
        >
          Lưu nháp
        </button>
        <button class="post-btn" onclick={handlePost} disabled={isSubmitting}>
          {isSubmitting ? "Đang đăng..." : "Đăng"}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Drafts Modal -->
<DraftsModal
  show={showDraftsModal}
  onClose={handleCloseDrafts}
  onEditDraft={handleEditDraft}
/>

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    z-index: 1000;
    padding: 60px 20px 20px;
    overflow-y: auto;
  }

  .modal-content {
    background: white;
    border-radius: 8px;
    width: 100%;
    max-width: 740px;
    padding: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .modal-header h2 {
    font-size: 20px;
    font-weight: 700;
    color: var(--blue--);
    margin: 0;
  }

  .drafts-indicator {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
    transition: background-color 0.2s ease;
  }

  .drafts-indicator:hover {
    background-color: rgba(0, 0, 0, 0.05);
    color: var(--blue--);
  }

  .draft-count {
    opacity: 0.6;
  }

  /* Community Selector */
  .community-selector {
    margin-bottom: 16px;
    position: relative;
  }

  .community-display-btn {
    background: rgba(214, 216, 222, 0.4);
    width: fit-content;
    border-radius: 16px;
    padding: 8px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    border: none;
    cursor: pointer;
    font-size: 14px;
    font-weight: 600;
    color: var(--blue--);
    transition: background 0.2s;
  }

  .community-display-btn:hover {
    background: rgba(214, 216, 222, 0.5);
  }

  .search-input-container {
    position: relative;
    width: 100%;
  }

  .search-icon {
    position: absolute;
    left: 16px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--blue--);
    pointer-events: none;
  }

  .community-search-input {
    width: 100%;
    background: rgba(214, 216, 222, 0.3);
    border: 2px solid var(--blue--);
    border-radius: 20px;
    padding: 10px 16px 10px 48px;
    font-size: 14px;
    color: var(--blue--);
    outline: none;
  }

  .community-search-input::placeholder {
    color: var(--grayfont);
  }

  .community-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 8px;
    background: white;
    border: 1px solid var(--lightgray--);
    border-radius: 8px;
    max-height: 400px;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 10;
  }

  .community-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border: none;
    background: white;
    cursor: pointer;
    text-align: left;
    transition: background 0.2s;
  }

  .community-item:hover {
    background: #f6f7f8;
  }

  .community-item-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
  }

  .community-item-icon img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .community-item-info {
    flex: 1;
  }

  .community-item-name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    margin-bottom: 4px;
  }

  .community-item-meta {
    font-size: 12px;
    color: var(--grayfont);
  }

  .no-results,
  .loading-communities {
    padding: 24px;
    text-align: center;
    color: var(--grayfont);
    font-size: 14px;
  }

  .no-communities-joined {
    padding: 24px;
    text-align: center;
    color: var(--grayfont);
    font-size: 14px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  .no-communities-joined svg {
    opacity: 0.5;
  }

  .no-communities-joined span {
    color: var(--grayfont);
  }

  .browse-communities-link {
    color: var(--primary);
    text-decoration: none;
    font-weight: 500;
    padding: 8px 16px;
    border-radius: 20px;
    background: var(--primary-light, rgba(0, 122, 255, 0.1));
    transition: background 0.2s;
  }

  .browse-communities-link:hover {
    background: var(--primary-hover, rgba(0, 122, 255, 0.2));
  }

  .community-icon {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .community-icon img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  /* Tabs */
  .tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid #edeff1;
    margin-bottom: 16px;
  }

  .tab-btn {
    padding: 12px 16px;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tab-btn.active {
    color: #1c6ea4;
    border-bottom-color: #1c6ea4;
  }

  .tab-btn:hover {
    color: #1c1c1c;
  }

  /* Input Groups */
  .input-group {
    margin-bottom: 16px;
  }

  .input-with-required {
    position: relative;
  }

  .required-mark {
    position: absolute;
    right: 16px;
    top: 50%;
    transform: translateY(-50%);
    color: #ff0000;
    font-size: 16px;
    font-weight: 600;
    pointer-events: none;
  }

  .title-input,
  .link-input {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid var(--lightgray--);
    border-radius: 16px;
    font-size: 14px;
    color: #1c1c1c;
    background: white;
  }

  .title-input:focus,
  .link-input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .title-input::placeholder,
  .link-input::placeholder {
    color: var(--grayfont);
  }

  /* Add Tags Button */
  .add-tags-btn {
    padding: 6px 12px;
    background: #f6f7f8;
    border: none;
    border-radius: 4px;
    font-size: 13px;
    font-weight: 500;
    color: #7c7c7c;
    cursor: pointer;
    margin-bottom: 16px;
    transition: background 0.2s;
  }

  .add-tags-btn:hover:not(:disabled) {
    background: #edeff1;
  }

  .add-tags-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Body Text Container */
  .body-text-container {
    position: relative;
    margin-bottom: 16px;
  }

  .body-textarea {
    width: 100%;
    min-height: 200px;
    padding: 16px;
    padding-bottom: 48px;
    border: 1px solid var(--lightgray--);
    border-radius: 12px;
    font-size: 14px;
    color: var(--grayfont);
    background: white;
    resize: vertical;
    font-family: inherit;
  }

  .body-textarea:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .body-textarea::placeholder {
    color: var(--grayfont);
  }

  .editor-tools {
    position: absolute;
    bottom: 12px;
    right: 12px;
    display: flex;
    gap: 8px;
  }

  .tool-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 1px solid #edeff1;
    border-radius: 4px;
    color: #7c7c7c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tool-btn:hover {
    background: #f6f7f8;
    color: #1c1c1c;
  }

  /* Media Upload Area */
  .media-upload-area {
    border: 2px dashed #000000;
    border-radius: 8px;
    padding: 48px 24px;
    text-align: center;
    margin-bottom: 16px;
    transition: all 0.2s;
  }

  .media-upload-area.dragging {
    border-color: var(--blue--);
    background: rgba(21, 48, 96, 0.05);
  }

  .upload-label {
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }

  .upload-label p {
    margin: 0;
    font-size: 14px;
    color: #7c7c7c;
  }

  .upload-hint {
    font-size: 12px;
    color: #999;
  }

  .media-previews {
    margin-top: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: flex-start;
  }

  .media-preview-item {
    position: relative;
    width: 120px;
    height: 120px;
    border-radius: 8px;
    overflow: hidden;
    background: #f6f7f8;
    border: 1px solid #edeff1;
  }

  .preview-image,
  .preview-video {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .video-indicator {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 36px;
    height: 36px;
    background: rgba(0, 0, 0, 0.6);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 14px;
    pointer-events: none;
  }

  .remove-media-btn {
    position: absolute;
    top: 4px;
    right: 4px;
    width: 24px;
    height: 24px;
    background: rgba(0, 0, 0, 0.7);
    border: none;
    border-radius: 50%;
    color: white;
    font-size: 14px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
    z-index: 1;
  }

  .remove-media-btn:hover {
    background: #c00;
  }

  .file-name {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 4px 6px;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    font-size: 10px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .media-preview {
    padding: 8px 12px;
    background: #f6f7f8;
    border-radius: 4px;
    font-size: 12px;
    color: #1c1c1c;
  }

  /* Error Message */
  .error-message {
    padding: 12px 16px;
    background: #fee;
    border: 1px solid #fcc;
    border-radius: 8px;
    color: #c00;
    font-size: 14px;
    margin-top: 16px;
  }

  /* Warning Message */
  .warning-message {
    padding: 12px 16px;
    background: #fff3cd;
    border: 1px solid #ffc107;
    border-radius: 8px;
    color: #856404;
    font-size: 14px;
    margin-bottom: 12px;
  }

  /* Poll Container */
  .poll-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 16px;
  }

  .poll-question-input {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid var(--lightgray--);
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    color: var(--grayfont);
  }

  .poll-question-input:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .poll-options {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .poll-option-row {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .poll-option-input {
    flex: 1;
    padding: 10px 14px;
    border: 1px solid var(--lightgray--);
    border-radius: 8px;
    font-size: 14px;
    color: var(--grayfont);
  }

  .poll-option-input:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .remove-option-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 1px solid #edeff1;
    border-radius: 4px;
    color: #c00;
    cursor: pointer;
    font-size: 16px;
    transition: all 0.2s;
  }

  .remove-option-btn:hover {
    background: #fee;
    border-color: #fcc;
  }

  .add-option-btn {
    padding: 8px 16px;
    background: transparent;
    border: 1px dashed var(--blue--);
    border-radius: 8px;
    color: var(--blue--);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .add-option-btn:hover {
    background: rgba(21, 48, 96, 0.05);
  }

  .poll-settings {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background: #f6f7f8;
    border-radius: 8px;
  }

  .setting-row {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 14px;
  }

  .duration-select {
    padding: 6px 12px;
    border: 1px solid var(--lightgray--);
    border-radius: 6px;
    background: white;
    color: var(--grayfont);
    font-size: 14px;
    cursor: pointer;
  }

  .duration-select:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .checkbox {
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  /* Action Buttons */
  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }

  .save-draft-btn,
  .post-btn {
    padding: 8px 24px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.2s;
    font-family: "Roboto", sans-serif;
    border: none;
  }

  .save-draft-btn:disabled,
  .post-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .save-draft-btn {
    background: #f6f7f8;
    color: #1c1c1c;
  }

  .save-draft-btn:hover:not(:disabled) {
    background: #edeff1;
  }

  .post-btn {
    background: var(--blue--);
    color: white;
  }

  .post-btn:hover:not(:disabled) {
    background: var(--darkblue--);
  }

  /* Mobile responsive */
  @media (max-width: 768px) {
    .modal-overlay {
      padding: var(--topbar-height) 0 0;
      align-items: stretch;
    }

    .modal-content {
      border-radius: 0;
      padding: 16px;
      max-width: 100%;
      min-height: 100vh;
    }

    .modal-header h2 {
      font-size: 18px;
    }

    .tab-button {
      padding: 8px 12px;
      font-size: 13px;
    }

    .post-input {
      font-size: 16px; /* Prevent zoom on iOS */
    }
  }

  @media (max-width: 480px) {
    .modal-content {
      padding: 12px;
    }

    .modal-header {
      margin-bottom: 16px;
    }

    .modal-header h2 {
      font-size: 16px;
    }

    .modal-actions {
      margin-top: 16px;
      gap: 8px;
    }

    .save-draft-btn,
    .post-btn {
      padding: 8px 16px;
      font-size: 13px;
    }
  }
</style>
