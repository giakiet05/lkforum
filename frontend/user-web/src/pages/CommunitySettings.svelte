<script lang="ts">
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import {
    getCommunities,
    updateCommunity,
  } from "../services/community-service";
  import type { CommunityResponse } from "../dtos/community-dto";
  import { authStore } from "../stores/auth-store";

  export interface Props {
    params?: { name?: string };
  }

  let { params }: Props = $props();

  const communityName = params?.name || "";

  let community = $state<CommunityResponse | null>(null);
  let isLoading = $state(true);
  let isSaving = $state(false);
  let error = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Settings state
  let postRequireApproval = $state(false);
  let joinRequireApproval = $state(false);
  let isPrivate = $state(false);
  let maxPostLength = $state(40000);

  // Appearance state
  let description = $state("");
  let avatarImage = $state<string>("");
  let bannerImage = $state<string>("");

  const currentUser = $derived($authStore.user);

  onMount(async () => {
    await loadCommunity();
  });

  async function loadCommunity() {
    try {
      isLoading = true;
      error = null;
      const response = await getCommunities({ limit: 100 });
      const found = response.communities.find((c) => c.name === communityName);

      if (found) {
        community = found;
        console.log("📥 Loaded community:", found);
        console.log("📋 Community settings:", found.setting);
        // Load current settings
        postRequireApproval = found.setting?.post_require_approval ?? false;
        joinRequireApproval = found.setting?.join_require_approval ?? false;
        isPrivate = found.setting?.is_private ?? false;
        maxPostLength = found.setting?.max_post_length ?? 40000;

        // Load appearance
        description = found.description || "";
        avatarImage = found.avatar || "";
        bannerImage = found.banner || "";
        console.log(
          "🔄 Local state - postRequireApproval:",
          postRequireApproval,
        );
      } else {
        error = "Không tìm thấy cộng đồng";
      }
    } catch (err) {
      console.error("Failed to load community:", err);
      error = "Không thể tải cài đặt cộng đồng";
    } finally {
      isLoading = false;
    }
  }

  function canManageSettings(): boolean {
    if (!community || !currentUser) return false;

    // Admin can manage everything
    if (currentUser.role === "admin") return true;

    // Check if user is creator
    if (community.create_by_id === currentUser.id) return true;

    // Check if user is moderator
    return (
      community.moderators?.some((mod) => mod.user_id === currentUser.id) ??
      false
    );
  }

  function handleAvatarUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      if (file.size > 2 * 1024 * 1024) {
        error = "Ảnh đại diện phải nhỏ hơn 2MB";
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        avatarImage = e.target?.result as string;
        error = null;
      };
      reader.readAsDataURL(file);
    }
  }

  function handleBannerUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      if (file.size > 100 * 1024 * 1024) {
        error = "Ảnh bìa phải nhỏ hơn 100MB";
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        bannerImage = e.target?.result as string;
        error = null;
      };
      reader.readAsDataURL(file);
    }
  }

  async function handleSaveSettings() {
    if (!community) return;

    try {
      isSaving = true;
      error = null;
      successMessage = null;

      const updateData = {
        id: community.id,
        description: description || undefined,
        avatar: avatarImage || undefined,
        banner: bannerImage || undefined,
        setting: {
          post_require_approval: postRequireApproval,
          join_require_approval: joinRequireApproval,
          is_private: isPrivate,
          max_post_length: maxPostLength,
        },
      };

      console.log("🔧 Saving settings:", updateData);
      const result = await updateCommunity(updateData);
      console.log("✅ Settings saved:", result);

      successMessage = "Đã lưu cài đặt thành công!";

      // Reload community to get fresh data
      await loadCommunity();
    } catch (err: any) {
      console.error("Failed to save settings:", err);
      error = err.message || "Không thể lưu cài đặt";
    } finally {
      isSaving = false;
    }
  }

  function handleBack() {
    push(`/lk/${communityName}`);
  }
</script>

