<script lang="ts">
  import type { CommunityResponse } from "../dtos/community-dto";

  interface Props {
    communities: CommunityResponse[];
    loading?: boolean;
    onBan: (communityId: string) => void;
    onUnban: (communityId: string) => void;
  }

  let { communities, loading = false, onBan, onUnban }: Props = $props();

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
  {:else if communities.length === 0}
    <div class="empty">Không có cộng đồng nào</div>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Avatar</th>
          <th>Tên</th>
          <th>Mô tả</th>
          <th>Thành viên</th>
          <th>Trạng thái</th>
          <th>Ngày tạo</th>
          <th>Thao tác</th>
        </tr>
      </thead>
      <tbody>
        {#each communities as community}
          <tr class:banned={community.is_banned}>
            <td>
              {#if community.avatar}
                <img
                  src={community.avatar}
                  alt={community.name}
                  class="avatar"
                />
              {:else}
                <div class="avatar-placeholder">{community.name[0]}</div>
              {/if}
            </td>
            <td>
              <div class="community-name">
                <strong>{community.name}</strong>
                <small>c/{community.name}</small>
              </div>
            </td>
            <td>
              <div class="description">
                {community.description || "Không có mô tả"}
              </div>
            </td>
            <td>{community.member_count || 0}</td>
            <td>
              {#if community.is_banned}
                <span class="status banned">Đã cấm</span>
              {:else}
                <span class="status active">Hoạt động</span>
              {/if}
            </td>
            <td>{formatDate(community.created_at)}</td>
            <td>
              <div class="actions">
                {#if community.is_banned}
                  <button
                    class="btn-unban"
                    onclick={() => onUnban(community.id)}>Gỡ cấm</button
                  >
                {:else}
                  <button class="btn-ban" onclick={() => onBan(community.id)}
                    >Cấm</button
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

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #4a70a9;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 1.25rem;
  }

  .community-name {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .community-name small {
    color: #6c757d;
  }

  .description {
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #6c757d;
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

  .loading,
  .empty {
    padding: 3rem;
    text-align: center;
    color: #6c757d;
  }
</style>
