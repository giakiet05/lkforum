<script lang="ts">
  type ConfirmModalProps = {
    show: boolean;
    title?: string;
    message: string;
    confirmText?: string;
    cancelText?: string;
    confirmVariant?: "danger" | "primary";
    onConfirm: () => void;
    onCancel: () => void;
  };

  let {
    show = false,
    title = "Xác nhận",
    message,
    confirmText = "Xác nhận",
    cancelText = "Hủy",
    confirmVariant = "danger",
    onConfirm,
    onCancel,
  }: ConfirmModalProps = $props();

  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      onCancel();
    }
  }

  function handleEscape(e: KeyboardEvent) {
    if (e.key === "Escape") {
      onCancel();
    }
  }

  $effect(() => {
    if (show) {
      document.addEventListener("keydown", handleEscape);
      document.body.style.overflow = "hidden";
    } else {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "unset";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "unset";
    };
  });
</script>

{#if show}
  <div
    class="confirm-overlay"
    role="dialog"
    aria-modal="true"
    onclick={handleOverlayClick}
    onkeydown={handleEscape}
  >
    <div class="confirm-modal">
      <div class="confirm-header">
        <h3>{title}</h3>
      </div>
      <div class="confirm-body">
        <p>{message}</p>
      </div>
      <div class="confirm-actions">
        <button class="btn-cancel" onclick={onCancel}>
          {cancelText}
        </button>
        <button
          class="btn-confirm"
          class:danger={confirmVariant === "danger"}
          class:primary={confirmVariant === "primary"}
          onclick={onConfirm}
        >
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .confirm-overlay {
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

  .confirm-modal {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 400px;
    padding: 24px;
    animation: slideIn 0.2s ease-out;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  .confirm-header h3 {
    margin: 0 0 16px 0;
    font-size: 18px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .confirm-body {
    margin-bottom: 24px;
  }

  .confirm-body p {
    margin: 0;
    font-size: 14px;
    color: #576f76;
    line-height: 1.5;
  }

  .confirm-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
  }

  .btn-cancel,
  .btn-confirm {
    padding: 10px 20px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-cancel {
    background: rgba(214, 216, 222, 0.5);
    color: #1c1c1c;
    border: none;
  }

  .btn-cancel:hover {
    background: rgba(214, 216, 222, 0.7);
  }

  .btn-confirm.danger {
    background: #ff4500;
    color: white;
    border: none;
  }

  .btn-confirm.danger:hover {
    background: #e03d00;
  }

  .btn-confirm.primary {
    background: #0079d3;
    color: white;
    border: none;
  }

  .btn-confirm.primary:hover {
    background: #0060a8;
  }
</style>