<div class="settings-page">
  {#if isLoading}
    <div class="loading">Đang tải cài đặt...</div>
  {:else if error && !community}
    <div class="error-page">
      <p>{error}</p>
      <button onclick={handleBack}>Quay lại cộng đồng</button>
    </div>
  {:else if community}
    {#if !canManageSettings()}
      <div class="error-page">
        <p>Bạn không có quyền quản lý cài đặt cộng đồng</p>
        <button onclick={handleBack}>Quay lại cộng đồng</button>
      </div>
    {:else}
      <!-- Header -->
      <div class="settings-header">
        <button class="back-btn" onclick={handleBack}>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <path
              d="M19 12H5M12 19l-7-7 7-7"
              stroke="currentColor"
              stroke-width="2"
            />
          </svg>
          Quay lại
        </button>
        <h1>Cài đặt cộng đồng: lk/{community.name}</h1>
      </div>

      <!-- Settings Form -->
      <div class="settings-container">
        {#if successMessage}
          <div class="success-message">{successMessage}</div>
        {/if}
        {#if error}
          <div class="error-message">{error}</div>
        {/if}

        <!-- Appearance Section -->
        <div class="settings-section">
          <h2>Giao diện</h2>

          <div class="setting-item-column">
            <label>
              <strong>Mô tả</strong>
              <textarea
                bind:value={description}
                placeholder="Mô tả về cộng đồng của bạn..."
                maxlength="500"
                rows="4"
                class="textarea-input"
              ></textarea>
              <span class="hint">{description.length}/500 ký tự</span>
            </label>
          </div>

          <div class="setting-item-column">
            <label>
              <strong>Ảnh đại diện cộng đồng</strong>
              {#if avatarImage}
                <div class="image-preview">
                  <img
                    src={avatarImage}
                    alt="Xem trước ảnh đại diện"
                    class="avatar-preview"
                  />
                  <span
                    class="remove-image-link"
                    onclick={() => (avatarImage = "")}>Xóa</span
                  >
                </div>
              {/if}
              <input
                type="file"
                accept="image/*"
                onchange={handleAvatarUpload}
                class="file-input"
              />
              <span class="hint">Tối đa 2MB, khuyến nghị ảnh vuông</span>
            </label>
          </div>

          <div class="setting-item-column">
            <label>
              <strong>Ảnh bìa cộng đồng</strong>
              {#if bannerImage}
                <div class="image-preview">
                  <img
                    src={bannerImage}
                    alt="Xem trước ảnh bìa"
                    class="banner-preview"
                  />
                  <span
                    class="remove-image-link"
                    onclick={() => (bannerImage = "")}>Xóa</span
                  >
                </div>
              {/if}
              <input
                type="file"
                accept="image/*"
                onchange={handleBannerUpload}
                class="file-input"
              />
              <span class="hint"
                >Tối đa 100MB, khuyến nghị ảnh ngang (1600x400)</span
              >
            </label>
          </div>
        </div>

        <div class="settings-section">
          <h2>Cài đặt kiểm duyệt</h2>

          <div class="setting-item">
            <div class="setting-info">
              <h3>Duyệt bài viết thủ công</h3>
              <p>
                {#if postRequireApproval}
                  <strong>Đã bật:</strong> Tất cả bài viết phải được quản trị viên
                  duyệt trước khi hiển thị trong cộng đồng.
                {:else}
                  <strong>Đã tắt:</strong> Bài viết được kiểm duyệt tự động bởi AI.
                  Bài viết sạch sẽ hiển thị ngay, bài viết nghi vấn sẽ vào hàng chờ
                  duyệt.
                {/if}
              </p>
            </div>
            <label class="toggle">
              <input type="checkbox" bind:checked={postRequireApproval} />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <h3>Duyệt thành viên thủ công</h3>
              <p>
                {#if joinRequireApproval}
                  Thành viên phải được duyệt trước khi tham gia
                {:else}
                  Bất kỳ ai cũng có thể tham gia mà không cần duyệt
                {/if}
              </p>
            </div>
            <label class="toggle">
              <input type="checkbox" bind:checked={joinRequireApproval} />
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>

        <div class="settings-section">
          <h2>Cài đặt quyền riêng tư</h2>

          <div class="setting-item">
            <div class="setting-info">
              <h3>Cộng đồng riêng tư</h3>
              <p>
                {#if isPrivate}
                  Chỉ thành viên được duyệt mới có thể xem bài viết
                {:else}
                  Bất kỳ ai cũng có thể xem bài viết
                {/if}
              </p>
            </div>
            <label class="toggle">
              <input type="checkbox" bind:checked={isPrivate} />
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>

        <div class="settings-section">
          <h2>Cài đặt bài viết</h2>

          <div class="setting-item">
            <div class="setting-info">
              <h3>Độ dài bài viết tối đa</h3>
              <p>
                Số ký tự tối đa cho phép trong một bài viết (mặc định: 40.000)
              </p>
            </div>
            <input
              type="number"
              bind:value={maxPostLength}
              min="1000"
              max="100000"
              step="1000"
              class="number-input"
            />
          </div>
        </div>

        <div class="actions">
          <button
            class="save-btn"
            onclick={handleSaveSettings}
            disabled={isSaving}
          >
            {isSaving ? "Đang lưu..." : "Lưu cài đặt"}
          </button>
          <button class="cancel-btn" onclick={handleBack} disabled={isSaving}>
            Hủy
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>

<style>
  .settings-page {
    max-width: 900px;
    margin: 0 auto;
    padding: 20px;
  }

  .loading {
    text-align: center;
    padding: 40px;
    color: #666;
  }

  .error-page {
    text-align: center;
    padding: 40px;
  }

  .error-page p {
    color: #d32f2f;
    margin-bottom: 20px;
  }

  .settings-header {
    margin-bottom: 30px;
  }

  .back-btn {
    background: none;
    border: none;
    color: #0079d3;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    padding: 8px 0;
    margin-bottom: 16px;
  }

  .back-btn:hover {
    text-decoration: underline;
  }

  .settings-header h1 {
    font-size: 24px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0;
  }

  .settings-container {
    background: white;
    border-radius: 8px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .success-message {
    background: #d4edda;
    color: #155724;
    padding: 12px 16px;
    border-radius: 4px;
    margin-bottom: 20px;
    border: 1px solid #c3e6cb;
  }

  .error-message {
    background: #f8d7da;
    color: #721c24;
    padding: 12px 16px;
    border-radius: 4px;
    margin-bottom: 20px;
    border: 1px solid #f5c6cb;
  }

  .settings-section {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid #edeff1;
  }

  .settings-section:last-of-type {
    border-bottom: none;
  }

  .settings-section h2 {
    font-size: 18px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0 0 16px 0;
  }

  .setting-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 0;
    gap: 20px;
  }

  .setting-info {
    flex: 1;
  }

  .setting-info h3 {
    font-size: 16px;
    font-weight: 500;
    color: #1c1c1c;
    margin: 0 0 4px 0;
  }

  .setting-info p {
    font-size: 14px;
    color: #7c7c7c;
    margin: 0;
    line-height: 1.4;
  }

  .setting-info p strong {
    color: #1c1c1c;
  }

  .setting-item-column {
    display: flex;
    flex-direction: column;
    padding: 16px 0;
    gap: 8px;
  }

  .setting-item-column label {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .textarea-input {
    width: 100%;
    padding: 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 14px;
    font-family: inherit;
    resize: vertical;
    min-height: 100px;
  }

  .textarea-input:focus {
    outline: none;
    border-color: #0079d3;
  }

  .file-input {
    padding: 8px 0;
    font-size: 14px;
  }

  .hint {
    font-size: 12px;
    color: #7c7c7c;
  }

  .image-preview {
    display: flex;
    align-items: center;
    gap: 12px;
    margin: 8px 0;
  }

  .avatar-preview {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid #edeff1;
  }

  .banner-preview {
    width: 200px;
    height: 50px;
    border-radius: 4px;
    object-fit: cover;
    border: 2px solid #edeff1;
  }

  .remove-image-link {
    color: #ff4500;
    font-size: 13px;
    cursor: pointer;
    font-weight: 600;
    text-decoration: underline;
  }

  .remove-image-link:hover {
    color: #d93900;
  }

  /* Toggle Switch */
  .toggle {
    position: relative;
    display: inline-block;
    width: 48px;
    height: 24px;
    flex-shrink: 0;
  }

  .toggle input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #ccc;
    transition: 0.3s;
    border-radius: 24px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }

  .toggle input:checked + .toggle-slider {
    background-color: #0079d3;
  }

  .toggle input:checked + .toggle-slider:before {
    transform: translateX(24px);
  }

  /* Number Input */
  .number-input {
    width: 150px;
    padding: 8px 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
  }

  .number-input:focus {
    outline: none;
    border-color: #0079d3;
  }

  /* Actions */
  .actions {
    display: flex;
    gap: 12px;
    margin-top: 24px;
    padding-top: 24px;
    border-top: 1px solid #edeff1;
  }

  .save-btn,
  .cancel-btn {
    padding: 10px 24px;
    border-radius: 9999px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .save-btn {
    background: rgba(214, 216, 222, 0.5);
    color: #1c1c1c;
    border: none;
  }

  .save-btn:hover:not(:disabled) {
    background: rgba(214, 216, 222, 0.7);
  }

  .save-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .cancel-btn {
    background: #0079d3;
    color: white;
    border: none;
  }

  .cancel-btn:hover:not(:disabled) {
    background: #0060a8;
  }

  .cancel-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .back-btn {
    background: none;
    border: none;
    color: #1c1c1c;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    padding: 8px 12px;
    margin-bottom: 16px;
    border-radius: 20px;
    background: rgba(214, 216, 222, 0.5);
    transition: background-color 0.2s;
  }

  .back-btn:hover {
    background: rgba(214, 216, 222, 0.7);
  }
</style>
