<script lang="ts">
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import { getCommunitiesByUserId } from "../services/community-service";
  import { authStore } from "../stores/auth-store";
  import type { CommunityResponse } from "../dtos/community-dto";
  import ConfirmModal from "../components/ConfirmModal.svelte";

  type Community = {
    id: string;
    name: string;
    description: string;
    avatar?: string;
    isOwner: boolean;
  };

  let allCommunities = $state<Community[]>([]);
  let isLoading = $state(true);
  let error = $state<string | null>(null);
  let searchQuery = $state("");
  let activeTab = $state<"owned" | "joined">("owned");

  async function loadCommunities() {
    try {
      isLoading = true;
      error = null;
      const userId = $authStore.user?.id;
      if (!userId) {
        error = "Bạn cần đăng nhập để xem cộng đồng";
        return;
      }

      const communities = await getCommunitiesByUserId(userId);
      allCommunities = communities.map((c: CommunityResponse) => ({
        id: c.id,
        name: c.name,
        description: c.description || "",
        avatar: c.avatar,
        isOwner: c.create_by_id === userId,
      }));
    } catch (err) {
      error = err instanceof Error ? err.message : "Không thể tải cộng đồng";
      console.error("Error loading communities:", err);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadCommunities();
  });

  let filteredCommunities = $derived(() => {
    let filtered = allCommunities;

    // Filter by tab
    if (activeTab === "owned") {
      filtered = filtered.filter((c) => c.isOwner);
    } else {
      filtered = filtered.filter((c) => !c.isOwner);
    }

    // Filter by search
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (c) =>
          c.name.toLowerCase().includes(query) ||
          c.description.toLowerCase().includes(query),
      );
    }

    return filtered;
  });

  let ownedCount = $derived(allCommunities.filter((c) => c.isOwner).length);
  let joinedCount = $derived(allCommunities.filter((c) => !c.isOwner).length);

  let showLeaveConfirm = $state(false);
  let communityToLeave = $state<Community | null>(null);

  function leaveCommunity(communityId: string) {
    const community = allCommunities.find((c) => c.id === communityId);
    if (community) {
      communityToLeave = community;
      showLeaveConfirm = true;
    }
  }

  function confirmLeaveCommunity() {
    if (communityToLeave) {
      // TODO: Call API to leave community
      allCommunities = allCommunities.filter(
        (c) => c.id !== communityToLeave!.id,
      );
      showLeaveConfirm = false;
      communityToLeave = null;
    }
  }

  function navigateToCommunity(communityName: string) {
    push(`/lk/${communityName}`);
  }
</script>

