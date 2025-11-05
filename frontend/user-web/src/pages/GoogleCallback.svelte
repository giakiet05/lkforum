<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { setAuth } from "../stores/auth-store";
  import {
    setAccessToken,
    setRefreshToken,
    setUser,
  } from "../services/storage-service";

  let status = $state<"loading" | "success" | "error">("loading");
  let errorMessage = $state("");

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const accessToken = urlParams.get("access_token");
    const refreshToken = urlParams.get("refresh_token");
    const userJson = urlParams.get("user");

    if (accessToken && refreshToken && userJson) {
      try {
        const user = JSON.parse(decodeURIComponent(userJson));

        // Lưu vào storage
        setAccessToken(accessToken);
        setRefreshToken(refreshToken);
        setUser(user);

        // Update auth store
        setAuth(user, accessToken);

        status = "success";

        // Redirect về home sau 1 giây
        setTimeout(() => {
          push("/");
        }, 1000);
      } catch (err) {
        console.error("Error parsing user data:", err);
        status = "error";
        errorMessage = "Không thể xử lý thông tin người dùng";
      }
    } else {
      status = "error";
      errorMessage = "Thiếu thông tin xác thực";
    }
  });
</script>

<div class="callback-container">
  {#if status === "loading"}
    <div class="spinner"></div>
    <h2>Đang xử lý đăng nhập...</h2>
  {:else if status === "success"}
    <div class="success-icon">✓</div>
    <h2>Đăng nhập thành công!</h2>
    <p>Đang chuyển hướng...</p>
  {:else}
    <div class="error-icon">✗</div>
    <h2>Đăng nhập thất bại</h2>
    <p>{errorMessage}</p>
    <button class="retry-btn" onclick={() => push("/")}>
      Quay về trang chủ
    </button>
  {/if}
</div>

<style>
  .callback-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 20px;
    text-align: center;
  }

  .spinner {
    width: 48px;
    height: 48px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid var(--blue--);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 24px;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .success-icon {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background-color: #4caf50;
    color: white;
    font-size: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 24px;
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
    margin-bottom: 24px;
  }

  h2 {
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 12px 0;
  }

  p {
    font-size: 16px;
    color: #7c7c7c;
    margin: 0 0 24px 0;
  }

  .retry-btn {
    padding: 10px 24px;
    background-color: var(--blue--);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .retry-btn:hover {
    opacity: 0.9;
  }
</style>
