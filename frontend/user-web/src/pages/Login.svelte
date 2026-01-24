<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import Button from "../components/Button.svelte";
  import type {
    CompleteRegistrationRequest,
    LoginRequest,
  } from "../dtos/auth-dto";

  export type Props = {
    mode?: "login" | "register";
    onSubmit: (data: LoginRequest | CompleteRegistrationRequest) => void;
    isLoading?: boolean;
    error?: string;
  };

  let {
    mode = "login",
    onSubmit,
    isLoading = false,
    error = "",
  }: Props = $props();

  const dispatch = createEventDispatcher();

  // form fields
  let identifier = $state(""); // username hoặc email tuỳ backend
  let email = $state(""); // email field cho register mode
  let password = $state("");
  let showPassword = $state(false);

  // Validator đơn giản
  function validate() {
    if (!identifier || !password) {
      error = "Vui lòng nhập tên đăng nhập và mật khẩu";
      return false;
    }
    if (mode === "register" && !email) {
      error = "Vui lòng nhập email";
      return false;
    }
    return true;
  }

  // Chuyển đổi hiển thị mật khẩu
  function toggleShowPasswordVisibility() {
    showPassword = !showPassword;
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      toggleShowPasswordVisibility();
    }
  }

  // Xử lý submit form
  async function handleSubmit() {
    if (!validate()) return;
    onSubmit({
      identifier: identifier,
      password: password,
    });
  }

  function handleGoogleLogin() {
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";
    window.location.href = `${API_BASE_URL}/api/auth/google/login`;
  }
</script>

<form
  onsubmit={(e) => {
    e.preventDefault();
    handleSubmit();
  }}
  class="login-form"
>
  <div class="input-group">
    <label for="identifier">Tên đăng nhập</label>
    <input
      id="identifier"
      type="text"
      bind:value={identifier}
      placeholder={mode === "register"
        ? "Chọn tên đăng nhập"
        : "Nhập tên đăng nhập"}
    />
  </div>

  {#if mode === "register"}
    <div class="input-group">
      <label for="email">Email</label>
      <input
        id="email"
        type="email"
        bind:value={email}
        placeholder="Nhập email của bạn"
      />
    </div>
  {/if}

  <div class="input-group password-group">
    <label for="password">Mật khẩu</label>
    <input
      id="password"
      type={showPassword ? "text" : "password"}
      bind:value={password}
      placeholder="Mật khẩu"
    />
    <button
      type="button"
      class="password-toggle-icon"
      onclick={toggleShowPasswordVisibility}
      onkeydown={handleKeydown}
      aria-label={showPassword ? "Ẩn mật khẩu" : "Hiện mật khẩu"}
    >
      {#if showPassword}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
      {:else}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
          <path
            d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
          />
          <path
            d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"
          />
          <line x1="2" x2="22" y1="2" y2="22" />
        </svg>
      {/if}
    </button>
  </div>

  {#if error}
    <div class="error" role="alert">{error}</div>
  {/if}

  {#if mode === "login"}
    <div class="forgot-password-link">
      <button
        type="button"
        class="text-button"
        onclick={() => dispatch("forgotPassword")}
      >
        Quên mật khẩu?
      </button>
    </div>
  {/if}

  <Button
    type="submit"
    label={isLoading
      ? "Đang xử lý..."
      : mode === "register"
        ? "Đăng Ký"
        : "Đăng Nhập"}
    variant="primary"
    disabled={isLoading}
  />

  <!-- Google login tạm ẩn
  <div class="separator">
    <span>HOẶC</span>
  </div>

  <Button
    label="Đăng nhập với Google"
    variant="google"
    onclick={handleGoogleLogin}
  />
  -->

  <div class="signup-link">
    {mode === "register" ? "Đã có tài khoản?" : "Chưa có tài khoản?"}
    <button
      type="button"
      class="text-button"
      onclick={() =>
        dispatch("switchMode", {
          mode: mode === "register" ? "login" : "register",
        })}
    >
      {mode === "register" ? "Đăng nhập" : "Đăng ký ngay"}
    </button>
  </div>
</form>

<style>
  .login-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .input-group {
    margin-bottom: 1rem;
  }

  .input-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #1a1a1b;
  }

  .input-group input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
    background: white;
    transition: border-color 0.2s;
  }

  .input-group input:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .password-group {
    position: relative;
    margin-bottom: 1rem;
  }

  .password-toggle-icon {
    position: absolute;
    top: 50%;
    right: 8px;
    transform: translateY(-50%);
    background: none;
    border: none;
    cursor: pointer;
    color: #878a8c;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    margin-top: 10px; /* Điều chỉnh vị trí so với input */
  }

  .password-toggle-icon:hover {
    color: #666;
  }

  .password-toggle-icon svg {
    width: 18px;
    height: 18px;
  }

  .password-group input {
    padding-right: 40px;
    height: 44px;
  }

  .separator {
    display: flex;
    align-items: center;
    text-align: center;
    color: #787c7e;
    margin: 1rem 0;
  }

  .separator::before,
  .separator::after {
    content: "";
    flex: 1;
    border-bottom: 1px solid #edeff1;
  }

  .separator:not(:empty)::before {
    margin-right: 1rem;
  }

  .separator:not(:empty)::after {
    margin-left: 1rem;
  }

  .signup-link {
    text-align: center;
    margin-top: 1rem;
    color: #787c7e;
    font-size: 12px;
  }

  .text-button {
    background: none;
    border: none;
    padding: 0;
    color: var(--primary-color);
    font-weight: 600;
    cursor: pointer;
    font-size: inherit;
    text-decoration: none;
  }

  .text-button:hover {
    text-decoration: underline;
  }

  .forgot-password-link {
    text-align: right;
    margin-top: -0.5rem;
    margin-bottom: 0.5rem;
  }

  .forgot-password-link .text-button {
    font-size: 12px;
    color: #787c7e;
  }

  .error {
    color: #ff4500;
    font-size: 12px;
    margin-bottom: 1rem;
    padding: 8px 12px;
    background: rgba(255, 69, 0, 0.1);
    border-radius: 4px;
  }
</style>
