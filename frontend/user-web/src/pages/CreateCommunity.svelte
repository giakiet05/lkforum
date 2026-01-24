<script lang="ts">
  import { push } from "svelte-spa-router";
  import { toastStore } from "../stores/toast-store";

  type Topic = {
    id: string;
    name: string;
    icon: string;
    subtopics: string[];
  };

  let communityName = $state("");
  let description = $state("");
  let communityType = $state<"public" | "restricted" | "private">("public");
  let isAdultContent = $state(false);
  let isLoading = $state(false);
  let error = $state("");
  let showTopicsModal = $state(false);
  let selectedTopics = $state<string[]>([]);
  let topicSearchQuery = $state("");
  let showStyleModal = $state(false);
  let bannerImage = $state<string>("");
  let iconImage = $state<string>("");
  let currentStep = $state(1); // 1: main form, 2: topics, 3: style

  // Mock topics data
  const allTopics: Topic[] = [
    {
      id: "anime",
      name: "Anime & Cosplay",
      icon: "🎭",
      subtopics: ["Anime & Manga", "Cosplay"],
    },
    {
      id: "art",
      name: "Art",
      icon: "🎨",
      subtopics: [
        "Performing Arts",
        "Architecture",
        "Design",
        "Art",
        "Filmmaking",
        "Digital Art",
        "Photography",
      ],
    },
    {
      id: "business",
      name: "Business & Finance",
      icon: "💼",
      subtopics: [
        "Personal Finance",
        "Crypto",
        "Economics",
        "Business News & Discussion",
        "Deals & Marketplace",
        "Startups & Entrepreneurship",
        "Real Estate",
        "Stocks & Investing",
      ],
    },
    {
      id: "collectibles",
      name: "Collectibles & Other Hobbies",
      icon: "⭐",
      subtopics: ["Model Building", "Collectibles", "Other Hobbies", "Toys"],
    },
    {
      id: "entertainment",
      name: "Entertainment",
      icon: "🎬",
      subtopics: [
        "Movies",
        "Television",
        "Streaming",
        "Music",
        "Podcasts & Streamers",
        "Celebrity",
      ],
    },
    {
      id: "food",
      name: "Food & Drink",
      icon: "🍕",
      subtopics: ["Food", "Drinks", "Cooking"],
    },
    {
      id: "gaming",
      name: "Gaming",
      icon: "🎮",
      subtopics: [
        "Video Games",
        "Board Games & TTRPGs",
        "eSports",
        "Gaming News & Discussion",
      ],
    },
    {
      id: "health",
      name: "Health & Fitness",
      icon: "💪",
      subtopics: ["Fitness", "Mental Health", "Wellness", "Health"],
    },
  ];

  const filteredTopics = $derived(() => {
    if (!topicSearchQuery.trim()) return allTopics;
    const query = topicSearchQuery.toLowerCase();
    return allTopics.filter(
      (topic) =>
        topic.name.toLowerCase().includes(query) ||
        topic.subtopics.some((sub) => sub.toLowerCase().includes(query))
    );
  });

  function toggleTopic(subtopic: string) {
    if (selectedTopics.includes(subtopic)) {
      selectedTopics = selectedTopics.filter((t) => t !== subtopic);
    } else if (selectedTopics.length < 3) {
      selectedTopics = [...selectedTopics, subtopic];
    }
  }

  function openTopicsModal() {
    showTopicsModal = true;
  }

  function closeTopicsModal() {
    showTopicsModal = false;
    topicSearchQuery = "";
  }

  function handleTopicsNext() {
    closeTopicsModal();
    showStyleModal = true;
    currentStep = 3;
  }

  function openStyleModal() {
    showStyleModal = true;
    currentStep = 3;
  }

  function closeStyleModal() {
    showStyleModal = false;
  }

  function handleStyleBack() {
    closeStyleModal();
    showTopicsModal = true;
    currentStep = 2;
  }

  function handleBannerUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      const reader = new FileReader();
      reader.onload = (e) => {
        bannerImage = e.target?.result as string;
      };
      reader.readAsDataURL(file);
    }
  }

  function handleIconUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      const reader = new FileReader();
      reader.onload = (e) => {
        iconImage = e.target?.result as string;
      };
      reader.readAsDataURL(file);
    }
  }

  function handleSubmit(e: Event) {
    e.preventDefault();

    // Validation
    if (!communityName.trim()) {
      error = "Tên cộng đồng là bắt buộc";
      return;
    }

    if (communityName.length < 3) {
      error = "Tên cộng đồng phải có ít nhất 3 ký tự";
      return;
    }

    if (communityName.length > 21) {
      error = "Tên cộng đồng phải dưới 21 ký tự";
      return;
    }

    // Check for valid characters (letters, numbers, underscores)
    if (!/^[a-zA-Z0-9_]+$/.test(communityName)) {
      error =
        "Tên cộng đồng chỉ được chứa chữ cái, số và dấu gạch dưới";
      return;
    }

    isLoading = true;
    error = "";

    // Mock create report - sau này sẽ call API
    setTimeout(() => {
      isLoading = false;
      toastStore.success(
        `Cộng đồng "lk/${communityName}" đã được tạo thành công!`
      );
      // Navigate to the new report page
      push(`/lk/${communityName}`);
    }, 1000);
  }

  function handleCancel() {
    push("/");
  }
