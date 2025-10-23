<script lang="ts">
  import type { Draft } from "../mocks/drafts.mock";
  import { mockDraftsList, totalDrafts } from "../mocks/drafts.mock";

  interface Props {
    show: boolean;
    onClose: () => void;
    onEditDraft?: (draftId: number) => void;
  }

  let { show, onClose, onEditDraft }: Props = $props();

  function handleEdit(draftId: number) {
    if (onEditDraft) {
      onEditDraft(draftId);
    }
  }

  function handleDelete(draftId: number) {
    console.log("Delete draft:", draftId);
    // TODO: Implement delete functionality
  }
</script>

{#if show}
  <div class="modal-overlay" onclick={onClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h2>
          Drafts <span class="draft-count"
            >{mockDraftsList.length}/{totalDrafts}</span
          >
        </h2>
      </div>

      <div class="drafts-list">
        {#each mockDraftsList as draft}
          <div class="draft-item">
            <div class="draft-info">
              <h3 class="draft-title">{draft.title}</h3>
              <p class="draft-time">{draft.editedTime}</p>
            </div>
            <div class="draft-actions">
              <button
                class="draft-action-btn edit-btn"
                onclick={() => handleEdit(draft.id)}
                title="Edit draft"
              >
                <img src="/write_icon.svg" alt="Edit" width="20" height="20" />
              </button>
              <button
                class="draft-action-btn delete-btn"
                onclick={() => handleDelete(draft.id)}
                title="Delete draft"
              >
                <img src="/bin_icon.svg" alt="Delete" width="20" height="20" />
              </button>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}

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
