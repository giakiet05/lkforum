<script lang="ts">
  import { onMount } from "svelte";
  import UserTable from "../components/UserTable.svelte";
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
    const reason = prompt("Lý do cấm:");
    if (!reason) return;

    try {
      await banUser(userId, reason);
      await loadUsers();
      toastStore.success("Đã cấm người dùng");
    } catch (error) {
      toastStore.error("Không thể cấm người dùng");
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
    if (!confirm("Bạn có chắc muốn xóa người dùng này?")) return;

    try {
      await deleteUser(userId);
      await loadUsers();
      toastStore.success("Đã xóa người dùng");
    } catch (error) {
      toastStore.error("Không thể xóa người dùng");
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
</style>
