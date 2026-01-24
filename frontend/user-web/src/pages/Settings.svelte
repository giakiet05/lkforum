<script lang="ts">
  import { onMount } from "svelte";
  import type { UserResponse, SettingsResponse } from "../dtos/user-dto";
  import {
    getMyProfile,
    updateProfile,
    changePassword,
    uploadAvatar,
    deleteAvatar,
    getSettings,
    updateSettings,
    getProvinces,
    getInterests,
    getGenders,
  } from "../services/user-service";
  import { ApiError } from "../errors/api-error";
  import { setAuth } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
  import ConfirmModal from "../components/ConfirmModal.svelte";

  let activeTab = $state<
    "account" | "privacy" | "notifications" | "appearance"
  >("account");

  let user = $state<UserResponse | null>(null);
  let settings = $state<SettingsResponse | null>(null);
  let provinces = $state<string[]>([]);
  let allInterests = $state<string[]>([]);
  let genders = $state<string[]>([]);

  let isLoadingUser = $state(true);
  let isLoadingSettings = $state(false);
  let isSaving = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Account form
  let editedBio = $state("");
  let editedGender = $state("");
  let editedDateOfBirth = $state("");
  let editedLocation = $state("");
  let editedInterests = $state<string[]>([]);
  let editedWebsite = $state("");
  let editedFacebook = $state("");
  let editedYouTube = $state("");
  let editedGitHub = $state("");

  // Password modal
  let showPasswordModal = $state(false);
  let oldPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");

  // Settings form
  let editedSettings = $state<SettingsResponse | null>(null);

  let avatarFileInput: HTMLInputElement;
  let isUploadingAvatar = $state(false);
  let isDeletingAvatar = $state(false);
  let showDeleteAvatarConfirm = $state(false);
  let showDeleteAccountModal = $state(false);
  let deleteAccountVerification = $state("");

  onMount(() => {
    loadUserProfile();
    loadMetadata();
  });

  async function loadUserProfile() {
    try {
      isLoadingUser = true;
      errorMessage = null;
      user = await getMyProfile();

      // Populate form
      editedBio = user.profile.bio || "";
      editedGender = user.profile.gender || "";
      editedLocation = user.profile.location || "";
      editedInterests = user.profile.interests || [];
      editedWebsite = user.profile.social_links?.website || "";
      editedFacebook = user.profile.social_links?.facebook || "";
      editedYouTube = user.profile.social_links?.youtube || "";
      editedGitHub = user.profile.social_links?.github || "";
    } catch (error) {
      console.error("Failed to load profile:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Không thể tải hồ sơ.";
      }
    } finally {
      isLoadingUser = false;
    }
  }

  async function loadMetadata() {
    try {
      const [provincesData, interestsData, gendersData] = await Promise.all([
        getProvinces(),
        getInterests(),
        getGenders(),
      ]);
      provinces = provincesData;
      allInterests = interestsData;
      genders = gendersData;
    } catch (error) {
      console.error("Failed to load metadata:", error);
    }
  }

  async function loadSettings() {
    if (settings) return;

    try {
      isLoadingSettings = true;
      settings = await getSettings();
      editedSettings = JSON.parse(JSON.stringify(settings));
    } catch (error) {
      console.error("Failed to load settings:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      }
    } finally {
      isLoadingSettings = false;
    }
  }

  $effect(() => {
    if (
      (activeTab === "privacy" ||
        activeTab === "notifications" ||
        activeTab === "appearance") &&
      !settings
    ) {
      loadSettings();
    }
  });

  async function handleSaveAccount() {
    if (!user) return;

    try {
      isSaving = true;
      errorMessage = null;
      successMessage = null;

      const payload: any = {};

      if (editedBio !== (user.profile.bio || ""))
        payload.bio = editedBio || null;
      if (editedGender !== (user.profile.gender || ""))
        payload.gender = editedGender || null;
      if (editedDateOfBirth) payload.date_of_birth = editedDateOfBirth;
      if (editedLocation !== (user.profile.location || ""))
        payload.location = editedLocation || null;
      if (
        JSON.stringify(editedInterests) !==
        JSON.stringify(user.profile.interests || [])
      ) {
        payload.interests = editedInterests.length > 0 ? editedInterests : null;
      }

      const currentLinks = user.profile.social_links || {};
      const newLinks: any = {};
      let hasLinkChanges = false;

      if (editedWebsite !== (currentLinks.website || "")) {
        newLinks.website = editedWebsite || null;
        hasLinkChanges = true;
      }
      if (editedFacebook !== (currentLinks.facebook || "")) {
        newLinks.facebook = editedFacebook || null;
        hasLinkChanges = true;
      }
      if (editedYouTube !== (currentLinks.youtube || "")) {
        newLinks.youtube = editedYouTube || null;
        hasLinkChanges = true;
      }
      if (editedGitHub !== (currentLinks.github || "")) {
        newLinks.github = editedGitHub || null;
        hasLinkChanges = true;
      }

      if (hasLinkChanges) payload.social_links = newLinks;

      if (Object.keys(payload).length === 0) {
        successMessage = "Không có thay đổi để lưu";
        return;
      }

      const updatedUser = await updateProfile(payload);
      updateUserState(updatedUser); // Đồng bộ với authStore
      successMessage = "Đã cập nhật hồ sơ!";
    } catch (error) {
      console.error("Failed to update profile:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Không thể cập nhật hồ sơ.";
      }
    } finally {
      isSaving = false;
    }
  }

  async function handleChangePassword() {
    if (!oldPassword || !newPassword || !confirmPassword) {
      errorMessage = "Vui lòng điền đầy đủ các trường mật khẩu";
      return;
    }

    if (newPassword !== confirmPassword) {
      errorMessage = "Mật khẩu mới không khớp";
      return;
    }

    if (newPassword.length < 6) {
      errorMessage = "Mật khẩu mới phải có ít nhất 6 ký tự";
      return;
    }

    try {
      isSaving = true;
      errorMessage = null;
      successMessage = null;

      await changePassword({
        old_password: oldPassword,
        new_password: newPassword,
      });

      successMessage = "Đã đổi mật khẩu thành công!";
      showPasswordModal = false;
      oldPassword = "";
      newPassword = "";
      confirmPassword = "";
    } catch (error) {
      console.error("Failed to change password:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Không thể đổi mật khẩu.";
      }
    } finally {
      isSaving = false;
    }
  }

  async function autoSaveSettings() {
    if (!editedSettings) return;

    try {
      // Save in background without blocking UI
      settings = await updateSettings(editedSettings);
      editedSettings = JSON.parse(JSON.stringify(settings));

      // Show brief success feedback
      successMessage = "Đã lưu";
      setTimeout(() => {
        if (successMessage === "Đã lưu") successMessage = null;
      }, 2000);
    } catch (error) {
      console.error("Failed to save settings:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Không thể lưu cài đặt.";
      }
    }
  }

  function handleDeleteAccount() {
    showDeleteAccountModal = true;
    deleteAccountVerification = "";
  }

  function confirmDeleteAccount() {
    if (deleteAccountVerification !== "DELETE") {
      toastStore.warning("Hủy xóa tài khoản. Không khớp mã xác nhận.");
      showDeleteAccountModal = false;
      return;
    }

    showDeleteAccountModal = false;

    // TODO: Call delete account API when available
    toastStore.info(
      "Tính năng xóa tài khoản chưa được triển khai.\n\nAPI yêu cầu: DELETE /api/users/me",
    );

    // Future implementation:
    // try {
    //   await deleteAccount();
    //   clearAuth();
    //   localStorage.clear();
    //   push("/");
    // } catch (error) {
    //   errorMessage = "Failed to delete account.";
    // }
  }

  // Helper function để update user state và đồng bộ với authStore + localStorage
  function updateUserState(updatedUser: UserResponse) {
    user = updatedUser;
    // Update authStore để sync với topbar và các component khác
    setAuth(updatedUser);
    // Update localStorage
    localStorage.setItem("user", JSON.stringify(updatedUser));
  }

  async function handleAvatarChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      errorMessage = "Vui lòng chọn tệp hình ảnh";
      return;
    }

    if (file.size > 100 * 1024 * 1024) {
      errorMessage = "Kích thước ảnh phải nhỏ hơn 100MB";
      return;
    }

    try {
      isUploadingAvatar = true;
      errorMessage = null;
      const updatedUser = await uploadAvatar(file);
      updateUserState(updatedUser); // Đồng bộ với authStore
      successMessage = "Đã tải ảnh đại diện lên!";
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Không thể tải ảnh đại diện.";
      }
    } finally {
      isUploadingAvatar = false;
      input.value = "";
    }
  }

  async function handleDeleteAvatar() {
    showDeleteAvatarConfirm = true;
  }

  async function confirmDeleteAvatar() {
    showDeleteAvatarConfirm = false;

    try {
      isDeletingAvatar = true;
      errorMessage = null;
      const updatedUser = await deleteAvatar();
      updateUserState(updatedUser); // Đồng bộ với authStore
      successMessage = "Đã xóa ảnh đại diện!";
    } catch (error) {
      console.error("Failed to delete avatar:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Không thể xóa ảnh đại diện.";
      }
    } finally {
      isDeletingAvatar = false;
    }
  }

  function toggleInterest(interest: string) {
    if (editedInterests.includes(interest)) {
      editedInterests = editedInterests.filter((i) => i !== interest);
    } else {
      if (editedInterests.length >= 10) {
        errorMessage = "Tối đa 10 sở thích";
        return;
      }
      editedInterests = [...editedInterests, interest];
    }
  }
