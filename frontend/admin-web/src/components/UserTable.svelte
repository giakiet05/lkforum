<script lang="ts">
  import type { UserResponse } from "../dtos/user-dto";

  interface Props {
    users: UserResponse[];
    loading?: boolean;
    onBan: (userId: string) => void;
    onUnban: (userId: string) => void;
    onDelete: (userId: string) => void;
  }

  let { users, loading = false, onBan, onUnban, onDelete }: Props = $props();

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString("vi-VN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });
  }
</script>

<div class="table-container">
  {#if loading}
    <div class="loading">Đang tải...</div>
  {:else if users.length === 0}
    <div class="empty">Không có người dùng nào</div>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Username</th>
          <th>Email</th>
          <th>Role</th>
          <th>Reputation</th>
          <th>Trạng thái</th>
          <th>Ngày tạo</th>
          <th>Thao tác</th>
        </tr>
      </thead>
      <tbody>
        {#each users as user}
          <tr class:banned={user.is_banned} class:deleted={user.deleted_at}>
            <td>{user.username}</td>
            <td>{user.email}</td>
            <td>
              <span class="badge" class:admin={user.role === "admin"}>
                {user.role}
              </span>
            </td>
            <td>{user.reputation}</td>
            <td>
              {#if user.deleted_at}
                <span class="status deleted">Đã xóa</span>
              {:else if user.is_banned}
                <span class="status banned">Đã cấm</span>
              {:else}
                <span class="status active">Hoạt động</span>
              {/if}
            </td>
            <td>{formatDate(user.created_at)}</td>
            <td>
              <div class="actions">
                {#if user.deleted_at}
                  <button class="btn-restore" onclick={() => {}}
                    >Khôi phục</button
                  >
                {:else if user.is_banned}
                  <button class="btn-unban" onclick={() => onUnban(user.id)}
                    >Gỡ cấm</button
                  >
                {:else}
                  <button class="btn-ban" onclick={() => onBan(user.id)}
                    >Cấm</button
                  >
                  <button class="btn-delete" onclick={() => onDelete(user.id)}
                    >Xóa</button
                  >
                {/if}
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .table-container {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: #f8f9fa;
  }

  th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #333;
    font-size: 0.875rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  td {
    padding: 1rem;
    border-top: 1px solid #e9ecef;
  }

  tr.banned {
    background-color: #fff5f5;
  }

  tr.deleted {
    background-color: #f5f5f5;
    opacity: 0.7;
  }

  .badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 600;
    background: #e9ecef;
    color: #495057;
  }

  .badge.admin {
    background: #4a70a9;
    color: white;
  }

  .status {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .status.active {
    background: #d4edda;
    color: #155724;
  }

  .status.banned {
    background: #f8d7da;
    color: #721c24;
  }

  .status.deleted {
    background: #e9ecef;
    color: #6c757d;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  button {
    padding: 0.375rem 0.75rem;
    border: none;
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .btn-ban {
    background: #dc3545;
    color: white;
  }

  .btn-ban:hover {
    background: #c82333;
  }

  .btn-unban {
    background: #28a745;
    color: white;
  }

  .btn-unban:hover {
    background: #218838;
  }

  .btn-delete {
    background: #6c757d;
    color: white;
  }

  .btn-delete:hover {
    background: #5a6268;
  }

  .btn-restore {
    background: #007bff;
    color: white;
  }

  .btn-restore:hover {
    background: #0056b3;
  }

  .loading,
  .empty {
    padding: 3rem;
    text-align: center;
    color: #6c757d;
  }
</style>
