<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";

  let errorMessage = $state("");
  let showLoginModal = $state(false);

  onMount(() => {
    // Đọc error message từ hash URL (e.g., #/auth/error?message=...)
    const hash = window.location.hash;
    const queryIndex = hash.indexOf("?");
    const queryString = queryIndex >= 0 ? hash.substring(queryIndex + 1) : "";
    const urlParams = new URLSearchParams(queryString);
    const message = urlParams.get("message");

    if (message) {
      errorMessage = decodeURIComponent(message);
    } else {
      errorMessage = "Đã xảy ra lỗi không xác định trong quá trình đăng nhập.";
    }
  });

  function handleLoginClick() {
    showLoginModal = true;
  }

  function handleCloseModal() {
    showLoginModal = false;
  }
</script>

<div class="error-container">
  <div class="error-card">
    <div class="error-icon">
      <svg
        width="32"
        height="32"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="15" y1="9" x2="9" y2="15"></line>
        <line x1="9" y1="9" x2="15" y2="15"></line>
      </svg>
    </div>
    <h1>Đăng nhập thất bại</h1>
    <p class="error-message">{errorMessage}</p>
    <div class="button-group">
      <button class="login-btn" onclick={handleLoginClick}>
        Đăng nhập lại
      </button>
      <button class="home-btn" onclick={() => push("/")}>
        Quay về trang chủ
      </button>
    </div>
  </div>
</div>

{#if showLoginModal}
  <div class="modal-overlay" onclick={handleCloseModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <button class="close-btn" onclick={handleCloseModal}>
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
      <h2>Đăng nhập</h2>
      <p class="modal-subtitle">Chọn phương thức đăng nhập</p>

      <a href="{API_BASE_URL}/api/auth/google" class="google-btn">
        <svg width="20" height="20" viewBox="0 0 24 24">
          <path
            fill="#4285F4"
            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
          />
          <path
            fill="#34A853"
            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
          />
          <path
            fill="#FBBC05"
            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
          />
          <path
            fill="#EA4335"
            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
          />
        </svg>
        Đăng nhập với Google
      </a>

      <div class="divider">
        <span>hoặc</span>
      </div>

      <button
        class="email-btn"
        onclick={() => {
          handleCloseModal();
          push("/login");
        }}
      >
        Đăng nhập bằng Email
      </button>

      <p class="register-link">
        Chưa có tài khoản?
        <button
          onclick={() => {
            handleCloseModal();
            push("/register");
          }}>Đăng ký ngay</button
        >
      </p>
    </div>
  </div>
{/if}

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
    border-radius: 12px;
    padding: 48px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    max-width: 480px;
    width: 100%;
    text-align: center;
  }

  .error-icon {
    width: 72px;
    height: 72px;
    border-radius: 50%;
    background-color: #fee2e2;
    color: #ef4444;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
  }

  h1 {
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 12px 0;
  }

  .error-message {
    font-size: 15px;
    color: #6b7280;
    margin: 0 0 32px 0;
    line-height: 1.6;
  }

  .button-group {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .login-btn {
    padding: 14px 24px;
    background-color: var(--primary-color, #4a70a9);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .login-btn:hover {
    background-color: #3d5d8f;
  }

  .home-btn {
    padding: 14px 24px;
    background-color: transparent;
    color: #6b7280;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .home-btn:hover {
    background-color: #f3f4f6;
  }

  /* Modal styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    background: white;
    border-radius: 12px;
    padding: 32px;
    max-width: 400px;
    width: 90%;
    position: relative;
    text-align: center;
  }

  .close-btn {
    position: absolute;
    top: 16px;
    right: 16px;
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    padding: 4px;
  }

  .close-btn:hover {
    color: #1f2937;
  }

  .modal-content h2 {
    font-size: 22px;
    font-weight: 600;
    margin: 0 0 8px 0;
    color: #1f2937;
  }

  .modal-subtitle {
    color: #6b7280;
    margin: 0 0 24px 0;
    font-size: 14px;
  }

  .google-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    width: 100%;
    padding: 12px 16px;
    background: white;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
    color: #374151;
    cursor: pointer;
    transition: all 0.2s;
    text-decoration: none;
  }

  .google-btn:hover {
    background: #f9fafb;
    border-color: #9ca3af;
  }

  .divider {
    display: flex;
    align-items: center;
    margin: 20px 0;
    color: #9ca3af;
    font-size: 13px;
  }

  .divider::before,
  .divider::after {
    content: "";
    flex: 1;
    height: 1px;
    background: #e5e7eb;
  }

  .divider span {
    padding: 0 12px;
  }

  .email-btn {
    width: 100%;
    padding: 12px 16px;
    background: var(--primary-color, #4a70a9);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .email-btn:hover {
    background: #3d5d8f;
  }

  .register-link {
    margin-top: 20px;
    color: #6b7280;
    font-size: 14px;
  }

  .register-link button {
    background: none;
    border: none;
    color: var(--primary-color, #4a70a9);
    font-weight: 500;
    cursor: pointer;
    text-decoration: underline;
  }

  .register-link button:hover {
    text-decoration: none;
  }
</style>
