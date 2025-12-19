<script lang="ts">
  import type { Report } from "../dtos/report-dto";

  interface Props {
    reports: Report[];
    loading?: boolean;
    onDelete: (reportId: string) => void;
    onView?: (reportId: string) => void;
  }

  let { reports, loading = false, onDelete, onView }: Props = $props();

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString("vi-VN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  function getTargetTypeLabel(type: string) {
    switch (type) {
      case "user":
        return "Người dùng";
      case "post":
        return "Bài viết";
      case "comment":
        return "Bình luận";
      default:
        return type;
    }
  }
</script>

<div class="table-container">
  {#if loading}
    <div class="loading">Đang tải...</div>
  {:else if reports.length === 0}
    <div class="empty">Không có báo cáo nào</div>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Loại</th>
          <th>Lý do</th>
          <th>Mô tả</th>
          <th>Người báo cáo</th>
          <th>Ngày tạo</th>
          <th>Thao tác</th>
        </tr>
      </thead>
      <tbody>
        {#each reports as report}
          <tr class:deleted={report.is_deleted}>
            <td>
              <span
                class="badge"
                class:user={report.target_type === "user"}
                class:post={report.target_type === "post"}
                class:comment={report.target_type === "comment"}
              >
                {getTargetTypeLabel(report.target_type)}
              </span>
            </td>
            <td>
              <strong>{report.reason}</strong>
            </td>
            <td>
              <div class="description">
                {report.description || "Không có mô tả"}
              </div>
            </td>
            <td>
              <code>{report.reporter_id.substring(0, 8)}...</code>
            </td>
            <td>{formatDate(report.created_at)}</td>
            <td>
              <div class="actions">
                {#if onView}
                  <button class="btn-view" onclick={() => onView?.(report.id)}
                    >Xem</button
                  >
                {/if}
                {#if !report.is_deleted}
                  <button class="btn-delete" onclick={() => onDelete(report.id)}
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
  }

  .badge.user {
    background: #d4edda;
    color: #155724;
  }

  .badge.post {
    background: #d1ecf1;
    color: #0c5460;
  }

  .badge.comment {
    background: #fff3cd;
    color: #856404;
  }

  .description {
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #6c757d;
  }

  code {
    background: #f8f9fa;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-family: monospace;
    font-size: 0.875rem;
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

  .btn-view {
    background: #007bff;
    color: white;
  }

  .btn-view:hover {
    background: #0056b3;
  }

  .btn-delete {
    background: #dc3545;
    color: white;
  }

  .btn-delete:hover {
    background: #c82333;
  }

  .loading,
  .empty {
    padding: 3rem;
    text-align: center;
    color: #6c757d;
  }
</style>
