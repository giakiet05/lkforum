<script lang="ts">
  import { onMount } from "svelte";
  import UserTable from "../components/UserTable.svelte";
  import ConfirmModal from "../components/ConfirmModal.svelte";
  import {
    getUsers,
    banUser,
    unbanUser,
    deleteUser,
  } from "../services/user-service";
  import { toastStore } from "../stores/toast-store";
  import type { UserResponse } from "../dtos/user-dto";

  let loading = $state(false);
  let users = $state<UserResponse[]>([]);

  // Modal states
  let showBanModal = $state(false);
  let showDeleteConfirm = $state(false);
  let targetUserId = $state<string | null>(null);
  let banReason = $state("");

  async function loadUsers() {
    loading = true;
    try {
      const response = await getUsers({ limit: 100 });
      users = response.users;
    } catch (error) {
      console.error("Failed to load users:", error);
      toastStore.error("Không thể tải danh sách người dùng");
    } finally {
      loading = false;
    }
  }

  async function handleBanUser(userId: string) {
    targetUserId = userId;
    banReason = "";
    showBanModal = true;
  }

  async function confirmBanUser() {
    if (!targetUserId || !banReason.trim()) return;
    showBanModal = false;

    try {
      await banUser(targetUserId, banReason);
      await loadUsers();
      toastStore.success("Đã cấm người dùng");
    } catch (error) {
      toastStore.error("Không thể cấm người dùng");
    } finally {
      targetUserId = null;
      banReason = "";
    }
  }

  async function handleUnbanUser(userId: string) {
    try {
      await unbanUser(userId);
      await loadUsers();
      toastStore.success("Đã gỡ cấm người dùng");
    } catch (error) {
      toastStore.error("Không thể gỡ cấm người dùng");
    }
  }

  async function handleDeleteUser(userId: string) {
    targetUserId = userId;
    showDeleteConfirm = true;
  }

  async function confirmDeleteUser() {
    if (!targetUserId) return;
    showDeleteConfirm = false;

    try {
      await deleteUser(targetUserId);
      await loadUsers();
      toastStore.success("Đã xóa người dùng");
    } catch (error) {
      toastStore.error("Không thể xóa người dùng");
    } finally {
      targetUserId = null;
    }
  }

  onMount(() => {
    loadUsers();
  });
</script>

<div class="users-page">
  <div class="page-header">
    <h1>Quản lý người dùng</h1>
  </div>

  <div class="page-content">
    {#if loading}
      <div class="loading">Đang tải...</div>
    {:else}
      <UserTable
        {users}
        onBan={handleBanUser}
        onUnban={handleUnbanUser}
        onDelete={handleDeleteUser}
      />
    {/if}
  </div>
</div>

<ConfirmModal
  show={showDeleteConfirm}
  title="Xác nhận xóa"
  message="Bạn có chắc muốn xóa người dùng này? Hành động này không thể hoàn tác."
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDeleteUser}
  onCancel={() => {
    showDeleteConfirm = false;
    targetUserId = null;
  }}
/>

<!-- Ban User Modal -->
{#if showBanModal}
  <div class="modal-overlay" onclick={() => (showBanModal = false)}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Cấm người dùng</h3>
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
            targetUserId = null;
          }}>Hủy</button
        >
        <button
          class="btn-confirm"
          onclick={confirmBanUser}
          disabled={!banReason.trim()}
        >
          Cấm
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .users-page {
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
