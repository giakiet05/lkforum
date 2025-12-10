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
  import { getNotifications } from "../services/notification-service";

  type TopbarProps = {
    user?: { name: string; avatar?: string; karma?: number };
    notificationCount?: number;
    onSearch?: (q: string) => void;
    onCreatePost?: () => void;
    onNotificationClick?: () => void;
    onLogout?: () => void;
  };

  let {
    user,
    notificationCount = 0,
    onSearch,
    onCreatePost,
    onNotificationClick,
    onLogout,
  }: TopbarProps = $props();

  let searchQuery = $state("");
  let showUserMenu = $state(false);
  let showNotifications = $state(false);
  let showAuthModal = $state(false);
  let showCreatePostModal = $state(false);
  let showChatPopup = $state(false);
  let dropdownElement = $state<HTMLDivElement | null>(null);
  let unreadNotificationCount = $state(0);

  // Load unread notification count on mount
  onMount(async () => {
    if (user) {
      try {
        const response = await getNotifications({ page: 1, limit: 1 });
        const unreadCount =
          response.notifications?.filter((n) => !n.is_read).length || 0;
        unreadNotificationCount = unreadCount;
      } catch (err) {
        console.error("Failed to load initial notification count:", err);
      }
    }
  });

  // Search dropdown state
  let showSearchDropdown = $state(false);
  let communityResults = $state<CommunityResponse[]>([]);
  let userResults = $state<UserResponse[]>([]);
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

      // Search both communities and users in parallel
      const [communitiesResponse, usersResponse] = await Promise.all([
        getCommunities({
          name: cleanQuery,
          limit: 3,
        }).catch(() => ({
          communities: [],
          pagination: { page: 1, page_size: 0, total_items: 0, total_pages: 0 },
        })),
        searchUsers(cleanQuery, 1, 3).catch(() => ({
          users: [],
          pagination: { page: 1, page_size: 0, total_items: 0, total_pages: 0 },
        })),
      ]);

      communityResults = communitiesResponse?.communities || [];
      userResults = usersResponse?.users || [];
      showSearchDropdown =
        communityResults.length > 0 || userResults.length > 0;
    } catch (error) {
      console.error("Search failed:", error);
      communityResults = [];
      userResults = [];
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
      onSearch?.(searchQuery);
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

  function handleCommunityClick(communityName: string) {
    console.log("🔍 Navigating to community:", communityName);
    push(`/lk/${communityName}`);
    searchQuery = "";
    communityResults = [];
    userResults = [];
    showSearchDropdown = false;
  }

  function handleUserClick(username: string) {
    push(`/profile/${username}`);
    searchQuery = "";
    communityResults = [];
    userResults = [];
    showSearchDropdown = false;
  }

  function handleSearchFocus() {
    if (
      searchQuery.trim() &&
      (communityResults.length > 0 || userResults.length > 0)
    ) {
      showSearchDropdown = true;
    }
  }

  function handleSearchBlur() {
    // Delay to allow click on dropdown items
    setTimeout(() => {
      showSearchDropdown = false;
    }, 200);
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
      <div class="brand" role="button" tabindex="0" onclick={() => {}}>
        <img src="/LKlogo.svg" alt="LKForum" class="brand-icon" />
        <span class="brand-name">LKForum</span>
      </div>
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
            placeholder="Search communities and users"
            bind:value={searchQuery}
            oninput={handleSearchInput}
            onkeydown={handleSearchKeydown}
            onfocus={handleSearchFocus}
            onblur={handleSearchBlur}
          />

          {#if showSearchDropdown && (communityResults.length > 0 || userResults.length > 0)}
            <div class="search-dropdown">
              {#if communityResults.length > 0}
                <div class="search-dropdown-header">Communities</div>
                {#each communityResults as community}
                  <button
                    class="search-result-item"
                    onclick={() => handleCommunityClick(community.name)}
                  >
                    <img
                      src={community.avatar || "/default-community.png"}
                      alt=""
                      class="result-avatar"
                    />
                    <div class="result-info">
                      <div class="result-name">c/{community.name}</div>
                      <div class="result-meta">
                        {community.member_count} members
                      </div>
                    </div>
                  </button>
                {/each}
              {/if}

              {#if userResults.length > 0}
                {#if communityResults.length > 0}
                  <div class="search-dropdown-divider"></div>
                {/if}
                <div class="search-dropdown-header">Users</div>
                {#each userResults as userResult}
                  <button
                    class="search-result-item"
                    onclick={() => handleUserClick(userResult.username)}
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
            </div>
          {/if}

          {#if isSearching}
            <div class="search-loading">Searching...</div>
          {/if}
        </div>
      </div>
    </div>

    <div class="topbar-right">
      <div class="topbar-actions">
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
          <span class="button-text">Create</span>
        </button>

        <button
          type="button"
          class="icon-button notification-btn"
          onclick={() => (showNotifications = !showNotifications)}
          title="Notifications"
          aria-label="Notifications"
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

        <!-- Messages Button -->
        <button
          class="icon-button"
          onclick={() => (showChatPopup = true)}
          title="Messages"
        >
          <img src="/message_icon.svg" alt="Messages" width="24" height="24" />
        </button>

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
                {#if user.avatar}
                  <img src={user.avatar} alt={user.name} class="avatar-image" />
                {:else}
                  <span class="avatar-fallback"
                    >{user.name.charAt(0).toUpperCase()}</span
                  >
                {/if}
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
                  Profile
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
                  Settings
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
                  Log Out
                </div>
              </div>
            {/if}
          </div>
        {:else}
          <button class="login-button" onclick={() => (showAuthModal = true)}
            >Log In</button
          >
          <AuthModal
            show={showAuthModal}
            onClose={() => (showAuthModal = false)}
          />
        {/if}
      </div>
    </div>
  </div>
</header>

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
    --topbar-search-background: rgba(33, 37, 41, 0.04);
    --topbar-search-border: rgba(33, 37, 41, 0.08);
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
    gap: 12px;
  }

  .topbar-center {
    flex: 1;
    display: flex;
    justify-content: center;
  }

  .topbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
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
    max-width: 680px;
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
    color: var(--muted-foreground);
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
    color: var(--topbar-foreground);
    outline: none;
  }

  .search-input:focus {
    background: var(--background);
    border-color: var(--topbar-accent);
  }

  .search-input::placeholder {
    color: var(--muted-foreground);
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
  }

  .icon-button:hover {
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

  @media (max-width: 640px) {
    .brand-name {
      display: none;
    }
    .topbar-search {
      max-width: none;
    }
    .button-text {
      display: none;
    }
  }
</style>
