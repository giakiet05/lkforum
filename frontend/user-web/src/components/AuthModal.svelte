<script lang="ts">
  import Modal from "./Modal.svelte";
  import Login from "../pages/Login.svelte";
  import {
    login,
    sendVerificationEmail,
    verifyEmail,
    completeRegistration,
    forgotPassword,
    verifyResetOTP,
    resetPassword,
  } from "../services/auth-service";
  import type {
    LoginRequest,
    CompleteRegistrationRequest,
  } from "../dtos/auth-dto";
  import Button from "./Button.svelte";
  import { ApiError } from "../errors/api-error";
  import { ApiErrorCode } from "../errors/error-codes";
  import { toastStore } from "../stores/toast-store";

  let { show = false, onClose }: { show: boolean; onClose: () => void } =
    $props();

  let activeTab = $state<"login" | "register" | "forgotPassword">("login");
  let step = $state<"email" | "otp" | "register" | "resetPassword">("email"); // Step 1: email → Step 2: otp → Step 3: register/resetPassword
  let isLoading = $state(false);
  let error = $state("");

  // Register form fields
  let email = $state("");
  let username = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let otp = $state(["", "", "", "", "", ""]); // 6 ô OTP
  let verificationToken = $state(""); // Lưu verification token sau khi verify OTP

  // Forgot password fields
  let resetToken = $state(""); // Token để reset password
  let newPassword = $state("");
  let confirmNewPassword = $state("");

  // OTP countdown timer
  let countdown = $state(60);
  let canResend = $state(false);
  let countdownInterval: number | null = null;

  function resetForm() {
    email = "";
    username = "";
    password = "";
    confirmPassword = "";
    newPassword = "";
    confirmNewPassword = "";
    resetToken = "";
    otp = ["", "", "", "", "", ""];
    step = "email";
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

      await login(data); // Gọi login, tự động lưu tokens và update authStore
      handleClose(); // Đóng modal, hoàn tất!
    } catch (err) {
      if (err instanceof ApiError) {
        // Hiển thị message từ backend (đã tiếng Việt)
        error = err.message;
      } else {
        error = "Đã xảy ra lỗi không mong muốn";
      }
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
    step = "email";
    error = "";

    // Xóa pending verification khi chuyển sang đăng nhập
    localStorage.removeItem("pending_verification_email");
  }

  // Step 1: Nhập Email → Backend gửi OTP
  async function handleEmailSubmit(e: Event) {
    e.preventDefault();

    if (!email) {
      error = "Vui lòng nhập email";
      return;
    }

    isLoading = true;
    error = "";

    try {
      await sendVerificationEmail(email);

      // Chuyển sang step nhập OTP
      step = "otp";
      error = "";
      localStorage.setItem("pending_verification_email", email);
      startCountdown();
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Không thể gửi mã OTP. Vui lòng thử lại.";
      }
    } finally {
      isLoading = false;
    }
  }

  // Step 2: Verify OTP → Chuyển sang step Register
  async function handleVerifyOTPStep(e: Event) {
    e.preventDefault();

    const otpCode = otp.join("");
    if (otpCode.length !== 6) {
      error = "Vui lòng nhập đầy đủ 6 chữ số";
      return;
    }

    isLoading = true;
    error = "";

    try {
      // Verify OTP trả về verification_token
      verificationToken = await verifyEmail(email, otpCode);

      // Chuyển sang step đăng ký (username + password)
      step = "register";
      error = "";
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Mã OTP không đúng. Vui lòng thử lại.";
      }
    } finally {
      isLoading = false;
    }
  }

  // Step 3: Register - nhập username, password → Hoàn tất đăng ký
  async function handleRegisterSubmit(e: Event) {
    e.preventDefault();

    if (!username || !password || !confirmPassword) {
      error = "Vui lòng điền đầy đủ thông tin";
      return;
    }
    if (username.length < 3 || username.length > 20) {
      error = "Tên đăng nhập phải từ 3-20 ký tự";
      return;
    }
    if (password.length < 6) {
      error = "Mật khẩu phải có ít nhất 6 ký tự";
      return;
    }
    if (password !== confirmPassword) {
      error = "Mật khẩu xác nhận không khớp!";
      return;
    }

    isLoading = true;
    error = "";

    try {
      const registrationData: CompleteRegistrationRequest = {
        username,
        password,
        verification_token: verificationToken,
      };

      console.log("🔍 Registration data:", {
        username,
        has_password: !!password,
        has_token: !!verificationToken,
      });
      await completeRegistration(registrationData);

      // Xóa pending verification
      localStorage.removeItem("pending_verification_email");

      // Đóng modal
      handleClose();
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Lỗi khi đăng ký";
      }
    } finally {
      isLoading = false;
    }
  }

  async function handleResendOTP() {
    if (!canResend) return;

    isLoading = true;
    error = "";

    try {
      if (activeTab === "forgotPassword") {
        await forgotPassword(email);
      } else {
        await sendVerificationEmail(email);
      }
      toastStore.success("Mã OTP mới đã được gửi đến email của bạn!");
      startCountdown(); // Bắt đầu đếm ngược lại
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Không thể gửi lại mã OTP";
      }
    } finally {
      isLoading = false;
    }
  }

  // --- Forgot Password Functions ---

  function handleSwitchToForgotPassword() {
    activeTab = "forgotPassword";
    step = "email";
    error = "";
  }

  async function handleForgotPasswordEmailSubmit(e: Event) {
    e.preventDefault();

    if (!email) {
      error = "Vui lòng nhập email";
      return;
    }

    isLoading = true;
    error = "";

    try {
      await forgotPassword(email);
      step = "otp";
      error = "";
      startCountdown();
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.error_code === "LOGIN_METHOD_MISMATCH") {
          error =
            "Tài khoản này đã đăng ký bằng Google. Vui lòng đăng nhập bằng Google hoặc liên hệ quản trị viên để hỗ trợ.";
        } else {
          error = err.message;
        }
      } else {
        error = "Không thể gửi mã OTP. Vui lòng thử lại.";
      }
    } finally {
      isLoading = false;
    }
  }

  async function handleVerifyResetOTPStep(e: Event) {
    e.preventDefault();

    const otpCode = otp.join("");
    if (otpCode.length !== 6) {
      error = "Vui lòng nhập đầy đủ 6 chữ số";
      return;
    }

    isLoading = true;
    error = "";

    try {
      const response = await verifyResetOTP(email, otpCode);
      resetToken = response.reset_token;
      step = "resetPassword";
      error = "";
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Mã OTP không đúng. Vui lòng thử lại.";
      }
    } finally {
      isLoading = false;
    }
  }

  async function handleResetPasswordSubmit(e: Event) {
    e.preventDefault();

    if (!newPassword || !confirmNewPassword) {
      error = "Vui lòng điền đầy đủ thông tin";
      return;
    }
    if (newPassword.length < 6) {
      error = "Mật khẩu phải có ít nhất 6 ký tự";
      return;
    }
    if (newPassword !== confirmNewPassword) {
      error = "Mật khẩu xác nhận không khớp!";
      return;
    }

    isLoading = true;
    error = "";

    try {
      await resetPassword(resetToken, newPassword);
      toastStore.success("Đổi mật khẩu thành công! Vui lòng đăng nhập lại.");
      handleClose();
    } catch (err) {
      if (err instanceof ApiError) {
        error = err.message;
      } else {
        error = "Không thể đổi mật khẩu";
      }
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
    : activeTab === "forgotPassword" && step === "email"
      ? "Quên mật khẩu"
      : activeTab === "forgotPassword" && step === "otp"
        ? "Xác thực OTP"
        : activeTab === "forgotPassword" && step === "resetPassword"
          ? "Đặt lại mật khẩu"
          : step === "email"
            ? "Tạo tài khoản mới"
            : step === "otp"
              ? "Xác thực Email"
              : "Hoàn tất đăng ký"}
>
  {#if activeTab === "login"}
    <Login
      mode="login"
      onSubmit={(data) => handleLogin(data as LoginRequest)}
      {isLoading}
      {error}
      on:switchMode={handleSwitchToRegister}
      on:forgotPassword={handleSwitchToForgotPassword}
    />
  {:else if activeTab === "register" && step === "email"}
    <!-- Step 1: Nhập Email -->
    <form on:submit={handleEmailSubmit} class="auth-form">
      <p class="form-description">
        Nhập email để bắt đầu đăng ký tài khoản LKForum.
      </p>

      <div class="input-group">
        <label for="email">Email</label>
        <input
          type="email"
          id="email"
          bind:value={email}
          placeholder="Nhập email của bạn"
          disabled={isLoading}
          required
        />
      </div>

      {#if error}
        <div class="error" role="alert">{error}</div>
      {/if}

      <Button
        type="submit"
        label={isLoading ? "Đang gửi..." : "Tiếp Tục"}
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
  {:else if activeTab === "register" && step === "otp"}
    <!-- Step 2: Verify OTP -->
    <form on:submit={handleVerifyOTPStep} class="auth-form otp-form">
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
            on:click={() => {
              step = "email";
              otp = ["", "", "", "", "", ""];
              error = "";
            }}
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
        <button type="button" class="link-btn" on:click={handleSwitchToLogin}>
          Đăng nhập
        </button>
      </div>
    </form>
  {:else if step === "register"}
    <!-- Step 3: Register Form (Username + Password) -->
    <form on:submit={handleRegisterSubmit} class="auth-form">
      <p class="form-description">
        Tạo tên đăng nhập và mật khẩu cho tài khoản của bạn.
      </p>

      <div class="input-group">
        <label for="username">Tên đăng nhập</label>
        <input
          type="text"
          id="username"
          bind:value={username}
          placeholder="Chọn một tên đăng nhập"
          disabled={isLoading}
          required
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
          required
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
          required
        />
      </div>

      {#if error}
        <div class="error" role="alert">{error}</div>
      {/if}

      <Button
        type="submit"
        label={isLoading ? "Đang xử lý..." : "Hoàn Tất"}
        variant="primary"
        disabled={isLoading}
      />

      <div class="otp-actions">
        <button
          type="button"
          class="back-btn"
          on:click={() => {
            step = "otp";
            username = "";
            password = "";
            confirmPassword = "";
            error = "";
          }}
          disabled={isLoading}
        >
          Quay lại
        </button>
      </div>

      <div class="switch-mode">
        Đã có tài khoản?
        <button type="button" class="link-btn" on:click={handleSwitchToLogin}>
          Đăng nhập
        </button>
      </div>
    </form>
  {:else if activeTab === "forgotPassword" && step === "email"}
    <!-- Forgot Password: Step 1 - Email -->
    <form on:submit={handleForgotPasswordEmailSubmit} class="auth-form">
      <p class="form-description">
        Nhập email đã đăng ký để nhận mã OTP khôi phục mật khẩu.
      </p>

      <div class="input-group">
        <label for="email">Email</label>
        <input
          type="email"
          id="email"
          bind:value={email}
          placeholder="Nhập email của bạn"
          disabled={isLoading}
          required
        />
      </div>

      {#if error}
        <div class="error" role="alert">{error}</div>
      {/if}

      <Button
        type="submit"
        label={isLoading ? "Đang gửi..." : "Gửi mã OTP"}
        variant="primary"
        disabled={isLoading}
      />

      <div class="switch-mode">
        Nhớ mật khẩu?
        <button type="button" class="link-btn" on:click={handleSwitchToLogin}>
          Đăng nhập
        </button>
      </div>
    </form>
  {:else if activeTab === "forgotPassword" && step === "otp"}
    <!-- Forgot Password: Step 2 - Verify OTP -->
    <form on:submit={handleVerifyResetOTPStep} class="auth-form otp-form">
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
            on:click={() => {
              step = "email";
              otp = ["", "", "", "", "", ""];
              error = "";
            }}
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
        Nhớ mật khẩu?
        <button type="button" class="link-btn" on:click={handleSwitchToLogin}>
          Đăng nhập
        </button>
      </div>
    </form>
  {:else if activeTab === "forgotPassword" && step === "resetPassword"}
    <!-- Forgot Password: Step 3 - Reset Password -->
    <form on:submit={handleResetPasswordSubmit} class="auth-form">
      <p class="form-description">Tạo mật khẩu mới cho tài khoản của bạn.</p>

      <div class="input-group">
        <label for="newPassword">Mật khẩu mới</label>
        <input
          type="password"
          id="newPassword"
          bind:value={newPassword}
          placeholder="Nhập mật khẩu mới"
          disabled={isLoading}
          required
        />
      </div>

      <div class="input-group">
        <label for="confirmNewPassword">Xác nhận mật khẩu mới</label>
        <input
          type="password"
          id="confirmNewPassword"
          bind:value={confirmNewPassword}
          placeholder="Nhập lại mật khẩu mới"
          disabled={isLoading}
          required
        />
      </div>

      {#if error}
        <div class="error" role="alert">{error}</div>
      {/if}

      <Button
        type="submit"
        label={isLoading ? "Đang xử lý..." : "Đặt lại mật khẩu"}
        variant="primary"
        disabled={isLoading}
      />

      <div class="switch-mode">
        Nhớ mật khẩu?
        <button type="button" class="link-btn" on:click={handleSwitchToLogin}>
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

  .form-description {
    color: #666;
    font-size: 14px;
    margin: -8px 0 8px 0;
    line-height: 1.5;
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

  /* Social Login */
  .social-divider {
    position: relative;
    text-align: center;
    margin: 20px 0;
  }

  .social-divider::before,
  .social-divider::after {
    content: "";
    position: absolute;
    top: 50%;
    width: 45%;
    height: 1px;
    background-color: #e0e0e0;
  }

  .social-divider::before {
    left: 0;
  }

  .social-divider::after {
    right: 0;
  }

  .social-divider span {
    background-color: white;
    padding: 0 12px;
    color: #7c7c7c;
    font-size: 14px;
  }
</style>
