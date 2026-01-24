<script lang="ts">
  import { toastStore, type Toast } from "../stores/toast-store";

  let toasts: Toast[] = $state([]);

  toastStore.subscribe((value) => {
    toasts = value;
  });

  function getIcon(type: Toast["type"]) {
    switch (type) {
      case "success":
        return "✓";
      case "error":
        return "✕";
      case "warning":
        return "⚠";
      case "info":
        return "ℹ";
    }
  }
</script>

<div class="toast-container">
  {#each toasts as toast (toast.id)}
    <div class="toast toast-{toast.type}" role="alert">
      <span class="toast-icon">{getIcon(toast.type)}</span>
      <span class="toast-message">{toast.message}</span>
      <button
        class="toast-close"
        onclick={() => toastStore.remove(toast.id)}
        aria-label="Đóng"
      >
        ✕
      </button>
    </div>
  {/each}
</div>

<style>
  .toast-container {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 10000;
    display: flex;
    flex-direction: column;
    gap: 12px;
    pointer-events: none;
  }

  .toast {
    pointer-events: auto;
    display: flex;
    align-items: center;
    gap: 12px;
    min-width: 300px;
    max-width: 500px;
    padding: 16px 20px;
    background: white;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    animation: slideIn 0.3s ease-out;
    border-left: 4px solid;
    font-size: 14px;
    line-height: 1.5;
  }

  @keyframes slideIn {
    from {
      transform: translateX(400px);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  .toast-success {
    border-left-color: #10b981;
    background: #f0fdf4;
  }

  .toast-error {
    border-left-color: #ef4444;
    background: #fef2f2;
  }

  .toast-warning {
    border-left-color: #f59e0b;
    background: #fffbeb;
  }

  .toast-info {
    border-left-color: #3b82f6;
    background: #eff6ff;
  }

  .toast-icon {
    flex-shrink: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    font-weight: bold;
    font-size: 16px;
  }

  .toast-success .toast-icon {
    background: #10b981;
    color: white;
  }

  .toast-error .toast-icon {
    background: #ef4444;
    color: white;
  }

  .toast-warning .toast-icon {
    background: #f59e0b;
    color: white;
  }

  .toast-info .toast-icon {
    background: #3b82f6;
    color: white;
  }

  .toast-message {
    flex: 1;
    color: #374151;
    font-weight: 500;
  }

  .toast-close {
    flex-shrink: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    color: #6b7280;
    cursor: pointer;
    border-radius: 4px;
    font-size: 18px;
    transition: all 0.2s;
  }

  .toast-close:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #374151;
  }

  @media (max-width: 640px) {
    .toast-container {
      top: 10px;
      right: 10px;
      left: 10px;
    }

    .toast {
      min-width: auto;
      width: 100%;
    }
  }
</style>
