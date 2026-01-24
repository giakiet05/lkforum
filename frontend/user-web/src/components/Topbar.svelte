<script lang="ts">
  import { onMount } from "svelte";
  import AuthModal from "./AuthModal.svelte";
  import CreatePostModal from "./CreatePostModal.svelte";
  import ChatPopup from "./ChatPopup.svelte";
  import NotificationsDropdown from "./NotificationsDropdown.svelte";
  import { push } from "svelte-spa-router";
  import { getCommunities } from "../services/community-service";
  import type { CommunityResponse } from "../dtos/community-dto";
  import { searchUsers } from "../services/user-service";
  import type { UserResponse } from "../dtos/user-dto";
  import { getPosts } from "../services/post-service";
  import type { PostResponse } from "../dtos/post-dto";
  import { getNotifications } from "../services/notification-service";
  import { totalUnreadCount } from "../stores/chat-store";
  import { generatePostUrl } from "../utils/slug";

  type TopbarProps = {
    user?: { name: string; avatar?: string; karma?: number };
    notificationCount?: number;
    onSearch?: (q: string) => void;
    onCreatePost?: () => void;
    onNotificationClick?: () => void;
    onLogout?: () => void;
    onMenuClick?: () => void;
    forceShowAuthModal?: boolean;
    onAuthModalClose?: () => void;
  };

  let {
    user,
    notificationCount = 0,
    onSearch,
    onCreatePost,
    onNotificationClick,
    onLogout,
    onMenuClick,
    forceShowAuthModal = false,
    onAuthModalClose,
  }: TopbarProps = $props();

  let searchQuery = $state("");
  let showUserMenu = $state(false);
  let showNotifications = $state(false);
  let showAuthModal = $state(false);
  let showCreatePostModal = $state(false);
  let showChatPopup = $state(false);
  let dropdownElement = $state<HTMLDivElement | null>(null);
  let unreadNotificationCount = $state(0);
  let unreadMessageCount = $state(0);

  // Subscribe to total unread count
  $effect(() => {
    unreadMessageCount = $totalUnreadCount;
  });

  // Watch forceShowAuthModal prop and open modal when true
  $effect(() => {
    if (forceShowAuthModal) {
      showAuthModal = true;
    }
  });

  // Load unread notification count on mount and when user changes
  async function loadUnreadNotificationCount() {
    if (user) {
      try {
        // Load more notifications to get accurate unread count
        const response = await getNotifications({ page: 1, pageSize: 100 });
        const unreadCount =
          response.notifications?.filter((n) => !n.is_read).length || 0;
        unreadNotificationCount = unreadCount;
        console.log(
          `🔔 [Topbar] Loaded unread notification count: ${unreadCount}`,
        );
      } catch (err) {
        console.error("Failed to load initial notification count:", err);
      }
    } else {
      unreadNotificationCount = 0;
    }
  }

  // Reload unread count when user changes (login/logout)
  $effect(() => {
    loadUnreadNotificationCount();
  });

  onMount(async () => {
    // Listen for open-chat event from Profile page
    const handleOpenChat = (event: CustomEvent) => {
      showChatPopup = true;
    };
    window.addEventListener("open-chat", handleOpenChat as EventListener);

    return () => {
      window.removeEventListener("open-chat", handleOpenChat as EventListener);
    };
  });

  // Search dropdown state
  let showSearchDropdown = $state(false);
  let communityResults = $state<CommunityResponse[]>([]);
  let userResults = $state<UserResponse[]>([]);
  let postResults = $state<PostResponse[]>([]);
  let isSearching = $state(false);
  let searchTimeout: number | null = null;

  async function performSearch(query: string) {
    if (!query.trim()) {
      communityResults = [];
      userResults = [];
      showSearchDropdown = false;
      return;
    }

    try {
      isSearching = true;
      const cleanQuery = query.replace(/^(c\/|lk\/|u\/)/i, "").trim();
      if (!cleanQuery) {
        communityResults = [];
        userResults = [];
        showSearchDropdown = false;
        isSearching = false;
        return;
      }

      // Search communities, users, and posts in parallel
      const [communitiesResponse, usersResponse, postsResponse] =
        await Promise.all([
          getCommunities({
            name: cleanQuery,
            limit: 3,
          }).catch(() => ({
            communities: [],
            pagination: {
              page: 1,
              page_size: 0,
              total_items: 0,
              total_pages: 0,
            },
          })),
          searchUsers(cleanQuery, 1, 3).catch(() => ({
            users: [],
            pagination: {
              page: 1,
              page_size: 0,
              total_items: 0,
              total_pages: 0,
            },
          })),
          getPosts({
            search: cleanQuery,
            limit: 3,
          }).catch(() => []),
        ]);

      communityResults = communitiesResponse?.communities || [];
      userResults = usersResponse?.users || [];
      postResults = postsResponse || [];
      showSearchDropdown =
        communityResults.length > 0 ||
        userResults.length > 0 ||
        postResults.length > 0;
    } catch (error) {
      console.error("Search failed:", error);
      communityResults = [];
      userResults = [];
      postResults = [];
      showSearchDropdown = false;
    } finally {
      isSearching = false;
    }
  }

  function handleSearchInput() {
    // Debounce search
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    searchTimeout = window.setTimeout(() => {
      performSearch(searchQuery);
    }, 300);
  }

  function handleSearch() {
    if (searchQuery.trim()) {
      // Navigate to search page to show full results
      push(`/search?q=${encodeURIComponent(searchQuery)}`);
      showSearchDropdown = false;
    }
  }

  function handleSearchKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSearch();
    } else if (e.key === "Escape") {
      showSearchDropdown = false;
    }
  }

  function handleCommunityClick(e: MouseEvent, communityName: string) {
    e.preventDefault();
    e.stopPropagation();
    push(`/lk/${communityName}`);
    searchQuery = "";
    communityResults = [];
    userResults = [];
    postResults = [];
    showSearchDropdown = false;
  }

  function handleUserClick(e: MouseEvent, username: string) {
    e.preventDefault();
    e.stopPropagation();
    push(`/profile/${username}`);
    searchQuery = "";
    communityResults = [];
    userResults = [];
    postResults = [];
    showSearchDropdown = false;
  }

  function handlePostClick(e: MouseEvent, postId: string, postTitle: string) {
    e.preventDefault();
    e.stopPropagation();
    console.log("🔍 Navigating to post:", postId, postTitle);

    if (!postId) {
      console.error("❌ Post ID is undefined!");
      return;
    }

    const url = generatePostUrl(postId, postTitle);
    console.log("🔗 Pushing URL:", url);
    push(url);

    searchQuery = "";
    communityResults = [];
    userResults = [];
    postResults = [];
    showSearchDropdown = false;
  }

  function handleSearchFocus() {
    if (
      searchQuery.trim() &&
      (communityResults.length > 0 ||
        userResults.length > 0 ||
        postResults.length > 0)
    ) {
      showSearchDropdown = true;
    }
  }

  function handleSearchBlur() {
    // Delay to allow click on dropdown items
    setTimeout(() => {
      showSearchDropdown = false;
    }, 300);
  }

  function toggleUserMenu() {
    showUserMenu = !showUserMenu;
  }

  function closeUserMenu() {
    showUserMenu = false;
  }

  function handleOverlayKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") closeUserMenu();
  }

  function handleOverlayClick(e: MouseEvent) {
    if (dropdownElement && dropdownElement.contains(e.target as Node)) {
      return;
    }
    closeUserMenu();
  }

  function handleLogoutClick() {
    console.log("Logout clicked!");
    onLogout?.();
    closeUserMenu();
  }

  function handleCreatePostClick() {
    if (!user) {
      // Chưa đăng nhập - hiển thị modal đăng nhập
      showAuthModal = true;
      return;
    }
    showCreatePostModal = true;
    onCreatePost?.();
  }

  function handleNavigation() {
    console.log("Navigation clicked!");
    closeUserMenu();
  }

  $effect(() => {
    if (showUserMenu) {
      document.addEventListener("keydown", handleOverlayKeydown);
      // Thêm click outside handler
      const handleClickOutside = (e: MouseEvent) => {
        const target = e.target as Node;
        if (dropdownElement && !dropdownElement.contains(target)) {
          const userButton = document.querySelector(".user-button");
          if (userButton && !userButton.contains(target)) {
            closeUserMenu();
          }
        }
      };
      setTimeout(() => {
        document.addEventListener("click", handleClickOutside);
      }, 0);

      return () => {
        document.removeEventListener("keydown", handleOverlayKeydown);
        document.removeEventListener("click", handleClickOutside);
      };
    }
  });
