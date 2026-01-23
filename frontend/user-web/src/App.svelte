<script lang="ts">
  import Router, { location } from "svelte-spa-router";
  import routes from "./routes";
  import Topbar from "./components/Topbar.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import ToastContainer from "./components/ToastContainer.svelte";
  import { authStore, getInitialAuthState } from "./stores/auth-store";
  import { logout, validateAuth } from "./services/auth-service";
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import { websocketService } from "./services/websocket-service";
  import {
    addOnlineUser,
    removeOnlineUser,
    setCurrentUserId,
    addMessage,
    setChannels,
    setMessages,
  } from "./stores/chat-store";
  import { getChannelsByUser } from "./services/channel-service";
  import { getMessages } from "./services/message-service";
  import type { MessageResponse } from "./dtos/message-dto";

  const sidebarItems = [
    {
      id: "home",
      label: "Trang chủ",
      to: "/",
      icon: `<svg viewBox=\"0 0 24 24\" width=\"20\" height=\"20\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">
      <path d=\"M3 11.5L12 4l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1v-8.5z\" fill=\"currentColor\"/>
    </svg>`,
    },
    {
      id: "popular",
      label: "Phổ biến",
      to: "/popular",
      icon: `<svg width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><g clip-path=\"url(#clip0_20_157)\"><path d=\"M23.2499 12.751L12.7769 23.25\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M17.25 12.751H23.25V18.75\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M18.75 0.75V5.25H12.75V11.25H6.75V17.25H0.75V23.25\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/></g><defs><clipPath id=\"clip0_20_157\"><rect width=\"24\" height=\"24\" fill=\"white\"/></clipPath></defs></svg>`,
    },
    {
      id: "explore",
      label: "Khám phá",
      to: "/explore",
      icon: `<svg width=\"20\" height=\"20\" viewBox=\"0 0 16 16\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">
      <path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M11.8957 3.33209L10.8016 4.06159L11.9615 6.9595L13.2837 6.57143L11.8957 3.33209ZM9.66815 4.81615L3.05234 9.2274L3.20634 9.53543L10.6762 7.3375L9.66815 4.81615ZM12.49 1.3335L15.1016 7.42709L10.1902 8.8715L12.1669 14.1267L10.9185 14.5949L8.9075 9.24884L8.0035 9.51415L6.13906 14.5954L4.91353 14.0702L6.4135 9.98215L2.51275 11.1297L1.3335 8.77118L12.49 1.3335Z\" fill=\"currentColor\"/>
    </svg>`,
    },
    {
      id: "all",
      label: "Tất cả",
      to: "/all",
      icon: `<svg width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><g clip-path=\"url(#clip0_20_165)\"><path d=\"M23.25 21.25C23.25 22.35 22.35 23.25 21.25 23.25H2.75C1.65 23.25 0.75 22.35 0.75 21.25V2.75C0.75 1.65 1.65 0.75 2.75 0.75H21.25C22.35 0.75 23.25 1.65 23.25 2.75V21.25Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M7.44971 18.4499C7.44971 19.0499 7.04971 19.4499 6.44971 19.4499H5.44971C4.84971 19.4499 4.44971 19.0499 4.44971 18.4499V15.6499C4.44971 15.0499 4.84971 14.6499 5.44971 14.6499H6.44971C7.04971 14.6499 7.44971 15.0499 7.44971 15.6499V18.4499Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M13.4502 18.4499C13.4502 19.0499 13.0502 19.4499 12.4502 19.4499H11.4502C10.8502 19.4499 10.4502 19.0499 10.4502 18.4499V6.6499C10.4502 6.0499 10.8502 5.6499 11.4502 5.6499H12.4502C13.0502 5.6499 13.4502 6.0499 13.4502 6.6499V18.4499Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/><path d=\"M19.4502 18.4499C19.4502 19.0499 19.0502 19.4499 18.4502 19.4499H17.4502C16.8502 19.4499 16.4502 19.0499 16.4502 18.4499V11.6499C16.4502 11.0499 16.8502 10.6499 17.4502 10.6499H18.4502C19.0502 10.6499 19.4502 11.0499 19.4502 11.6499V18.4499Z\" stroke=\"currentColor\" stroke-opacity=\"0.7\" stroke-width=\"1.5\" stroke-miterlimit=\"10\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/></g><defs><clipPath id=\"clip0_20_165\"><rect width=\"24\" height=\"24\" fill=\"white\"/></clipPath></defs></svg>`,
    },
  ];

  let isSidebarCompact = $state(false);
  let isSidebarOpen = $state(false);
  let shouldShowAuthModal = $state(false);
  let isSearchPage = $derived($location === "/search");

  let topbarUser = $state<
    { name: string; avatar?: string; karma?: number } | undefined
  >(undefined);

  let isAuthChecking = $state(true);

  // Validate auth khi app khởi động
  onMount(async () => {
    console.log("App mounted - starting auth check");
    try {
      await validateAuth();
      console.log("Auth validation completed");
    } catch (error) {
      console.error("Auth validation error:", error);
    } finally {
      isAuthChecking = false;
      console.log("isAuthChecking set to false");
    }
  });

  // Listen event từ token.ts khi refresh token thành công
  onMount(() => {
    const handleAuthRefreshed = () => {
      authStore.set(getInitialAuthState());
    };

    const handleAuthUnauthorized = () => {
      logout();
      shouldShowAuthModal = true;
    };

    window.addEventListener("auth:refreshed", handleAuthRefreshed);
    window.addEventListener("auth:unauthorized", handleAuthUnauthorized);

    return () => {
      window.removeEventListener("auth:refreshed", handleAuthRefreshed);
      window.removeEventListener("auth:unauthorized", handleAuthUnauthorized);
    };
  });

  // Subscribe to authStore để update realtime khi login/logout
  authStore.subscribe((state) => {
    if (state.user) {
      console.log("User data in authStore:", state.user);
      topbarUser = {
        name: state.user.username || state.user.email || "User",
        avatar: state.user.profile?.avatar?.url,
        karma: state.user.reputation || 0,
      };
      console.log("Topbar user:", topbarUser);
    } else {
      topbarUser = undefined;
    }
  });

  // Define presence handlers outside so they can be reused
  const handlePresenceOnline = (payload: { user_id: string }) => {
    console.log("🟢 User online:", payload.user_id);
    addOnlineUser(payload.user_id);
  };

  const handlePresenceOffline = (payload: { user_id: string }) => {
    console.log("⚫ User offline:", payload.user_id);
    removeOnlineUser(payload.user_id);
  };

  // Global handler for incoming messages (so badge shows even when ChatPopup is closed)
  const handleGlobalMessage = (payload: any) => {
    const message: MessageResponse = payload.message || payload;
    console.log("📨 [App] Global message received:", message);
    addMessage(message.channel_id, message);
  };

  // Load channels and messages on login to get unread count
  async function loadUserChannelsAndMessages(userId: string) {
    try {
      console.log("📥 [App] Loading channels for unread count...");
      const response = await getChannelsByUser(userId);
      const channels = response.channels || [];
      setChannels(channels);

      // Load messages for each channel to calculate unread count
      for (const channel of channels) {
        try {
          const messages = await getMessages({ channel_id: channel.id });
          setMessages(channel.id, messages);
        } catch (err) {
          console.error(
            "Failed to load messages for channel:",
            channel.id,
            err,
          );
        }
      }
      console.log("✅ [App] Channels and messages loaded for unread count");
    } catch (err) {
      console.error("❌ [App] Failed to load channels:", err);
    }
  }

  // Connect WebSocket when user is authenticated
  $effect(() => {
    const user = $authStore.user;
    if (user && !isAuthChecking) {
      // Set current user ID for chat store (needed to filter own messages from unread count)
      setCurrentUserId(user.id);

      // Load channels to get initial unread count
      loadUserChannelsAndMessages(user.id);

      // Always register presence handlers first
      websocketService.on("presence_online", handlePresenceOnline);
      websocketService.on("presence_offline", handlePresenceOffline);
      // Register global message handler to update badge even when ChatPopup is closed
      websocketService.on("send_message", handleGlobalMessage);
      websocketService.on("ack_message", handleGlobalMessage);
      console.log("✅ Presence and message handlers registered");

      // Connect WebSocket if not already connected
      if (!websocketService.isConnected()) {
        websocketService
          .connect()
          .then(() => {
            console.log("✅ WebSocket connected for user:", user.username);
          })
          .catch((error) => {
            console.error("❌ Failed to connect WebSocket:", error);
          });
      } else {
        console.log("✅ WebSocket already connected for user:", user.username);
      }
    } else {
      // Clear current user ID when logged out
      setCurrentUserId(null);
    }

    // Cleanup function
    return () => {
      websocketService.off("presence_online", handlePresenceOnline);
      websocketService.off("presence_offline", handlePresenceOffline);
      websocketService.off("send_message", handleGlobalMessage);
      websocketService.off("ack_message", handleGlobalMessage);
    };
  });

  function handleLogout() {
    logout();
    // Redirect to home and refresh the page
    window.location.href = "/#/";
    window.location.reload();
  }

  function handleNavigate(item: any) {
    push(item.to);
  }