<div class="manage-communities-page">
  <div class="page-header">
    <h1>Quản lý cộng đồng</h1>
  </div>

  <div class="page-content">
    <!-- Tabs -->
    <div class="tabs">
      <button
        class="tab"
        class:active={activeTab === "owned"}
        onclick={() => (activeTab = "owned")}
      >
        Cộng đồng của tôi
        <span class="tab-count">{ownedCount}</span>
      </button>
      <button
        class="tab"
        class:active={activeTab === "joined"}
        onclick={() => (activeTab = "joined")}
      >
        Đã tham gia
        <span class="tab-count">{joinedCount}</span>
      </button>
    </div>
    <!-- Search Bar -->
    <div class="search-bar">
      <svg
        class="search-icon"
        width="20"
        height="20"
        viewBox="0 0 20 20"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <circle
          cx="8.5"
          cy="8.5"
          r="5.5"
          stroke="currentColor"
          stroke-width="2"
        />
        <path
          d="M12.5 12.5L17 17"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        />
      </svg>
      <input
        type="text"
        placeholder="Tìm kiếm cộng đồng"
        bind:value={searchQuery}
      />
    </div>

    <!-- Communities List -->
    {#if isLoading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Đang tải...</p>
      </div>
    {:else if error}
      <div class="error-state">
        <p>{error}</p>
        <button onclick={loadCommunities}>Thử lại</button>
      </div>
    {:else}
      <div class="communities-list">
        {#each filteredCommunities() as community (community.id)}
          <div class="community-card">
            <div class="community-main">
              <button
                class="community-info"
                onclick={() => navigateToCommunity(community.name)}
              >
                <div class="community-avatar">
                  {#if community.avatar}
                    <img src={community.avatar} alt={community.name} />
                  {:else}
                    <span class="avatar-placeholder"
                      >{community.name.charAt(0).toUpperCase()}</span
                    >
                  {/if}
                </div>
                <div class="community-details">
                  <h3 class="community-name">lk/{community.name}</h3>
                  <p class="community-description">{community.description}</p>
                </div>
              </button>

              <div class="community-actions">
                <button class="joined-btn">Đã tham gia</button>
              </div>
            </div>

            <!-- Leave button (appears on hover) -->
            <button
              class="leave-btn"
              onclick={() => leaveCommunity(community.id)}
            >
              Rời khỏi lk/{community.name}
            </button>
          </div>
        {/each}

        {#if filteredCommunities().length === 0}
          <div class="empty-state">
            <p>Không tìm thấy cộng đồng nào</p>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<ConfirmModal
  show={showLeaveConfirm}
  title="Xác nhận rời cộng đồng"
  message={`Bạn có chắc chắn muốn rời khỏi lk/${communityToLeave?.name || ""}?`}
  confirmText="Rời"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmLeaveCommunity}
  onCancel={() => {
    showLeaveConfirm = false;
    communityToLeave = null;
  }}
/>

<style>
  .manage-communities-page {
    max-width: 1000px;
    margin: 0 auto;
    padding: 20px;
  }

  .page-header {
    margin-bottom: 24px;
  }

  .page-header h1 {
    font-size: 28px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0;
  }

  .page-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  /* Tabs */
  .tabs {
    display: flex;
    gap: 0;
    border-bottom: 2px solid #ededed;
    margin-bottom: 8px;
  }

  .tab {
    padding: 12px 24px;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
    color: #7c7c7c;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .tab:hover {
    color: #1c1c1c;
  }

  .tab.active {
    color: var(--blue--);
    border-bottom-color: var(--blue--);
  }

  .tab-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 24px;
    height: 20px;
    padding: 0 6px;
    background: #f0f0f0;
    border-radius: 10px;
    font-size: 12px;
    font-weight: 600;
    color: #7c7c7c;
  }

  .tab.active .tab-count {
    background: var(--blue--);
    color: white;
  }

  /* Search Bar */
  .search-bar {
    position: relative;
    width: 100%;
  }

  .search-icon {
    position: absolute;
    left: 16px;
    top: 50%;
    transform: translateY(-50%);
    color: #878a8c;
  }

  .search-bar input {
    width: 100%;
    padding: 12px 16px 12px 48px;
    border: 1px solid #ccc;
    border-radius: 24px;
    font-size: 14px;
    background: white;
    transition: border-color 0.2s;
  }

  .search-bar input:focus {
    outline: none;
    border-color: #0079d3;
  }

  .search-bar input::placeholder {
    color: #878a8c;
  }

  /* Loading & Error States */
  .loading-state,
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 20px;
    gap: 16px;
  }

  .spinner {
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

  .loading-state p,
  .error-state p {
    font-size: 16px;
    color: #878a8c;
    margin: 0;
  }

  .error-state button {
    padding: 8px 16px;
    background: var(--blue--);
    color: white;
    border: none;
    border-radius: 20px;
    cursor: pointer;
    font-weight: 600;
  }

  /* Communities List */
  .communities-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .community-card {
    background: white;
    border: 1px solid #ccc;
    border-radius: 8px;
    padding: 16px;
    transition: all 0.2s;
    position: relative;
  }

  .community-card:hover {
    border-color: #878a8c;
  }

  .community-card:hover .leave-btn {
    opacity: 1;
    pointer-events: auto;
  }

  .community-main {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  .community-info {
    flex: 1;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    text-align: left;
    min-width: 0;
  }

  .community-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    overflow: hidden;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .community-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    color: white;
    font-size: 20px;
    font-weight: 700;
  }

  .community-details {
    flex: 1;
    min-width: 0;
  }

  .community-name {
    font-size: 16px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0 0 4px 0;
  }

  .community-description {
    font-size: 14px;
    color: #7c7c7c;
    margin: 0;
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .community-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
  }

  .joined-btn {
    padding: 8px 24px;
    background: transparent;
    border: 1px solid var(--blue--);
    border-radius: 24px;
    color: var(--blue--);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .joined-btn:hover {
    background: var(--blue--);
    color: white;
  }

  .leave-btn {
    position: absolute;
    bottom: 16px;
    left: 50%;
    transform: translateX(-50%);
    padding: 8px 16px;
    background: #1c1c1c;
    color: white;
    border: none;
    border-radius: 24px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    opacity: 0;
    pointer-events: none;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .leave-btn:hover {
    background: #000;
  }

  .empty-state {
    text-align: center;
    padding: 48px 20px;
    color: #878a8c;
  }

  .empty-state p {
    font-size: 16px;
    margin: 0;
  }
</style>
