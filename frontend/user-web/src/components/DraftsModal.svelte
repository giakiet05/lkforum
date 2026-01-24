<script lang="ts">
  import { getDrafts, deleteDraft } from "../services/draft-service";
  import type { DraftSummaryResponse } from "../dtos/draft-dto";
  import { toastStore } from "../stores/toast-store";
  import ConfirmModal from "./ConfirmModal.svelte";

  interface Props {
    show: boolean;
    onClose: () => void;
    onEditDraft?: (draftId: string) => void;
  }

  let { show, onClose, onEditDraft }: Props = $props();

  let drafts = $state<DraftSummaryResponse[]>([]);
  let totalDrafts = $state(0);
  let isLoading = $state(false);
  let showDeleteConfirm = $state(false);
  let draftToDelete = $state<string | null>(null);

  // Load drafts when modal opens
  $effect(() => {
    if (show) {
      loadDrafts();
    }
  });

  async function loadDrafts() {
    try {
      isLoading = true;
      const response = await getDrafts(1, 20);
      drafts = response.drafts;
      totalDrafts = response.pagination.total;
    } catch (error) {
      console.error("Failed to load drafts:", error);
    } finally {
      isLoading = false;
    }
  }

  function handleEdit(draftId: string) {
    if (onEditDraft) {
      onEditDraft(draftId);
    }
  }

  async function handleDelete(draftId: string) {
    draftToDelete = draftId;
    showDeleteConfirm = true;
  }

  async function confirmDelete() {
    if (!draftToDelete) return;
    showDeleteConfirm = false;

    try {
      await deleteDraft(draftToDelete);
      // Reload drafts after deletion
      await loadDrafts();
    } catch (error) {
      console.error("Failed to delete draft:", error);
      toastStore.error("Không thể xóa bản nháp. Vui lòng thử lại.");
    } finally {
      draftToDelete = null;
    }
  }

  function formatTime(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 60) {
      return `${diffMins} phút trước`;
    } else if (diffHours < 24) {
      return `${diffHours} giờ trước`;
    } else if (diffDays < 7) {
      return `${diffDays} ngày trước`;
    } else {
      return date.toLocaleDateString();
    }
  }
</script>

{#if show}
  <div class="modal-overlay" onclick={onClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h2>
          Bản nháp <span class="draft-count">{drafts.length}/{totalDrafts}</span
          >
        </h2>
      </div>

      {#if isLoading}
        <div class="loading">Đang tải bản nháp...</div>
      {:else if drafts.length === 0}
        <div class="empty-state">
          <p>Chưa có bản nháp</p>
        </div>
      {:else}
        <div class="drafts-list">
          {#each drafts as draft}
            <div class="draft-item">
              <div class="draft-info">
                <h3 class="draft-title">
                  {draft.title || "Bản nháp chưa đặt tên"}
                </h3>
                <p class="draft-time">{formatTime(draft.updated_at)}</p>
              </div>
              <div class="draft-actions">
                <button
                  class="draft-action-btn edit-btn"
                  onclick={() => handleEdit(draft.id)}
                  title="Sửa bản nháp"
                >
                  <img
                    src="/write_icon.svg"
                    alt="Edit"
                    width="20"
                    height="20"
                  />
                </button>
                <button
                  class="draft-action-btn delete-btn"
                  onclick={() => handleDelete(draft.id)}
                  title="Xóa bản nháp"
                >
                  <img
                    src="/bin_icon.svg"
                    alt="Delete"
                    width="20"
                    height="20"
                  />
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
{/if}

<ConfirmModal
  show={showDeleteConfirm}
  title="Xác nhận xóa"
  message="Bạn có chắc muốn xóa bản nháp này?"
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDelete}
  onCancel={() => {
    showDeleteConfirm = false;
    draftToDelete = null;
  }}
/>

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    z-index: 1000;
    padding: 60px 20px 20px;
    overflow-y: auto;
  }

  .modal-content {
    background: white;
    border-radius: 8px;
    width: 100%;
    max-width: 640px;
    padding: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .modal-header {
    margin-bottom: 24px;
  }

  .modal-header h2 {
    font-size: 20px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0;
  }

  .draft-count {
    opacity: 0.6;
  }

  .loading,
  .empty-state {
    text-align: center;
    padding: 40px 20px;
    color: #666;
  }

  .empty-state p {
    margin: 0;
    font-size: 16px;
  }

  .drafts-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .draft-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    background: #f6f7f8;
    border-radius: 8px;
    transition: background 0.2s;
  }

  .draft-item:hover {
    background: #edeff1;
  }

  .draft-info {
    flex: 1;
  }

  .draft-title {
    font-size: 16px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0 0 4px 0;
  }

  .draft-time {
    font-size: 13px;
    color: var(--grayfont);
    margin: 0;
  }

  .draft-actions {
    display: flex;
    gap: 8px;
  }

  .draft-action-btn {
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .draft-action-btn:hover {
    background: #f6f7f8;
    border-color: var(--lightgray--);
  }

  .delete-btn {
    background: #fff0f0;
    border-color: #ffcccc;
  }

  .delete-btn:hover {
    background: #ffe0e0;
    border-color: #ff9999;
  }

  .delete-btn img {
    filter: brightness(0) saturate(100%) invert(23%) sepia(89%) saturate(7471%)
      hue-rotate(357deg) brightness(95%) contrast(118%);
  }

  .draft-action-btn img {
    display: block;
  }
</style>
