<script lang="ts">
  import Router, { location } from "svelte-spa-router";
  import { routes } from "./routes";
  import { isAuthenticated } from "./stores/auth-store";
  import { push } from "svelte-spa-router";
  import ToastContainer from "./components/ToastContainer.svelte";
  import Topbar from "./components/Topbar.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import { tokenManager } from "./auth/token";

  let isSidebarCompact = $state(false);
  let topbarUser = $state<{ name: string; avatar?: string } | undefined>(
    undefined,
  );
  let showLayout = $derived($isAuthenticated && $location !== "/login");

  // Redirect to login if not authenticated
  $effect(() => {
    if (!$isAuthenticated && $location !== "/login") {
      push("/login");
    }
  });

  // Load user info
  $effect(() => {
    if ($isAuthenticated) {
      topbarUser = { name: "Admin" };
    } else {
      topbarUser = undefined;
    }
  });

  const sidebarItems = [
    {
      id: "overview",
      label: "Tổng quan",
      to: "/dashboard",
      icon: `<svg viewBox="0 0 24 24" width="20" height="20" fill="none"><path d="M3 11.5L12 4l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1v-8.5z" fill="currentColor"/></svg>`,
    },
    {
      id: "users",
      label: "Người dùng",
      to: "/users",
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/><path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    },
    {
      id: "communities",
      label: "Cộng đồng",
      to: "/communities",
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none"><rect x="3" y="3" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><rect x="14" y="3" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><rect x="14" y="14" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><rect x="3" y="14" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    },
    {
      id: "reports",
      label: "Báo cáo",
      to: "/reports",
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="9" x2="12" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="17" x2="12.01" y2="17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    },
  ];

  function handleLogout() {
    console.log("[App] handleLogout called");
    console.trace("[App] handleLogout trace");
    tokenManager.clearTokens();
    isAuthenticated.set(false);
    push("/login");
  }

  function handleNavigate(item: any) {
    push(item.to);
  }
</script>

{#if showLayout}
  <div class="app-layout">
    <Topbar user={topbarUser} onLogout={handleLogout} />
    <Sidebar
      items={sidebarItems}
      bind:compact={isSidebarCompact}
      onNavigate={handleNavigate}
    />
    <main class="main-content" data-compact={isSidebarCompact}>
      <Router {routes} />
    </main>
  </div>
{:else}
  <Router {routes} />
{/if}

<ToastContainer />
