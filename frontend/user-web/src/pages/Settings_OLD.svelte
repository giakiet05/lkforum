<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
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

  let activeTab: "account" | "privacy" | "notifications" | "appearance" =
    "account";

  let user = $state<UserResponse | null>(null);
  let settings = $state<SettingsResponse | null>(null);
  let provinces = $state<string[]>([]);
  let interests = $state<string[]>([]);
  let genders = $state<string[]>([]);

  let isLoadingUser = $state(true);
  let isLoadingSettings = $state(false);
  let isLoadingMetadata = $state(false);
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

  // Password form
  let showPasswordModal = $state(false);
  let oldPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");

  // Settings form (reactive copies)
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

      // Populate form with current data
      editedBio = user.profile.bio || "";
      editedGender = user.profile.gender || "";
      editedDateOfBirth = ""; // Backend returns age, not DOB
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
        errorMessage = "Failed to load profile. Please try again.";
      }
    } finally {
      isLoadingUser = false;
    }
  }

  async function loadMetadata() {
    try {
      isLoadingMetadata = true;
      const [provincesData, interestsData, gendersData] = await Promise.all([
        getProvinces(),
        getInterests(),
        getGenders(),
      ]);
      provinces = provincesData;
      interests = interestsData;
      genders = gendersData;
    } catch (error) {
      console.error("Failed to load metadata:", error);
    } finally {
      isLoadingMetadata = false;
    }
  }

  async function loadSettings() {
    if (settings) return; // Already loaded

    try {
      isLoadingSettings = true;
      settings = await getSettings();
      editedSettings = JSON.parse(JSON.stringify(settings)); // Deep copy
    } catch (error) {
      console.error("Failed to load settings:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      }
    } finally {
      isLoadingSettings = false;
    }
  }

  // Watch tab changes to load settings when needed
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

      if (editedBio !== (user.profile.bio || "")) {
        payload.bio = editedBio || null;
      }
      if (editedGender !== (user.profile.gender || "")) {
        payload.gender = editedGender || null;
      }
      if (editedDateOfBirth) {
        payload.date_of_birth = editedDateOfBirth;
      }
      if (editedLocation !== (user.profile.location || "")) {
        payload.location = editedLocation || null;
      }
      if (
        JSON.stringify(editedInterests) !==
        JSON.stringify(user.profile.interests || [])
      ) {
        payload.interests = editedInterests.length > 0 ? editedInterests : null;
      }

      // Social links
      const currentLinks = user.profile.social_links || {};
      const newLinks: any = {};
      let hasChanges = false;

      if (editedWebsite !== (currentLinks.website || "")) {
        newLinks.website = editedWebsite || null;
        hasChanges = true;
      }
      if (editedFacebook !== (currentLinks.facebook || "")) {
        newLinks.facebook = editedFacebook || null;
        hasChanges = true;
      }
      if (editedYouTube !== (currentLinks.youtube || "")) {
        newLinks.youtube = editedYouTube || null;
        hasChanges = true;
      }
      if (editedGitHub !== (currentLinks.github || "")) {
        newLinks.github = editedGitHub || null;
        hasChanges = true;
      }

      if (hasChanges) {
        payload.social_links = newLinks;
      }

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
        errorMessage = "Failed to update profile. Please try again.";
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

      await changePassword(oldPassword, newPassword);

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
        errorMessage = "Failed to change password. Please try again.";
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
      editedSettings = JSON.parse(JSON.stringify(settings)); // Update copy

      successMessage = "Settings saved successfully!";
    } catch (error) {
      console.error("Failed to save settings:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to save settings. Please try again.";
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
      successMessage = "Avatar uploaded successfully!";
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to upload avatar. Please try again.";
      }
    } finally {
      isUploadingAvatar = false;
      input.value = "";
    }
  }

  async function handleDeleteAvatar() {
    if (!confirm("Are you sure you want to delete your avatar?")) return;

    try {
      isDeletingAvatar = true;
      errorMessage = null;
      user = await deleteAvatar();
      successMessage = "Avatar deleted successfully!";
    } catch (error) {
      console.error("Failed to delete avatar:", error);
      if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to delete avatar. Please try again.";
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
        errorMessage = "You can select up to 10 interests";
        return;
      }
      editedInterests = [...editedInterests, interest];
    }
  }

  function handleDeleteAccount() {
    // TODO: Implement delete account functionality
    console.log("Delete account");
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
      <button class="alert-close" onclick={() => (errorMessage = null)}>
        ×
      </button>
    </div>
  {/if}

  {#if successMessage}
    <div class="alert alert-success">
      {successMessage}
      <button class="alert-close" onclick={() => (successMessage = null)}>
        ×
      </button>
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
          on:click={() => (activeTab = "account")}
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
          on:click={() => (activeTab = "privacy")}
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
          class:active={activeTab === "email"}
          on:click={() => (activeTab = "email")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <rect
              x="3"
              y="5"
              width="14"
              height="10"
              rx="2"
              stroke="currentColor"
              stroke-width="1.5"
            />
            <path
              d="M3 7L10 11L17 7"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          Email
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "notification"}
          on:click={() => (activeTab = "notification")}
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
          Notification
        </button>
      </div>

      <div class="settings-main">
        {#if activeTab === "account"}
          <div class="settings-section">
            <h2>Account Settings</h2>
            <p class="section-description">
              Manage your account information and preferences
            </p>

            <div class="form-group">
              <label for="display-name">Display Name</label>
              <input
                type="text"
                id="display-name"
                value={user.displayName}
                placeholder="Enter your display name"
              />
            </div>

            <div class="form-group">
              <label for="username">Username</label>
              <div class="input-with-prefix">
                <span class="input-prefix">u/</span>
                <input
                  type="text"
                  id="username"
                  value={user.name}
                  placeholder="Enter your username"
                />
              </div>
              <p class="input-hint">Your username cannot be changed once set</p>
            </div>

            <div class="form-group">
              <label for="bio">Bio</label>
              <textarea
                id="bio"
                rows="4"
                value={user.bio}
                placeholder="Tell us about yourself"
              ></textarea>
              <p class="input-hint">Brief description for your profile</p>
            </div>

            <div class="form-group">
              <label>Profile Picture</label>
              <div class="avatar-upload">
                <div class="avatar-preview">
                  <img src={user.avatar} alt="Profile" />
                </div>
                <div class="avatar-actions">
                  <button class="btn-secondary">Change Avatar</button>
                  <button class="btn-text">Remove</button>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button class="btn-primary" on:click={handleSaveAccount}>
                Save Changes
              </button>
              <button class="btn-secondary">Cancel</button>
            </div>

            <div class="password-section">
              <h3>Password</h3>
              <button class="btn-secondary" on:click={handleChangePassword}>
                Change Password
              </button>
            </div>
          </div>
        {:else if activeTab === "privacy"}
          <div class="settings-section">
            <h2>Privacy Settings</h2>
            <p class="section-description">
              Privacy settings will be implemented here
            </p>
          </div>
        {:else if activeTab === "email"}
          <div class="settings-section">
            <h2>Email Settings</h2>
            <p class="section-description">
              Email settings will be implemented here
            </p>
          </div>
        {:else if activeTab === "notification"}
          <div class="settings-section">
            <h2>Notification Settings</h2>
            <p class="section-description">
              Notification settings will be implemented here
            </p>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .settings-page {
    background-color: white;
    min-height: 100vh;
    padding-top: 72px;
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

  .settings-section {
    margin-bottom: 2rem;
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

  .form-group input[type="text"],
  .form-group textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    font-size: 0.95rem;
    color: #1a1a1b;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s ease;
  }

  .form-group input[type="text"]:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #153060;
    box-shadow: 0 0 0 3px rgba(21, 48, 96, 0.1);
  }

  .input-with-prefix {
    display: flex;
    align-items: center;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.2s ease;
  }

  .input-with-prefix:focus-within {
    border-color: #153060;
    box-shadow: 0 0 0 3px rgba(21, 48, 96, 0.1);
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

  .input-with-prefix input:focus {
    outline: none;
    box-shadow: none;
  }

  .input-hint {
    margin: 0.5rem 0 0 0;
    font-size: 0.85rem;
    color: #7c7c7c;
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

  .btn-primary:hover {
    background-color: #0d2144;
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

  .btn-secondary:hover {
    background-color: #e9e9e9;
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

  .btn-text:hover {
    color: #1a1a1b;
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

  .password-section .btn-secondary {
    font-weight: 500;
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
  }
</style>
