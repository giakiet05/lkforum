<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import Post from "../components/Post.svelte";
  import type { PostResponse } from "../dtos/post-dto";
  import type { UserResponse } from "../dtos/user-dto";
  import {
    getMyProfile,
    uploadAvatar,
    uploadCover,
  } from "../services/user-service";
  import { ApiError } from "../errors/api-error";

  let user = $state<UserResponse | null>(null);
  let isLoadingUser = $state(true);
  let errorMessage = $state<string | null>(null);

  // TODO: Replace with API call for posts
  let posts = $state<PostResponse[]>([]);
  let isLoadingPosts = $state(false);

  let activeTab = $state<"posts" | "comments" | "upvoted" | "downvoted">(
    "posts"
  );
  let sortBy = $state<"hot" | "newest" | "oldest">("hot");

  let avatarFileInput: HTMLInputElement;
  let coverFileInput: HTMLInputElement;
  let isUploadingAvatar = $state(false);
  let isUploadingCover = $state(false);

  onMount(() => {
    loadUserProfile();
  });

  async function loadUserProfile() {
    try {
      isLoadingUser = true;
      errorMessage = null;
      user = await getMyProfile();
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

  async function handleAvatarChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      alert("Please select an image file");
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert("Image size must be less than 5MB");
      return;
    }

    try {
      isUploadingAvatar = true;
      user = await uploadAvatar(file);
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      if (error instanceof ApiError) {
        alert(error.message);
      } else {
        alert("Failed to upload avatar. Please try again.");
      }
    } finally {
      isUploadingAvatar = false;
      input.value = ""; // Reset input
    }
  }

  async function handleCoverChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      alert("Please select an image file");
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert("Image size must be less than 5MB");
      return;
    }

    try {
      isUploadingCover = true;
      user = await uploadCover(file);
    } catch (error) {
      console.error("Failed to upload cover:", error);
      if (error instanceof ApiError) {
        alert(error.message);
      } else {
        alert("Failed to upload cover. Please try again.");
      }
    } finally {
      isUploadingCover = false;
      input.value = ""; // Reset input
    }
  }

  function handleCreatePost() {
    push("/submit");
  }

  function handleEditProfile() {
    push("/settings");
  }

  function handleSettings() {
    push("/settings");
  }

  function formatDate(dateString?: string): string {
    if (!dateString) return "Unknown";
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "long",
      day: "numeric",
      year: "numeric",
    });
  }
</script>

