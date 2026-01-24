<script lang="ts">
  import { push } from "svelte-spa-router";
  import CreateCommunityModal from "./CreateCommunityModal.svelte";
  import { getCommunitiesByUserId } from "../services/community-service";
  import { authStore } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
  import type { CommunityResponse } from "../dtos/community-dto";
  import type { UserResponse } from "../dtos/user-dto";
  import { onMount } from "svelte";

  type Props = {
    compact?: boolean;
  };

  let { compact = false }: Props = $props();

  let userCommunities = $state<CommunityResponse[]>([]);
  let isLoadingCommunities = $state(false);
  let user = $state<UserResponse | null>(null);

  let isExpanded = $state(true);
  let showCreateModal = $state(false);

  authStore.subscribe((state) => {
    user = state.user;
    if (user) {
      loadUserCommunities();
    }
  });

  onMount(() => {
    if (user) {
      loadUserCommunities();
    }
  });

  async function loadUserCommunities() {
    if (!user?.id) return;

    try {
      isLoadingCommunities = true;

      // Get communities directly from new endpoint
      userCommunities = await getCommunitiesByUserId(user.id);
      console.log("User communities:", userCommunities);
    } catch (error) {
      console.error("Failed to load user communities:", error);
      userCommunities = [];
    } finally {
      isLoadingCommunities = false;
    }
  }
  function toggleExpand() {
    isExpanded = !isExpanded;
  }

  function navigateToCommunity(communityName: string) {
    push(`/lk/${communityName}`);
  }

  function handleCreateCommunity() {
    if (!user) {
      toastStore.warning("Vui lòng đăng nhập để tạo cộng đồng");
      return;
    }
    showCreateModal = true;
  }

  function handleCloseModal() {
    showCreateModal = false;
    // Reload communities after creating a new one
    if (user) {
      loadUserCommunities();
    }
  }

  function handleManageCommunities() {
    push("/communities/manage");
  }
</script>

<div class="communities-section" class:compact>
  {#if !compact}
    <button class="section-header" onclick={toggleExpand}>
      <span class="section-title">CỘNG ĐỒNG</span>
      <span class="expand-icon" class:expanded={isExpanded}>
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
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
  {:else}
    <!-- Compact mode: Show community avatars -->
    <div class="compact-communities">
      {#if isLoadingCommunities}
        <div class="compact-loading">⏳</div>
      {:else if userCommunities.length > 0}
        {#each userCommunities.slice(0, 4) as community (community.id)}
          <button
            class="compact-community-avatar"
            onclick={() => navigateToCommunity(community.name)}
            title="lk/{community.name}"
          >
            {#if community.avatar}
              <img src={community.avatar} alt={community.name} />
            {:else}
              <span class="avatar-placeholder">
                {community.name.charAt(0).toUpperCase()}
              </span>
            {/if}
          </button>
        {/each}
        {#if userCommunities.length > 4}
          <button
            class="compact-more"
            title="+{userCommunities.length - 4} more"
          >
            +{userCommunities.length - 4}
          </button>
        {/if}
      {:else}
        <button
          class="compact-add"
          onclick={handleCreateCommunity}
          title="Tạo cộng đồng"
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 4V16M4 10H16"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>
      {/if}
    </div>
  {/if}

  {#if isExpanded && !compact}
    <div class="communities-content">
      <!-- Create Community -->
      <button class="action-button" onclick={handleCreateCommunity}>
        <span class="action-icon">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 4V16M4 10H16"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </span>
        <span class="action-label">Tạo cộng đồng</span>
      </button>

      <!-- User's Communities -->
      {#if isLoadingCommunities}
        <div class="loading-message">
          <span class="loading-spinner">⏳</span>
          <span>Đang tải...</span>
        </div>
      {:else if userCommunities.length > 0}
        <div class="user-communities-list">
          {#each userCommunities as community (community.id)}
            <button
              class="community-item"
              onclick={() => navigateToCommunity(community.name)}
            >
              {#if community.avatar}
                <img
                  src={community.avatar}
                  alt={community.name}
                  class="community-avatar"
                />
              {:else}
                <span class="community-icon">📁</span>
              {/if}
              <span class="community-name">lk/{community.name}</span>
            </button>
          {/each}
        </div>
      {/if}

      <!-- Manage Communities -->
      <button class="action-button" onclick={handleManageCommunities}>
        <span class="action-icon">
          <img src="/setting_icon.svg" alt="Settings" width="20" height="20" />
        </span>
        <span class="action-label">Quản lý cộng đồng</span>
      </button>
    </div>
  {/if}
</div>

<CreateCommunityModal show={showCreateModal} onClose={handleCloseModal} />

<style>
  .communities-section {
    margin-top: 12px;
    border-top: 1px solid hsl(var(--sidebar-border));
    padding-top: 12px;
  }

  .communities-section.compact {
    padding-top: 8px;
  }

  .section-header {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 8px 12px;
    background: transparent;
    border: none;
    cursor: pointer;
    border-radius: 6px;
    transition: background 0.15s ease;
  }

  .section-header:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .section-title {
    font-size: 13px;
    font-weight: 600;
    color: #a8a8a8;
    letter-spacing: 0.8px;
  }

  .section-icon {
    font-size: 20px;
  }

  .expand-icon {
    width: 16px;
    height: 16px;
    color: #878a8c;
    transition: transform 0.2s ease;
  }

  .expand-icon.expanded {
    transform: rotate(0deg);
  }

  .expand-icon:not(.expanded) {
    transform: rotate(-90deg);
  }

  .communities-content {
    padding: 4px 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .action-button {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    background: transparent;
    border: none;
    cursor: pointer;
    border-radius: 6px;
    transition: background 0.15s ease;
    color: #1c1c1c;
    font-size: 14px;
    font-weight: 500;
  }

  .action-button:hover {
    background: rgba(0, 0, 0, 0.08);
  }

  .action-icon {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #1c1c1c;
  }

  .action-label {
    flex: 1;
    text-align: left;
  }

  .loading-message {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    font-size: 13px;
    color: #878a8c;
  }

  .loading-spinner {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .user-communities-list {
    margin: 4px 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding-bottom: 8px;
    border-bottom: 1px solid hsl(var(--sidebar-border));
  }

  .community-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    background: transparent;
    border: none;
    cursor: pointer;
    border-radius: 6px;
    transition: background 0.15s ease;
  }

  .community-item:hover {
    background: rgba(0, 0, 0, 0.08);
  }

  .community-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
  }

  .community-icon {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    flex-shrink: 0;
    background: #f0f0f0;
  }

  .community-name {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
    color: #1c1c1c;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    text-align: left;
  }

  /* Compact mode styles */
  .compact-communities {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 8px 4px;
  }

  .compact-community-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    padding: 0;
    cursor: pointer;
    overflow: hidden;
    transition: all 0.2s ease;
    position: relative;
  }

  .compact-community-avatar:hover {
    transform: scale(1.1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  }

  .compact-community-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    font-weight: 600;
    font-size: 16px;
  }

  .compact-loading {
    font-size: 20px;
    animation: spin 1s linear infinite;
  }

  .compact-more {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: 2px dashed #d0d0d0;
    background: #f5f5f5;
    color: #666;
    font-size: 11px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .compact-more:hover {
    background: #e8e8e8;
    border-color: #b0b0b0;
    transform: scale(1.05);
  }

  .compact-add {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: 2px dashed var(--blue--);
    background: rgba(21, 48, 96, 0.05);
    color: var(--blue--);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .compact-add:hover {
    background: rgba(21, 48, 96, 0.1);
    transform: scale(1.05);
  }
</style>
