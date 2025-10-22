<script lang="ts">
  import Modal from "./Modal.svelte";
  import Login from "../pages/Login.svelte";
  import { login } from "../services/auth-service";
  import type { LoginRequest } from "../dtos/auth-dto";
  import { setAuth } from "../stores/auth-store";
  import Button from "./Button.svelte";

  let { show = false, onClose }: { show: boolean; onClose: () => void } =
    $props();

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

  let activeTab = $state<"login" | "register">("login");
  let step = $state<"register" | "verify">("register");
  let isLoading = $state(false);
  let error = $state("");

  // Register form fields
  let email = $state("");
  let username = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let otp = $state(["", "", "", "", "", ""]); // 6 ô OTP

  // OTP countdown timer
  let countdown = $state(60);
  let canResend = $state(false);
  let countdownInterval: number | null = null;

  function resetForm() {
    email = "";
    username = "";
    password = "";
    confirmPassword = "";
    otp = ["", "", "", "", "", ""];
    step = "register";
    activeTab = "login";
    error = "";
    if (countdownInterval) clearInterval(countdownInterval);
    countdown = 60;
    canResend = false;

    // Xóa pending verification khi reset form
    localStorage.removeItem("pending_verification_email");
  }

  // Wrapper để reset form khi đóng modal
  function handleClose() {
    resetForm();
    onClose();
  }

  // Check for pending email verification on mount
  $effect(() => {
    if (show) {
      const pendingEmail = localStorage.getItem("pending_verification_email");
      if (pendingEmail) {
        email = pendingEmail;
        activeTab = "register";
        step = "verify";
        error =
          "Bạn chưa xác thực email. Vui lòng nhập mã OTP đã gửi đến email của bạn.";
        startCountdown();
      }
    }
  });

  function startCountdown() {
    countdown = 60;
    canResend = false;

    if (countdownInterval) clearInterval(countdownInterval);

    countdownInterval = setInterval(() => {
      countdown--;
      if (countdown <= 0) {
        if (countdownInterval) clearInterval(countdownInterval);
        canResend = true;
      }
    }, 1000);
  }

  function formatCountdown(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  }

  const handleLogin = async (data: LoginRequest) => {
    try {
      isLoading = true;
      error = "";
      await login(data);
      handleClose();
    } catch (err) {
      error = err instanceof Error ? err.message : "Login failed";
    } finally {
      isLoading = false;
    }
  };

  function handleSwitchToRegister() {
    activeTab = "register";
    error = "";
  }

  function handleSwitchToLogin() {
    activeTab = "login";
    step = "register";
    error = "";

    // Xóa pending verification khi chuyển sang đăng nhập
    localStorage.removeItem("pending_verification_email");
  }

  // Step 1: Register
  async function handleRegisterSubmit(e: Event) {
    e.preventDefault();

    if (!email || !username || !password || !confirmPassword) {
      error = "Vui lòng điền đầy đủ thông tin";
      return;
    }
    if (password !== confirmPassword) {
      error = "Mật khẩu xác nhận không khớp!";
      return;
    }

    isLoading = true;
    error = "";

    try {
      const res = await fetch(`${API_BASE_URL}/api/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, username, password }),
      });

      if (!res.ok) {
        const errObj = await res
          .json()
          .catch(() => ({ error: `HTTP ${res.status}` }));
        throw errObj.error || errObj.message || "Lỗi đăng ký";
      }

      step = "verify";
      error = "";

      // Lưu email để tracking pending verification
      localStorage.setItem("pending_verification_email", email);

      startCountdown(); // Bắt đầu đếm ngược
    } catch (err: any) {
      error = typeof err === "string" ? err : err.message || "Lỗi khi đăng ký";
    } finally {
      isLoading = false;
    }
  }

  // Step 2: Verify OTP
  async function handleVerifyOTP(e: Event) {
    e.preventDefault();

    const otpCode = otp.join("");
    if (otpCode.length !== 6) {
      error = "Vui lòng nhập đầy đủ 6 chữ số";
      return;
    }

    isLoading = true;
    error = "";

    try {
      const res = await fetch(`${API_BASE_URL}/api/auth/verify-email`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, otp: otpCode }),
      });

      if (!res.ok) {
        const errObj = await res
          .json()
          .catch(() => ({ error: "Mã OTP không đúng" }));
        throw errObj.error || errObj.message || "Mã OTP không đúng";
      }

      const data = await res.json();

      // Lưu tokens và user
      localStorage.setItem("access_token", data.access_token);
      localStorage.setItem("refresh_token", data.refresh_token);
      localStorage.setItem("user", JSON.stringify(data.user));

      // Update authStore
      setAuth(data.user, data.access_token);

      // Xóa pending verification
      localStorage.removeItem("pending_verification_email");

      // Đóng modal và reset form
      handleClose();
    } catch (err: any) {
      error =
        typeof err === "string" ? err : err.message || "Mã OTP không đúng";
    } finally {
      isLoading = false;
    }
  }

  async function handleResendOTP() {
    if (!canResend) return;

    isLoading = true;
    error = "";

    try {
      const res = await fetch(
        `${API_BASE_URL}/api/auth/resend-verification-email`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email }),
        }
      );

      if (!res.ok) throw new Error("Không thể gửi lại mã OTP");

      alert("Mã OTP mới đã được gửi đến email của bạn!");
      startCountdown(); // Bắt đầu đếm ngược lại
    } catch (err: any) {
      error = "Không thể gửi lại mã OTP";
    } finally {
      isLoading = false;
    }
  }

  function handleBackToRegister() {
    step = "register";
    otp = ["", "", "", "", "", ""];
    error = "";
    if (countdownInterval) clearInterval(countdownInterval);
    countdown = 60;
    canResend = false;

    // Xóa pending verification nếu quay lại
    localStorage.removeItem("pending_verification_email");
  }

  // Handle OTP input with auto-focus next box
  function handleOtpInput(index: number, e: Event) {
    const input = e.target as HTMLInputElement;
    const value = input.value;

    // Only allow numbers
    if (value && !/^\d$/.test(value)) {
      input.value = "";
      return;
    }

    otp[index] = value;

    // Auto focus next input
    if (value && index < 5) {
      const nextInput =
        input.parentElement?.nextElementSibling?.querySelector("input");
      nextInput?.focus();
    }
  }

  function handleOtpKeydown(index: number, e: KeyboardEvent) {
    const input = e.target as HTMLInputElement;

    // Backspace - go to previous input
    if (e.key === "Backspace" && !input.value && index > 0) {
      const prevInput =
        input.parentElement?.previousElementSibling?.querySelector("input");
      prevInput?.focus();
    }

    // Arrow keys navigation
    if (e.key === "ArrowLeft" && index > 0) {
      const prevInput =
        input.parentElement?.previousElementSibling?.querySelector("input");
      prevInput?.focus();
    }
    if (e.key === "ArrowRight" && index < 5) {
      const nextInput =
        input.parentElement?.nextElementSibling?.querySelector("input");
      nextInput?.focus();
    }
  }

  function handleOtpPaste(e: ClipboardEvent) {
    e.preventDefault();
    const pastedData = e.clipboardData?.getData("text").slice(0, 6) || "";
    const digits = pastedData
      .split("")
      .filter((char) => /\d/.test(char))
      .slice(0, 6);

    digits.forEach((digit, i) => {
      otp[i] = digit;
    });

    // Focus last filled input or first empty
    const focusIndex = Math.min(digits.length, 5);
    const inputs = document.querySelectorAll(".otp-input");
    (inputs[focusIndex] as HTMLInputElement)?.focus();
  }
</script>

<Modal
  {show}
  onClose={handleClose}
  title={activeTab === "login"
    ? "Đăng nhập vào LKForum"
    : step === "register"
      ? "Tạo tài khoản mới"
      : "Xác thực Email"}
>
  {#if activeTab === "login"}
    <Login
      mode="login"
      onSubmit={(data) => handleLogin(data as LoginRequest)}
      {isLoading}
      {error}
      on:switchMode={handleSwitchToRegister}
    />
  {:else if step === "register"}
    <!-- Register Form -->
    <form on:submit={handleRegisterSubmit} class="auth-form">
      <div class="input-group">
        <label for="email">Email</label>
        <input
          type="email"
          id="email"
          bind:value={email}
          placeholder="Nhập email của bạn"
          disabled={isLoading}
        />
      </div>

      <div class="input-group">
        <label for="username">Username</label>
        <input
          type="text"
          id="username"
          bind:value={username}
          placeholder="Chọn một tên đăng nhập"
          disabled={isLoading}
        />
      </div>

      <div class="input-group">
        <label for="password">Mật khẩu</label>
        <input
          type="password"
          id="password"
          bind:value={password}
          placeholder="Tạo mật khẩu"
          disabled={isLoading}
        />
      </div>

      <div class="input-group">
        <label for="confirmPassword">Xác nhận mật khẩu</label>
        <input
          type="password"
          id="confirmPassword"
          bind:value={confirmPassword}
          placeholder="Nhập lại mật khẩu"
          disabled={isLoading}
        />
      </div>

      {#if error}
        <div class="error" role="alert">{error}</div>
      {/if}

      <Button
        type="submit"
        label={isLoading ? "Đang xử lý..." : "Đăng Ký"}
        variant="primary"
        disabled={isLoading}
      />

      <div class="switch-mode">
        Đã có tài khoản?
        <button type="button" class="link-btn" on:click={handleSwitchToLogin}>
          Đăng nhập
        </button>
      </div>
    </form>
  {:else}
    <!-- Verify OTP Form -->
    <form on:submit={handleVerifyOTP} class="auth-form otp-form">
      <p class="otp-instruction">
        Chúng tôi đã gửi mã OTP đến email <strong>{email}</strong>
      </p>
      <p class="otp-hint">Vui lòng nhập mã gồm 6 chữ số</p>

      <div class="otp-inputs">
        {#each otp as _, i}
          <div class="otp-box">
            <input
              type="text"
              class="otp-input"
              maxlength="1"
              value={otp[i]}
              on:input={(e) => handleOtpInput(i, e)}
              on:keydown={(e) => handleOtpKeydown(i, e)}
              on:paste={i === 0 ? handleOtpPaste : undefined}
              disabled={isLoading}
            />
          </div>
        {/each}
      </div>

      {#if error}
        <div class="error" role="alert">{error}</div>
      {/if}

      <Button
        type="submit"
        label={isLoading ? "Đang xác thực..." : "Xác Nhận"}
        variant="primary"
        disabled={isLoading}
      />

      <div class="otp-timer">
        {#if !canResend}
          <span class="timer-text">{formatCountdown(countdown)}</span>
        {/if}
      </div>

      <div class="otp-actions">
        <p class="resend-text">
          Chưa nhận được OTP?
          <button
            type="button"
            class="back-btn"
            on:click={handleBackToRegister}
            disabled={isLoading}
          >
            Quay lại
          </button>
          <button
            type="button"
            class="resend-btn"
            on:click={handleResendOTP}
            disabled={!canResend || isLoading}
          >
            Gửi lại
          </button>
        </p>
      </div>

      <div class="switch-mode">
        Đã có tài khoản?
        <button
          type="button"
          class="login-link-btn"
          on:click={handleSwitchToLogin}
        >
          Đăng nhập
        </button>
      </div>
    </form>
  {/if}
</Modal>

<style>
  :global(.modal-container) {
    width: 420px !important;
    padding: 32px !important;
  }

  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .input-group label {
    font-size: 14px;
    font-weight: 500;
    color: #1a1a1b;
  }

  .input-group input {
    padding: 10px 12px;
    border: 1px solid #ccc;
    border-radius: 6px;
    font-size: 14px;
    transition: border-color 0.2s;
  }

  .input-group input:focus {
    outline: none;
    border-color: var(--darkblue--);
  }

  .input-group input:disabled {
    background-color: #f6f7f8;
    cursor: not-allowed;
  }

  .error {
    padding: 10px;
    background-color: #fee;
    border: 1px solid #fcc;
    border-radius: 6px;
    color: #c00;
    font-size: 14px;
  }

  .switch-mode {
    text-align: center;
    font-size: 14px;
    color: #666;
    margin-top: 8px;
  }

  .link-btn {
    background: none;
    border: none;
    color: var(--primary-color);
    cursor: pointer;
    font-weight: 600;
    padding: 0;
    text-decoration: none;
  }

  .link-btn:hover {
    text-decoration: underline;
  }

  .link-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* OTP Form Styles */
  .otp-form {
    text-align: center;
  }

  .otp-instruction {
    font-size: 14px;
    color: #1a1a1b;
    margin: 0 0 8px 0;
  }

  .otp-hint {
    font-size: 13px;
    color: #7c7c7c;
    margin: 0 0 24px 0;
  }

  .otp-inputs {
    display: flex;
    justify-content: center;
    gap: 8px;
    margin-bottom: 20px;
  }

  .otp-box {
    width: 48px;
    height: 56px;
  }

  .otp-input {
    width: 100%;
    height: 100%;
    text-align: center;
    font-size: 24px;
    font-weight: 600;
    border: 2px solid #ccc;
    border-radius: 8px;
    transition: border-color 0.2s;
  }

  .otp-input:focus {
    outline: none;
    border-color: var(--darkblue--);
    box-shadow: 0 0 0 3px rgba(21, 77, 113, 0.1);
  }

  .otp-input:disabled {
    background-color: #f6f7f8;
    cursor: not-allowed;
  }

  .otp-actions {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    margin-top: 16px;
    font-size: 14px;
  }

  .otp-timer {
    text-align: center;
    margin-top: 12px;
    min-height: 24px;
  }

  .timer-text {
    font-size: 16px;
    font-weight: 600;
    color: var(--error--);
  }

  .resend-text {
    margin: 0;
    color: #000;
    font-size: 14px;
  }

  .back-btn {
    background: none;
    border: none;
    color: #000;
    text-decoration: underline;
    cursor: pointer;
    padding: 0;
    margin-left: 4px;
    font-size: 14px;
  }

  .back-btn:hover {
    opacity: 0.7;
  }

  .back-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .resend-btn {
    background: none;
    border: none;
    color: var(--error--);
    cursor: pointer;
    padding: 0;
    margin-left: 4px;
    font-size: 14px;
    font-weight: 500;
  }

  .resend-btn:hover {
    opacity: 0.7;
  }

  .resend-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .login-link-btn {
    background: none;
    border: none;
    color: var(--primary-color);
    cursor: pointer;
    padding: 0;
    margin-left: 4px;
    font-size: 14px;
    font-weight: 600;
    text-decoration: none;
  }

  .login-link-btn:hover {
    text-decoration: underline;
  }
</style>
