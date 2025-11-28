<script lang="ts">
  import DraftsModal from "./DraftsModal.svelte";
  import { mockDraftsDetails } from "../mocks/drafts.mock";
  import { createPost, uploadPostImages } from "../services/post-service";
  import { getCommunities } from "../services/community-service";
  import type { CommunityResponse } from "../dtos/community-dto";

  interface Props {
    show: boolean;
    onClose: () => void;
    communityName?: string; // Nếu có thì auto-fill community
    onPostCreated?: () => void; // Callback khi đăng bài thành công
  }

  let { show, onClose, communityName, onPostCreated }: Props = $props();

  let activeTab = $state<"text" | "images" | "link">("text");
  let selectedCommunity = $state(communityName || "");
  let title = $state("");
  let tags = $state<string[]>([]);
  let bodyText = $state("");
  let linkUrl = $state("");
  let mediaFiles = $state<File[]>([]);
  let isDragging = $state(false);
  let showCommunitySearch = $state(false);
  let communitySearchQuery = $state("");
  let showDraftsModal = $state(false);
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);
  let allCommunities = $state<CommunityResponse[]>([]);

  // Load communities from API
  $effect(() => {
    if (show && allCommunities.length === 0) {
      loadCommunities();
    }
  });

  async function loadCommunities() {
    try {
      const response = await getCommunities({ limit: 100 });
      allCommunities = response.communities;
    } catch (error) {
      console.error("Failed to load communities:", error);
    }
  }

  $effect(() => {
    if (communityName) {
      selectedCommunity = communityName;
    }
  });

  const filteredCommunities = $derived(
    allCommunities.filter((c) =>
      c.name.toLowerCase().includes(communitySearchQuery.toLowerCase())
    )
  );

  // Get selected community details
  const selectedCommunityData = $derived(
    allCommunities.find((c) => c.name === selectedCommunity)
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
    showCommunitySearch = false;
    communitySearchQuery = "";
    errorMessage = null;
    isSubmitting = false;
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

  function handleEditDraft(draftId: number) {
    console.log("Edit draft:", draftId);

    // Load draft data from mockDraftsDetails
    const draft = mockDraftsDetails[draftId];
    if (draft) {
      title = draft.title;
      selectedCommunity = draft.community;
      activeTab = draft.tab;
      tags = draft.tags || [];

      if (draft.tab === "text") {
        bodyText = draft.bodyText || "";
      } else if (draft.tab === "link") {
        linkUrl = draft.linkUrl || "";
      }

      console.log("Loaded draft:", draft);
    }

    showDraftsModal = false;
  }

  function handleSaveDraft() {
    // TODO: Backend doesn't have drafts API yet
    console.log("Save draft:", { title, bodyText, tags, selectedCommunity });
    alert("Draft saved locally (backend API not available yet)");
  }

  async function handlePost() {
    // Validation
    if (!title.trim()) {
      errorMessage = "Title is required!";
      return;
    }

    const targetCommunity = selectedCommunity || communityName;
    if (!targetCommunity) {
      errorMessage = "Please select a community!";
      return;
    }

    // Backend doesn't support link posts yet
    if (activeTab === "link") {
      errorMessage =
        "Link posts are not supported yet. Please use Text or Images tab.";
      return;
    }

    try {
      isSubmitting = true;
      errorMessage = null;

      // Get community ID from name
      let community: CommunityResponse | undefined;

      if (allCommunities.length === 0) {
        const response = await getCommunities({ limit: 100 });
        allCommunities = response.communities;
      }

      community = allCommunities.find((c) => c.name === targetCommunity);

      if (!community) {
        errorMessage = "Community not found";
        isSubmitting = false;
        return;
      }

      // Create post
      const post = await createPost({
        community_id: community.id,
        title: title.trim(),
        type: "text", // Backend only accepts text/poll for now
        text: bodyText.trim() || undefined,
      });

      // Upload images if any (2-step flow)
      if (activeTab === "images" && mediaFiles.length > 0) {
        await uploadPostImages(post.id, mediaFiles);
      }

      console.log("✅ Post created successfully:", post);

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

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragging = false;

    const files = Array.from(e.dataTransfer?.files || []);
    const imageFiles = files.filter(
      (f) => f.type.startsWith("image/") || f.type.startsWith("video/")
    );
    mediaFiles = [...mediaFiles, ...imageFiles];
  }

  function handleFileSelect(e: Event) {
    const input = e.target as HTMLInputElement;
    const files = Array.from(input.files || []);
    mediaFiles = [...mediaFiles, ...files];
  }
</script>

{#if show}
  <div class="modal-overlay" onclick={handleClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h2>Create post</h2>
        <span class="drafts-indicator" onclick={handleOpenDrafts}>
          Drafts <span class="draft-count">1</span>
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
                : "Select a community"}</span
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
              placeholder="Select a community"
              class="community-search-input"
            />
          </div>
        {/if}

        <!-- Dropdown list of communities -->
        {#if showCommunitySearch}
          <div class="community-dropdown">
            {#if filteredCommunities.length > 0}
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
                      {community.member_count} members
                    </div>
                  </div>
                </button>
              {/each}
            {:else}
              <div class="no-results">No communities found</div>
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
          Text
        </button>
        <button
          class="tab-btn"
          class:active={activeTab === "images"}
          onclick={() => (activeTab = "images")}
        >
          Images & Video
        </button>
        <button
          class="tab-btn"
          class:active={activeTab === "link"}
          onclick={() => (activeTab = "link")}
        >
          Link
        </button>
      </div>

      <!-- Title Input -->
      <div class="input-group input-with-required">
        <input
          type="text"
          placeholder="Title"
          bind:value={title}
          class="title-input"
        />
        {#if !title}
          <span class="required-mark">*</span>
        {/if}
      </div>

      <!-- Add Tags Button -->
      <!-- TODO: Backend doesn't support tags yet -->
      <button class="add-tags-btn" disabled title="Tags not supported yet"
        >Add tags</button
      >

      <!-- Content Area based on active tab -->
      {#if activeTab === "text"}
        <div class="body-text-container">
          <textarea
            placeholder="Body text (optional)"
            bind:value={bodyText}
            class="body-textarea"
          ></textarea>
          <div class="editor-tools">
            <button class="tool-btn" title="Add picture">
              <img
                src="/picture_icon.svg"
                alt="Picture"
                width="20"
                height="20"
              />
            </button>
            <button class="tool-btn" title="Add link">
              <img src="/link_icon.svg" alt="Link" width="20" height="20" />
            </button>
            <button class="tool-btn" title="Add video">
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
            <p>Drag or upload media</p>
          </label>
          {#if mediaFiles.length > 0}
            <div class="media-previews">
              {#each mediaFiles as file}
                <div class="media-preview">{file.name}</div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "link"}
        <div class="warning-message">
          ⚠️ Link posts are not supported yet. Please use Text or Images tab.
        </div>
        <div class="input-group input-with-required">
          <input
            type="url"
            placeholder="Link URL"
            bind:value={linkUrl}
            class="link-input"
            disabled
          />
          {#if !linkUrl}
            <span class="required-mark">*</span>
          {/if}
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
          Save draft
        </button>
        <button class="post-btn" onclick={handlePost} disabled={isSubmitting}>
          {isSubmitting ? "Posting..." : "Post"}
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

  .no-results {
    padding: 24px;
    text-align: center;
    color: var(--grayfont);
    font-size: 14px;
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

  .media-previews {
    margin-top: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
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
</style>
