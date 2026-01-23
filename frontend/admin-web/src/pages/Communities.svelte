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

  // Ban modal state
  let showBanModal = $state(false);
  let targetCommunityId = $state<string | null>(null);
  let banReason = $state("");

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
    targetCommunityId = communityId;
    banReason = "";
    showBanModal = true;
  }

  async function confirmBanCommunity() {
    if (!targetCommunityId || !banReason.trim()) return;
    showBanModal = false;

    try {
      await banCommunity(targetCommunityId, banReason);
      await loadCommunities();
      toastStore.success("Đã cấm cộng đồng");
    } catch (error) {
      toastStore.error("Không thể cấm cộng đồng");
    } finally {
      targetCommunityId = null;
      banReason = "";
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

<!-- Ban Community Modal -->
{#if showBanModal}
  <div class="modal-overlay" onclick={() => (showBanModal = false)}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Cấm cộng đồng</h3>
      </div>
      <div class="modal-body">
        <label for="ban-reason">Lý do cấm:</label>
        <input
          type="text"
          id="ban-reason"
          bind:value={banReason}
          placeholder="Nhập lý do cấm..."
        />
      </div>
      <div class="modal-actions">
        <button
          class="btn-cancel"
          onclick={() => {
            showBanModal = false;
            targetCommunityId = null;
          }}>Hủy</button
        >
        <button
          class="btn-confirm"
          onclick={confirmBanCommunity}
          disabled={!banReason.trim()}
        >
          Cấm
        </button>
      </div>
    </div>
  </div>
{/if}

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

  /* Ban Modal Styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10000;
  }

  .modal-content {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 400px;
    padding: 24px;
  }

  .modal-header h3 {
    margin: 0 0 16px 0;
    font-size: 18px;
    font-weight: 600;
  }

  .modal-body {
    margin-bottom: 20px;
  }

  .modal-body label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
  }

  .modal-body input {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
  }

  .modal-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
  }

  .btn-cancel,
  .btn-confirm {
    padding: 10px 20px;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    border: none;
  }

  .btn-cancel {
    background: #f0f0f0;
    color: #333;
  }

  .btn-confirm {
    background: #ff4500;
    color: white;
  }

  .btn-confirm:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