</script>

{#if isAuthChecking}
  <div class="loading-screen">
    <div class="loading-spinner"></div>
    <p>Đang kiểm tra đăng nhập...</p>
  </div>
{:else}
  <div class="app-layout" class:search-page={isSearchPage}>
    {#if !isSearchPage}
      <Topbar
        user={topbarUser}
        onLogout={handleLogout}
        onMenuClick={() => (isSidebarOpen = !isSidebarOpen)}
        forceShowAuthModal={shouldShowAuthModal}
        onAuthModalClose={() => (shouldShowAuthModal = false)}
      />

      <Sidebar
        items={sidebarItems}
        onNavigate={handleNavigate}
        bind:compact={isSidebarCompact}
        isOpen={isSidebarOpen}
        onClose={() => (isSidebarOpen = false)}
      />
    {/if}

    <main
      class="main-content"
      data-compact={isSidebarCompact}
      class:full-width={isSearchPage}
    >
      <Router {routes} />
    </main>
  </div>
  <ToastContainer />
{/if}

<style>
  :root {
    --sidebar-width: 256px;
    --sidebar-compact-width: 64px;
    --topbar-height: 56px;
  }

  .loading-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background-color: white;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #ff4500;
    border-radius: 50%;
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

  .loading-screen p {
    margin-top: 16px;
    color: #666;
    font-size: 14px;
  }

  .app-layout {
    position: relative;
    min-height: 100vh;
    background-color: white;
  }

  .app-layout.search-page {
    padding-top: 0;
  }

  .main-content {
    margin-left: var(--sidebar-width);
    transition: margin-left 0.2s ease;
    padding-top: var(--topbar-height);
    padding-right: 0;
    padding-left: 0;
    padding-bottom: 0;
    min-height: 100vh;
    box-sizing: border-box;
    overflow-y: auto;
  }

  .main-content.full-width {
    margin-left: 0;
    padding-top: 0;
  }

  .main-content[data-compact="true"] {
    margin-left: var(--sidebar-compact-width);
  }

  @media (max-width: 1024px) {
    .main-content {
      margin-left: 0;
      padding-top: var(--topbar-height);
    }

    .main-content[data-compact="true"] {
      margin-left: 0;
    }
  }
</style>
