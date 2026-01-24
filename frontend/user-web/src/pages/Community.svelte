<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import Post from "../components/Post.svelte";
  import CreatePostModal from "../components/CreatePostModal.svelte";
  import { toastStore } from "../stores/toast-store";
  import type { PostResponse } from "../dtos/post-dto";
  import type { CommunityResponse } from "../dtos/community-dto";
  import {
    getCommunityByName,
    addModerators,
    updateCommunity,
  } from "../services/community-service";
  import { getPosts } from "../services/post-service";
  import { authStore } from "../stores/auth-store";
  import {
    checkMembership as checkMembershipAPI,
    createMembership,
    deleteMembership,
  } from "../services/membership-service";
  import { getUserByUsername } from "../services/user-service";

  type CommunityProps = {
    params?: { name: string };
  };

  let { params = { name: "sveltejs" } }: CommunityProps = $props();

  let activeSort = $state<"hot" | "new" | "top">("hot");
  let activeTimeFrame = $state<
    "hour" | "day" | "week" | "month" | "year" | "all"
  >("day");
  let isJoined = $state(false);
  let isCheckingMembership = $state(false);
  let isTogglingMembership = $state(false);
  let expandedRules = $state(new Set<number>());
  let showCreatePostModal = $state(false);
  let showInviteModModal = $state(false);
  let inviteUsername = $state("");
  let invitePermission = $state("Everything");
  let inviteCanEdit = $state(true);

  // Callback when new post is created
  function handlePostCreated() {
    console.log("📨 handlePostCreated called - reloading posts...");
    loadPosts(true); // Reload posts to show new post
  }

  // API data
  let community = $state<CommunityResponse | null>(null);
  let posts = $state<PostResponse[]>([]);
  let isLoadingCommunity = $state(true);
  let isLoadingPosts = $state(true);
  let loadingMorePosts = $state(false);
  let communityError = $state<string | null>(null);
  let postsError = $state<string | null>(null);
  let currentPage = $state(1);
  const postsPerPage = 20;
  let hasMorePosts = $state(true);
  let postsSentinel: HTMLDivElement | null = null;

  // Check if current user is the creator of this community
  const isCreator = $derived(() => {
    const currentUser = $authStore.user;
    if (!currentUser || !community) return false;
    return community.create_by_id === currentUser.id;
  });

  // Check if current user is a moderator of this community
  const isModerator = $derived(() => {
    const currentUser = $authStore.user;
    if (!currentUser || !community) return false;
    return (
      community.moderators?.some((mod) => mod.user_id === currentUser.id) ||
      false
    );
  });

  // Check if current user can access mod tools (creator or moderator)
  const canUseModeTools = $derived(() => isCreator() || isModerator());

  // Upload state
  let isUploadingBanner = $state(false);
  let isUploadingAvatar = $state(false);
  let bannerFileInput: HTMLInputElement | undefined = $state();
  let avatarFileInput: HTMLInputElement | undefined = $state();

  async function handleBannerChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file || !community) return;

    if (!file.type.startsWith("image/")) {
      toastStore.warning("Vui lòng chọn tệp hình ảnh");
      return;
    }

    if (file.size > 100 * 1024 * 1024) {
      toastStore.warning("Kích thước ảnh phải nhỏ hơn 100MB");
      return;
    }

    try {
      isUploadingBanner = true;
      const reader = new FileReader();
      reader.onload = async (e) => {
        const base64 = e.target?.result as string;
        const updated = await updateCommunity({
          id: community!.id,
          banner: base64,
        });
        community = updated;
        toastStore.success("Đã cập nhật ảnh bìa!");
      };
      reader.readAsDataURL(file);
    } catch (error) {
      console.error("Failed to upload banner:", error);
      toastStore.error("Không thể tải ảnh bìa. Vui lòng thử lại.");
    } finally {
      isUploadingBanner = false;
      input.value = "";
    }
  }

  async function handleAvatarChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file || !community) return;

    if (!file.type.startsWith("image/")) {
      toastStore.warning("Vui lòng chọn tệp hình ảnh");
      return;
    }

    if (file.size > 2 * 1024 * 1024) {
      toastStore.warning("Kích thước ảnh phải nhỏ hơn 2MB");
      return;
    }

    try {
      isUploadingAvatar = true;
      const reader = new FileReader();
      reader.onload = async (e) => {
        const base64 = e.target?.result as string;
        const updated = await updateCommunity({
          id: community!.id,
          avatar: base64,
        });
        community = updated;
        toastStore.success("Đã cập nhật ảnh đại diện!");
      };
      reader.readAsDataURL(file);
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      toastStore.error("Không thể tải ảnh đại diện. Vui lòng thử lại.");
    } finally {
      isUploadingAvatar = false;
      input.value = "";
    }
  }

  // Load community data
  async function loadCommunity() {
    try {
      isLoadingCommunity = true;
      communityError = null;

      community = await getCommunityByName(params.name);
    } catch (error) {
      console.error("Failed to load community:", error);
      communityError =
        error instanceof Error ? error.message : "Failed to load community";
      community = null;
    } finally {
      isLoadingCommunity = false;
    }
  }

  // Load posts for this community
  async function loadPosts(reset = true) {
    if (!community) {
      console.warn("⚠️ loadPosts called but community is null");
      return;
    }

    if (reset) {
      isLoadingPosts = true;
      currentPage = 1;
      posts = [];
      hasMorePosts = true;
    } else {
      loadingMorePosts = true;
    }

    try {
      console.log(
        `🔍 Loading posts for community: ${community.name} (${community.id}), page: ${currentPage}`,
      );
      postsError = null;

      const fetchedPosts = await getPosts({
        community_id: community.id,
        sort: activeSort,
        time: activeSort === "top" ? activeTimeFrame : undefined,
        page: currentPage,
        limit: postsPerPage,
      });

      console.log("📦 Fetched posts response:", fetchedPosts);
      console.log(`✅ Loaded ${fetchedPosts.length} posts`);

      if (reset) {
        posts = fetchedPosts;
      } else {
        posts = [...posts, ...fetchedPosts];
      }

      hasMorePosts = fetchedPosts.length === postsPerPage;
    } catch (error) {
      console.error("❌ Failed to load posts:", error);
      postsError =
        error instanceof Error ? error.message : "Failed to load posts";
      if (reset) {
        posts = [];
      }
    } finally {
      isLoadingPosts = false;
      loadingMorePosts = false;
    }
  }

  async function loadMorePosts() {
    if (loadingMorePosts || !hasMorePosts) return;
    currentPage++;
    await loadPosts(false);
  }

  // Check if user is member of this community
  async function checkMembership() {
    const currentUser = $authStore.user;
    if (!currentUser || !community) {
      isJoined = false;
      return;
    }

    // Creators are automatically members
    if (isCreator()) {
      isJoined = true;
      console.log("✅ You are the creator - automatically joined");
      return;
    }

    try {
      isCheckingMembership = true;
      isJoined = await checkMembershipAPI(currentUser.id, community.id);
      console.log(
        `✅ Membership status: ${isJoined ? "Joined" : "Not joined"}`,
      );
    } catch (error) {
      console.error("❌ Failed to check membership:", error);
      isJoined = false;
    } finally {
      isCheckingMembership = false;
    }
  }

  // Reload posts when sort or time frame changes
  $effect(() => {
    if (community && !isLoadingCommunity) {
      loadPosts(true);
    }
  });

  // Reload community when params.name changes (user navigates to different community)
  $effect(() => {
    if (params.name) {
      loadCommunity();
    }
  });

  // Check membership when community loads
  $effect(() => {
    if (community && !isLoadingCommunity) {
      checkMembership();
    }
  });

  // Intersection Observer for infinite scroll
  $effect(() => {
    if (!postsSentinel || isLoadingPosts) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMorePosts && !loadingMorePosts) {
          loadMorePosts();
        }
      },
      { threshold: 0.1 },
    );

    observer.observe(postsSentinel);

    return () => {
      observer.disconnect();
    };
  });

  /* Old mock data
  const posts: PostResponse[] = [
    {
      id: "1",
      type: "text",
      community: params.name,
      author: "user123",
      time: "4 hours ago",
      title: "Welcome to " + params.name + "!",
      upvotes: 123,
      downvotes: 5,
      commentsCount: 42,
      content:
        "This is the community page. Join us to see more content and participate in discussions!",
    },
    {
      id: "2",
      type: "image",
      community: params.name,
      author: "photographer",
      time: "8 hours ago",
      title: "Check out this amazing photo!",
      upvotes: 456,
      downvotes: 12,
      commentsCount: 89,
      images: ["/GirlFromNowhere.jpg"],
    },
  ];
  */

  async function toggleJoin() {
    const currentUser = $authStore.user;
    if (!currentUser || !community) {
      toastStore.warning("Vui lòng đăng nhập để tham gia cộng đồng");
      return;
    }

    // Prevent creator from leaving their own community
    if (isCreator()) {
      toastStore.warning("Bạn không thể rời khỏi cộng đồng do bạn tạo!");
      return;
    }

    try {
      isTogglingMembership = true;

      if (isJoined) {
        // Leave community
        await deleteMembership(currentUser.id, community.id);
        console.log(`✅ Left community: ${community.name}`);
        isJoined = false;
        // Decrement member count locally
        if (community.member_count > 0) {
          community.member_count--;
        }
        toastStore.success(`Đã rời khỏi cộng đồng ${community.name}!`);
      } else {
        // Join community
        const membership = await createMembership(currentUser.id, community.id);
        console.log(`✅ Joined community: ${community.name}`, membership);

        // Check if membership is pending approval
        if (membership?.status === "pending") {
          toastStore.success(
            "Đã gửi yêu cầu tham gia! Vui lòng chờ quản trị viên duyệt.",
          );
        } else {
          isJoined = true;
          // Increment member count locally
          community.member_count++;
          toastStore.success(`Đã tham gia cộng đồng ${community.name}!`);
        }
      }
    } catch (error: any) {
      console.error("❌ Failed to toggle membership:", error);
      // Handle specific error cases
      if (error?.message?.includes("đang chờ duyệt")) {
        toastStore.warning("Bạn đã gửi yêu cầu tham gia và đang chờ duyệt.");
      } else {
        toastStore.error(
          error instanceof Error
            ? error.message
            : "Failed to update membership",
        );
      }
    } finally {
      isTogglingMembership = false;
    }
  }

  function toggleRule(ruleIndex: number) {
    if (expandedRules.has(ruleIndex)) {
      expandedRules.delete(ruleIndex);
    } else {
      expandedRules.add(ruleIndex);
    }
    expandedRules = new Set(expandedRules);
  }

  function handleModTools() {
    push(`/lk/${params.name}/mod`);
  }

  function handleSettings() {
    push(`/lk/${params.name}/settings`);
  }

  function handleOpenInviteModModal() {
    showInviteModModal = true;
    inviteUsername = "";
    invitePermission = "Everything";
    inviteCanEdit = true;
  }

  function handleCloseInviteModModal() {
    showInviteModModal = false;
    inviteUsername = "";
    invitePermission = "Everything";
    inviteCanEdit = true;
  }

  async function handleInviteMod() {
    if (!inviteUsername.trim() || !community) {
      toastStore.warning("Vui lòng nhập tên người dùng!");
      return;
    }

    console.log("Invite mod:", {
      username: inviteUsername,
      permission: invitePermission,
      canEdit: inviteCanEdit,
    });

    try {
      console.log("🔍 Looking up user:", inviteUsername.trim());
      const user = await getUserByUsername(inviteUsername.trim());
      console.log("✅ User found:", user);

      console.log("📤 Adding moderator to community:", community.id);
      await addModerators({
        id: community.id,
        added_moderator: [user.id], // Just send user ID
      });

      console.log("✅ Moderator added successfully!");
      toastStore.success(`Đã thêm ${inviteUsername} làm quản trị viên!`);
      handleCloseInviteModModal();

      // Reload community to get updated moderators
      await loadCommunity();
    } catch (error) {
      console.error("❌ Failed to add moderator:", error);
      toastStore.error(
        "Không thể thêm quản trị viên. Vui lòng kiểm tra tên người dùng và thử lại.",
      );
    }
  }

  onMount(() => {
    window.scrollTo(0, 0);
    // Don't call loadCommunity here - let $effect handle it
  });