{#if isLoadingUser}
  <div class="loading-state">
    <div class="spinner"></div>
    <p>Loading profile...</p>
  </div>
{:else if errorMessage}
  <div class="error-state">
    <p>{errorMessage}</p>
    <button class="retry-btn" onclick={loadUserProfile}>Retry</button>
  </div>
{:else if user}
  <div class="profile-page">
    <input
      type="file"
      accept="image/*"
      bind:this={avatarFileInput}
      onchange={handleAvatarChange}
      style="display: none;"
    />
    <input
      type="file"
      accept="image/*"
      bind:this={coverFileInput}
      onchange={handleCoverChange}
      style="display: none;"
    />

    <div class="profile-header">
      <div class="cover-image-wrapper">
        {#if user.profile.cover?.url}
          <img src={user.profile.cover.url} alt="Cover" class="cover-image" />
        {:else}
          <div class="cover-placeholder"></div>
        {/if}
        <button
          class="change-cover-btn"
          onclick={() => coverFileInput.click()}
          disabled={isUploadingCover}
        >
          {#if isUploadingCover}
            <div class="mini-spinner"></div>
          {:else}
            <img src="/change_profile_image.png" alt="Change cover" />
          {/if}
        </button>
      </div>
      <div class="profile-info-bar">
        <div class="profile-details">
          <div class="avatar-wrapper">
            {#if user.profile.avatar?.url}
              <img
                src={user.profile.avatar.url}
                alt={user.username}
                class="profile-avatar"
              />
            {:else}
              <div class="avatar-placeholder">
                {user.username[0].toUpperCase()}
              </div>
            {/if}
            <button
              class="change-avatar-btn"
              onclick={() => avatarFileInput.click()}
              disabled={isUploadingAvatar}
            >
              {#if isUploadingAvatar}
                <div class="mini-spinner"></div>
              {:else}
                <img src="/change_profile_image.png" alt="Change avatar" />
              {/if}
            </button>
          </div>
          <div class="profile-text">
            <h1 class="username">{user.username}</h1>
            <p class="user-handle">u/{user.username}</p>
          </div>
        </div>
        <div class="profile-actions">
          <button class="action-btn primary" onclick={handleCreatePost}>
            <i class="fas fa-plus"></i>
            + Create Post
          </button>
          <button class="action-btn secondary" onclick={handleEditProfile}>
            Edit Profile
          </button>
          <button class="action-btn tertiary" onclick={handleSettings}>
            <img src="/dot.png" alt="Settings" class="settings-icon" />
          </button>
        </div>
      </div>
    </div>

    <div class="profile-content">
      <div class="profile-main-content">
        <div class="profile-tabs">
          <button
            class="tab-btn"
            class:active={activeTab === "posts"}
            onclick={() => (activeTab = "posts")}>Posts</button
          >
          <button
            class="tab-btn"
            class:active={activeTab === "comments"}
            onclick={() => (activeTab = "comments")}>Comments</button
          >
          <button
            class="tab-btn"
            class:active={activeTab === "upvoted"}
            onclick={() => (activeTab = "upvoted")}>Upvoted</button
          >
          <button
            class="tab-btn"
            class:active={activeTab === "downvoted"}
            onclick={() => (activeTab = "downvoted")}>Downvoted</button
          >

          <div class="sort-options">
            <select bind:value={sortBy}>
              <option value="hot">Hot</option>
              <option value="newest">Newest</option>
              <option value="oldest">Oldest</option>
            </select>
          </div>
        </div>

        <div class="tab-content">
          {#if activeTab === "posts"}
            {#if isLoadingPosts}
              <div class="loading-posts">
                <div class="spinner"></div>
                <p>Loading posts...</p>
              </div>
            {:else if posts.length === 0}
              <div class="empty-state">
                <p>No posts yet</p>
              </div>
            {:else}
              <div class="post-list">
                {#each posts as post}
                  <Post {post} />
                {/each}
              </div>
            {/if}
          {:else if activeTab === "comments"}
            <p>Comments will be shown here.</p>
          {:else if activeTab === "upvoted"}
            <p>Upvoted posts will be shown here.</p>
          {:else if activeTab === "downvoted"}
            <p>Downvoted posts will be shown here.</p>
          {/if}
        </div>
      </div>
      <div class="profile-sidebar">
        <div class="user-card">
          <div class="user-card-body">
            <h3>About</h3>
            <p class="bio">
              {user.profile.bio || "No bio yet."}
            </p>
            <div class="user-stats">
              <div class="stat">
                <span class="stat-value"
                  >{user.profile.stats?.post_count ?? 0}</span
                >
                <span class="stat-label">Posts</span>
              </div>
              <div class="stat">
                <span class="stat-value"
                  >{user.profile.stats?.comment_count ?? 0}</span
                >
                <span class="stat-label">Comments</span>
              </div>
              <div class="stat">
                <span class="stat-value">{user.reputation ?? 0}</span>
                <span class="stat-label">Reputation</span>
              </div>
            </div>
            <div class="cake-day">
              <i class="fas fa-birthday-cake"></i>
              Member since: {formatDate(user.profile.stats?.member_since)}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .profile-page {
    background-color: white;
    min-height: 100vh;
  }

  .profile-header {
    background-color: white;
    border-bottom: 1px solid #e6e6e6;
    margin-bottom: 1rem;
  }

  .cover-image-wrapper {
    height: 200px;
    overflow: hidden;
    background-color: #f6f7f8;
    position: relative;
    width: calc(100% - 48px);
    border-radius: 8px;
    margin: 8px 24px;
  }

  .cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .change-cover-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    width: 40px;
    height: 40px;
    border: none;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.9);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.2s ease;
  }

  .change-cover-btn:hover {
    background-color: rgba(255, 255, 255, 1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    transform: scale(1.05);
  }

  .change-cover-btn img {
    width: 20px;
    height: 20px;
  }

  .profile-info-bar {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    padding: 0 24px;
    max-width: 100%;
    margin: 0 auto;
    padding-bottom: 1rem;
  }

  .profile-details {
    display: flex;
    align-items: flex-end;
    transform: translateY(-30px);
  }

  .avatar-wrapper {
    position: relative;
  }

  .profile-avatar {
    width: 150px;
    height: 150px;
    border-radius: 50%;
    border: 4px solid white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .change-avatar-btn {
    position: absolute;
    bottom: 5px;
    right: 5px;
    width: 36px;
    height: 36px;
    border: 2px solid white;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.95);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
    transition: all 0.2s ease;
  }

  .change-avatar-btn:hover {
    background-color: rgba(255, 255, 255, 1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    transform: scale(1.1);
  }

  .change-avatar-btn img {
    width: 18px;
    height: 18px;
  }

  .profile-text {
    margin-left: 1.5rem;
    margin-bottom: 0.5rem;
  }

  .username {
    font-size: 2rem;
    font-weight: bold;
    margin: 0;
    color: #1a1a1b;
  }

  .user-handle {
    color: #7c7c7c;
    margin: 0.25rem 0 0 0;
    font-size: 0.9rem;
    font-weight: 400;
  }

  .profile-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .action-btn {
    border: none;
    padding: 0.65rem 1.5rem;
    border-radius: 24px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.95rem;
    transition: all 0.2s ease;
  }

  .action-btn.primary {
    background-color: #f6f7f8;
    color: #1a1a1b;
  }

  .action-btn.primary:hover {
    background-color: #e9e9e9;
  }

  .action-btn.secondary {
    background-color: #153060;
    color: white;
  }

  .action-btn.secondary:hover {
    background-color: #0d2144;
  }

  .action-btn.tertiary {
    background-color: #f6f7f8;
    color: #1a1a1b;
    padding: 0.65rem;
  }

  .action-btn.tertiary:hover {
    background-color: #e9e9e9;
  }

  .profile-content {
    display: flex;
    gap: 1.5rem;
    padding: 0 24px;
    max-width: 100%;
    margin: 0 auto;
  }

  .profile-main-content {
    flex-grow: 1;
  }

  .profile-tabs {
    background-color: white;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    font-family: "Roboto", sans-serif;
  }

  .tab-btn {
    background: none;
    border: none;
    padding: 0.75rem 1.25rem;
    font-weight: 500;
    cursor: pointer;
    color: #1a1a1b;
    transition: all 0.2s;
    font-family: "Roboto", sans-serif;
  }

  .tab-btn:hover {
    opacity: 0.7;
  }

  .tab-btn.active {
    color: #00008b;
    font-weight: 600;
    position: relative;
  }

  .tab-btn.active::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 2px;
    background-color: #00008b;
  }

  .sort-options {
    padding: 0 1rem;
    position: relative;
  }

  .sort-options::after {
    content: "";
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    width: 20px;
    height: 20px;
    background-image: url("/Sort.jpg");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    pointer-events: none;
    opacity: 0.8;
  }

  .sort-options select {
    padding: 0.5rem 2rem 0.5rem 2.75rem;
    border: none;
    border-radius: 4px;
    background-color: #f8f9fa;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 1em;
    color: #1a1a1b;
    font-size: 0.9rem;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    font-weight: 400;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    transition: all 0.2s ease;
  }

  /* Style for the dropdown menu */
  .sort-options select:not(:focus) {
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .sort-options select:hover {
    background-color: #f0f1f2;
  }

  .sort-options select:focus {
    outline: none;
    background-color: #fff;
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.06);
  }

  /* Style for the dropdown options */
  .sort-options select option {
    padding: 0.75rem 1rem;
    background-color: white;
    color: #1a1a1b;
    font-size: 0.9rem;
    font-weight: 400;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .sort-options select option:hover {
    background-color: #f8f9fa;
    color: #00008b;
  }

  .sort-options select option:checked {
    background-color: #f0f1f2;
    font-weight: 500;
  }

  /* Style for the dropdown container when opened */
  .sort-options select:focus {
    border-radius: 4px;
  }

  @media screen and (-webkit-min-device-pixel-ratio: 0) {
    .sort-options select {
      border-radius: 4px !important;
    }

    .sort-options select:focus {
      border: none;
    }

    .sort-options select option:checked {
      background: #f0f1f2 linear-gradient(0deg, #f0f1f2 0%, #f0f1f2 100%);
      font-weight: 500;
    }

    .sort-options select option:hover {
      background: #e8f0fe linear-gradient(0deg, #e8f0fe 0%, #e8f0fe 100%);
      color: #00008b;
    }
  }

  .tab-content {
    padding-top: 1rem;
  }

  .profile-sidebar {
    width: 300px;
  }

  .user-card {
    background-color: white;
    border-radius: 4px;
    border: 1px solid #e6e6e6;
  }

  .user-card-body {
    padding: 1rem;
  }

  .user-card h3 {
    margin: 0 0 1rem 0;
    color: #1a1a1b;
    font-size: 1rem;
    font-weight: 500;
  }

  .bio {
    color: #7c7c7c;
    font-size: 0.9rem;
    margin: 0 0 1rem 0;
    line-height: 1.4;
  }

  .user-stats {
    display: flex;
    justify-content: space-between;
    margin: 1rem 0;
    padding: 1rem 0;
    border-top: 1px solid #e6e6e6;
    border-bottom: 1px solid #e6e6e6;
  }

  .stat {
    text-align: center;
  }

  .stat-value {
    font-weight: 600;
    display: block;
    color: #1a1a1b;
    font-size: 1.1rem;
  }

  .stat-label {
    color: #7c7c7c;
    font-size: 0.8rem;
    margin-top: 0.25rem;
  }

  .cake-day {
    color: #7c7c7c;
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .cake-day i {
    color: #0079d3;
  }

  .settings-icon {
    width: 20px;
    height: 20px;
  }

  /* Loading States */
  .loading-state,
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    padding: 2rem;
  }

  .loading-state p,
  .error-state p {
    margin-top: 1rem;
    color: #7c7c7c;
  }

  .spinner {
    border: 3px solid #f3f3f3;
    border-top: 3px solid #153060;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
  }

  .mini-spinner {
    border: 2px solid #f3f3f3;
    border-top: 2px solid #153060;
    border-radius: 50%;
    width: 16px;
    height: 16px;
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

  .retry-btn {
    margin-top: 1rem;
    padding: 0.75rem 1.5rem;
    background-color: #153060;
    color: white;
    border: none;
    border-radius: 24px;
    font-weight: 600;
    cursor: pointer;
  }

  .retry-btn:hover {
    background-color: #0d2144;
  }

  .loading-posts {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 3rem;
  }

  .empty-state {
    text-align: center;
    padding: 3rem;
    color: #7c7c7c;
  }

  /* Placeholder styles */
  .cover-placeholder {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .avatar-placeholder {
    width: 150px;
    height: 150px;
    border-radius: 50%;
    border: 4px solid white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    background-color: #153060;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3rem;
    font-weight: bold;
  }

  /* Disabled button state */
  .change-avatar-btn:disabled,
  .change-cover-btn:disabled {
    cursor: not-allowed;
    opacity: 0.7;
  }
</style>
