<script lang="ts">
  import { onMount } from "svelte";
  import CommunityTable from "../components/CommunityTable.svelte";
  import {
    getCommunities,
    banCommunity,
    unbanCommunity,
  } from "../services/community-service";
  import { toastStore } from "../stores/toast-store";
  import type { CommunityResponse } from "../dtos/community-dto";

  let loading = $state(false);
  let communities = $state<CommunityResponse[]>([]);

  async function loadCommunities() {
    loading = true;
    try {
      const response = await getCommunities({ limit: 100 });
      communities = response.communities;
    } catch (error) {
      console.error("Failed to load communities:", error);
      toastStore.error("Không thể tải danh sách cộng đồng");
    } finally {
      loading = false;
    }
  }

  async function handleBanCommunity(communityId: string) {
    const reason = prompt("Lý do cấm:");
    if (!reason) return;

    try {
      await banCommunity(communityId, reason);
      await loadCommunities();
      toastStore.success("Đã cấm cộng đồng");
    } catch (error) {
      toastStore.error("Không thể cấm cộng đồng");
    }
  }

  async function handleUnbanCommunity(communityId: string) {
    try {
      await unbanCommunity(communityId);
      await loadCommunities();
      toastStore.success("Đã gỡ cấm cộng đồng");
    } catch (error) {
      toastStore.error("Không thể gỡ cấm cộng đồng");
    }
  }

  onMount(() => {
    loadCommunities();
  });
</script>

<div class="communities-page">
  <div class="page-header">
    <h1>Quản lý cộng đồng</h1>
  </div>

  <div class="page-content">
    {#if loading}
      <div class="loading">Đang tải...</div>
    {:else}
      <CommunityTable
        {communities}
        onBan={handleBanCommunity}
        onUnban={handleUnbanCommunity}
      />
    {/if}
  </div>
</div>

<style>
  .communities-page {
    height: 100%;
  }

  .page-header {
    margin-bottom: 2rem;
  }

  .page-header h1 {
    font-size: 1.75rem;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0;
  }

  .page-content {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .loading {
    padding: 3rem;
    text-align: center;
    color: #6c757d;
    font-size: 1rem;
  }
</style>
