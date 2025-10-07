<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let label: string;
  export let variant: "primary" | "google" = "primary";
  export let disabled: boolean = false;
  export let type: "button" | "submit" = "button";

  const dispatch = createEventDispatcher();

  function handleClick() {
    dispatch("click");
  }
</script>

<button
  {type}
  {disabled}
  class="btn"
  class:primary={variant === "primary"}
  class:google={variant === "google"}
  on:click={handleClick}
>
  {label}

  {#if variant === "google"}
    <img
      src="https://www.svgrepo.com/show/475656/google-color.svg"
      alt="Google icon"
      class="google-icon"
    />
  {/if}
</button>

<style>
  /* Style của button giữ nguyên như cũ, rất tốt */
  .btn {
    width: 100%;
    padding: 1rem;
    border: none;
    border-radius: 8px;
    font-size: 1.1em;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.75rem;
  }
  .btn:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }
  .primary {
    background-color: var(--primary-color, #154d71); /* Thêm màu mặc định */
    color: white;
  }
  .btn:focus-visible {
    /* Thay thế outline mặc định của trình duyệt */
    outline: 3px solid var(--primary-color-hover);
    outline-offset: 2px; /* Tạo một khoảng cách nhỏ giữa nút và outline */
  }
  .primary:hover:not(:disabled) {
    background-color: #1c6ea4;
  }
  .google {
    background-color: var(--btn-google-bg, #e9f1ff); /* Thêm màu mặc định */
    color: #333;
  }
  .google:hover:not(:disabled) {
    background-color: #cccccc;
  }
  .google-icon {
    width: 24px;
    height: 24px;
  }
</style>
