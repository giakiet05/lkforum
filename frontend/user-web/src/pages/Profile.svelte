<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import Post from "../components/Post.svelte";
  import CreatePostModal from "../components/CreatePostModal.svelte";
  import type { PostResponse } from "../dtos/post-dto";
  import type { UserResponse } from "../dtos/user-dto";
  import {
    getMyProfile,
    getUserByUsername,
    uploadAvatar,
    uploadCover,
  } from "../services/user-service";
  import {
    getMyPosts,
    getSavedPosts,
    getPostsByUserId,
  } from "../services/post-service";
  import {
    getChannelBetweenUsers,
    createChannel,
  } from "../services/channel-service";
  import { ApiError } from "../errors/api-error";
  import { setAuth, authStore } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
  import { setActiveChannel, chatStore } from "../stores/chat-store";

  // Create post modal state
  let showCreatePostModal = $state(false);

  // Route params
  let { params = {} }: { params?: { username?: string } } = $props();

  let user = $state<UserResponse | null>(null);
  let isLoadingUser = $state(true);
  let errorMessage = $state<string | null>(null);
  let isPrivateProfile = $state(false);
  let privateProfileData = $state<{
    username?: string;
    avatar_url?: string;
    cover_url?: string;
  } | null>(null);

  let posts = $state<PostResponse[]>([]);
  let savedPosts = $state<PostResponse[]>([]);
  let isLoadingPosts = $state(false);
  let isLoadingSaved = $state(false);

  // Flags to prevent infinite loading when user has no data
  let hasLoadedPosts = $state(false);
  let hasLoadedSaved = $state(false);
  let postsError = $state<string | null>(null);
  let savedError = $state<string | null>(null);

  let activeTab = $state<"posts" | "saved">("posts");
  let sortBy = $state<"hot" | "newest" | "oldest">("hot");

  let avatarFileInput = $state<HTMLInputElement | undefined>();
  let coverFileInput = $state<HTMLInputElement | undefined>();
  let isUploadingAvatar = $state(false);
  let isUploadingCover = $state(false);
  let isCreatingChannel = $state(false);

  // Watch for params changes to reload profile
  $effect(() => {
    // Reset states when params change
    hasLoadedPosts = false;
    hasLoadedSaved = false;
    posts = [];
    savedPosts = [];
    postsError = null;
    savedError = null;
    activeTab = "posts";

    // Load new profile
    loadUserProfile();
  });

  onMount(() => {
    // Listen for auth:unauthorized event to redirect
    const handleUnauthorized = () => {
      push("/");
    };
    window.addEventListener("auth:unauthorized", handleUnauthorized);

    return () => {
      window.removeEventListener("auth:unauthorized", handleUnauthorized);
    };
  });

  // Watch for tab changes to load data
  $effect(() => {
    if (user && activeTab === "posts" && !hasLoadedPosts && !postsError) {
      loadUserPosts();
    } else if (
      user &&
      activeTab === "saved" &&
      !hasLoadedSaved &&
      !savedError
    ) {
      loadSavedPosts();
    }
  });

  async function loadUserProfile() {
    try {
      isLoadingUser = true;
      errorMessage = null;
      isPrivateProfile = false;
      privateProfileData = null;

      // If username param exists, load that user's profile
      // Otherwise load current user's profile
      if (params.username) {
        user = await getUserByUsername(params.username);
      } else {
        user = await getMyProfile();
      }
    } catch (error) {
      console.error("Failed to load profile:", error);
      if (error instanceof ApiError && error.code === "FORBIDDEN") {
        // Profile is private - show limited info
        isPrivateProfile = true;
        console.log("🔒 Private profile detected");
        console.log("Error data:", error.data);
        // Extract basic profile data from error response
        const profileData = error.data || {};
        console.log("Profile data:", profileData);
        privateProfileData = {
          username: profileData.username || params.username || "User",
          avatar_url: profileData.profile?.avatar?.url,
          cover_url: profileData.profile?.cover?.url,
        };
        console.log("Private profile data:", privateProfileData);
      } else if (error instanceof ApiError) {
        errorMessage = error.message;
      } else {
        errorMessage = "Failed to load profile. Please try again.";
      }
    } finally {
      isLoadingUser = false;
    }
  }

  // Check if viewing own profile
  const isOwnProfile = $derived(!params.username);

  async function handleAvatarChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      toastStore.warning("Vui lòng chọn tệp hình ảnh");
      return;
    }

    // Validate file size (max 100MB)
    if (file.size > 100 * 1024 * 1024) {
      toastStore.warning("Kích thước ảnh phải nhỏ hơn 100MB");
      return;
    }

    try {
      isUploadingAvatar = true;
      const updatedUser = await uploadAvatar(file);
      user = updatedUser;
      // Đồng bộ với authStore và localStorage
      setAuth(updatedUser);
      localStorage.setItem("user", JSON.stringify(updatedUser));
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      if (error instanceof ApiError) {
        toastStore.error(error.message);
      } else {
        toastStore.error("Không thể tải ảnh đại diện. Vui lòng thử lại.");
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
      toastStore.warning("Vui lòng chọn tệp hình ảnh");
      return;
    }

    // Validate file size (max 100MB)
    if (file.size > 100 * 1024 * 1024) {
      toastStore.warning("Kích thước ảnh phải nhỏ hơn 100MB");
      return;
    }

    try {
      isUploadingCover = true;
      const updatedUser = await uploadCover(file);
      user = updatedUser;
      // Đồng bộ với authStore và localStorage
      setAuth(updatedUser);
      localStorage.setItem("user", JSON.stringify(updatedUser));
    } catch (error) {
      console.error("Failed to upload cover:", error);
      if (error instanceof ApiError) {
        toastStore.error(error.message);
      } else {
        toastStore.error("Không thể tải ảnh bìa. Vui lòng thử lại.");
      }
    } finally {
      isUploadingCover = false;
      input.value = ""; // Reset input
    }
  }

  function handleCreatePost() {
    console.log(
      "Profile: handleCreatePost clicked, showCreatePostModal =",
      showCreatePostModal,
    );
    showCreatePostModal = true;
    console.log(
      "Profile: after setting, showCreatePostModal =",
      showCreatePostModal,
    );
  }

  function handleEditProfile() {
    push("/settings");
  }

  function handleSettings() {
    push("/settings");
  }

  // Handle send message to user
  async function handleSendMessage() {
    if (!user || isOwnProfile) return;

    const currentUser = $authStore.user;
    if (!currentUser) {
      toastStore.warning("Vui lòng đăng nhập để gửi tin nhắn");
      return;
    }

    try {
      isCreatingChannel = true;

      // Check if channel already exists
      let channel = await getChannelBetweenUsers(currentUser.id, user.id);

      // If not, create new channel
      if (!channel) {
        channel = await createChannel(
          user.id,
          user.username,
          user.profile.avatar?.url || "",
        );
      }

      // Set active channel and open chat popup
      setActiveChannel(channel.id);

      // Trigger chat popup to open by dispatching custom event
      window.dispatchEvent(
        new CustomEvent("open-chat", { detail: { channelId: channel.id } }),
      );

      toastStore.success(`Đang mở trò chuyện với ${user.username}`);
    } catch (error) {
      console.error("Failed to create/open channel:", error);
      if (error instanceof ApiError) {
        toastStore.error(error.message);
      } else {
        toastStore.error("Không thể mở trò chuyện. Vui lòng thử lại.");
      }
    } finally {
      isCreatingChannel = false;
    }
  }

  async function loadUserPosts() {
    if (!user) return;

    try {
      isLoadingPosts = true;
      postsError = null;

      // If viewing other user's profile, use their ID
      // Otherwise use getMyPosts for current user
      if (isOwnProfile) {
        posts = await getMyPosts({ page: 1, limit: 20 });
      } else {
        posts = await getPostsByUserId(user.id, 1, 20);
      }
    } catch (error) {
      console.error("Failed to load posts:", error);
      if (error instanceof ApiError) {
        postsError = error.message;
      } else {
        postsError = "Không thể tải bài viết. Vui lòng thử lại.";
      }
      posts = [];
    } finally {
      isLoadingPosts = false;
      hasLoadedPosts = true;
    }
  }

  async function loadSavedPosts() {
    if (!user) return;

    try {
      isLoadingSaved = true;
      savedError = null;
      savedPosts = await getSavedPosts(1, 20);
    } catch (error) {
      console.error("Failed to load saved posts:", error);
      if (error instanceof ApiError) {
        savedError = error.message;
      } else {
        savedError = "Không thể tải bài viết đã lưu. Vui lòng thử lại.";
      }
      savedPosts = [];
    } finally {
      isLoadingSaved = false;
      hasLoadedSaved = true;
    }
  }

  function formatDate(dateString?: string): string {
    if (!dateString) return "Không xác định";
    const date = new Date(dateString);
    return date.toLocaleDateString("vi-VN", {
      month: "long",
      day: "numeric",
      year: "numeric",
    });
  }

  function formatSocialLink(url: string, maxLength: number = 30): string {
    try {
      const urlObj = new URL(url);
      let displayText = urlObj.hostname + urlObj.pathname;

      // Remove www. prefix
      displayText = displayText.replace(/^www\./, "");

      // Remove trailing slash
      displayText = displayText.replace(/\/$/, "");

      // Truncate if too long
      if (displayText.length > maxLength) {
        return displayText.substring(0, maxLength - 3) + "...";
      }

      return displayText;
    } catch {
      // If URL parsing fails, just truncate the string
      if (url.length > maxLength) {
        return url.substring(0, maxLength - 3) + "...";
      }
      return url;
    }
  }
