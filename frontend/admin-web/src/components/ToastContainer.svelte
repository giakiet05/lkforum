<script lang="ts">
  import { toastStore } from "../stores/toast-store";
  import type { Toast } from "../stores/toast-store";

  let toasts: Toast[] = $state([]);

  toastStore.subscribe((value) => {
    toasts = value;
  });

  function getIcon(type: string): string {
    switch (type) {
      case "success":
        return "✓";
      case "error":
        return "✕";
      case "warning":
        return "⚠";
      case "info":
        return "ℹ";
      default:
        return "ℹ";
    }
  }
</script>

<div class="toast-container">
  {#each toasts as toast (toast.id)}
    <div class="toast toast-{toast.type}">
      <span class="toast-icon">{getIcon(toast.type)}</span>
      <span class="toast-message">{toast.message}</span>
      <button class="toast-close" onclick={() => toastStore.remove(toast.id)}
        >✕</button
      >
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
    gap: 10px;
    max-width: 400px;
  }

  .toast {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    background: white;
    animation: slideIn 0.3s ease-out;
    min-width: 300px;
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

  .toast-icon {
    font-size: 20px;
    font-weight: bold;
    flex-shrink: 0;
  }

  .toast-message {
    flex: 1;
    font-size: 14px;
    color: #333;
  }

  .toast-close {
    background: none;
    border: none;
    font-size: 18px;
    cursor: pointer;
    color: #666;
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .toast-close:hover {
    color: #333;
  }

  .toast-success {
    border-left: 4px solid #10b981;
  }

  .toast-success .toast-icon {
    color: #10b981;
  }

  .toast-error {
    border-left: 4px solid #ef4444;
  }

  .toast-error .toast-icon {
    color: #ef4444;
  }

  .toast-warning {
    border-left: 4px solid #f59e0b;
  }

  .toast-warning .toast-icon {
    color: #f59e0b;
  }

  .toast-info {
    border-left: 4px solid #3b82f6;
  }

  .toast-info .toast-icon {
    color: #3b82f6;
  }

  @media (max-width: 640px) {
    .toast-container {
      left: 10px;
      right: 10px;
      max-width: none;
    }

    .toast {
      min-width: auto;
    }
  }
</style>