</script>

<div class="settings-page">
  <input
    type="file"
    accept="image/*"
    bind:this={avatarFileInput}
    onchange={handleAvatarChange}
    style="display: none;"
  />

  {#if errorMessage}
    <div class="alert alert-error">
      {errorMessage}
      <button class="alert-close" onclick={() => (errorMessage = null)}
        >×</button
      >
    </div>
  {/if}

  {#if successMessage}
    <div class="alert alert-success">
      {successMessage}
      <button class="alert-close" onclick={() => (successMessage = null)}
        >×</button
      >
    </div>
  {/if}

  <div class="settings-container">
    <div class="settings-header">
      <h1>Cài đặt</h1>
      <p class="settings-description">
        Quản lý cài đặt tài khoản và tùy chọn của bạn
      </p>
    </div>

    <div class="settings-content">
      <div class="settings-sidebar">
        <button
          class="settings-tab"
          class:active={activeTab === "account"}
          onclick={() => (activeTab = "account")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <circle
              cx="10"
              cy="6"
              r="3"
              stroke="currentColor"
              stroke-width="1.5"
            />
            <path
              d="M4 18C4 14.6863 6.68629 12 10 12C13.3137 12 16 14.6863 16 18"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
          Tài khoản
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "privacy"}
          onclick={() => (activeTab = "privacy")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 2L4 5V9C4 13 6.5 16.5 10 18C13.5 16.5 16 13 16 9V5L10 2Z"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          Quyền riêng tư
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "notifications"}
          onclick={() => (activeTab = "notifications")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 3C7.79 3 6 4.79 6 7V10L4 12V13H16V12L14 10V7C14 4.79 12.21 3 10 3Z"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M8.5 16C8.5 16.8284 9.17157 17.5 10 17.5C10.8284 17.5 11.5 16.8284 11.5 16"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
          Thông báo
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "appearance"}
          onclick={() => activeTab === "appearance"}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <circle
              cx="10"
              cy="10"
              r="3"
              stroke="currentColor"
              stroke-width="1.5"
            />
            <path
              d="M10 1v2M10 17v2M3.93 3.93l1.41 1.41M14.66 14.66l1.41 1.41M1 10h2M17 10h2M3.93 16.07l1.41-1.41M14.66 5.34l1.41-1.41"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
          Giao diện
        </button>
      </div>

      <div class="settings-main">
        {#if isLoadingUser}
          <div class="loading-state">
            <div class="spinner"></div>
            <p>Đang tải...</p>
          </div>
        {:else if activeTab === "account" && user}
          <div class="settings-section">
            <h2>Cài đặt tài khoản</h2>
            <p class="section-description">Quản lý thông tin tài khoản</p>

            <div class="form-group">
              <span class="form-label">Ảnh đại diện</span>
              <div class="avatar-upload">
                <div class="avatar-preview">
                  {#if user.profile.avatar?.url}
                    <img src={user.profile.avatar.url} alt="Avatar" />
                  {:else}
                    <div class="avatar-placeholder">
                      {user.username[0].toUpperCase()}
                    </div>
                  {/if}
                </div>
                <div class="avatar-actions">
                  <button
                    class="btn-secondary"
                    onclick={() => avatarFileInput.click()}
                    disabled={isUploadingAvatar || isDeletingAvatar}
                  >
                    {isUploadingAvatar ? "Đang tải..." : "Đổi ảnh đại diện"}
                  </button>
                  {#if user.profile.avatar?.url}
                    <button
                      class="btn-text btn-remove"
                      onclick={handleDeleteAvatar}
                      disabled={isUploadingAvatar || isDeletingAvatar}
                    >
                      {isDeletingAvatar ? "Đang xóa..." : "Xóa"}
                    </button>
                  {/if}
                </div>
              </div>
            </div>

            <div class="form-group">
              <label for="username-label">Tên người dùng</label>
              <div class="input-with-prefix">
                <span class="input-prefix">u/</span>
                <input
                  type="text"
                  id="username-label"
                  value={user.username}
                  disabled
                />
              </div>
              <p class="input-hint">Tên người dùng không thể thay đổi</p>
            </div>

            <div class="form-group">
              <label for="email-label">Email</label>
              <input
                type="email"
                id="email-label"
                value={user.email}
                disabled
              />
              <p class="input-hint">Email không thể thay đổi</p>
            </div>

            <div class="form-group">
              <label for="bio">Tiểu sử</label>
              <textarea
                id="bio"
                rows="4"
                bind:value={editedBio}
                placeholder="Kể về bản thân bạn"
                maxlength="500"
              ></textarea>
              <p class="input-hint">{editedBio.length}/500 ký tự</p>
            </div>

            <div class="form-group">
              <label for="gender">Giới tính</label>
              <select id="gender" bind:value={editedGender}>
                <option value="">Chọn giới tính</option>
                {#each genders as gender}
                  <option value={gender}
                    >{gender === "male"
                      ? "Nam"
                      : gender === "female"
                        ? "Nữ"
                        : "Không tiết lộ"}</option
                  >
                {/each}
              </select>
            </div>

            <div class="form-group">
              <label for="dob">Ngày sinh</label>
              <input type="date" id="dob" bind:value={editedDateOfBirth} />
              <p class="input-hint">Bạn phải ít nhất 13 tuổi</p>
            </div>

            <div class="form-group">
              <label for="location">Địa điểm</label>
              <select id="location" bind:value={editedLocation}>
                <option value="">Chọn tỉnh/thành</option>
                {#each provinces as province}
                  <option value={province}>{province}</option>
                {/each}
              </select>
            </div>

            <div class="form-group">
              <label for="interests-label">Sở thích (tối đa 10)</label>
              <div class="interests-grid" id="interests-label">
                {#each allInterests as interest}
                  <button
                    type="button"
                    class="interest-btn"
                    class:selected={editedInterests.includes(interest)}
                    onclick={() => toggleInterest(interest)}
                  >
                    {interest}
                  </button>
                {/each}
              </div>
              <p class="input-hint">{editedInterests.length}/10 đã chọn</p>
            </div>

            <div class="form-group">
              <label for="social-website">Liên kết xã hội</label>
              <input
                type="url"
                id="social-website"
                bind:value={editedWebsite}
                placeholder="URL Website"
              />
            </div>
            <div class="form-group">
              <label for="social-facebook" class="sr-only">Facebook</label>
              <input
                type="text"
                id="social-facebook"
                bind:value={editedFacebook}
                placeholder="Tên đăng nhập hoặc URL Facebook"
              />
            </div>
            <div class="form-group">
              <label for="social-youtube" class="sr-only">YouTube</label>
              <input
                type="text"
                id="social-youtube"
                bind:value={editedYouTube}
                placeholder="URL kênh YouTube"
              />
            </div>
            <div class="form-group">
              <label for="social-github" class="sr-only">GitHub</label>
              <input
                type="text"
                id="social-github"
                bind:value={editedGitHub}
                placeholder="Tên đăng nhập GitHub"
              />
            </div>

            <div class="form-actions">
              <button
                class="btn-primary"
                onclick={handleSaveAccount}
                disabled={isSaving}
              >
                {isSaving ? "Đang lưu..." : "Lưu thay đổi"}
              </button>
            </div>

            <div class="password-section">
              <h3>Mật khẩu</h3>
              <button
                class="btn-secondary"
                onclick={() => (showPasswordModal = true)}>Đổi mật khẩu</button
              >
            </div>

            <!-- Delete Account Section -->
            <div class="danger-zone">
              <h3>Vùng nguy hiểm</h3>
              <div class="danger-content">
                <div>
                  <h4>Xóa tài khoản</h4>
                  <p>
                    Xóa vĩnh viễn tài khoản và tất cả nội dung. Hành động này
                    không thể hoàn tác.
                  </p>
                </div>
                <button
                  class="btn-danger"
                  onclick={() => handleDeleteAccount()}
                >
                  Xóa tài khoản
                </button>
              </div>
            </div>
          </div>
        {:else if activeTab === "privacy" && editedSettings}
          <div class="settings-section">
            <h2>Cài đặt quyền riêng tư</h2>
            <p class="section-description">Kiểm soát tùy chọn quyền riêng tư</p>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Hiển thị hồ sơ</h4>
                <p>Cho phép người khác xem hồ sơ của bạn</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.show_profile}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Hiển thị Email</h4>
                <p>Hiển thị email trên hồ sơ của bạn</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.show_email}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Hiển thị lịch sử bài viết</h4>
                <p>Cho phép người khác xem lịch sử bài viết của bạn</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.show_post_history}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Cho phép tin nhắn trực tiếp</h4>
                <p>Cho phép người khác gửi tin nhắn cho bạn</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.allow_direct_messages}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Cho phép nhắc đến</h4>
                <p>
                  Cho phép người khác nhắc đến bạn trong bài viết và bình luận
                </p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.allow_mentions}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
        {:else if activeTab === "notifications" && editedSettings}
          <div class="settings-section">
            <h2>Cài đặt thông báo</h2>
            <p class="section-description">Quản lý cách bạn nhận thông báo</p>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Thông báo trong ứng dụng</h4>
                <p>Hiển thị thông báo trên website</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.in_app_enabled}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Thông báo Email</h4>
                <p>Gửi thông báo đến email của bạn</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.email_enabled}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <h3 class="subsection-title">Thông báo cho tôi khi:</h3>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Có người bình luận bài viết của tôi</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_comment}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Có người nhắc đến tôi</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_mention}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Có người upvote bài viết của tôi</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_upvote}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Có người gửi tin nhắn cho tôi</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_message}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
        {:else if activeTab === "appearance" && editedSettings}
          <div class="settings-section">
            <h2>Cài đặt giao diện</h2>
            <p class="section-description">Tùy chỉnh giao diện</p>

            <div class="form-group">
              <label for="theme">Giao diện</label>
              <select
                id="theme"
                bind:value={editedSettings.appearance.theme}
                onchange={autoSaveSettings}
              >
                <option value="light">Sáng</option>
                <option value="dark">Tối</option>
                <option value="auto">Tự động</option>
              </select>
            </div>

            <div class="form-group">
              <label for="fontSize">Cỡ chữ</label>
              <select
                id="fontSize"
                bind:value={editedSettings.appearance.font_size}
                onchange={autoSaveSettings}
              >
                <option value="small">Nhỏ</option>
                <option value="medium">Vừa</option>
                <option value="large">Lớn</option>
              </select>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Cho phép nội dung NSFW</h4>
                <p>Hiển thị nội dung giới hạn độ tuổi</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.content.allow_nsfw}
                  onchange={autoSaveSettings}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
        {:else if isLoadingSettings}
          <div class="loading-state">
            <div class="spinner"></div>
            <p>Đang tải cài đặt...</p>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<!-- Password Change Modal -->
{#if showPasswordModal}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="modal-overlay" onclick={() => (showPasswordModal = false)}>
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>Đổi mật khẩu</h2>
      <div class="form-group">
        <label for="oldPassword">Mật khẩu hiện tại</label>
        <input type="password" id="oldPassword" bind:value={oldPassword} />
      </div>
      <div class="form-group">
        <label for="newPassword">Mật khẩu mới</label>
        <input type="password" id="newPassword" bind:value={newPassword} />
      </div>
      <div class="form-group">
        <label for="confirmPassword">Xác nhận mật khẩu mới</label>
        <input
          type="password"
          id="confirmPassword"
          bind:value={confirmPassword}
        />
      </div>
      <div class="modal-actions">
        <button
          class="btn-primary"
          onclick={handleChangePassword}
          disabled={isSaving}
        >
          {isSaving ? "Đang đổi..." : "Đổi mật khẩu"}
        </button>
        <button
          class="btn-secondary"
          onclick={() => (showPasswordModal = false)}>Hủy</button
        >
      </div>
    </div>
  </div>
{/if}

<ConfirmModal
  show={showDeleteAvatarConfirm}
  title="Xác nhận xóa"
  message="Bạn có chắc muốn xóa ảnh đại diện?"
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDeleteAvatar}
  onCancel={() => (showDeleteAvatarConfirm = false)}
