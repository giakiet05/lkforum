<script lang="ts">
  import { onDestroy } from "svelte";

  type ModalProps = {
    show?: boolean;
    title: string;
    onClose: () => void;
  };

  let { show = false, title, onClose }: ModalProps = $props();

  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }

  function handleEscape(e: KeyboardEvent) {
    if (e.key === "Escape") {
      onClose();
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

    // Cleanup
    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "unset";
    };
  });
</script>

{#if show}
  <div class="modal-overlay" role="dialog" onclick={handleOverlayClick}>
    <div class="modal-container">
      <div class="modal-header">
        <h2>{title}</h2>
        <button class="close-button" onclick={onClose}>×</button>
      </div>
      <div class="modal-content">
        <slot />
      </div>
    </div>
  </div>
{/if}

<style>
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
    z-index: 1000;
  }

  .modal-container {
    background-color: white;
    border-radius: 8px;
    width: 95%;
    max-width: 400px;
    max-height: 90vh;
    overflow-y: auto;
    padding: 20px;
    animation: modal-appear 0.2s ease-out;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .modal-header h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 500;
  }

  .close-button {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    padding: 0;
    color: #878a8c;
  }

  .close-button:hover {
    color: #1a1a1b;
  }

  .modal-content {
    position: relative;
  }

  @keyframes modal-appear {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
