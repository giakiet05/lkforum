<script lang="ts">
  import { push } from "svelte-spa-router";
  import CreateCommunityModal from "./CreateCommunityModal.svelte";
  import type { SidebarCommunity } from "../mocks/sidebar-communities.mock";
  import { mockSidebarCommunities } from "../mocks/sidebar-communities.mock";

  type Props = {
    compact?: boolean;
  };

  let { compact = false }: Props = $props();

  let communities = $state<SidebarCommunity[]>(mockSidebarCommunities);

  let isExpanded = $state(true);
  let showCreateModal = $state(false);

  function toggleExpand() {
    isExpanded = !isExpanded;
  }

  function toggleFavorite(communityId: string) {
    const community = communities.find((c) => c.id === communityId);
    if (community) {
      community.isFavorite = !community.isFavorite;
    }
  }

  function navigateToCommunity(communityName: string) {
    push(`/lk/${communityName}`);
  }

  function handleCreateCommunity() {
    showCreateModal = true;
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

      <!-- Manage Communities -->
      <button class="action-button" onclick={handleManageCommunities}>
        <span class="action-icon">
          <img src="/setting_icon.svg" alt="Settings" width="20" height="20" />
        </span>
        <span class="action-label">Manage Communities</span>
      </button>

      <!-- Communities List -->
      <div class="communities-list">
        {#each communities as community (community.id)}
          <div class="community-item">
            <button
              class="community-link"
              onclick={() => navigateToCommunity(community.name)}
            >
              <span class="community-icon">{community.icon || "📁"}</span>
              <span class="community-name">lk/{community.name}</span>
            </button>
            <button
              class="favorite-button"
              class:active={community.isFavorite}
              onclick={() => toggleFavorite(community.id)}
            >
              {#if community.isFavorite}
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 16 16"
                  fill="currentColor"
                >
                  <path
                    d="M8 1.5l1.5 4.5h4.5l-3.5 2.5 1.5 4.5L8 10.5 4 13l1.5-4.5L2 6h4.5L8 1.5z"
                  />
                </svg>
              {:else}
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                >
                  <path
                    d="M8 1.5l1.5 4.5h4.5l-3.5 2.5 1.5 4.5L8 10.5 4 13l1.5-4.5L2 6h4.5L8 1.5z"
                    stroke-width="1.5"
                  />
                </svg>
              {/if}
            </button>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>

<CreateCommunityModal
  show={showCreateModal}
  onClose={() => (showCreateModal = false)}
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

  .communities-list {
    margin-top: 8px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .community-item {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 4px;
    padding: 0 4px;
  }

  .community-link {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 8px;
    background: transparent;
    border: none;
    cursor: pointer;
    border-radius: 6px;
    transition: background 0.15s ease;
    min-width: 0;
  }

  .community-link:hover {
    background: rgba(0, 0, 0, 0.08);
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
  }

  .community-name {
    font-size: 14px;
    font-weight: 500;
    color: #1c1c1c;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    text-align: left;
  }

  .favorite-button {
    width: 20px;
    height: 20px;
    padding: 0;
    background: transparent;
    border: none;
    cursor: pointer;
    color: #d3d3d3;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: color 0.15s ease;
  }

  .favorite-button:hover {
    color: #ffd700;
  }

  .favorite-button.active {
    color: #ffd700;
  }

  .favorite-button.active:hover {
    color: #ffed4e;
  }
</style>