</script>

<header class="topbar">
  <div class="topbar-container">
    <div class="topbar-left">
      <!-- Hamburger menu for mobile -->
      <button
        class="menu-button"
        onclick={() => onMenuClick?.()}
        aria-label="Menu"
      >
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
          <path
            d="M3 12H21M3 6H21M3 18H21"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
      </button>

      <div
        class="brand"
        role="button"
        tabindex="0"
        onclick={() => (window.location.href = "/")}
        style="cursor: pointer;"
      >
        <img src="/LKlogo.svg" alt="LKForum" class="brand-icon" />
        <span class="brand-name">LKForum</span>
      </div>

      <!-- Search icon for mobile -->
      <button
        class="search-icon-button"
        onclick={() => push("/search")}
        aria-label="Tìm kiếm"
      >
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
          <circle
            cx="8.5"
            cy="8.5"
            r="5.5"
            stroke="currentColor"
            stroke-width="2"
          />
          <path
            d="M12.5 12.5L17 17"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
      </button>
    </div>

    <div class="topbar-center">
      <div class="topbar-search" role="search">
        <div class="search-wrapper">
          <svg
            class="search-icon"
            width="20"
            height="20"
            viewBox="0 0 20 20"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <circle
              cx="8.5"
              cy="8.5"
              r="5.5"
              stroke="currentColor"
              stroke-width="2"
            />
            <path
              d="M12.5 12.5L17 17"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
          <input
            class="search-input"
            type="text"
            placeholder="Tìm cộng đồng và người dùng"
            bind:value={searchQuery}
            oninput={handleSearchInput}
            onkeydown={handleSearchKeydown}
            onfocus={handleSearchFocus}
            onblur={handleSearchBlur}
          />

          {#if showSearchDropdown && (communityResults.length > 0 || userResults.length > 0 || postResults.length > 0)}
            <div class="search-dropdown">
              {#if communityResults.length > 0}
                <div class="search-dropdown-header">Cộng đồng</div>
                {#each communityResults as community}
                  <button
                    class="search-result-item"
                    onmousedown={(e) => handleCommunityClick(e, community.name)}
                  >
                    <img
                      src={community.avatar || "/default-community.png"}
                      alt=""
                      class="result-avatar"
                    />
                    <div class="result-info">
                      <div class="result-name">lk/{community.name}</div>
                      <div class="result-meta">
                        {community.member_count} thành viên
                      </div>
                    </div>
                  </button>
                {/each}
              {/if}

              {#if userResults.length > 0}
                {#if communityResults.length > 0}
                  <div class="search-dropdown-divider"></div>
                {/if}
                <div class="search-dropdown-header">Người dùng</div>
                {#each userResults as userResult}
                  <button
                    class="search-result-item"
                    onmousedown={(e) => handleUserClick(e, userResult.username)}
                  >
                    {#if userResult.profile?.avatar?.url}
                      <img
                        src={userResult.profile.avatar.url}
                        alt=""
                        class="result-avatar"
                      />
                    {:else}
                      <div class="result-avatar fallback">
                        {userResult.username.charAt(0).toUpperCase()}
                      </div>
                    {/if}
                    <div class="result-info">
                      <div class="result-name">u/{userResult.username}</div>
                      <div class="result-meta">
                        {userResult.reputation} karma • {userResult.title}
                      </div>
                    </div>
                  </button>
                {/each}
              {/if}

              {#if postResults.length > 0}
                {#if communityResults.length > 0 || userResults.length > 0}
                  <div class="search-dropdown-divider"></div>
                {/if}
                <div class="search-dropdown-header">Bài viết</div>
                {#each postResults as post}
                  <button
                    class="search-result-item post-result"
                    onmousedown={(e) =>
                      handlePostClick(e, post.id, post.title || "post")}
                  >
                    <div class="result-info">
                      <div class="result-name">
                        {post.title || "Bài viết không tiêu đề"}
                      </div>
                      <div class="result-meta">
                        lk/{post.community.name} • {post.author.username}
                      </div>
                    </div>
                  </button>
                {/each}
              {/if}
            </div>
          {/if}

          {#if isSearching}
            <div class="search-loading">Đang tìm kiếm...</div>
          {/if}
        </div>
      </div>
    </div>

    <div class="topbar-right">
      <div class="topbar-actions">
        {#if user}
          <button
            type="button"
            class="create-button"
            onclick={handleCreatePostClick}
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path
                d="M8 3V13M3 8H13"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              />
            </svg>
            <span class="button-text">Tạo</span>
          </button>

          <button
            type="button"
            class="icon-button message-btn"
            onclick={() => (showChatPopup = !showChatPopup)}
            title="Tin nhắn"
            aria-label="Tin nhắn"
          >
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <path
                d="M17 9C17 13.4183 13.4183 17 9 17C7.87087 17 6.79301 16.7625 5.81818 16.3362L3 17L3.66379 14.1818C3.23749 13.207 3 12.1291 3 11C3 6.58172 6.58172 3 11 3C15.4183 3 19 6.58172 19 11"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            {#if unreadMessageCount > 0}
              <span class="notification-badge">{unreadMessageCount}</span>
            {/if}
          </button>

          <button
            type="button"
            class="icon-button notification-btn"
            onclick={() => (showNotifications = !showNotifications)}
            title="Thông báo"
            aria-label="Thông báo"
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
            {#if unreadNotificationCount > 0}
              <span class="notification-badge">{unreadNotificationCount}</span>
            {/if}

            <NotificationsDropdown
              show={showNotifications}
              onClose={() => (showNotifications = false)}
              onUnreadCountChange={(count) => (unreadNotificationCount = count)}
            />
          </button>
        {:else}
          <!-- Login button -->
          <button class="login-button" onclick={() => (showAuthModal = true)}>
            Đăng nhập
          </button>
        {/if}

        {#if user}
          <div class="user-menu-wrapper">
            <div
              class="user-button"
              role="button"
              tabindex="0"
              onclick={toggleUserMenu}
              onkeydown={(e) =>
                (e.key === "Enter" || e.key === " ") && toggleUserMenu()}
              aria-haspopup="true"
              aria-expanded={showUserMenu}
            >
              <div class="user-avatar">
                <img
                  src={user.avatar || "/user.jpg"}
                  alt={user.name}
                  class="avatar-image"
                />
              </div>
              <div class="user-info">
                <span class="user-name">{user.name}</span>
                {#if user.karma !== undefined}
                  <span class="user-karma">{user.karma} karma</span>
                {/if}
              </div>
              <svg
                class="chevron-icon"
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
            </div>

            {#if showUserMenu}
              <div
                class="user-dropdown"
                role="menu"
                aria-label="User menu"
                bind:this={dropdownElement}
              >
                <div
                  class="dropdown-item"
                  role="menuitem"
                  onclick={() => {
                    console.log("Profile clicked!");
                    closeUserMenu();
                    push("/profile");
                  }}
                >
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <circle
                      cx="8"
                      cy="5"
                      r="2.5"
                      stroke="currentColor"
                      stroke-width="1.5"
                    />
                    <path
                      d="M3 14C3 11.7909 4.79086 10 7 10H9C11.2091 10 13 11.7909 13 14"
                      stroke="currentColor"
                      stroke-width="1.5"
                      stroke-linecap="round"
                    />
                  </svg>
                  Hồ sơ
                </div>

                <div
                  class="dropdown-item"
                  role="menuitem"
                  onclick={() => {
                    console.log("Settings clicked!");
                    closeUserMenu();
                    push("/settings");
                  }}
                >
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <circle
                      cx="8"
                      cy="8"
                      r="2"
                      stroke="currentColor"
                      stroke-width="1.5"
                    />
                    <path
                      d="M8 1V3M8 13V15M15 8H13M3 8H1M12.5 3.5L11 5M5 11L3.5 12.5M12.5 12.5L11 11M5 5L3.5 3.5"
                      stroke="currentColor"
                      stroke-width="1.5"
                      stroke-linecap="round"
                    />
                  </svg>
                  Cài đặt
                </div>

                <div class="dropdown-separator" role="separator"></div>

                <div
                  class="dropdown-item"
                  role="menuitem"
                  onclick={() => {
                    console.log("Logout clicked!");
                    handleLogoutClick();
                  }}
                >
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path
                      d="M6 14H3C2.44772 14 2 13.5523 2 13V3C2 2.44772 2.44772 2 3 2H6M11 11L14 8M14 8L11 5M14 8H6"
                      stroke="currentColor"
                      stroke-width="1.5"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                  Đăng xuất
                </div>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
</header>

<!-- Auth Modal for both authenticated and unauthenticated users -->
<AuthModal
  show={showAuthModal}
  onClose={() => {
    showAuthModal = false;
    if (onAuthModalClose) {
      onAuthModalClose();
    }
  }}
/>

<CreatePostModal
  show={showCreatePostModal}
  onClose={() => (showCreatePostModal = false)}
/>

<ChatPopup show={showChatPopup} onClose={() => (showChatPopup = false)} />

<style>
  :root {
    --topbar-height: 56px;
    --topbar-background: #ffffff;
    --topbar-border: #e6e9ee;
    --topbar-foreground: #213547;
    --topbar-accent: #ff8a00;
    --topbar-accent-hover: #ff7a00;
    --topbar-search-background: #ececec;
    --topbar-search-border: #ececec;
    --topbar-search-focus-background: #ececec;
    --topbar-search-focus-border: #4a70a9;
    --topbar-search-icon: #000000;
    --topbar-search-text: #000000;
    --muted-foreground: #9aa4b2;
    --background: #ffffff;
    --border: #e6e9ee;
  }

  .topbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: var(--topbar-height);
    background: var(--topbar-background);
    border-bottom: 1px solid var(--topbar-border);
    z-index: 200;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  }

  .topbar-container {
    max-width: 100%;
    height: 100%;
    padding: 0 16px;
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .topbar-left {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 0 0 auto;
  }

  .menu-button {
    display: none;
    width: 40px;
    height: 40px;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    color: var(--topbar-foreground);
    transition: background 0.2s ease;
  }

  .menu-button:hover {
    background: var(--topbar-search-background);
  }

  .topbar-center {
    flex: 1;
    display: flex;
    justify-content: center;
    max-width: 600px;
    margin: 0 auto;
    padding: 0 16px;
  }

  .topbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 0 0 auto;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    color: var(--topbar-foreground);
  }

  .brand-icon {
    width: 40px;
    height: 40px;
    object-fit: contain;
    display: block;
  }

  .brand-name {
    font-size: 18px;
    font-weight: 700;
    color: var(--topbar-foreground);
  }

  .topbar-search {
    width: 100%;
    max-width: 500px;
  }

  .search-wrapper {
    position: relative;
    width: 100%;
  }

  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--topbar-search-icon);
    pointer-events: none;
  }

  .search-input {
    width: 100%;
    padding: 8px 16px 8px 40px;
    background: var(--topbar-search-background);
    border: 1px solid var(--topbar-search-border);
    border-radius: 20px;
    height: 40px;
    font-size: 14px;
    color: var(--topbar-search-text);
    outline: none;
    transition:
      background 0.2s ease,
      border-color 0.2s ease;
  }

  .search-input:focus {
    background: var(--topbar-search-focus-background);
    border-color: var(--topbar-search-focus-border);
    border-width: 2px;
  }

  .search-input::placeholder {
    color: var(--topbar-search-text);
    opacity: 0.7;
  }

  .search-dropdown {
    position: absolute;
    top: calc(100% + 8px);
    left: 0;
    right: 0;
    background: var(--background);
    border: 1px solid var(--border);
    border-radius: 12px;
    box-shadow:
      0 8px 24px rgba(0, 0, 0, 0.12),
      0 2px 6px rgba(0, 0, 0, 0.08);
    overflow: hidden;
    max-height: 480px;
    overflow-y: auto;
    z-index: 1000;
  }

  .search-dropdown-header {
    padding: 12px 16px 8px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #000000;
    background: var(--background);
    position: sticky;
    top: 0;
    z-index: 1;
  }

  .search-dropdown-divider {
    height: 1px;
    background: var(--border);
    margin: 8px 0;
  }

  .search-result-item {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 10px 16px;
    background: transparent;
    border: none;
    cursor: pointer;
    transition: background-color 0.15s ease;
    text-align: left;
  }

  .search-result-item:hover {
    background: var(--topbar-search-background);
  }

  .result-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    object-fit: cover;
    background: var(--topbar-search-background);
    flex-shrink: 0;
    border: 1px solid var(--border);
  }

  .result-avatar.fallback {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--topbar-accent);
    color: white;
    font-weight: 700;
    font-size: 15px;
  }

  .result-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 3px;
    min-width: 0;
  }

  .result-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--topbar-foreground);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.3;
  }

  .result-meta {
    font-size: 12px;
    color: var(--muted-foreground);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2;
  }

  .post-result .result-info {
    gap: 4px;
  }

  .post-result .result-name {
    white-space: normal;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    line-height: 1.4;
  }

  .search-loading {
    position: absolute;
    top: calc(100% + 8px);
    left: 0;
    right: 0;
    padding: 12px 16px;
    background: var(--background);
    border: 1px solid var(--border);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    font-size: 14px;
    color: var(--muted-foreground);
    text-align: center;
  }

  .topbar-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .create-button {
    display: flex;
    align-items: center;
    gap: 6px;
    background: rgba(214, 216, 222, 0.4);
    color: #000000;
    border: none;
    border-radius: 20px;
    padding: 0 12px;
    height: 36px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .create-button:hover {
    background-color: rgba(214, 216, 222, 0.6);
  }

  .button-text {
    display: none;
  }

  .icon-button {
    position: relative;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    background: transparent;
    color: var(--topbar-foreground);
  }

  .icon-button:hover {
    background: var(--topbar-search-background);
  }

  .menu-dots {
    color: var(--topbar-foreground);
  }

  .search-icon-button {
    display: none;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border: none;
    background: transparent;
    color: var(--topbar-foreground);
    cursor: pointer;
    border-radius: 8px;
    transition: background 0.2s;
  }

  .search-icon-button:hover {
    background: var(--topbar-search-background);
  }

  .notification-badge {
    position: absolute;
    top: 4px;
    right: 4px;
    min-width: 18px;
    height: 18px;
    padding: 0 4px;
    background: #ff4444;
    color: white;
    font-size: 10px;
    font-weight: 700;
    border-radius: 9px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .user-menu-wrapper {
    position: relative;
  }

  .user-button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 8px 4px 4px;
    border: 1px solid var(--topbar-border);
    border-radius: 8px;
    cursor: pointer;
    background: transparent;
    transition: background-color 0.2s;
  }

  .user-button:hover {
    background: var(--topbar-search-background);
  }

  .user-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    overflow: hidden;
    background: var(--topbar-accent);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .avatar-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-fallback {
    color: white;
    font-weight: 600;
    font-size: 14px;
  }

  .user-info {
    display: none;
    flex-direction: column;
    align-items: flex-start;
  }

  .user-name {
    font-size: 13px;
    font-weight: 600;
    color: var(--topbar-foreground);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-karma {
    font-size: 11px;
    color: var(--muted-foreground);
  }

  .chevron-icon {
    display: none;
    color: var(--muted-foreground);
  }

  .user-dropdown {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    min-width: 200px;
    background: var(--background);
    border: 1px solid var(--border);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    padding: 4px;
    z-index: 302;
    pointer-events: auto;
  }

  .dropdown-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    background: transparent;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    color: var(--topbar-foreground);
    text-align: left;
    text-decoration: none;
    transition: background-color 0.2s;
  }

  .dropdown-item:hover {
    background: var(--topbar-search-background);
  }

  .dropdown-separator {
    height: 1px;
    background: var(--border);
    margin: 4px 0;
  }

  .login-button {
    padding: 8px 16px;
    border: none;
    border-radius: 20px;
    background: var(--blue--);
    cursor: pointer;
    text-decoration: none;
    color: white;
    font-weight: 500;
    transition: opacity 0.2s;
  }

  .login-button:hover {
    opacity: 0.9;
  }

  @media (min-width: 640px) {
    .button-text {
      display: inline;
    }
    .user-info {
      display: flex;
    }
    .chevron-icon {
      display: block;
    }
  }

  @media (max-width: 1024px) {
    .menu-button {
      display: flex;
    }
  }

  @media (max-width: 640px) {
    .topbar-container {
      padding: 0 8px;
      gap: 8px;
    }

    .menu-button {
      display: flex;
    }

    .search-icon-button {
      display: flex;
    }

    .brand-name {
      display: none;
    }

    .brand-icon {
      width: 32px;
      height: 32px;
    }

    .topbar-center {
      display: none;
    }

    .topbar-right {
      gap: 4px;
      margin-left: auto;
    }

    .button-text {
      display: none;
    }

    .topbar-btn {
      padding: 6px;
      min-width: 36px;
    }

    .user-avatar {
      width: 32px;
      height: 32px;
    }
  }

  @media (max-width: 480px) {
    .topbar-container {
      padding: 0 4px;
      gap: 4px;
    }

    .topbar-btn {
      padding: 4px;
      min-width: 32px;
    }

    .create-post-btn {
      padding: 6px 12px;
      font-size: 13px;
    }
  }
</style>
