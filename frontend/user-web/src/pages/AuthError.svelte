<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";

  let errorMessage = $state("");

  onMount(() => {
    // Đọc error message từ query params
    const urlParams = new URLSearchParams(window.location.search);
    const message = urlParams.get("message");

    if (message) {
      errorMessage = decodeURIComponent(message);
    } else {
      errorMessage = "Đã xảy ra lỗi không xác định";
    }
  });
</script>

<div class="error-container">
  <div class="error-card">
    <div class="error-icon">✗</div>
    <h1>Đăng nhập thất bại</h1>
    <p class="error-message">{errorMessage}</p>
    <button class="home-btn" onclick={() => push("/")}>
      Quay về trang chủ
    </button>
  </div>
</div>

<style>
  .error-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background-color: #f6f7f8;
    padding: 20px;
  }

  .error-card {
    background: white;
    border-radius: 8px;
    padding: 40px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    max-width: 480px;
    width: 100%;
    text-align: center;
  }

  .error-icon {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background-color: #f44336;
    color: white;
    font-size: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
  }

  h1 {
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 16px 0;
  }

  .error-message {
    font-size: 16px;
    color: #7c7c7c;
    margin: 0 0 32px 0;
    line-height: 1.5;
  }

  .home-btn {
    padding: 12px 24px;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .home-btn:hover {
    opacity: 0.9;
  }
</style>
