<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { completeGoogleSetup } from "../services/auth-service";
  import { ApiError } from "../errors/api-error";
  import Button from "../components/Button.svelte";

  let setupToken = $state("");
  let username = $state("");
  let isLoading = $state(false);
  let error = $state("");
  let status = $state<"form" | "error">("form");

  onMount(() => {
    // Đọc setup_token từ query params trong hash
    // URL format: /#/auth/google-setup?setup_token=...
    // window.location.hash = "#/auth/google-setup?setup_token=..."
    const hash = window.location.hash; // "#/auth/google-setup?setup_token=..."

    // Tìm vị trí dấu ?
    const queryIndex = hash.indexOf("?");
    if (queryIndex === -1) {
      status = "error";
      error = "Không tìm thấy mã xác thực. Vui lòng thử đăng nhập lại.";
      return;
    }

    // Lấy phần query string sau dấu ?
    const queryString = hash.substring(queryIndex + 1); // "setup_token=..."
    const urlParams = new URLSearchParams(queryString);
    const token = urlParams.get("setup_token");

    console.log("Full URL:", window.location.href);
    console.log("Hash:", hash);
    console.log("Query string:", queryString);
    console.log("Parsed token:", token ? "Found" : "Not found");

    if (token) {
      setupToken = token;
    } else {
      status = "error";
      error = "Không tìm thấy mã xác thực. Vui lòng thử đăng nhập lại.";
    }
  });

  async function handleSubmit(e: Event) {
    e.preventDefault();

    // Validate username
    if (!username || username.trim().length < 3) {
      error = "Tên người dùng phải có ít nhất 3 ký tự";
      return;
    }

    if (username.length > 20) {
      error = "Tên người dùng không được vượt quá 20 ký tự";
      return;
    }

    try {
      isLoading = true;
      error = "";

      await completeGoogleSetup(setupToken, username.trim());

      // Redirect về home sau khi thành công
      push("/");
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Đã xảy ra lỗi không mong muốn. Vui lòng thử lại.";
      }
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="setup-container">
  {#if status === "error"}
    <div class="error-card">
      <div class="error-icon">✗</div>
      <h2>Có lỗi xảy ra</h2>
      <p>{error}</p>
      <button class="retry-btn" onclick={() => push("/")}>
        Quay về trang chủ
      </button>
    </div>
  {:else}
    <div class="setup-card">
      <h1>Hoàn tất đăng ký</h1>
      <p class="subtitle">
        Chào mừng bạn! Vui lòng chọn tên người dùng để hoàn tất đăng ký.
      </p>

      <form onsubmit={handleSubmit}>
        <div class="input-group">
          <label for="username">Tên người dùng</label>
          <input
            id="username"
            type="text"
            bind:value={username}
            placeholder="Nhập tên người dùng (3-20 ký tự)"
            disabled={isLoading}
            minlength="3"
            maxlength="20"
            required
          />
          <span class="input-hint">
            Tên người dùng sẽ hiển thị công khai trên profile của bạn
          </span>
        </div>

        {#if error}
          <div class="error" role="alert">{error}</div>
        {/if}

        <Button
          type="submit"
          label={isLoading ? "Đang xử lý..." : "Hoàn tất"}
          variant="primary"
          disabled={isLoading || !username}
        />
      </form>
    </div>
  {/if}
</div>

<style>
  .setup-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background-color: #f6f7f8;
    padding: 20px;
  }

  .setup-card {
    background: white;
    border-radius: 8px;
    padding: 40px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    max-width: 480px;
    width: 100%;
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

  h1 {
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 8px 0;
    text-align: center;
  }

  h2 {
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 12px 0;
  }

  .subtitle {
    font-size: 14px;
    color: #7c7c7c;
    text-align: center;
    margin: 0 0 32px 0;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 1rem;
  }

  .input-group label {
    font-size: 14px;
    font-weight: 500;
    color: #1a1a1b;
  }

  .input-group input {
    width: 100%;
    padding: 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
    background: white;
    transition: border-color 0.2s;
    box-sizing: border-box;
  }

  .input-group input:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .input-group input:disabled {
    background-color: #f6f7f8;
    cursor: not-allowed;
  }

  .input-hint {
    font-size: 12px;
    color: #878a8c;
  }

  .error {
    color: #ff4500;
    font-size: 12px;
    padding: 8px 12px;
    background: rgba(255, 69, 0, 0.1);
    border-radius: 4px;
    margin-bottom: 1rem;
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

  .error-card p {
    font-size: 16px;
    color: #7c7c7c;
    margin: 0 0 24px 0;
  }

  .retry-btn {
    padding: 10px 24px;
    background-color: var(--primary-color);
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