</script>

<div class="create-report-page">
  <div class="container">
    <div class="header">
      <h1>Tạo cộng đồng</h1>
      <button class="close-btn" onclick={handleCancel}>×</button>
    </div>

    <form onsubmit={handleSubmit}>
      <!-- Community Name -->
      <div class="form-section">
        <label for="name" class="label">Tên</label>
        <p class="help-text">
          Tên cộng đồng bao gồm chữ hoa không thể thay đổi.
        </p>
        <div class="input-wrapper">
          <span class="prefix">lk/</span>
          <input
            type="text"
            id="name"
            bind:value={communityName}
            placeholder="ten_cong_dong"
            maxlength="21"
            disabled={isLoading}
          />
        </div>
        <div class="char-count">
          {communityName.length} / 21
        </div>
      </div>

      <!-- Description -->
      <div class="form-section">
        <label for="description" class="label">Mô tả (tùy chọn)</label>
        <textarea
          id="description"
          bind:value={description}
          placeholder="Cộng đồng của bạn về điều gì?"
          rows="4"
          maxlength="500"
          disabled={isLoading}
        ></textarea>
        <div class="char-count">
          {description.length} / 500
        </div>
      </div>

      <!-- Community Type -->
      <div class="form-section">
        <label class="label">Loại cộng đồng</label>

        <label class="radio-option">
          <input
            type="radio"
            name="type"
            value="public"
            bind:group={communityType}
            disabled={isLoading}
          />
          <div class="radio-content">
            <div class="radio-header">
              <svg
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  d="M10 2a8 8 0 100 16 8 8 0 000-16zM9 13V9h2v4H9zm0-6V5h2v2H9z"
                />
              </svg>
              <span class="radio-title">Công khai</span>
            </div>
            <p class="radio-description">
              Bất kỳ ai cũng có thể xem, đăng và bình luận
            </p>
          </div>
        </label>

        <label class="radio-option">
          <input
            type="radio"
            name="type"
            value="restricted"
            bind:group={communityType}
            disabled={isLoading}
          />
          <div class="radio-content">
            <div class="radio-header">
              <svg
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  d="M10 2a8 8 0 00-8 8v5a3 3 0 003 3h10a3 3 0 003-3v-5a8 8 0 00-8-8zm0 2a6 6 0 016 6v5a1 1 0 01-1 1H5a1 1 0 01-1-1v-5a6 6 0 016-6z"
                />
              </svg>
              <span class="radio-title">Hạn chế</span>
            </div>
            <p class="radio-description">
              Bất kỳ ai cũng có thể xem, nhưng chỉ người được duyệt mới đăng được
            </p>
          </div>
        </label>

        <label class="radio-option">
          <input
            type="radio"
            name="type"
            value="private"
            bind:group={communityType}
            disabled={isLoading}
          />
          <div class="radio-content">
            <div class="radio-header">
              <svg
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  d="M10 2a4 4 0 00-4 4v2H5a2 2 0 00-2 2v6a2 2 0 002 2h10a2 2 0 002-2v-6a2 2 0 00-2-2h-1V6a4 4 0 00-4-4zm2 6V6a2 2 0 10-4 0v2h4z"
                />
              </svg>
              <span class="radio-title">Riêng tư</span>
            </div>
            <p class="radio-description">
              Chỉ người được duyệt mới có thể xem và gửi
            </p>
          </div>
        </label>
      </div>

      <!-- Adult Content -->
      <div class="form-section">
        <label class="checkbox-option">
          <input
            type="checkbox"
            bind:checked={isAdultContent}
            disabled={isLoading}
          />
          <div class="checkbox-content">
            <span class="checkbox-title">Cộng đồng 18+</span>
            <p class="checkbox-description">
              Phải trên 18 tuổi để xem và đóng góp
            </p>
          </div>
        </label>
      </div>

      <!-- Topics Section -->
      <div class="form-section">
        <label class="label">Chủ đề (Tùy chọn)</label>
        <p class="help-text">
          Thêm tối đa 3 chủ đề để giúp mọi người tìm thấy cộng đồng của bạn.
        </p>
        <button
          type="button"
          class="add-topics-btn"
          onclick={openTopicsModal}
          disabled={isLoading}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 4V16M4 10H16"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
          Thêm chủ đề
        </button>
        {#if selectedTopics.length > 0}
          <div class="selected-topics">
            {#each selectedTopics as topic}
              <span class="topic-tag">
                {topic}
                <button
                  type="button"
                  class="remove-topic"
                  onclick={() => toggleTopic(topic)}
                >
                  ×
                </button>
              </span>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Error Message -->
      {#if error}
        <div class="error-message">{error}</div>
      {/if}

      <!-- Actions -->
      <div class="actions">
        <button
          type="button"
          class="btn btn-secondary"
          onclick={handleCancel}
          disabled={isLoading}
        >
          Hủy
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          disabled={isLoading || !communityName.trim()}
        >
          {isLoading ? "Đang tạo..." : "Tạo cộng đồng"}
        </button>
      </div>
    </form>
  </div>
</div>

<!-- Topics Modal -->
{#if showTopicsModal}
  <div class="modal-overlay" onclick={closeTopicsModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h2>Thêm chủ đề</h2>
        <button class="modal-close-btn" onclick={closeTopicsModal}>×</button>
      </div>

      <p class="modal-subtitle">
        Thêm tối đa 3 chủ đề để giúp mọi người tìm thấy cộng đồng của bạn.
      </p>

      <!-- Search Input -->
      <div class="search-wrapper">
        <svg
          class="search-icon"
          width="20"
          height="20"
          viewBox="0 0 20 20"
          fill="none"
        >
          <path
            d="M9 17A8 8 0 1 0 9 1a8 8 0 0 0 0 16zM18 18l-4-4"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
        <input
          type="text"
          placeholder="Lọc chủ đề"
          bind:value={topicSearchQuery}
        />
      </div>

      <div class="topics-counter">Chủ đề {selectedTopics.length}/3</div>

      <!-- Topics List -->
      <div class="topics-list">
        {#each filteredTopics() as topic (topic.id)}
          <div class="topic-category">
            <div class="category-header">
              <span class="category-icon">{topic.icon}</span>
              <span class="category-name">{topic.name}</span>
            </div>
            <div class="subtopics">
              {#each topic.subtopics as subtopic}
                <button
                  type="button"
                  class="subtopic-tag"
                  class:selected={selectedTopics.includes(subtopic)}
                  class:disabled={!selectedTopics.includes(subtopic) &&
                    selectedTopics.length >= 3}
                  onclick={() => toggleTopic(subtopic)}
                  disabled={!selectedTopics.includes(subtopic) &&
                    selectedTopics.length >= 3}
                >
                  {subtopic}
                  {#if selectedTopics.includes(subtopic)}
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                      <path
                        d="M2 6L5 9L10 3"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  {/if}
                </button>
              {/each}
            </div>
          </div>
        {/each}
      </div>

      <!-- Modal Actions -->
      <div class="modal-actions">
        <button
          type="button"
          class="btn btn-secondary"
          onclick={closeTopicsModal}
        >
          Quay lại
        </button>
        <button
          type="button"
          class="btn btn-primary"
          onclick={handleTopicsNext}
          disabled={selectedTopics.length === 0}
        >
          Tiếp theo
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Style Modal -->
{#if showStyleModal}
  <div class="modal-overlay" onclick={closeStyleModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h2>Trang trí cộng đồng</h2>
        <button class="modal-close-btn" onclick={closeStyleModal}>×</button>
      </div>

      <p class="modal-subtitle">
        Thêm hình ảnh sẽ thu hút sự chú ý của thành viên mới và giúp thiết lập văn hóa cộng đồng của bạn! Bạn có thể cập nhật bất cứ lúc nào.
      </p>

      <div class="style-content">
        <!-- Banner Upload -->
        <div class="upload-section">
          <label class="upload-label">Ảnh bìa</label>
          <label class="upload-button">
            <input
              type="file"
              accept="image/*"
              onchange={handleBannerUpload}
              style="display: none;"
            />
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <path
                d="M4 16V14C4 12.8954 4.89543 12 6 12H14C15.1046 12 16 12.8954 16 14V16M13 7L10 4M10 4L7 7M10 4V12"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            Thêm
          </label>
        </div>

        <!-- Icon Upload -->
        <div class="upload-section">
          <label class="upload-label">Biểu tượng</label>
          <label class="upload-button">
            <input
              type="file"
              accept="image/*"
              onchange={handleIconUpload}
              style="display: none;"
            />
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <path
                d="M4 16V14C4 12.8954 4.89543 12 6 12H14C15.1046 12 16 12.8954 16 14V16M13 7L10 4M10 4L7 7M10 4V12"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            Thêm
          </label>
        </div>

        <!-- Preview Card -->
        <div class="preview-card">
          <div
            class="preview-banner"
            style={bannerImage ? `background-image: url(${bannerImage})` : ""}
          ></div>
          <div class="preview-content">
            <div class="preview-icon-wrapper">
              {#if iconImage}
                <img src={iconImage} alt="Icon" class="preview-icon-img" />
              {:else}
                <div class="preview-icon-placeholder">
                  <svg
                    width="32"
                    height="32"
                    viewBox="0 0 32 32"
                    fill="currentColor"
                  >
                    <circle cx="16" cy="16" r="16" fill="#ff4500" />
                    <text
                      x="16"
                      y="20"
                      text-anchor="middle"
                      fill="white"
                      font-size="16"
                      font-weight="bold"
                    >
                      lk/
                    </text>
                  </svg>
                </div>
              {/if}
            </div>
            <div class="preview-info">
              <h3>lk/{communityName || "report"}</h3>
              <p>1 thành viên · 1 trực tuyến</p>
            </div>
          </div>
          <div class="preview-description">
            {description || "Mô tả cộng đồng"}
          </div>
        </div>

        <!-- Progress Dots -->
        <div class="progress-dots">
          <span class="dot"></span>
          <span class="dot"></span>
          <span class="dot active"></span>
        </div>
      </div>

      <!-- Modal Actions -->
      <div class="modal-actions">
        <button
          type="button"
          class="btn btn-secondary"
          onclick={handleStyleBack}
        >
          Quay lại
        </button>
        <button type="button" class="btn btn-primary" onclick={handleSubmit}>
          {isLoading ? "Đang tạo..." : "Tiếp theo"}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .create-report-page {
    min-height: 100vh;
    background: #dae0e6;
    padding: 40px 20px;
  }

  .container {
    max-width: 640px;
    margin: 0 auto;
    background: white;
    border-radius: 4px;
    padding: 0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 24px;
    border-bottom: 1px solid #edeff1;
  }

  .header h1 {
    font-size: 18px;
    font-weight: 600;
    margin: 0;
    color: #1c1c1c;
  }

  .close-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    font-size: 28px;
    color: #878a8c;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    padding: 0;
    line-height: 1;
  }

  .close-btn:hover {
    background: #f6f7f8;
  }

  form {
    padding: 24px;
  }

  .form-section {
    margin-bottom: 24px;
  }

  .label {
    display: block;
    font-size: 16px;
    font-weight: 600;
    color: #1c1c1c;
    margin-bottom: 8px;
  }

  .help-text {
    font-size: 12px;
    color: #7c7c7c;
    margin: 0 0 8px 0;
  }

  .input-wrapper {
    display: flex;
    align-items: center;
    border: 1px solid #ccc;
    border-radius: 4px;
    overflow: hidden;
    background: white;
  }

  .input-wrapper:focus-within {
    border-color: var(--blue--);
  }

  .prefix {
    padding: 10px 12px;
    background: #f6f7f8;
    color: #1c1c1c;
    font-size: 14px;
    font-weight: 500;
    border-right: 1px solid #ccc;
  }

  input[type="text"] {
    flex: 1;
    border: none;
    padding: 10px 12px;
    font-size: 14px;
    color: #1c1c1c;
  }

  input[type="text"]:focus {
    outline: none;
  }

  textarea {
    width: 100%;
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 10px 12px;
    font-size: 14px;
    font-family: inherit;
    resize: vertical;
    color: #1c1c1c;
  }

  textarea:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .char-count {
    text-align: right;
    font-size: 12px;
    color: #7c7c7c;
    margin-top: 4px;
  }

  .radio-option,
  .checkbox-option {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    margin-bottom: 8px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .radio-option:hover,
  .checkbox-option:hover {
    background: #f6f7f8;
    border-color: #ccc;
  }

  .radio-option input[type="radio"],
  .checkbox-option input[type="checkbox"] {
    margin-top: 2px;
    cursor: pointer;
  }

  .radio-content,
  .checkbox-content {
    flex: 1;
  }

  .radio-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
  }

  .radio-header svg {
    color: #878a8c;
  }

  .radio-title,
  .checkbox-title {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .radio-description,
  .checkbox-description {
    font-size: 12px;
    color: #7c7c7c;
    margin: 0;
  }

  .error-message {
    padding: 12px;
    background: #fee;
    border: 1px solid #fcc;
    border-radius: 4px;
    color: #c00;
    font-size: 14px;
    margin-bottom: 16px;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-top: 16px;
    border-top: 1px solid #edeff1;
  }

  .btn {
    padding: 10px 24px;
    border-radius: 24px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: transparent;
    color: var(--blue--);
    border: 1px solid var(--blue--);
  }

  .btn-secondary:hover:not(:disabled) {
    background: rgba(21, 48, 96, 0.08);
  }

  .btn-primary {
    background: var(--blue--);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--darkblue--);
  }

  /* Topics Section */
  .add-topics-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: transparent;
    border: 1px solid var(--blue--);
    border-radius: 24px;
    color: var(--blue--);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .add-topics-btn:hover:not(:disabled) {
    background: rgba(21, 48, 96, 0.08);
  }

  .add-topics-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .selected-topics {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 12px;
  }

  .topic-tag {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: #f6f7f8;
    border-radius: 16px;
    font-size: 13px;
    color: #1c1c1c;
    font-weight: 500;
  }

  .remove-topic {
    width: 16px;
    height: 16px;
    border: none;
    background: transparent;
    color: #878a8c;
    font-size: 18px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    line-height: 1;
    border-radius: 50%;
    transition: all 0.2s;
  }

  .remove-topic:hover {
    background: #e0e0e0;
    color: #1c1c1c;
  }

  /* Modal Styles */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 20px;
  }

  .modal-content {
    background: white;
    border-radius: 8px;
    width: 100%;
    max-width: 700px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 24px;
    border-bottom: 1px solid #edeff1;
  }

  .modal-header h2 {
    font-size: 18px;
    font-weight: 600;
    margin: 0;
    color: #1c1c1c;
  }

  .modal-close-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    font-size: 28px;
    color: #878a8c;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    padding: 0;
    line-height: 1;
  }

  .modal-close-btn:hover {
    background: #f6f7f8;
  }

  .modal-subtitle {
    padding: 12px 24px 0;
    font-size: 14px;
    color: #7c7c7c;
    margin: 0;
  }

  .search-wrapper {
    position: relative;
    margin: 16px 24px;
  }

  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: #878a8c;
  }

  .search-wrapper input {
    width: 100%;
    padding: 10px 12px 10px 40px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    background: #f6f7f8;
    font-size: 14px;
    color: #1c1c1c;
  }

  .search-wrapper input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .topics-counter {
    padding: 0 24px 12px;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .topics-list {
    flex: 1;
    overflow-y: auto;
    padding: 0 24px 16px;
  }

  .topic-category {
    margin-bottom: 20px;
  }

  .category-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }

  .category-icon {
    font-size: 18px;
  }

  .category-name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .subtopics {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .subtopic-tag {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 14px;
    background: #f6f7f8;
    border: 1px solid #edeff1;
    border-radius: 16px;
    font-size: 13px;
    color: #1c1c1c;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .subtopic-tag:hover:not(:disabled) {
    background: #e9ecef;
    border-color: #ccc;
  }

  .subtopic-tag.selected {
    background: var(--blue--);
    color: white;
    border-color: var(--blue--);
  }

  .subtopic-tag.selected:hover {
    background: var(--darkblue--);
    border-color: var(--darkblue--);
  }

  .subtopic-tag.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .subtopic-tag svg {
    flex-shrink: 0;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 24px;
    border-top: 1px solid #edeff1;
  }

  /* Style Modal */
  .style-content {
    padding: 24px;
    flex: 1;
    overflow-y: auto;
  }

  .upload-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0;
    border-bottom: 1px solid #edeff1;
  }

  .upload-label {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .upload-button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    background: #f6f7f8;
    border: 1px solid #edeff1;
    border-radius: 4px;
    color: #1c1c1c;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .upload-button:hover {
    background: #e9ecef;
  }

  .upload-button svg {
    color: #878a8c;
  }

  .preview-card {
    margin-top: 24px;
    border: 1px solid #edeff1;
    border-radius: 8px;
    overflow: hidden;
    background: white;
  }

  .preview-banner {
    width: 100%;
    height: 80px;
    background: linear-gradient(135deg, #ffc9aa 0%, #ffaec0 100%);
    background-size: cover;
    background-position: center;
  }

  .preview-content {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    position: relative;
  }

  .preview-icon-wrapper {
    width: 72px;
    height: 72px;
    margin-top: -36px;
    background: white;
    border-radius: 50%;
    border: 4px solid white;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .preview-icon-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .preview-icon-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f6f7f8;
  }

  .preview-info h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0 0 4px 0;
  }

  .preview-info p {
    font-size: 12px;
    color: #7c7c7c;
    margin: 0;
  }

  .preview-description {
    padding: 0 16px 16px;
    font-size: 14px;
    color: #1c1c1c;
    line-height: 1.5;
  }

  .progress-dots {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-top: 24px;
  }

  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #ccc;
  }

  .dot.active {
    background: #1c1c1c;
  }
</style>
