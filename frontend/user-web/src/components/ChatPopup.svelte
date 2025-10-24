<script lang="ts">
  import { push } from "svelte-spa-router";
  import { mockConversations } from "../mocks/conversations.mock";
  import type { Conversation, Message } from "../mocks/conversations.mock";

  type ChatPopupProps = {
    show: boolean;
    onClose: () => void;
  };

  let { show, onClose }: ChatPopupProps = $props();

  let conversations = $state(mockConversations);
  let selectedConversation = $state<Conversation | null>(conversations[0]);
  let messageInput = $state("");
  let showChatMenu = $state(false);

  function handleSelectConversation(conversation: Conversation) {
    selectedConversation = conversation;
    if (selectedConversation) {
      selectedConversation.unreadCount = 0;
    }
  }

  function handleSendMessage() {
    if (!messageInput.trim() || !selectedConversation) return;

    const newMessage: Message = {
      id: `m${Date.now()}`,
      senderId: "me",
      senderName: "You",
      senderAvatar: "/avatar.jpg",
      content: messageInput,
      timestamp: new Date().toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
      }),
      isRead: false,
      isSent: true,
    };

    selectedConversation.messages.push(newMessage);
    selectedConversation.lastMessage = `You: ${messageInput}`;
    selectedConversation.lastMessageTime = "Just now";

    messageInput = "";

    setTimeout(() => {
      const messagesArea = document.querySelector(".popup-messages-area");
      if (messagesArea) {
        messagesArea.scrollTop = messagesArea.scrollHeight;
      }
    }, 0);
  }

  function toggleMute() {
    if (selectedConversation) {
      selectedConversation.isMuted = !selectedConversation.isMuted;
    }
  }

  function toggleChatMenu() {
    showChatMenu = !showChatMenu;
  }

  function handleExpand() {
    push("/messages");
    onClose();
  }
</script>

