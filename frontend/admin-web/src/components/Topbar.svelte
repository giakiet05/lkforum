<script lang="ts">
  interface TopbarProps {
    user?: { name: string; avatar?: string };
    onLogout?: () => void;
  }

  let { user, onLogout }: TopbarProps = $props();

  let showUserMenu = $state(false);

  function toggleUserMenu() {
    showUserMenu = !showUserMenu;
  }

  function handleLogout() {
    showUserMenu = false;
    onLogout?.();
  }

  // Close dropdown when clicking outside
  function handleClickOutside(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest(".user-menu-container")) {
      showUserMenu = false;
    }
  }

  $effect(() => {
    if (showUserMenu) {
      document.addEventListener("click", handleClickOutside);
      return () => {
        document.removeEventListener("click", handleClickOutside);
      };
    }
  });
</script>

<header class="topbar">
  <div class="topbar-container">
    <div class="topbar-left">
      <div class="brand">
        <svg
          width="32"
          height="32"
          viewBox="0 0 32 32"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <rect width="32" height="32" rx="8" fill="#FF4500" />
          <text
            x="16"
            y="22"
            text-anchor="middle"
            font-size="18"
            font-weight="bold"
            fill="white">A</text
          >
        </svg>
        <span class="brand-text">LKForum Admin</span>
      </div>
    </div>

    <div class="topbar-center"></div>

    <div class="topbar-right">
      {#if user}
        <div class="user-menu-container">
          <button class="user-button" onclick={toggleUserMenu}>
            {#if user.avatar}
              <img src={user.avatar} alt={user.name} class="user-avatar" />
            {:else}
              <div class="user-avatar-placeholder">
                {user.name.charAt(0).toUpperCase()}
              </div>
            {/if}
            <span class="user-name">{user.name}</span>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path
                d="M4 6L8 10L12 6"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>

          {#if showUserMenu}
            <div class="user-dropdown">
              <button class="dropdown-item" onclick={handleLogout}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                  <polyline
                    points="16 17 21 12 16 7"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                  <line
                    x1="21"
                    y1="12"
                    x2="9"
                    y2="12"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
                Đăng xuất
              </button>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
</header>

<style>
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
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    padding: 0 16px;
    max-width: 100%;
  }

  .topbar-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 18px;
    font-weight: 700;
    color: var(--topbar-foreground);
  }

  .brand-text {
    white-space: nowrap;
  }

  .topbar-center {
    flex: 1;
  }

  .topbar-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-menu-container {
    position: relative;
  }

  .user-button {
    all: unset;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    border-radius: 20px;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  .user-button:hover {
    background: var(--topbar-hover);
  }

  .user-avatar,
  .user-avatar-placeholder {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }

  .user-avatar-placeholder {
    background: var(--primary-color);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 14px;
  }

  .user-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--topbar-foreground);
  }

  .user-dropdown {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    min-width: 180px;
    padding: 8px;
    z-index: 1000;
  }

  .dropdown-item {
    all: unset;
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 10px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    color: var(--text-color);
    transition: background 0.15s ease;
  }

  .dropdown-item:hover {
    background: var(--topbar-hover);
  }

  .dropdown-item svg {
    flex-shrink: 0;
  }
</style>
