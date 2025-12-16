<script lang="ts">
  import { push } from "svelte-spa-router";
  import { createCommunity } from "../services/community-service";
  import { toastStore } from "../stores/toast-store";
  import type {
    CreateCommunityRequest,
    CommunitySetting,
    CommunityResponse,
  } from "../dtos/community-dto";
  import { authStore } from "../stores/auth-store";

  interface Props {
    show: boolean;
    onClose: () => void;
  }

  let { show = false, onClose }: Props = $props();

  let currentStep = $state(1); // 1: form, 2: style (removed topics step)
  let communityName = $state("");
  let description = $state("");
  let communityType = $state<"public" | "restricted" | "private">("public");
  let isAdultContent = $state(false);
  let bannerImage = $state<string>("");
  let iconImage = $state<string>("");
  let isLoading = $state(false);
  let error = $state("");

  // Get current user from auth store
  let user = $state<any>(null);

  authStore.subscribe((state) => {
    user = state.user;
  });

  function handleBannerUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      // Check file size (max 5MB)
      if (file.size > 5 * 1024 * 1024) {
        error = "Banner image must be less than 5MB";
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        bannerImage = e.target?.result as string;
        error = ""; // Clear error
      };
      reader.readAsDataURL(file);
    }
  }

  function handleIconUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      // Check file size (max 2MB)
      if (file.size > 2 * 1024 * 1024) {
        error = "Icon image must be less than 2MB";
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        iconImage = e.target?.result as string;
        error = ""; // Clear error
      };
      reader.readAsDataURL(file);
    }
  }

  function validateStep1() {
    if (!communityName.trim()) {
      error = "Community name is required";
      return false;
    }
    if (communityName.length < 3) {
      error = "Community name must be at least 3 characters";
      return false;
    }
    if (communityName.length > 50) {
      error = "Community name must be less than 50 characters";
      return false;
    }
    if (!/^[a-zA-Z0-9_]+$/.test(communityName)) {
      error =
        "Community name can only contain letters, numbers, and underscores";
      return false;
    }
    if (!description.trim()) {
      error = "Mô tả cộng đồng là bắt buộc";
      return false;
    }
    if (description.length > 500) {
      error = "Mô tả không được vượt quá 500 ký tự";
      return false;
    }
    error = "";
    return true;
  }

  function handleStep1Next() {
    if (validateStep1()) {
      currentStep = 2; // Go directly to style step
    }
  }

  function handleBack() {
    if (currentStep > 1) {
      currentStep--;
    }
  }

  function handleClose() {
    // Reset all states
    currentStep = 1;
    communityName = "";
    description = "";
    communityType = "public";
    bannerImage = "";
    iconImage = "";
    error = "";
    onClose();
  }

  /**
   * Map community type to backend CommunitySetting
   */
  function getCommunitySettings(): CommunitySetting {
    switch (communityType) {
      case "public":
        return {
          is_private: false,
          post_require_approval: false,
          join_require_approval: false,
          max_post_length: 40000, // default
        };
      case "restricted":
        return {
          is_private: false,
          post_require_approval: true,
          join_require_approval: true,
          max_post_length: 40000,
        };
      case "private":
        return {
          is_private: true,
          post_require_approval: true,
          join_require_approval: true,
          max_post_length: 40000,
        };
      default:
        return {
          is_private: false,
          post_require_approval: false,
          join_require_approval: false,
          max_post_length: 40000,
        };
    }
  }

  async function handleSubmit() {
    try {
      isLoading = true;
      error = "";

      // Validate base64 image sizes before sending
      if (iconImage && iconImage.length > 500000) {
        error =
          "Icon image is too large. Please use a smaller image (max ~375KB).";
        currentStep = 1; // Go back to show error
        return;
      }
      if (bannerImage && bannerImage.length > 1000000) {
        error =
          "Banner image is too large. Please use a smaller image (max ~750KB).";
        currentStep = 1; // Go back to show error
        return;
      }

      const requestData: CreateCommunityRequest = {
        name: communityName,
        description: description, // Now required, not optional
        avatar: iconImage || undefined,
        banner: bannerImage || undefined,
        setting: getCommunitySettings(),
        creator_name: user?.username,
        creator_avatar: user?.profile?.avatar?.url,
        is_18_plus: isAdultContent,
      };

      const result = await createCommunity(requestData);

      toastStore.success(`Community "lk/${result.name}" created successfully!`);
      handleClose();
      push(`/lk/${result.name}`);
    } catch (err: any) {
      // Handle specific error cases
      if (
        err.message?.includes("already exists") ||
        err.message?.includes("đã tồn tại")
      ) {
        error =
          "Community name already exists. Please choose a different name.";
      } else if (
        err.message?.includes("server") ||
        err.message?.includes("500")
      ) {
        error =
          "Server error occurred. Try reducing image sizes or try again later.";
      } else {
        error = err.message || "Failed to create community. Please try again.";
      }

      // Go back to step 1 to show error
      currentStep = 1;
      console.error("Create community error:", err);
    } finally {
      isLoading = false;
    }
  }
</script>