/>

<!-- Delete Account Modal -->
{#if showDeleteAccountModal}
  <div class="modal-overlay" onclick={() => (showDeleteAccountModal = false)}>
    <div
      class="modal-content delete-account-modal"
      onclick={(e) => e.stopPropagation()}
    >
      <div class="modal-header">
        <h3>Xóa tài khoản</h3>
        <button
          class="close-btn"
          onclick={() => (showDeleteAccountModal = false)}>×</button
        >
      </div>
      <div class="modal-body">
        <div class="warning-text">
          <p><strong>⚠️ Cảnh báo:</strong> Hành động này không thể hoàn tác!</p>
          <p>Việc xóa tài khoản sẽ xóa vĩnh viễn:</p>
          <ul>
            <li>Tất cả bài viết và bình luận</li>
            <li>Hồ sơ và cài đặt</li>
            <li>Nội dung đã lưu</li>
            <li>Lịch sử hoạt động</li>
          </ul>
        </div>
        <div class="verification-input">
          <label for="delete-verification"
            >Nhập <strong>DELETE</strong> để xác nhận:</label
          >
          <input
            type="text"
            id="delete-verification"
            bind:value={deleteAccountVerification}
            placeholder="DELETE"
          />
        </div>
      </div>
      <div class="modal-footer">
        <button
          class="btn-cancel"
          onclick={() => (showDeleteAccountModal = false)}>Hủy</button
        >
        <button
          class="btn-danger"
          onclick={confirmDeleteAccount}
          disabled={deleteAccountVerification !== "DELETE"}
        >
          Xóa tài khoản
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .settings-page {
    background-color: white;
    min-height: 100vh;
    padding-top: 72px;
  }

  .alert {
    position: fixed;
    top: 80px;
    right: 20px;
    padding: 1rem 3rem 1rem 1.5rem;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 1000;
    animation: slideIn 0.3s ease;
    max-width: 400px;
  }

  .alert-error {
    background-color: #fee;
    color: #c00;
    border: 1px solid #fcc;
  }

  .alert-success {
    position: fixed;
    top: 1rem;
    right: 1rem;
    background-color: #10b981;
    color: white;
    border: none;
    padding: 0.75rem 2rem 0.75rem 1rem;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
    animation:
      slideInRight 0.3s ease,
      fadeOut 0.3s ease 1.7s;
    z-index: 1000;
    font-size: 0.9rem;
    font-weight: 500;
  }

  @keyframes slideInRight {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  @keyframes fadeOut {
    from {
      opacity: 1;
    }
    to {
      opacity: 0;
    }
  }

  .alert-close {
    position: absolute;
    top: 0.5rem;
    right: 0.75rem;
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: inherit;
    opacity: 0.6;
  }

  .alert-close:hover {
    opacity: 1;
  }

  @keyframes slideIn {
    from {
      transform: translateX(100%);
    }
    to {
      transform: translateX(0);
    }
  }

  .settings-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1.5rem 2rem;
  }

  .settings-header {
    margin-bottom: 2rem;
  }

  .settings-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: #1a1a1b;
    margin: 0 0 0.5rem 0;
    font-family: "Roboto", sans-serif;
  }

  .settings-description {
    color: #7c7c7c;
    margin: 0;
    font-size: 0.95rem;
  }

  .settings-content {
    display: flex;
    gap: 2rem;
  }

  .settings-sidebar {
    width: 240px;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .settings-tab {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: none;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    color: #7c7c7c;
    font-size: 0.95rem;
    font-weight: 500;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s ease;
    text-align: left;
  }

  .settings-tab:hover {
    background-color: #f6f7f8;
    color: #1a1a1b;
  }

  .settings-tab.active {
    background-color: #f0f1f2;
    color: #153060;
    font-weight: 600;
  }

  .settings-tab svg {
    flex-shrink: 0;
  }

  .settings-main {
    flex: 1;
    background-color: white;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    padding: 2rem;
  }

  .settings-section h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 0.5rem 0;
    font-family: "Roboto", sans-serif;
  }

  .section-description {
    color: #7c7c7c;
    margin: 0 0 2rem 0;
    font-size: 0.9rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #1a1a1b;
    font-size: 0.9rem;
  }

  .form-label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #1a1a1b;
    font-size: 0.9rem;
  }

  .form-group input[type="text"],
  .form-group input[type="email"],
  .form-group input[type="url"],
  .form-group input[type="date"],
  .form-group input[type="password"],
  .form-group textarea,
  .form-group select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    font-size: 0.95rem;
    color: #1a1a1b;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s ease;
  }

  .form-group input:focus,
  .form-group textarea:focus,
  .form-group select:focus {
    outline: none;
    border-color: #153060;
    box-shadow: 0 0 0 3px rgba(21, 48, 96, 0.1);
  }

  .form-group input:disabled {
    background-color: #f6f7f8;
    cursor: not-allowed;
  }

  .input-with-prefix {
    display: flex;
    align-items: center;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    overflow: hidden;
  }

  .input-prefix {
    padding: 0.75rem;
    background-color: #f6f7f8;
    color: #7c7c7c;
    font-weight: 500;
    border-right: 1px solid #e6e6e6;
  }

  .input-with-prefix input {
    flex: 1;
    border: none;
    padding: 0.75rem;
  }

  .input-hint {
    margin: 0.5rem 0 0 0;
    font-size: 0.85rem;
    color: #7c7c7c;
  }

  .interests-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 0.5rem;
  }

  .interest-btn {
    padding: 0.5rem 1rem;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    background: white;
    color: #1a1a1b;
    cursor: pointer;
    font-size: 0.85rem;
    transition: all 0.2s ease;
  }

  .interest-btn:hover {
    border-color: #153060;
  }

  .interest-btn.selected {
    background-color: #153060;
    color: white;
    border-color: #153060;
  }

  .avatar-upload {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .avatar-preview {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    overflow: hidden;
    border: 3px solid #e6e6e6;
  }

  .avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background-color: #153060;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2rem;
    font-weight: bold;
  }

  .avatar-actions {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 2rem;
    border-top: 1px solid #e6e6e6;
  }

  .btn-primary {
    padding: 0.75rem 2rem;
    background-color: #153060;
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
  }

  .btn-primary:hover:not(:disabled) {
    background-color: #0d2144;
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-secondary {
    padding: 0.75rem 2rem;
    background-color: #f6f7f8;
    color: #1a1a1b;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
  }

  .btn-secondary:hover:not(:disabled) {
    background-color: #e9e9e9;
  }

  .btn-secondary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-text {
    padding: 0.5rem 1rem;
    background: none;
    color: #7c7c7c;
    border: none;
    font-weight: 500;
    cursor: pointer;
    font-size: 0.9rem;
    transition: color 0.2s ease;
  }

  .btn-text:hover:not(:disabled) {
    color: #1a1a1b;
  }

  .btn-text:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-remove {
    color: #ff4444;
  }

  .btn-remove:hover:not(:disabled) {
    color: #ff0000;
  }

  .password-section {
    margin-top: 3rem;
    padding-top: 2rem;
    border-top: 1px solid #e6e6e6;
  }

  .password-section h3 {
    font-size: 1rem;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 1rem 0;
  }

  .danger-zone {
    margin-top: 3rem;
    padding: 2rem;
    background-color: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 12px;
  }

  .danger-zone h3 {
    font-size: 1rem;
    font-weight: 600;
    color: #dc2626;
    margin: 0 0 1rem 0;
  }

  .danger-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 2rem;
  }

  .danger-content h4 {
    margin: 0 0 0.5rem 0;
    font-size: 0.95rem;
    font-weight: 600;
    color: #1a1a1b;
  }

  .danger-content p {
    margin: 0;
    font-size: 0.9rem;
    color: #7c7c7c;
    line-height: 1.5;
  }

  .btn-danger {
    padding: 0.75rem 1.5rem;
    background-color: #dc2626;
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    white-space: nowrap;
    transition: background-color 0.2s ease;
  }

  .btn-danger:hover {
    background-color: #b91c1c;
  }

  .setting-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 0;
    border-bottom: 1px solid #e6e6e6;
  }

  .setting-info h4 {
    margin: 0 0 0.25rem 0;
    font-size: 0.95rem;
    font-weight: 600;
    color: #1a1a1b;
  }

  .setting-info p {
    margin: 0;
    font-size: 0.85rem;
    color: #7c7c7c;
  }

  .subsection-title {
    font-size: 1rem;
    font-weight: 600;
    color: #1a1a1b;
    margin: 2rem 0 1rem 0;
  }

  .toggle {
    position: relative;
    display: inline-block;
    width: 50px;
    height: 26px;
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
    border-radius: 26px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 4px;
    bottom: 4px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }

  .toggle input:checked + .toggle-slider {
    background-color: #153060;
  }

  .toggle input:checked + .toggle-slider:before {
    transform: translateX(24px);
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    background: white;
    padding: 2rem;
    border-radius: 12px;
    max-width: 500px;
    width: 90%;
  }

  .modal-content h2 {
    margin: 0 0 1.5rem 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a1a1b;
  }

  .modal-actions {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
  }

  .spinner {
    border: 3px solid #f3f3f3;
    border-top: 3px solid #153060;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .loading-state p {
    margin-top: 1rem;
    color: #7c7c7c;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border-width: 0;
  }

  @media (max-width: 768px) {
    .settings-content {
      flex-direction: column;
    }

    .settings-sidebar {
      width: 100%;
      flex-direction: row;
      overflow-x: auto;
    }

    .settings-tab {
      white-space: nowrap;
    }

    .interests-grid {
      grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    }
  }
</style>