</script>

<!-- Create Post Modal -->
<CreatePostModal
  show={showCreatePostModal}
  onClose={() => (showCreatePostModal = false)}
/>

{#if isLoadingUser}
  <div class="loading-state">
    <div class="spinner"></div>
    <p>Đang tải hồ sơ...</p>
  </div>
{:else if isPrivateProfile && privateProfileData}
  <div class="profile-page private-profile">
    <div class="profile-header">
      <!-- Cover Image -->
      <div class="cover-image-wrapper">
        {#if privateProfileData.cover_url}
          <img
            src={privateProfileData.cover_url}
            alt="Cover"
            class="cover-image"
          />
        {:else}
          <div class="cover-placeholder"></div>
        {/if}
      </div>

      <!-- Profile Info Bar with Avatar -->
      <div class="profile-info-bar">
        <div class="profile-details">
          <div class="avatar-wrapper">
            <img
              src={privateProfileData.avatar_url || "/user.jpg"}
              alt={privateProfileData.username}
              class="profile-avatar"
            />
            {#if user && $chatStore.onlineUsers.has(user.id)}
              <span class="online-indicator"></span>
            {/if}
          </div>
          <div class="profile-text">
            <h1 class="username">{privateProfileData.username}</h1>
            <p class="user-handle">u/{privateProfileData.username}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Private Account Message -->
    <div class="private-message">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="48"
        height="48"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        class="lock-icon"
      >
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
        <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
      </svg>
      <h2>Đây là tài khoản riêng tư</h2>
      <p>
        Chỉ những người được chấp thuận mới có thể xem nội dung của tài khoản
        này.
      </p>
    </div>
  </div>
{:else if errorMessage}
  <div class="error-state">
    <p>{errorMessage}</p>
    <button class="retry-btn" onclick={loadUserProfile}>Thử lại</button>
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
        {#if isOwnProfile}
          <div class="cover-actions">
            <button
              class="change-cover-btn"
              onclick={() => coverFileInput?.click()}
              disabled={isUploadingCover}
              title="Đổi ảnh bìa"
            >
              {#if isUploadingCover}
                <div class="mini-spinner"></div>
              {:else}
                <img src="/change_profile_image.png" alt="Đổi ảnh bìa" />
              {/if}
            </button>
          </div>
        {/if}
      </div>
      <div class="profile-info-bar">
        <div class="profile-details">
          <div class="avatar-wrapper">
            <img
              src={user.profile.avatar?.url || "/user.jpg"}
              alt={user.username}
              class="profile-avatar"
            />
            {#if isOwnProfile}
              <button
                class="change-avatar-btn"
                onclick={() => avatarFileInput?.click()}
                disabled={isUploadingAvatar}
                title="Đổi ảnh đại diện"
              >
                {#if isUploadingAvatar}
                  <div class="mini-spinner"></div>
                {:else}
                  <img src="/change_profile_image.png" alt="Đổi ảnh đại diện" />
                {/if}
              </button>
            {/if}
          </div>
          <div class="profile-text">
            <h1 class="username">{user.username}</h1>
            <p class="user-handle">u/{user.username}</p>
          </div>
        </div>
        {#if isOwnProfile}
          <div class="profile-actions">
            <button class="action-btn primary" onclick={handleCreatePost}>
              <i class="fas fa-plus"></i>
              + Tạo bài viết
            </button>
            <button class="action-btn secondary" onclick={handleEditProfile}>
              Sửa hồ sơ
            </button>
            <button class="action-btn tertiary" onclick={handleSettings}>
              <img src="/dot.png" alt="Settings" class="settings-icon" />
            </button>
          </div>
        {:else}
          <div class="profile-actions">
            <button
              class="action-btn primary"
              onclick={handleSendMessage}
              disabled={isCreatingChannel}
            >
              {#if isCreatingChannel}
                <div class="mini-spinner"></div>
              {:else}
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 20 20"
                  fill="none"
                  style="margin-right: 6px;"
                >
                  <path
                    d="M17 9C17 13.4183 13.4183 17 9 17C7.87087 17 6.79301 16.7625 5.81818 16.3362L3 17L3.66379 14.1818C3.23749 13.207 3 12.1291 3 11C3 6.58172 6.58172 3 11 3C15.4183 3 19 6.58172 19 11"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
                Nhắn tin
              {/if}
            </button>
          </div>
        {/if}
      </div>
    </div>

    <div class="profile-content">
      <div class="profile-main-content">
        <div class="profile-tabs">
          <button
            class="tab-btn"
            class:active={activeTab === "posts"}
            onclick={() => (activeTab = "posts")}>Bài viết</button
          >
          {#if isOwnProfile}
            <button
              class="tab-btn"
              class:active={activeTab === "saved"}
              onclick={() => (activeTab = "saved")}>Đã lưu</button
            >
          {/if}

          <div class="sort-options">
            <select bind:value={sortBy}>
              <option value="hot">Nổi bật</option>
              <option value="newest">Mới nhất</option>
              <option value="oldest">Cũ nhất</option>
            </select>
          </div>
        </div>

        <div class="tab-content">
          {#if activeTab === "posts"}
            {#if isLoadingPosts}
              <div class="loading-posts">
                <div class="spinner"></div>
                <p>Đang tải bài viết...</p>
              </div>
            {:else if postsError}
              <div class="error-state">
                <p>{postsError}</p>
                <button class="retry-btn" onclick={loadUserPosts}
                  >Thử lại</button
                >
              </div>
            {:else if posts.length === 0}
              <div class="empty-state">
                <p>Chưa có bài viết</p>
              </div>
            {:else}
              <div class="post-list">
                {#each posts as post}
                  <Post {post} />
                {/each}
              </div>
            {/if}
          {:else if activeTab === "saved"}
            {#if isLoadingSaved}
              <div class="loading-posts">
                <div class="spinner"></div>
                <p>Đang tải bài viết đã lưu...</p>
              </div>
            {:else if savedError}
              <div class="error-state">
                <p>{savedError}</p>
              </div>
            {:else if savedPosts.length === 0}
              <div class="empty-state">
                <p>Chưa lưu bài viết nào</p>
                <p class="note">Các bài viết bạn lưu sẽ xuất hiện ở đây</p>
              </div>
            {:else}
              <div class="posts-list">
                {#each savedPosts as post}
                  <Post {post} />
                {/each}
              </div>
            {/if}
          {/if}
        </div>
      </div>
      <div class="profile-sidebar">
        <div class="user-card">
          <div class="user-card-body">
            <h3>Giới thiệu</h3>
            <p class="bio">
              {user.profile.bio || "Chưa có tiểu sử."}
            </p>

            <!-- Personal Info -->
            {#if user.profile.gender || user.profile.age || user.profile.location}
              <div class="personal-info">
                {#if user.profile.gender}
                  <div class="info-item">
                    <img
                      src="/gender_icon.svg"
                      alt="Gender"
                      width="16"
                      height="16"
                    />
                    <span
                      >{user.profile.gender === "male"
                        ? "Nam"
                        : user.profile.gender === "female"
                          ? "Nữ"
                          : "Không tiết lộ"}</span
                    >
                  </div>
                {/if}
                {#if user.profile.age}
                  <div class="info-item">
                    <img
                      src="/Calendar_duotone.svg"
                      alt="Age"
                      width="16"
                      height="16"
                    />
                    <span>{user.profile.age} tuổi</span>
                  </div>
                {/if}
                {#if user.profile.location}
                  <div class="info-item">
                    <img
                      src="/location_icon.svg"
                      alt="Location"
                      width="16"
                      height="16"
                    />
                    <span>{user.profile.location}</span>
                  </div>
                {/if}
              </div>
            {/if}

            <!-- Interests -->
            {#if user.profile.interests && user.profile.interests.length > 0}
              <div class="interests-section">
                <h4>Sở thích</h4>
                <div class="interests-tags">
                  {#each user.profile.interests as interest}
                    <span class="interest-tag">{interest}</span>
                  {/each}
                </div>
              </div>
            {/if}

            <!-- Social Links -->
            {#if user.profile.social_links && (user.profile.social_links.website || user.profile.social_links.facebook || user.profile.social_links.youtube || user.profile.social_links.github)}
              <div class="social-links-section">
                <h4>Liên kết</h4>
                <div class="social-links">
                  {#if user.profile.social_links.website}
                    <a
                      href={user.profile.social_links.website}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="social-link social-website"
                      title={user.profile.social_links.website}
                    >
                      <svg
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="currentColor"
                      >
                        <path
                          d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0zM4.5 7.5a.5.5 0 0 0 0 1h5.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 1 0-.708.708L10.293 7.5H4.5z"
                        />
                      </svg>
                      <span class="link-text"
                        >{formatSocialLink(
                          user.profile.social_links.website,
                        )}</span
                      >
                    </a>
                  {/if}
                  {#if user.profile.social_links.facebook}
                    <a
                      href={user.profile.social_links.facebook}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="social-link social-facebook"
                      title={user.profile.social_links.facebook}
                    >
                      <svg
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="currentColor"
                      >
                        <path
                          d="M16 8.049c0-4.446-3.582-8.05-8-8.05C3.58 0-.002 3.603-.002 8.05c0 4.017 2.926 7.347 6.75 7.951v-5.625h-2.03V8.05H6.75V6.275c0-2.017 1.195-3.131 3.022-3.131.876 0 1.791.157 1.791.157v1.98h-1.009c-.993 0-1.303.621-1.303 1.258v1.51h2.218l-.354 2.326H9.25V16c3.824-.604 6.75-3.934 6.75-7.951z"
                        />
                      </svg>
                      <span class="link-text"
                        >{formatSocialLink(
                          user.profile.social_links.facebook,
                        )}</span
                      >
                    </a>
                  {/if}
                  {#if user.profile.social_links.youtube}
                    <a
                      href={user.profile.social_links.youtube}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="social-link social-youtube"
                      title={user.profile.social_links.youtube}
                    >
                      <svg
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="currentColor"
                      >
                        <path
                          d="M8.051 1.999h.089c.822.003 4.987.033 6.11.335a2.01 2.01 0 0 1 1.415 1.42c.101.38.172.883.22 1.402l.01.104.022.26.008.104c.065.914.073 1.77.074 1.957v.075c-.001.194-.01 1.108-.082 2.06l-.008.105-.009.104c-.05.572-.124 1.14-.235 1.558a2.007 2.007 0 0 1-1.415 1.42c-1.16.312-5.569.334-6.18.335h-.142c-.309 0-1.587-.006-2.927-.052l-.17-.006-.087-.004-.171-.007-.171-.007c-1.11-.049-2.167-.128-2.654-.26a2.007 2.007 0 0 1-1.415-1.419c-.111-.417-.185-.986-.235-1.558L.09 9.82l-.008-.104A31.4 31.4 0 0 1 0 7.68v-.123c.002-.215.01-.958.064-1.778l.007-.103.003-.052.008-.104.022-.26.01-.104c.048-.519.119-1.023.22-1.402a2.007 2.007 0 0 1 1.415-1.42c.487-.13 1.544-.21 2.654-.26l.17-.007.172-.006.086-.003.171-.007A99.788 99.788 0 0 1 7.858 2h.193zM6.4 5.209v4.818l4.157-2.408L6.4 5.209z"
                        />
                      </svg>
                      <span class="link-text"
                        >{formatSocialLink(
                          user.profile.social_links.youtube,
                        )}</span
                      >
                    </a>
                  {/if}
                  {#if user.profile.social_links.github}
                    <a
                      href={user.profile.social_links.github}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="social-link social-github"
                      title={user.profile.social_links.github}
                    >
                      <svg
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="currentColor"
                      >
                        <path
                          d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z"
                        />
                      </svg>
                      <span class="link-text"
                        >{formatSocialLink(
                          user.profile.social_links.github,
                        )}</span
                      >
                    </a>
                  {/if}
                </div>
              </div>
            {/if}

            <div class="user-stats">
              <div class="stat">
                <span class="stat-value"
                  >{user.profile.stats?.post_count ?? 0}</span
                >
                <span class="stat-label">Bài viết</span>
              </div>
              <div class="stat">
                <span class="stat-value"
                  >{user.profile.stats?.comment_count ?? 0}</span
                >
                <span class="stat-label">Bình luận</span>
              </div>
              <div class="stat">
                <span class="stat-value">{user.reputation ?? 0}</span>
                <span class="stat-label">Danh tiếng</span>
              </div>
            </div>
            <div class="cake-day">
              <i class="fas fa-birthday-cake"></i>
              {user.profile.stats?.member_since || "Thành viên"}
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

  .cover-actions {
    position: absolute;
    top: 1rem;
    right: 1rem;
    display: flex;
    gap: 0.5rem;
  }

  .cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .change-cover-btn {
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

  .online-indicator {
    position: absolute;
    bottom: 12px;
    right: 12px;
    width: 24px;
    height: 24px;
    background: #46d160;
    border: 4px solid white;
    border-radius: 50%;
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

  .personal-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin: 1rem 0;
    padding: 1rem 0;
    border-top: 1px solid #e6e6e6;
  }

  .info-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #7c7c7c;
    font-size: 0.9rem;
  }

  .info-item img {
    flex-shrink: 0;
    opacity: 0.7;
  }

  .interests-section {
    margin: 1rem 0;
    padding: 1rem 0;
    border-top: 1px solid #e6e6e6;
  }

  .interests-section h4 {
    margin: 0 0 0.75rem 0;
    color: #1a1a1b;
    font-size: 0.9rem;
    font-weight: 600;
  }

  .interests-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .interest-tag {
    padding: 0.35rem 0.75rem;
    background-color: #f6f7f8;
    color: #153060;
    border-radius: 16px;
    font-size: 0.85rem;
    font-weight: 500;
  }

  .social-links-section {
    margin: 1rem 0;
    padding: 1rem 0;
    border-top: 1px solid #e6e6e6;
  }

  .social-links-section h4 {
    margin: 0 0 0.75rem 0;
    color: #1a1a1b;
    font-size: 0.9rem;
    font-weight: 600;
  }

  .social-links {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .social-link {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    color: #153060;
    text-decoration: none;
    border-radius: 4px;
    font-size: 0.9rem;
    font-weight: 500;
    transition: all 0.2s ease;
  }

  .social-link:hover {
    background-color: #f6f7f8;
  }

  .social-link svg {
    flex-shrink: 0;
  }

  /* Social media brand colors */
  .social-facebook svg {
    color: #1877f2;
  }

  .social-youtube svg {
    color: #ff0000;
  }

  .social-github svg {
    color: #181717;
  }

  .social-website svg {
    color: #0066cc;
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

  /* Private Profile Styles */
  .private-message {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 4rem 2rem;
    margin: 2rem 0;
    background-color: white;
    border-radius: 8px;
  }

  .lock-icon {
    color: #7c7c7c;
    margin-bottom: 1.5rem;
  }

  .private-message h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a1a1b;
    margin-bottom: 0.75rem;
  }

  .private-message p {
    font-size: 1rem;
    color: #7c7c7c;
    max-width: 400px;
    line-height: 1.5;
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

  .empty-state .note {
    font-size: 0.85rem;
    color: #999;
    font-style: italic;
    margin-top: 0.5rem;
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
