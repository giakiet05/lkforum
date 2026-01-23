<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import {
    getNotifications,
    markNotificationAsRead,
    deleteNotification,
    markAllAsRead,
  } from "../services/notification-service";
  import { activateModerator } from "../services/community-service";
  import { toastStore } from "../stores/toast-store";
  import type { NotificationResponse } from "../dtos/notification-dto";

  type Props = {
    show?: boolean;
    onClose?: () => void;
    onUnreadCountChange?: (count: number) => void;
  };

  let { show = false, onClose, onUnreadCountChange }: Props = $props();

  let notifications = $state<NotificationResponse[]>([]);
  let isLoading = $state(false);
  let error = $state<string | null>(null);

  // Load notifications when dropdown opens
  $effect(() => {
    if (show) {
      loadNotifications();
    }
  });

  async function loadNotifications() {
    try {
      isLoading = true;
      error = null;
      const response = await getNotifications({ page: 1, pageSize: 20 });
      notifications = response.notifications || [];
      // Notify parent of unread count
      const unreadCount = notifications.filter((n) => !n.is_read).length;
      onUnreadCountChange?.(unreadCount);
    } catch (err) {
      console.error("Failed to load notifications:", err);
      error = "Không thể tải thông báo";
    } finally {
      isLoading = false;
    }
  }

  async function handleMarkAsRead(notificationId: string) {
    try {
      await markNotificationAsRead(notificationId);
      // Update local state
      notifications = notifications.map((n) =>
        n.id === notificationId ? { ...n, is_read: true } : n,
      );
      // Update parent unread count
      const unreadCount = notifications.filter((n) => !n.is_read).length;
      onUnreadCountChange?.(unreadCount);
    } catch (err) {
      console.error("Failed to mark as read:", err);
    }
  }

  async function handleDelete(notificationId: string) {
    try {
      await deleteNotification(notificationId);
      // Remove from local state
      notifications = notifications.filter((n) => n.id !== notificationId);
      // Update parent unread count
      const unreadCount = notifications.filter((n) => !n.is_read).length;
      onUnreadCountChange?.(unreadCount);
    } catch (err) {
      console.error("Failed to delete notification:", err);
    }
  }

  async function handleMarkAllAsRead() {
    try {
      console.log("🔔 [NotificationsDropdown] Calling markAllAsRead API...");
      await markAllAsRead();
      console.log("✅ [NotificationsDropdown] markAllAsRead API success");
      // Update all notifications as read
      notifications = notifications.map((n) => ({ ...n, is_read: true }));
      // Update parent unread count to 0
      onUnreadCountChange?.(0);
    } catch (err) {
      console.error(
        "❌ [NotificationsDropdown] Failed to mark all as read:",
        err,
      );
    }
  }

  async function handleAcceptModeratorInvite(
    notificationId: string,
    communityId: string,
    notificationLink?: string,
  ) {
    try {
      await activateModerator(communityId);
      await handleDelete(notificationId);

      // Navigate to community page and reload to see updated moderator status
      if (notificationLink) {
        onClose?.();
        push(notificationLink);
        // Force page reload to refresh community data
        setTimeout(() => {
          window.location.reload();
        }, 100);
      } else {
        toastStore.success("Bạn đã chấp nhận làm moderator!");
        loadNotifications();
      }
    } catch (err) {
      console.error("Failed to accept moderator invite:", err);
      toastStore.error("Không thể chấp nhận lời mời");
    }
  }

  async function handleDeclineModeratorInvite(notificationId: string) {
    try {
      await handleDelete(notificationId);
      toastStore.info("Bạn đã từ chối lời mời");
    } catch (err) {
      console.error("Failed to decline moderator invite:", err);
    }
  }

  async function handleNotificationClick(notification: NotificationResponse) {
    console.log("🔔 Notification clicked:", notification);
    console.log("🔗 Link:", notification.link);

    // Mark as read when clicked
    if (!notification.is_read) {
      await handleMarkAsRead(notification.id);
    }

    // Navigate to link if exists
    if (notification.link) {
      let frontendLink = notification.link;

      // Handle message notification - link format: /channels/:channelId
      if (
        notification.type === "new_message" &&
        notification.link.startsWith("/channels/")
      ) {
        const channelId = notification.link.replace("/channels/", "");
        frontendLink = `/messages?channel=${channelId}`;
        console.log("💬 Message notification - navigating to:", frontendLink);
        onClose?.();
        push(frontendLink);
        return;
      }

      // Convert backend link format to frontend routing
      // Backend: /posts/:id or /posts/:id#comment-:commentId
      // Frontend: /post/:id or /post/:id#comment-:commentId (no 's')
      frontendLink = notification.link.replace(/^\/posts\//, "/post/");

      console.log("📍 Navigating to:", frontendLink);

      // Extract post ID from link (ignore comment anchor)
      const currentPostId = window.location.hash.match(/\/post\/([^#?]+)/)?.[1];
      const targetPostId = frontendLink.match(/\/post\/([^#?]+)/)?.[1];

      console.log(
        "🔍 Current post:",
        currentPostId,
        "Target post:",
        targetPostId,
      );

      // Force reload if same post ID
      if (currentPostId && targetPostId && currentPostId === targetPostId) {
        console.log("♻️ Same post - forcing reload");
        window.location.href = window.location.origin + "/#" + frontendLink;
        window.location.reload();
      } else {
        // Different post - navigate normally
        console.log("➡️ Different post - navigating");
        onClose?.();
        push(frontendLink);
      }
    }
  }

  function extractCommunityIdFromMetadata(
    notification: NotificationResponse,
  ): string | null {
    // Assuming metadata contains community_id
    try {
      const metadata = (notification as any).metadata;
      return metadata?.community_id || null;
    } catch {
      return null;
    }
  }

  const unreadCount = $derived(notifications.filter((n) => !n.is_read).length);
</script>

{#if show}
  <div class="notifications-dropdown">
    <div class="dropdown-header">
      <h3>Thông báo</h3>
      {#if unreadCount > 0}
        <button class="mark-all-btn" onclick={handleMarkAllAsRead}>
          Đánh dấu tất cả đã đọc
        </button>
      {/if}
    </div>

    <div class="notifications-list">
      {#if isLoading}
        <div class="loading">
          <div class="spinner"></div>
          <p>Đang tải...</p>
        </div>
      {:else if error}
        <div class="error">{error}</div>
      {:else if notifications.length === 0}
        <div class="empty">
          <p>Không có thông báo nào</p>
        </div>
      {:else}
        {#each notifications as notification (notification.id)}
          <div class="notification-item" class:unread={!notification.is_read}>
            <!-- Special handling for moderator invitation -->
            {#if notification.type === "system" && notification.message.includes("moderator")}
              {@const communityId =
                extractCommunityIdFromMetadata(notification)}
              <div
                class="notification-content"
                onclick={() => handleNotificationClick(notification)}
              >
                <div class="notification-message">{notification.message}</div>
                <div class="notification-time">
                  {new Date(notification.created_at).toLocaleString("vi-VN")}
                </div>
              </div>
              {#if communityId}
                <div class="notification-actions">
                  <button
                    class="accept-btn"
                    onclick={() =>
                      handleAcceptModeratorInvite(
                        notification.id,
                        communityId,
                        notification.link,
                      )}
                  >
                    Chấp nhận
                  </button>
                  <button
                    class="decline-btn"
                    onclick={() =>
                      handleDeclineModeratorInvite(notification.id)}
                  >
                    Từ chối
                  </button>
                </div>
              {/if}
            {:else}
              <div
                class="notification-content"
                onclick={() => handleNotificationClick(notification)}
              >
                <div class="notification-message">{notification.message}</div>
                <div class="notification-time">
                  {new Date(notification.created_at).toLocaleString("vi-VN")}
                </div>
              </div>
            {/if}

            <button
              class="delete-btn"
              onclick={() => handleDelete(notification.id)}
              title="Xóa thông báo"
            >
              ×
            </button>
          </div>
        {/each}
      {/if}
    </div>
  </div>
{/if}

<style>
  .notifications-dropdown {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 8px;
    width: 400px;
    max-height: 500px;
    background: white;
    border-radius: 8px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    z-index: 1000;
  }

  .dropdown-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    border-bottom: 1px solid #e0e0e0;
  }

  .dropdown-header h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
  }

  .mark-all-btn {
    padding: 6px 12px;
    font-size: 12px;
    color: #0079d3;
    background: none;
    border: 1px solid #0079d3;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .mark-all-btn:hover {
    background: #0079d3;
    color: white;
  }

  .notifications-list {
    max-height: 450px;
    overflow-y: auto;
  }

  .loading,
  .error,
  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
    color: #666;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #0079d3;
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

  .notification-item {
    position: relative;
    padding: 16px;
    border-bottom: 1px solid #f0f0f0;
    transition: background 0.2s;
  }

  .notification-item:hover {
    background: #f8f8f8;
  }

  .notification-item.unread {
    background: #f0f7ff;
  }

  .notification-item.unread::before {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: #0079d3;
  }

  .notification-content {
    cursor: pointer;
    margin-right: 30px;
  }

  .notification-message {
    font-size: 14px;
    line-height: 1.5;
    color: #1c1c1c;
    margin-bottom: 4px;
  }

  .notification-time {
    font-size: 12px;
    color: #787c7e;
  }

  .notification-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
  }

  .accept-btn,
  .decline-btn {
    flex: 1;
    padding: 8px 16px;
    font-size: 13px;
    font-weight: 600;
    border-radius: 20px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .accept-btn {
    background: #0079d3;
    color: white;
    border: none;
  }

  .accept-btn:hover {
    background: #005fa3;
  }

  .decline-btn {
    background: white;
    color: #1c1c1c;
    border: 1px solid #e0e0e0;
  }

  .decline-btn:hover {
    background: #f0f0f0;
  }

  .delete-btn {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: none;
    font-size: 20px;
    color: #999;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .delete-btn:hover {
    background: #f0f0f0;
    color: #1c1c1c;
  }
</style>
