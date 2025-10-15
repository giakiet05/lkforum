<script lang="ts">
  import { run } from "svelte/legacy";

  // Import các thành phần cần thiết
  import Sidebar from "../components/Sidebar.svelte"; // Updated import path
  import Topbar from "../components/Topbar.svelte";
  import { authStore } from "../stores/auth-store";
  import { logout } from "../services/auth-service";
  import { push } from "svelte-spa-router";
  import { get } from "svelte/store";
  import type { User } from "../models/user";

  const sidebarItems = [
    {
      id: "home",
      label: "Home",
      to: "/",
      icon: `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M3 11.5L12 4l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1v-8.5z" fill="currentColor"/>
    </svg>`,
    },
    { id: "popular", label: "Popular", to: "/popular", icon: "🔥" },
    {
      id: "community",
      label: "Community",
      icon: `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 2a5 5 0 100 10 5 5 0 000-10zm0 12c-4 0-8 2-8 4v2h16v-2c0-2-4-4-8-4z" fill="currentColor"/>
    </svg>`,
      children: [
        { id: "r-all", label: "r/all", to: "/r/all", icon: "🌐" },
        { id: "r-news", label: "r/news", to: "/r/news", icon: "📰" },
        { id: "r-dev", label: "r/dev", to: "/r/dev", icon: "💻" },
        // thêm sau: { id: 'r-your', label: 'r/your', to: '/r/your', icon: '...' }
      ],
    },
    { id: "explore", label: "Explore", to: "/explore", icon: "🧭" },
    { id: "all", label: "All", to: "/all", icon: "📋" },
  ];

  // Logic của trang Home giữ nguyên
  let user: User | null = null;
  // read authStore once on mount
  run(() => {
    const state = get(authStore);
    user = state.user;
  });

  // track whether sidebar is compact (bind from Sidebar)
  let isSidebarCompact = false;

  // pass a simple user object to Topbar (read from store on mount)
  let topbarUser:
    | { name: string; avatar?: string; karma?: number }
    | undefined = undefined;
  run(() => {
    const state = get(authStore);
    if (state.user) {
      topbarUser = {
        name: state.user.username || state.user.email || "User",
        avatar: "../avatar.jpg",
        karma: 20,
      };
    }
  });

  function handleLogout() {
    // logout(); // Tạm thời comment lại
    alert("Đăng xuất!");
    // push("/login");
  }

  function handleNavigate(item: any) {
    push(item.to);
  }
</script>

<div class="app-layout">
  <!-- Sidebar is fixed; main content reserves space using CSS variable -->
  <Topbar user={topbarUser} onLogout={handleLogout} />

  <Sidebar
    items={sidebarItems}
    onNavigate={handleNavigate}
    bind:compact={isSidebarCompact}
  />

  <main class="main-content" data-compact={isSidebarCompact}>
    <h1>Home Page</h1>
    {#if user}
      <p>Chào mừng trở lại!</p>
      <div class="user-info">
        <p><strong>Username:</strong> {user.username}</p>
        <p><strong>Email:</strong> {user.email}</p>
      </div>

      <button onclick={handleLogout}>Logout</button>
    {:else}
      <p>You are not logged in.</p>
      <a href="/#/login">Go to Login</a>
    {/if}
  </main>
</div>

<style>
  :root {
    --sidebar-width: 256px;
    --sidebar-compact-width: 64px;
  }

  .app-layout {
    position: relative;
    min-height: 100vh;
    background-color: #f0f2f5;
  }

  /* main-content reserves space equal to the sidebar width to avoid overlap */
  .main-content {
    margin-left: var(--sidebar-width); /* reserves space for the sidebar */
    transition: margin-left 0.2s ease;
    padding-top: calc(var(--topbar-height) + 1rem); /* reserve topbar space */
    padding-right: 2rem;
    padding-left: 2rem;
    padding-bottom: 2rem;
    min-height: 100vh;
    box-sizing: border-box;
    overflow-y: auto;
  }

  /* compact sidebar state (bind:compact -> data-compact attribute) */
  .main-content[data-compact="true"] {
    margin-left: var(--sidebar-compact-width);
  }

  /* small screens: sidebar becomes overlaying, so remove left margin */
  @media (max-width: 768px) {
    .main-content {
      margin-left: 0;
    }
  }

  .user-info {
    background-color: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    max-width: 400px;
    margin: 1rem 0;
  }

  .user-info p {
    margin: 0.5rem 0;
  }

  button {
    /* Style đơn giản cho nút logout */
    padding: 0.6rem 1.2rem;
    border: none;
    background-color: #dc3545;
    color: white;
    border-radius: 5px;
    cursor: pointer;
  }
</style>