{#if show}
  <div class="chat-popup">
    <!-- Popup Header -->
    <div class="popup-header">
      <h2>Tin nhắn</h2>
      <div class="popup-header-actions">
        <button class="header-icon-btn" onclick={handleExpand} title="Expand">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M13 3h4v4M7 17H3v-4M17 3l-7 7M3 17l7-7"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>
        <button class="header-icon-btn" onclick={onClose} title="Close">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M15 5L5 15M5 5l10 10"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>
      </div>
    </div>

    <div class="popup-container">
      <!-- Left Side - Conversations List -->
      <div class="popup-conversations">
        <div class="popup-conversations-list">
          {#each conversations as conversation}
            <button
              class="popup-conversation-item"
              class:active={selectedConversation?.id === conversation.id}
              class:unread={conversation.unreadCount > 0}
              onclick={() => handleSelectConversation(conversation)}
            >
              <div class="popup-avatar-wrapper">
                <img
                  src={conversation.userAvatar}
                  alt={conversation.userName}
                  class="popup-avatar"
                />
                {#if conversation.isOnline}
                  <span class="popup-online-indicator"></span>
                {/if}
              </div>
              <div class="popup-conversation-info">
                <div class="popup-conversation-top">
                  <span class="popup-conversation-name"
                    >{conversation.userName}</span
                  >
                  <span class="popup-conversation-time"
                    >{conversation.lastMessageTime}</span
                  >
                </div>
                <div class="popup-conversation-bottom">
                  <p class="popup-conversation-preview">
                    {conversation.lastMessage}
                  </p>
                  {#if conversation.unreadCount > 0}
                    <span class="popup-unread-badge"
                      >{conversation.unreadCount}</span
                    >
                  {/if}
                </div>
              </div>
            </button>
          {/each}
        </div>
      </div>

      <!-- Right Side - Chat Detail -->
      <div class="popup-chat-detail">
        {#if selectedConversation}
          <!-- Chat Header -->
          <div class="popup-chat-header">
            <div class="popup-chat-user-info">
              <img
                src={selectedConversation.userAvatar}
                alt={selectedConversation.userName}
                class="popup-chat-avatar"
              />
              <div class="popup-chat-user-details">
                <h3>{selectedConversation.userName}</h3>
                {#if selectedConversation.isOnline}
                  <span class="popup-status-online">Active</span>
                {/if}
              </div>
            </div>
            <div class="popup-chat-actions">
              <div class="popup-chat-menu-wrapper">
                <button
                  class="popup-action-icon-btn"
                  onclick={toggleChatMenu}
                  title="More options"
                >
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"
                    />
                  </svg>
                </button>
                {#if showChatMenu}
                  <div class="popup-chat-dropdown">
                    <button
                      class="popup-dropdown-item"
                      class:muted={selectedConversation.isMuted}
                      onclick={() => {
                        toggleMute();
                        showChatMenu = false;
                      }}
                    >
                      <img
                        src="/muted_icon.svg"
                        alt="Mute"
                        width="18"
                        height="18"
                      />
                      <span
                        >{selectedConversation.isMuted
                          ? "Unmute"
                          : "Mute"}</span
                      >
                    </button>
                    <button
                      class="popup-dropdown-item"
                      onclick={() => {
                        showChatMenu = false;
                      }}
                    >
                      <img
                        src="/notification_icon.svg"
                        alt="Notification"
                        width="18"
                        height="18"
                      />
                      <span>Notifications</span>
                    </button>
                  </div>
                {/if}
              </div>
            </div>
          </div>

          <!-- Messages Area -->
          <div class="popup-messages-area">
            <div class="popup-messages-wrapper">
              {#each selectedConversation.messages as message}
                <div class="popup-message-row" class:sent={message.isSent}>
                  {#if !message.isSent}
                    <img
                      src={message.senderAvatar}
                      alt={message.senderName}
                      class="popup-message-avatar"
                    />
                  {/if}
                  <div class="popup-message-bubble" class:sent={message.isSent}>
                    <p>{message.content}</p>
                    <span class="popup-message-time">{message.timestamp}</span>
                  </div>
                </div>
              {/each}
            </div>
          </div>

          <!-- Message Input -->
          <div class="popup-message-input">
            <button class="popup-emoji-btn" title="Emoji">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <circle
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="2"
                />
                <path
                  d="M8 14s1.5 2 4 2 4-2 4-2"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <circle cx="9" cy="9" r="1" fill="currentColor" />
                <circle cx="15" cy="9" r="1" fill="currentColor" />
              </svg>
            </button>
            <input
              type="text"
              placeholder="Nhắn tin..."
              bind:value={messageInput}
              onkeydown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  handleSendMessage();
                }
              }}
            />
            <button
              class="popup-send-btn"
              onclick={handleSendMessage}
              disabled={!messageInput.trim()}
              title="Send"
            >
              <img src="/send_icon.svg" alt="Send" width="20" height="20" />
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .chat-popup {
    position: fixed;
    bottom: 0;
    right: 80px;
    width: 680px;
    height: 520px;
    background: white;
    border-radius: 12px 12px 0 0;
    display: flex;
    flex-direction: column;
    box-shadow: 0 -2px 16px rgba(0, 0, 0, 0.15);
    z-index: 9999;
    overflow: hidden;
  }

  /* Popup Header */
  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #edeff1;
    background: white;
  }

  .popup-header h2 {
    margin: 0;
    font-size: 16px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .popup-header-actions {
    display: flex;
    gap: 8px;
  }

  .header-icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .header-icon-btn:hover {
    background: #f6f7f8;
  }

  /* Container */
  .popup-container {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  /* Conversations List */
  .popup-conversations {
    width: 240px;
    border-right: 1px solid #edeff1;
    display: flex;
    flex-direction: column;
    background: white;
  }

  .popup-conversations-list {
    flex: 1;
    overflow-y: auto;
  }

  .popup-conversation-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 16px;
    width: 100%;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    transition: background 0.2s;
    border-left: 3px solid transparent;
  }

  .popup-conversation-item:hover {
    background: #f6f7f8;
  }

  .popup-conversation-item.active {
    background: #e8f4fd;
    border-left-color: var(--blue--);
  }

  .popup-avatar-wrapper {
    position: relative;
    flex-shrink: 0;
  }

  .popup-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }

  .popup-online-indicator {
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 12px;
    height: 12px;
    background: #31a24c;
    border: 2px solid white;
    border-radius: 50%;
  }

  .popup-conversation-info {
    flex: 1;
    min-width: 0;
  }

  .popup-conversation-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }

  .popup-conversation-name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .popup-conversation-item.unread .popup-conversation-name {
    font-weight: 700;
  }

  .popup-conversation-item.unread .popup-conversation-preview {
    color: #1c1c1c;
    font-weight: 600;
  }

  .popup-conversation-time {
    font-size: 11px;
    color: #7c7c7c;
  }

  .popup-conversation-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
  }

  .popup-conversation-preview {
    flex: 1;
    margin: 0;
    font-size: 12px;
    color: #7c7c7c;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .popup-unread-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    background: var(--error--);
    color: white;
    font-size: 10px;
    font-weight: 700;
    border-radius: 9px;
  }

  /* Chat Detail */
  .popup-chat-detail {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: white;
  }

  .popup-chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #edeff1;
  }

  .popup-chat-user-info {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .popup-chat-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .popup-chat-user-details h3 {
    margin: 0 0 2px 0;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .popup-status-online {
    font-size: 11px;
    color: #31a24c;
  }

  .popup-chat-actions {
    position: relative;
  }

  .popup-chat-menu-wrapper {
    position: relative;
  }

  .popup-action-icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .popup-action-icon-btn:hover {
    background: #f6f7f8;
  }

  .popup-chat-dropdown {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    min-width: 180px;
    z-index: 100;
    overflow: hidden;
  }

  .popup-dropdown-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 10px 14px;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    font-size: 13px;
    color: #1c1c1c;
    transition: background 0.2s;
  }

  .popup-dropdown-item:hover {
    background: #f6f7f8;
  }

  .popup-dropdown-item.muted {
    color: #ff4444;
  }

  .popup-dropdown-item.muted:hover {
    background: #ffe0e0;
  }

  .popup-dropdown-item img {
    opacity: 0.7;
  }

  .popup-dropdown-item span {
    font-weight: 500;
  }

  /* Messages Area */
  .popup-messages-area {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    background: #f6f7f8;
  }

  .popup-messages-wrapper {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .popup-message-row {
    display: flex;
    align-items: flex-end;
    gap: 6px;
  }

  .popup-message-row.sent {
    flex-direction: row-reverse;
  }

  .popup-message-avatar {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
  }

  .popup-message-bubble {
    max-width: 70%;
    padding: 8px 12px;
    background: white;
    border-radius: 16px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .popup-message-bubble.sent {
    background: var(--blue--);
    color: white;
  }

  .popup-message-bubble p {
    margin: 0 0 3px 0;
    font-size: 13px;
    line-height: 1.4;
    word-wrap: break-word;
  }

  .popup-message-time {
    font-size: 10px;
    opacity: 0.7;
  }

  /* Message Input */
  .popup-message-input {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    border-top: 1px solid #edeff1;
    background: white;
  }

  .popup-emoji-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .popup-emoji-btn:hover {
    background: #f6f7f8;
  }

  .popup-message-input input {
    flex: 1;
    padding: 8px 14px;
    border: 1px solid #edeff1;
    border-radius: 18px;
    background: #f6f7f8;
    font-size: 13px;
    color: #1c1c1c;
  }

  .popup-message-input input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .popup-send-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    transition: background 0.2s;
  }

  .popup-send-btn:hover:not(:disabled) {
    background: #f6f7f8;
  }

  .popup-send-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
</style>
