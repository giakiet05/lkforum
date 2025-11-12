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
        errorMessage = "Failed to load profile.";
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
        successMessage = "No changes to save";
        return;
      }

      user = await updateProfile(payload);
      successMessage = "Profile updated successfully!";
    } catch (error) {
      console.error("Failed to update profile:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to update profile.";
      }
    } finally {
      isSaving = false;
    }
  }

  async function handleChangePassword() {
    if (!oldPassword || !newPassword || !confirmPassword) {
      errorMessage = "All password fields are required";
      return;
    }

    if (newPassword !== confirmPassword) {
      errorMessage = "New passwords do not match";
      return;
    }

    if (newPassword.length < 6) {
      errorMessage = "New password must be at least 6 characters";
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

      successMessage = "Password changed successfully!";
      showPasswordModal = false;
      oldPassword = "";
      newPassword = "";
      confirmPassword = "";
    } catch (error) {
      console.error("Failed to change password:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to change password.";
      }
    } finally {
      isSaving = false;
    }
  }

  async function handleSaveSettings() {
    if (!editedSettings) return;

    try {
      isSaving = true;
      errorMessage = null;
      successMessage = null;

      settings = await updateSettings(editedSettings);
      editedSettings = JSON.parse(JSON.stringify(settings));

      successMessage = "Settings saved successfully!";
    } catch (error) {
      console.error("Failed to save settings:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to save settings.";
      }
    } finally {
      isSaving = false;
    }
  }

  async function handleAvatarChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      errorMessage = "Please select an image file";
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      errorMessage = "Image size must be less than 5MB";
      return;
    }

    try {
      isUploadingAvatar = true;
      errorMessage = null;
      user = await uploadAvatar(file);
      successMessage = "Avatar uploaded!";
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to upload avatar.";
      }
    } finally {
      isUploadingAvatar = false;
      input.value = "";
    }
  }

  async function handleDeleteAvatar() {
    if (!confirm("Delete your avatar?")) return;

    try {
      isDeletingAvatar = true;
      errorMessage = null;
      user = await deleteAvatar();
      successMessage = "Avatar deleted!";
    } catch (error) {
      console.error("Failed to delete avatar:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to delete avatar.";
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
        errorMessage = "Max 10 interests";
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
      <h1>Settings</h1>
      <p class="settings-description">
        Manage your account settings and preferences
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
          Account
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
          Privacy
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
          Notifications
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
          Appearance
        </button>
      </div>

      <div class="settings-main">
        {#if isLoadingUser}
          <div class="loading-state">
            <div class="spinner"></div>
            <p>Loading...</p>
          </div>
        {:else if activeTab === "account" && user}
          <div class="settings-section">
            <h2>Account Settings</h2>
            <p class="section-description">Manage your account information</p>

            <div class="form-group">
              <label for="username-label">Username</label>
              <div class="input-with-prefix">
                <span class="input-prefix">u/</span>
                <input
                  type="text"
                  id="username-label"
                  value={user.username}
                  disabled
                />
              </div>
              <p class="input-hint">Your username cannot be changed</p>
            </div>

            <div class="form-group">
              <label for="email-label">Email</label>
              <input
                type="email"
                id="email-label"
                value={user.email}
                disabled
              />
              <p class="input-hint">Your email cannot be changed</p>
            </div>

            <div class="form-group">
              <label for="bio">Bio</label>
              <textarea
                id="bio"
                rows="4"
                bind:value={editedBio}
                placeholder="Tell us about yourself"
                maxlength="500"
              ></textarea>
              <p class="input-hint">{editedBio.length}/500 characters</p>
            </div>

            <div class="form-group">
              <label for="gender">Gender</label>
              <select id="gender" bind:value={editedGender}>
                <option value="">Select gender</option>
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
              <label for="dob">Date of Birth</label>
              <input type="date" id="dob" bind:value={editedDateOfBirth} />
              <p class="input-hint">You must be at least 13 years old</p>
            </div>

            <div class="form-group">
              <label for="location">Location</label>
              <select id="location" bind:value={editedLocation}>
                <option value="">Select province</option>
                {#each provinces as province}
                  <option value={province}>{province}</option>
                {/each}
              </select>
            </div>

            <div class="form-group">
              <label for="interests-label">Interests (max 10)</label>
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
              <p class="input-hint">{editedInterests.length}/10 selected</p>
            </div>

            <div class="form-group">
              <label for="social-website">Social Links</label>
              <input
                type="url"
                id="social-website"
                bind:value={editedWebsite}
                placeholder="Website URL"
              />
            </div>
            <div class="form-group">
              <label for="social-facebook" class="sr-only">Facebook</label>
              <input
                type="text"
                id="social-facebook"
                bind:value={editedFacebook}
                placeholder="Facebook username or URL"
              />
            </div>
            <div class="form-group">
              <label for="social-youtube" class="sr-only">YouTube</label>
              <input
                type="text"
                id="social-youtube"
                bind:value={editedYouTube}
                placeholder="YouTube channel URL"
              />
            </div>
            <div class="form-group">
              <label for="social-github" class="sr-only">GitHub</label>
              <input
                type="text"
                id="social-github"
                bind:value={editedGitHub}
                placeholder="GitHub username"
              />
            </div>

            <div class="form-group">
              <span class="form-label">Profile Picture</span>
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
                    {isUploadingAvatar ? "Uploading..." : "Change Avatar"}
                  </button>
                  {#if user.profile.avatar?.url}
                    <button
                      class="btn-text"
                      onclick={handleDeleteAvatar}
                      disabled={isUploadingAvatar || isDeletingAvatar}
                    >
                      {isDeletingAvatar ? "Deleting..." : "Remove"}
                    </button>
                  {/if}
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button
                class="btn-primary"
                onclick={handleSaveAccount}
                disabled={isSaving}
              >
                {isSaving ? "Saving..." : "Save Changes"}
              </button>
            </div>

            <div class="password-section">
              <h3>Password</h3>
              <button
                class="btn-secondary"
                onclick={() => (showPasswordModal = true)}
                >Change Password</button
              >
            </div>
          </div>
        {:else if activeTab === "privacy" && editedSettings}
          <div class="settings-section">
            <h2>Privacy Settings</h2>
            <p class="section-description">Control your privacy preferences</p>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Show Profile</h4>
                <p>Make your profile visible to others</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.show_profile}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Show Email</h4>
                <p>Display your email on your profile</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.show_email}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Show Post History</h4>
                <p>Allow others to see your post history</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.show_post_history}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Allow Direct Messages</h4>
                <p>Let others send you direct messages</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.allow_direct_messages}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Allow Mentions</h4>
                <p>Let others mention you in posts and comments</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.privacy.allow_mentions}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="form-actions">
              <button
                class="btn-primary"
                onclick={handleSaveSettings}
                disabled={isSaving}
              >
                {isSaving ? "Saving..." : "Save Settings"}
              </button>
            </div>
          </div>
        {:else if activeTab === "notifications" && editedSettings}
          <div class="settings-section">
            <h2>Notification Settings</h2>
            <p class="section-description">
              Manage how you receive notifications
            </p>

            <div class="setting-item">
              <div class="setting-info">
                <h4>In-App Notifications</h4>
                <p>Show notifications in the website</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.in_app_enabled}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Email Notifications</h4>
                <p>Send notifications to your email</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.email_enabled}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <h3 class="subsection-title">Notify me when:</h3>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Someone comments on my post</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_comment}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Someone mentions me</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_mention}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Someone upvotes my post</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_upvote}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Someone sends me a message</h4>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.notifications.notify_on_message}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="form-actions">
              <button
                class="btn-primary"
                onclick={handleSaveSettings}
                disabled={isSaving}
              >
                {isSaving ? "Saving..." : "Save Settings"}
              </button>
            </div>
          </div>
        {:else if activeTab === "appearance" && editedSettings}
          <div class="settings-section">
            <h2>Appearance Settings</h2>
            <p class="section-description">Customize the look and feel</p>

            <div class="form-group">
              <label for="theme">Theme</label>
              <select id="theme" bind:value={editedSettings.appearance.theme}>
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="auto">Auto</option>
              </select>
            </div>

            <div class="form-group">
              <label for="fontSize">Font Size</label>
              <select
                id="fontSize"
                bind:value={editedSettings.appearance.font_size}
              >
                <option value="small">Small</option>
                <option value="medium">Medium</option>
                <option value="large">Large</option>
              </select>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h4>Allow NSFW Content</h4>
                <p>Show age-restricted content</p>
              </div>
              <label class="toggle">
                <input
                  type="checkbox"
                  bind:checked={editedSettings.content.allow_nsfw}
                />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div class="form-actions">
              <button
                class="btn-primary"
                onclick={handleSaveSettings}
                disabled={isSaving}
              >
                {isSaving ? "Saving..." : "Save Settings"}
              </button>
            </div>
          </div>
        {:else if isLoadingSettings}
          <div class="loading-state">
            <div class="spinner"></div>
            <p>Loading settings...</p>
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
      <h2>Change Password</h2>
      <div class="form-group">
        <label for="oldPassword">Current Password</label>
        <input type="password" id="oldPassword" bind:value={oldPassword} />
      </div>
      <div class="form-group">
        <label for="newPassword">New Password</label>
        <input type="password" id="newPassword" bind:value={newPassword} />
      </div>
      <div class="form-group">
        <label for="confirmPassword">Confirm New Password</label>
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
          {isSaving ? "Changing..." : "Change Password"}
        </button>
        <button
          class="btn-secondary"
          onclick={() => (showPasswordModal = false)}>Cancel</button
        >
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
    background-color: #efe;
    color: #060;
    border: 1px solid #cfc;
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
