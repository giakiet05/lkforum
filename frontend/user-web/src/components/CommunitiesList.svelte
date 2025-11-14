<script lang="ts">
  import { push } from "svelte-spa-router";
  import CreateCommunityModal from "./CreateCommunityModal.svelte";
  import { getCommunities } from "../services/community-service";
  import { authStore } from "../stores/auth-store";
  import type { CommunityResponse } from "../dtos/community-dto";
  import type { UserResponse } from "../dtos/user-dto";
  import { onMount } from "svelte";

  type Props = {
    compact?: boolean;
  };

  let { compact = false }: Props = $props();

  const STORAGE_KEY = "user_communities";

  let userCommunities = $state<CommunityResponse[]>([]);
  let isLoadingCommunities = $state(false);
  let user = $state<UserResponse | null>(null);

  let isExpanded = $state(true);
  let showCreateModal = $state(false);

  authStore.subscribe((state) => {
    user = state.user;
    if (user) {
      loadCommunitiesFromStorage();
    }
  });

  onMount(() => {
    if (user) {
      loadCommunitiesFromStorage();
    }
  });

  function loadCommunitiesFromStorage() {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        userCommunities = JSON.parse(stored);
      }
    } catch (error) {
      console.error("Failed to load communities from storage:", error);
      userCommunities = [];
    }
  }

  function saveCommunitiesToStorage(communities: CommunityResponse[]) {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(communities));
    } catch (error) {
      console.error("Failed to save communities to storage:", error);
    }
  }

  function addCommunity(community: CommunityResponse) {
    userCommunities = [community, ...userCommunities];
    saveCommunitiesToStorage(userCommunities);
  }

  function toggleExpand() {
    isExpanded = !isExpanded;
  }

  function navigateToCommunity(communityName: string) {
    push(`/lk/${communityName}`);
  }

  function handleCreateCommunity() {
    showCreateModal = true;
  }

  function handleCloseModal() {
    showCreateModal = false;
  }

  function handleCommunityCreated(community: CommunityResponse) {
    addCommunity(community);
  }

  function handleManageCommunities() {
    push("/communities/manage");
  }
</script>

<div class="communities-section" class:compact>
  <button class="section-header" onclick={toggleExpand}>
    {#if !compact}
      <span class="section-title">COMMUNITIES</span>
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
    {:else}
      <span class="section-icon">👥</span>
    {/if}
  </button>

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
        <span class="action-label">Create Community</span>
      </button>

      <!-- User's Communities -->
      {#if isLoadingCommunities}
        <div class="loading-message">
          <span class="loading-spinner">⏳</span>
          <span>Loading...</span>
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
        <span class="action-label">Manage Communities</span>
      </button>
    </div>
  {/if}
</div>

<CreateCommunityModal
  show={showCreateModal}
  onClose={handleCloseModal}
  onCommunityCreated={handleCommunityCreated}
/>

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
</style>