</script>

<div class="community-page">
  {#if isLoadingCommunity}
    <div class="loading-container">
      <p>Đang tải cộng đồng...</p>
    </div>
  {:else if communityError || !community}
    <div class="error-container">
      <p>{communityError || "Không tìm thấy cộng đồng"}</p>
      <button onclick={() => push("/")}>Quay về trang chủ</button>
    </div>
  {:else}
    <!-- Hidden file inputs -->
    <input
      type="file"
      accept="image/*"
      bind:this={bannerFileInput}
      onchange={handleBannerChange}
      style="display: none;"
    />
    <input
      type="file"
      accept="image/*"
      bind:this={avatarFileInput}
      onchange={handleAvatarChange}
      style="display: none;"
    />

    <!-- Banner -->
    <div class="community-banner">
      {#if community.banner}
        <img src={community.banner} alt="Community banner" />
      {:else}
        <div class="banner-placeholder"></div>
      {/if}
      {#if canUseModeTools()}
        <button
          class="change-banner-btn"
          onclick={() => bannerFileInput?.click()}
          disabled={isUploadingBanner}
          title="Đổi ảnh bìa"
        >
          {#if isUploadingBanner}
            <div class="mini-spinner"></div>
          {:else}
            <img src="/change_profile_image.png" alt="Change banner" />
          {/if}
        </button>
      {/if}
    </div>

    <!-- Community Header -->
    <div class="community-header">
      <div class="community-header-content">
        <div class="community-info">
          <div class="community-icon-wrapper">
            <img
              src={community.avatar || "/community_logo.jpg"}
              alt="Community icon"
              class="community-icon"
            />
            {#if canUseModeTools()}
              <button
                class="change-avatar-btn"
                onclick={() => avatarFileInput?.click()}
                disabled={isUploadingAvatar}
                title="Đổi ảnh đại diện"
              >
                {#if isUploadingAvatar}
                  <div class="mini-spinner"></div>
                {:else}
                  <img src="/change_profile_image.png" alt="Change avatar" />
                {/if}
              </button>
            {/if}
          </div>
          <div class="community-title">
            <div class="title-row">
              <h1>lk/{community.name}</h1>
              {#if community.is_18_plus}
                <span class="badge-18plus">18+</span>
              {/if}
            </div>
            <p class="community-name">{community.name}</p>
          </div>
        </div>

        <div class="community-actions">
          <!-- Join/Leave Button -->
          <button
            class="join-btn"
            class:joined={isJoined}
            disabled={isCheckingMembership || isTogglingMembership}
            onclick={toggleJoin}
          >
            {#if isTogglingMembership}
              {isJoined ? "Đang rời..." : "Đang tham gia..."}
            {:else if isJoined}
              Đã tham gia
            {:else}
              Tham gia
            {/if}
          </button>

          <!-- Create Post Button - Only show if joined -->
          {#if isJoined}
            <button
              class="create-post-action-btn"
              title="Tạo bài viết"
              onclick={() => {
                if (!$authStore.user) {
                  toastStore.error("Vui lòng đăng nhập để tạo bài viết");
                  return;
                }
                showCreatePostModal = true;
              }}
            >
              <svg viewBox="0 0 20 20" fill="currentColor">
                <path
                  d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
                />
              </svg>
              Tạo bài viết
            </button>
          {/if}

          <!-- Notification Bell Button -->
          <button class="action-btn" title="Thông báo">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path
                d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z"
              />
            </svg>
          </button>

          <!-- Mod Tools Button - Only show if creator or moderator -->
          {#if canUseModeTools()}
            <button
              class="mod-tools-btn"
              title="Công cụ quản trị"
              onclick={handleModTools}
            >
              <svg viewBox="0 0 20 20" fill="currentColor">
                <path
                  fill-rule="evenodd"
                  d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z"
                  clip-rule="evenodd"
                />
              </svg>
              Quản trị
            </button>

            <!-- More Options Button -->
            <button
              class="action-btn"
              title="Thêm tùy chọn"
              onclick={handleSettings}
            >
              <svg viewBox="0 0 20 20" fill="currentColor">
                <path
                  d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"
                />
              </svg>
            </button>
          {/if}
        </div>
      </div>
    </div>

    <div class="community-container">
      <!-- Main Content -->
      <div class="community-main">
        <!-- Sort Bar -->
        <div class="sorting-bar">
          <button
            class="sort-btn"
            class:active={activeSort === "hot"}
            onclick={() => (activeSort = "hot")}
          >
            Nổi bật
          </button>
          <button
            class="sort-btn"
            class:active={activeSort === "new"}
            onclick={() => (activeSort = "new")}
          >
            Mới nhất
          </button>
          <button
            class="sort-btn"
            class:active={activeSort === "top"}
            onclick={() => (activeSort = "top")}
          >
            Hàng đầu
          </button>
        </div>

        <!-- Posts -->
        <div class="post-list">
          {#if isLoadingPosts}
            <div class="loading-posts">
              <p>Đang tải bài viết...</p>
            </div>
          {:else if postsError}
            <div class="error-posts">
              <p>{postsError}</p>
            </div>
          {:else if posts.length === 0}
            <div class="no-posts">
              <p>Chưa có bài viết nào. Hãy là người đầu tiên đăng bài!</p>
            </div>
          {:else}
            {#each posts as post}
              <Post {post} />
            {/each}

            <!-- Sentinel for infinite scroll -->
            {#if hasMorePosts}
              <div class="posts-sentinel" bind:this={postsSentinel}>
                {#if loadingMorePosts}
                  <div class="loading-more">
                    <div class="spinner"></div>
                    <p>Đang tải thêm...</p>
                  </div>
                {/if}
              </div>
            {:else if posts.length > 0}
              <div class="end-message">
                <p>Bạn đã xem hết bài viết</p>
              </div>
            {/if}
          {/if}
        </div>
      </div>

      <!-- Sidebar -->
      <div class="community-sidebar">
        <div class="about-card">
          <h3>Giới thiệu cộng đồng</h3>
          <p class="about-description">
            {community.description || "Chưa có mô tả"}
          </p>

          <div class="community-stats">
            <div class="stat">
              <div class="stat-value">
                {community.member_count?.toLocaleString() || "0"}
              </div>
              <div class="stat-label">Thành viên</div>
            </div>
            <div class="stat">
              <!-- TODO: Backend doesn't track online count yet -->
              <div class="stat-value">-</div>
              <div class="stat-label">Trực tuyến</div>
            </div>
          </div>

          <!-- TODO: Backend doesn't have created_at field in CommunityResponse -->
        </div>

        <!-- Rules Card -->
        <div class="rules-card">
          <h3>Nội quy cộng đồng</h3>
          {#if community.rules && community.rules.length > 0}
            <div class="rules-accordion">
              {#each community.rules as rule, index}
                <div class="rule-item">
                  <button
                    class="rule-header"
                    onclick={() => toggleRule(index)}
                    aria-expanded={expandedRules.has(index)}
                  >
                    <span class="rule-number">{index + 1}.</span>
                    <span class="rule-title">{rule.title}</span>
                    <span
                      class="rule-toggle"
                      class:expanded={expandedRules.has(index)}
                    >
                      <svg
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="none"
                      >
                        <path
                          d="M4 6L8 10L12 6"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </span>
                  </button>
                  {#if expandedRules.has(index)}
                    <div class="rule-content">
                      <p style="color: #666;">{rule.description}</p>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {:else}
            <p class="no-rules">Chưa có nội quy</p>
          {/if}
        </div>

        <!-- Moderators Card -->
        <div class="moderators-card">
          <div class="moderators-header">
            <h3>Quản trị viên</h3>
            {#if isCreator()}
              <button
                class="invite-mod-btn"
                title="Mời quản trị viên"
                onclick={handleOpenInviteModModal}
              >
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    d="M8 9a3 3 0 100-6 3 3 0 000 6zM8 11a6 6 0 016 6H2a6 6 0 016-6zm10-5a1 1 0 10-2 0v1h-1a1 1 0 100 2h1v1a1 1 0 102 0v-1h1a1 1 0 100-2h-1V6z"
                  />
                </svg>
                Mời quản trị
              </button>
            {/if}
          </div>
          <div class="moderator-list">
            {#if community.moderators && community.moderators.filter((m) => m.is_active).length > 0}
              {#each community.moderators.filter((m) => m.is_active) as moderator}
                <div class="moderator">
                  <img
                    src={moderator.avatar?.url || "/user.jpg"}
                    alt="Moderator"
                    class="mod-avatar"
                  />
                  <span class="mod-name">u/{moderator.username}</span>
                </div>
              {/each}
            {:else}
              <p class="no-moderators">Chưa có quản trị viên</p>
            {/if}
          </div>
          {#if community.moderators && community.moderators.filter((m) => m.is_active).length > 3}
            <button class="view-all-mods-btn">Xem tất cả quản trị viên</button>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

{#if community}
  <CreatePostModal
    show={showCreatePostModal}
    onClose={() => {
      showCreatePostModal = false;
    }}
    onPostCreated={handlePostCreated}
    communityName={community.name}
  />
{/if}

{#if showInviteModModal}
  <div class="modal-overlay" onclick={handleCloseInviteModModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>Mời quản trị viên</h2>

      <div class="form-group">
        <div class="search-input-wrapper">
          <img
            src="/searchuser_icon.svg"
            alt="Search"
            class="search-icon"
            width="20"
            height="20"
          />
          <input
            type="text"
            placeholder="Tìm kiếm người dùng"
            bind:value={inviteUsername}
            class="search-input"
          />
        </div>
        <p class="hint">Nhập tên người dùng để tìm kiếm</p>
      </div>

      <div class="form-group">
        <label>Quyền hạn</label>
        <select bind:value={invitePermission} class="permission-select">
          <option value="Everything">Toàn bộ</option>
          <option value="Manage Posts & Comments"
            >Quản lý bài viết & bình luận</option
          >
          <option value="Manage Users">Quản lý người dùng</option>
          <option value="Manage Settings">Quản lý cài đặt</option>
        </select>
      </div>

      <div class="form-group">
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={inviteCanEdit} />
          <span>Bạn có thể chỉnh sửa quản trị viên này</span>
        </label>
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseInviteModModal}>
          Hủy
        </button>
        <button class="action-btn-primary" onclick={handleInviteMod}>
          Mời
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .community-page {
    min-height: 100vh;
    background-color: white;
  }

  /* Loading and Error States */
  .loading-container,
  .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    padding: 40px;
    text-align: center;
  }

  .error-container p {
    color: #ff4500;
    margin-bottom: 16px;
    font-size: 16px;
  }

  .error-container button {
    padding: 8px 16px;
    background-color: #0079d3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 600;
  }

  .error-container button:hover {
    background-color: #0060a8;
  }

  .loading-posts,
  .error-posts,
  .no-posts {
    padding: 40px;
    text-align: center;
    background-color: white;
    border-radius: 4px;
    border: 1px solid #ccc;
  }

  .error-posts p {
    color: #ff4500;
  }

  .no-posts p {
    color: #7c7c7c;
    font-size: 14px;
  }

  .no-moderators {
    padding: 12px;
    color: #7c7c7c;
    font-size: 13px;
    text-align: center;
  }

  /* Banner */
  .community-banner {
    width: calc(100% - 48px);
    height: 200px;
    overflow: hidden;
    background: linear-gradient(to bottom, #33a8ff, #0079d3);
    border-radius: 8px;
    margin: 8px 24px;
  }

  .community-banner img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .banner-placeholder {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  /* Community Header */
  .community-header {
    background-color: white;
    border-bottom: 1px solid #edeff1;
    padding: 16px 0;
  }

  .community-header-content {
    max-width: 100%;
    padding: 0 24px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }

  .community-info {
    display: flex;
    align-items: center;
    gap: 16px;
    flex: 1;
  }

  .community-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    border: 1.5px solid #ccc;
    background: white;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn:hover {
    background: #f6f7f8;
    border-color: #999;
  }

  .action-btn svg {
    width: 20px;
    height: 20px;
  }

  .create-post-action-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    background: white;
    color: #1c1c1c;
    border: 1.5px solid #ccc;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s;
  }

  .create-post-action-btn:hover {
    background: #f6f7f8;
    border-color: #999;
  }

  .create-post-action-btn svg {
    width: 18px;
    height: 18px;
  }

  .mod-tools-btn {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;
    white-space: nowrap;
  }

  .mod-tools-btn:hover {
    background: var(--darkblue--);
  }

  .settings-btn {
    background: white;
    color: var(--blue--);
    border: 1px solid var(--blue--);
    padding: 8px 16px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;
    white-space: nowrap;
  }

  .settings-btn:hover {
    background: #f6f7f8;
  }

  .settings-btn svg {
    width: 16px;
    height: 16px;
  }

  .community-icon-wrapper {
    position: relative;
  }

  .community-icon {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    border: 4px solid white;
    background: white;
    margin-top: -60px;
    object-fit: cover;
  }

  .change-avatar-btn {
    position: absolute;
    bottom: 5px;
    right: 5px;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: transform 0.2s;
  }

  .change-avatar-btn:hover {
    transform: scale(1.1);
  }

  .change-avatar-btn img {
    width: 20px;
    height: 20px;
  }

  .change-banner-btn {
    position: absolute;
    bottom: 10px;
    right: 10px;
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: transform 0.2s;
  }

  .change-banner-btn:hover {
    transform: scale(1.1);
  }

  .change-banner-btn img {
    width: 22px;
    height: 22px;
  }

  .mini-spinner {
    width: 16px;
    height: 16px;
    border: 2px solid #ccc;
    border-top-color: #0079d3;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .community-banner {
    position: relative;
  }

  .community-title h1 {
    font-size: 28px;
    font-weight: 700;
    margin: 0;
    color: #1c1c1c;
  }

  .title-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .badge-18plus {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 4px 12px;
    background: linear-gradient(135deg, #ff4757 0%, #ff6348 100%);
    color: white;
    font-size: 13px;
    font-weight: 700;
    border-radius: 4px;
    border: 2px solid #ff3838;
    box-shadow: 0 2px 4px rgba(255, 69, 87, 0.3);
    letter-spacing: 0.5px;
  }

  .community-name {
    font-size: 14px;
    color: #7c7c7c;
    margin: 4px 0 0 0;
  }

  .join-btn {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 8px 24px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .join-btn:hover {
    background: var(--darkblue--);
  }

  .join-btn.joined {
    background: #edeff1;
    color: #1c1c1c;
    border: 1px solid #edeff1;
  }

  .join-btn.joined:hover {
    background: #d7dadc;
  }

  /* Container */
  .community-container {
    max-width: 100%;
    margin: 0 auto;
    padding: 24px;
    display: grid;
    grid-template-columns: 1fr 312px;
    gap: 24px;
  }

  /* Main Content */
  .community-main {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .sorting-bar {
    background-color: #f6f7f8;
    border-radius: 4px;
    padding: 4px;
    display: flex;
    gap: 0px;
  }

  .sort-btn {
    flex: 1;
    padding: 8px 12px;
    border: none;
    background-color: transparent;
    color: #878a8c;
    font-weight: bold;
    border-radius: 4px;
    cursor: pointer;
    transition:
      background-color 0.2s,
      color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .sort-btn.active {
    background-color: white;
    color: var(--blue--);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .sort-btn:not(.active):hover {
    background-color: #e9ebee;
  }

  .post-list {
    display: flex;
    flex-direction: column;
  }

  .posts-sentinel {
    min-height: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .loading-more {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 24px;
    text-align: center;
  }

  .loading-more p {
    color: #5a5a5a;
    font-size: 14px;
    margin: 12px 0 0 0;
  }

  .loading-more .spinner {
    border: 3px solid #f3f3f3;
    border-top: 3px solid var(--blue--);
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

  .end-message {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32px 24px;
    text-align: center;
  }

  .end-message p {
    color: #878a8c;
    font-size: 14px;
    font-weight: 500;
  }

  /* Sidebar */
  .community-sidebar {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .about-card,
  .rules-card,
  .moderators-card {
    background: var(--table-bg);
    border: none;
    border-radius: 4px;
    padding: 12px;
    color: var(--grayfont);
  }

  .about-card h3,
  .rules-card h3,
  .moderators-card h3 {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #1c1c1c;
    margin: 0 0 8px 0;
    padding-bottom: 12px;
    border-bottom: 1px solid #edeff1;
  }

  .about-description {
    font-size: 14px;
    line-height: 21px;
    color: var(--grayfont);
    margin: 12px 0;
  }

  .community-stats {
    display: flex;
    gap: 24px;
    padding: 12px 0;
    border-top: 1px solid #edeff1;
    border-bottom: 1px solid #edeff1;
    margin: 12px 0;
  }

  .stat {
    display: flex;
    flex-direction: column;
  }

  .stat-value {
    font-size: 16px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .stat-label {
    font-size: 12px;
    color: var(--grayfont);
  }

  .community-created {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--grayfont);
    font-size: 14px;
    margin: 12px 0;
  }

  .community-created img {
    width: 20px;
    height: 20px;
  }

  /* Rules */
  .rules-accordion {
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .rule-item {
    border-bottom: 1px solid #edeff1;
  }

  .rule-item:last-child {
    border-bottom: none;
  }

  .rule-header {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 0;
    background: transparent;
    border: none;
    cursor: pointer;
    text-align: left;
    transition: background 0.2s;
  }

  .rule-header:hover {
    background: #f6f7f8;
  }

  .rule-number {
    font-size: 14px;
    font-weight: 700;
    color: #1c1c1c;
    min-width: 24px;
  }

  .rule-title {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
    color: var(--grayfont);
  }

  .rule-toggle {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    transition: transform 0.2s ease;
    color: var(--grayfont);
    transform: rotate(-90deg);
  }

  .rule-toggle.expanded {
    transform: rotate(0deg);
  }

  .rule-content {
    padding: 0 0 12px 32px;
    animation: slideDown 0.2s ease;
  }

  .rule-content p {
    margin: 0;
    font-size: 13px;
    line-height: 20px;
    color: var(--grayfont);
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-4px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .no-rules {
    font-size: 14px;
    color: var(--grayfont);
    text-align: center;
    padding: 16px 0;
    margin: 0;
  }

  /* Moderators */
  .moderators-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .moderators-header h3 {
    margin: 0;
    padding: 0;
    border: none;
  }

  .invite-mod-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    background: transparent;
    border: 1px solid #edeff1;
    border-radius: 9999px;
    font-size: 12px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
    font-family: "Roboto", sans-serif;
  }

  .invite-mod-btn:hover {
    background: #f6f7f8;
    border-color: #d7dadc;
  }

  .invite-mod-btn svg {
    width: 16px;
    height: 16px;
  }

  .moderator-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 12px;
  }

  .moderator {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .mod-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }

  .mod-name {
    font-size: 14px;
    color: var(--grayfont);
    font-weight: 500;
  }

  .mod-name:hover {
    text-decoration: underline;
    cursor: pointer;
  }

  .view-all-mods-btn {
    width: 100%;
    padding: 8px 16px;
    background: var(--button-secondary-bg);
    border: none;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
    font-family: "Roboto", sans-serif;
  }

  .view-all-mods-btn:hover {
    background: rgba(214, 216, 222, 0.6);
  }

  /* Modal Styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    background: white;
    padding: 32px;
    border-radius: 12px;
    max-width: 500px;
    width: 90%;
    max-height: 90vh;
    overflow-y: auto;
  }

  .modal-content h2 {
    margin: 0 0 24px 0;
    font-size: 20px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .form-group {
    margin-bottom: 16px;
  }

  .form-group label {
    display: block;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    margin-bottom: 8px;
  }

  .search-input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-input-wrapper .search-icon {
    position: absolute;
    left: 16px;
    pointer-events: none;
    opacity: 0.6;
  }

  .search-input {
    width: 100%;
    padding: 12px 16px 12px 48px;
    border: 1px solid #edeff1;
    border-radius: 8px;
    font-size: 14px;
    background: #f6f7f8;
    color: #1c1c1c;
  }

  .search-input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .permission-select {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #edeff1;
    border-radius: 8px;
    font-size: 14px;
    background: #f6f7f8;
    color: #1c1c1c;
    cursor: pointer;
  }

  .permission-select:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    user-select: none;
  }

  .checkbox-label input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
    accent-color: var(--blue--);
  }

  .checkbox-label span {
    font-size: 14px;
    color: #1c1c1c;
  }

  .hint {
    font-size: 12px;
    color: var(--grayfont);
    margin: 6px 0 0 0;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }

  .btn-cancel {
    padding: 10px 20px;
    background: var(--button-secondary-bg);
    color: var(--blue--);
    border: none;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-cancel:hover {
    background: rgba(214, 216, 222, 0.6);
  }

  .action-btn-primary {
    padding: 10px 20px;
    background: var(--blue--);
    color: white;
    border: none;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn-primary:hover {
    filter: brightness(0.85);
  }

  /* Responsive */
  @media (max-width: 960px) {
    .community-container {
      grid-template-columns: 1fr;
      gap: 16px;
    }

    .community-sidebar {
      order: -1;
    }
  }

  @media (max-width: 768px) {
    .community-page {
      padding: 0;
    }

    .community-banner {
      height: 120px;
    }

    .community-header {
      padding: 12px 16px;
    }

    .community-avatar {
      width: 60px;
      height: 60px;
      margin-top: -30px;
    }

    .community-info {
      gap: 8px;
    }

    .community-name {
      font-size: 18px;
    }

    .community-description {
      font-size: 13px;
    }

    .community-container {
      padding: 0;
      gap: 12px;
    }

    .community-sidebar {
      display: none;
    }

    .sort-tabs {
      gap: 8px;
      padding: 0;
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
    }

    .sort-tab {
      padding: 6px 12px;
      font-size: 13px;
      white-space: nowrap;
    }
  }

  @media (max-width: 480px) {
    .community-banner {
      height: 100px;
    }

    .community-header {
      padding: 8px 12px;
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .community-avatar {
      width: 50px;
      height: 50px;
      margin-top: -25px;
    }

    .community-name {
      font-size: 16px;
    }

    .join-button,
    .settings-button {
      padding: 6px 16px;
      font-size: 13px;
    }
  }
</style>
