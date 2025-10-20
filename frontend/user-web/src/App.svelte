<script lang="ts">
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import Topbar from "./components/Topbar.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import { authStore } from "./stores/auth-store";
  import { logout } from "./services/auth-service";
  import { push } from "svelte-spa-router";
  import { get } from "svelte/store";
  import type { User } from "./models/user";
  import { run } from "svelte/legacy";

  const sidebarItems = [
    {
      id: "home",
      label: "Home",
      to: "/",
      icon: `<svg viewBox=\"0 0 24 24\" width=\"20\" height=\"20\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">
      <path d=\"M3 11.5L12 4l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1v-8.5z\" fill=\"currentColor\"/>
    </svg>`,
    },
    {
      id: "popular",
      label: "Popular",
      to: "/popular",
      icon: `<svg width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><g clip-path=\"url(#clip0_20_157)\"><path d=\"M23.2499 12.751L12.7769 23.25\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M17.25 12.751H23.25V18.75\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M18.75 0.75V5.25H12.75V11.25H6.75V17.25H0.75V23.25\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/></g><defs><clipPath id=\"clip0_20_157\"><rect width=\"24\" height=\"24\" fill=\"white\"/></clipPath></defs></svg>`,
    },
    {
      id: "community",
      label: "Community",
      icon: `<svg viewBox=\"0 0 24 24\" width=\"20\" height=\"20\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">
      <path d=\"M12 2a5 5 0 100 10 5 5 0 000-10zm0 12c-4 0-8 2-8 4v2h16v-2c0-2-4-4-8-4z\" fill=\"currentColor\"/>
    </svg>`,
      children: [
        { id: "lk-all", label: "lk/all", to: "/lk/all", icon: "🌐" },
        { id: "lk-news", label: "lk/news", to: "/lk/news", icon: "📰" },
        { id: "lk-dev", label: "lk/dev", to: "/lk/dev", icon: "💻" },
      ],
    },
    { id: "explore", label: "Explore", to: "/explore", icon: "🧭" },
    {
      id: "all",
      label: "All",
      to: "/all",
      icon: `<svg width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><g clip-path=\"url(#clip0_20_165)\"><path d=\"M23.25 21.25C23.25 22.35 22.35 23.25 21.25 23.25H2.75C1.65 23.25 0.75 22.35 0.75 21.25V2.75C0.75 1.65 1.65 0.75 2.75 0.75H21.25C22.35 0.75 23.25 1.65 23.25 2.75V21.25Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M7.44971 18.4499C7.44971 19.0499 7.04971 19.4499 6.44971 19.4499H5.44971C4.84971 19.4499 4.44971 19.0499 4.44971 18.4499V15.6499C4.44971 15.0499 4.84971 14.6499 5.44971 14.6499H6.44971C7.04971 14.6499 7.44971 15.0499 7.44971 15.6499V18.4499Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M13.4502 18.4499C13.4502 19.0499 13.0502 19.4499 12.4502 19.4499H11.4502C10.8502 19.4499 10.4502 19.0499 10.4502 18.4499V6.6499C10.4502 6.0499 10.8502 5.6499 11.4502 5.6499H12.4502C13.0502 5.6499 13.4502 6.0499 13.4502 6.6499V18.4499Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M19.4502 18.4499C19.4502 19.0499 19.0502 19.4499 18.4502 19.4499H17.4502C16.8502 19.4499 16.4502 19.0499 16.4502 18.4499V11.6499C16.4502 11.0499 16.8502 10.6499 17.4502 10.6499H18.4502C19.0502 10.6499 19.4502 11.0499 19.4502 11.6499V18.4499Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/></g><defs><clipPath id=\"clip0_20_165\"><rect width=\"24\" height=\"24\" fill=\"white\"/></clipPath></defs></svg>`,
    },
  ];

  let user: User | null = null;
  run(() => {
    const state = get(authStore);
    user = state.user;
  });

  let isSidebarCompact = false;

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
    alert("Đăng xuất!");
    logout();
  }

  function handleNavigate(item: any) {
    push(item.to);
  }
</script>

<div class="app-layout">
  <Topbar user={topbarUser} onLogout={handleLogout} />

  <Sidebar
    items={sidebarItems}
    onNavigate={handleNavigate}
    bind:compact={isSidebarCompact}
  />

  <main class="main-content" data-compact={isSidebarCompact}>
    <Router {routes} />
  </main>
</div>

<style>
  :root {
    --sidebar-width: 256px;
    --sidebar-compact-width: 64px;
    --topbar-height: 56px;
  }

  .app-layout {
    position: relative;
    min-height: 100vh;
    background-color: white;
  }

  .main-content {
    margin-left: var(--sidebar-width);
    transition: margin-left 0.2s ease;
    padding-top: calc(var(--topbar-height) + 1rem);
    padding-right: 2rem;
    padding-left: 2rem;
    padding-bottom: 2rem;
    min-height: 100vh;
    box-sizing: border-box;
    overflow-y: auto;
  }

  .main-content[data-compact="true"] {
    margin-left: var(--sidebar-compact-width);
  }

  @media (max-width: 768px) {
    .main-content {
      margin-left: 0;
    }
  }
</style>