{#if show}
  <div class="modal-overlay" onclick={handleClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <!-- Step 1: Create Community Form -->
      {#if currentStep === 1}
        <div class="modal-header">
          <h2>Create a community</h2>
          <button class="modal-close-btn" onclick={handleClose}>×</button>
        </div>

        <div class="modal-body">
          <!-- Community Name -->
          <div class="form-section">
            <label for="name" class="label">Name</label>
            <p class="help-text">
              Community names including capitalization cannot be changed.
            </p>
            <div class="input-wrapper">
              <span class="prefix">lk/</span>
              <input
                type="text"
                id="name"
                bind:value={communityName}
                placeholder="community_name"
                maxlength="50"
              />
            </div>
            <div class="char-count">{communityName.length} / 50</div>
          </div>

          <!-- Description -->
          <div class="form-section">
            <label for="description" class="label">Description (optional)</label
            >
            <textarea
              id="description"
              bind:value={description}
              placeholder="What is your community about?"
              rows="4"
              maxlength="500"
            ></textarea>
            <div class="char-count">{description.length} / 500</div>
          </div>

          <!-- Community Type -->
          <div class="form-section">
            <label class="label">Community type</label>

            <label class="radio-option">
              <input
                type="radio"
                name="type"
                value="public"
                bind:group={communityType}
              />
              <div class="radio-content">
                <div class="radio-header">
                  <img
                    src="/material-symbols_public.svg"
                    alt="Public"
                    width="20"
                    height="20"
                  />
                  <span class="radio-title">Public</span>
                </div>
                <p class="radio-description">
                  Anyone can view, post, and comment to this community
                </p>
              </div>
            </label>

            <label class="radio-option">
              <input
                type="radio"
                name="type"
                value="restricted"
                bind:group={communityType}
              />
              <div class="radio-content">
                <div class="radio-header">
                  <img
                    src="/carbon_navaid-private.svg"
                    alt="Restricted"
                    width="20"
                    height="20"
                  />
                  <span class="radio-title">Restricted</span>
                </div>
                <p class="radio-description">
                  Anyone can view this community, but only approved users can
                  post
                </p>
              </div>
            </label>

            <label class="radio-option">
              <input
                type="radio"
                name="type"
                value="private"
                bind:group={communityType}
              />
              <div class="radio-content">
                <div class="radio-header">
                  <img
                    src="/simple-icons_privateinternetaccess.svg"
                    alt="Private"
                    width="20"
                    height="20"
                  />
                  <span class="radio-title">Private</span>
                </div>
                <p class="radio-description">
                  Only approved users can view and submit to this community
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
                <div class="checkbox-header">
                  <img
                    src="/stash_sensitive.svg"
                    alt="18+"
                    width="20"
                    height="20"
                  />
                  <span class="checkbox-title">18+ year old community</span>
                </div>
                <p class="checkbox-description">
                  Must be over 18 to view and contribute
                </p>
              </div>
            </label>
          </div>

          {#if error}
            <div class="error-message">{error}</div>
          {/if}

          <div class="progress-dots">
            <span class="dot active"></span>
            <span class="dot"></span>
          </div>
        </div>

        <div class="modal-actions">
          <button type="button" class="btn btn-secondary" onclick={handleClose}
            >Cancel</button
          >
          <button
            type="button"
            class="btn btn-primary"
            onclick={handleStep1Next}
            disabled={!communityName.trim()}
          >
            Next
          </button>
        </div>
      {/if}

      <!-- Step 2: Style -->
      {#if currentStep === 2}
        <div class="modal-header">
          <h2>Style your community</h2>
          <button class="modal-close-btn" onclick={handleClose}>×</button>
        </div>

        <p class="modal-subtitle">
          Adding visual flair will catch new members attention and help
          establish your community's culture! You can update this at any time.
        </p>

        <div class="style-content">
          <!-- Banner Upload -->
          <div class="upload-section">
            <label class="upload-label">Banner</label>
            <label class="upload-button">
              <input
                type="file"
                accept="image/*"
                onchange={handleBannerUpload}
                style="display: none;"
              />
              <img
                src="/hugeicons_image-upload.svg"
                alt="Upload"
                width="20"
                height="20"
              />
              Add
            </label>
          </div>

          <!-- Icon Upload -->
          <div class="upload-section">
            <label class="upload-label">Icon</label>
            <label class="upload-button">
              <input
                type="file"
                accept="image/*"
                onchange={handleIconUpload}
                style="display: none;"
              />
              <img
                src="/hugeicons_image-upload.svg"
                alt="Upload"
                width="20"
                height="20"
              />
              Add
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
                      width="64"
                      height="64"
                      viewBox="0 0 64 64"
                      fill="currentColor"
                    >
                      <circle cx="32" cy="32" r="32" fill="#ff4500" />
                      <text
                        x="32"
                        y="40"
                        text-anchor="middle"
                        fill="white"
                        font-size="20"
                        font-weight="bold"
                      >
                        lk/
                      </text>
                    </svg>
                  </div>
                {/if}
              </div>
              <div class="preview-info">
                <h3>lk/{communityName || "community"}</h3>
                <p>1 member · 1 online</p>
              </div>
            </div>
            <div class="preview-description">
              {description || "Community description"}
            </div>
          </div>

          <div class="progress-dots">
            <span class="dot"></span>
            <span class="dot active"></span>
          </div>
        </div>

        <div class="modal-actions">
          <button type="button" class="btn btn-secondary" onclick={handleBack}
            >Back</button
          >
          <button
            type="button"
            class="btn btn-primary"
            onclick={handleSubmit}
            disabled={isLoading}
          >
            {isLoading ? "Creating..." : "Create Community"}
          </button>
        </div>
      {/if}
    </div>
  </div>
{/if}

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
    max-height: 85vh;
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
    flex-shrink: 0;
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
    flex-shrink: 0;
  }

  .modal-body {
    padding: 24px;
    overflow-y: auto;
    flex: 1;
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

  .checkbox-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
  }

  .checkbox-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
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

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 24px;
    border-top: 1px solid #edeff1;
    flex-shrink: 0;
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

  /* Step 2: Topics */
  .search-wrapper {
    position: relative;
    margin: 16px 0;
    padding: 0 24px;
    flex-shrink: 0;
  }

  .search-icon {
    position: absolute;
    left: 36px;
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
    box-sizing: border-box;
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
    flex-shrink: 0;
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

  /* Step 3: Style */
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
</style>
