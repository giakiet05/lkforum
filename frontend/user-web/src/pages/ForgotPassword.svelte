<script lang="ts">
  import { push } from "svelte-spa-router";
  import Button from "../components/Button.svelte";
  import {
    forgotPassword,
    verifyResetOTP,
    resetPassword,
  } from "../services/auth-service";

  let step = $state<"email" | "otp" | "newPassword">("email");
  let email = $state("");
  let otp = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");
  let resetToken = $state("");
  let showPassword = $state(false);
  let showConfirmPassword = $state(false);
  let isLoading = $state(false);
  let error = $state("");
  let successMessage = $state("");

  // Step 1: Gửi OTP đến email
  async function handleSendOTP() {
    if (!email.trim()) {
      error = "Vui lòng nhập email";
      return;
    }

    try {
      isLoading = true;
      error = "";
      await forgotPassword(email);
      successMessage = "Mã OTP đã được gửi đến email của bạn";
      step = "otp";
    } catch (err: any) {
      error = err.message || "Không thể gửi OTP. Vui lòng thử lại";
    } finally {
      isLoading = false;
    }
  }

  // Step 2: Xác thực OTP
  async function handleVerifyOTP() {
    if (!otp.trim()) {
      error = "Vui lòng nhập mã OTP";
      return;
    }

    try {
      isLoading = true;
      error = "";
      const response = await verifyResetOTP(email, otp);
      resetToken = response.reset_token;
      successMessage = "Xác thực thành công! Vui lòng nhập mật khẩu mới";
      step = "newPassword";
    } catch (err: any) {
      error = err.message || "Mã OTP không đúng";
    } finally {
      isLoading = false;
    }
  }

  // Step 3: Đổi mật khẩu
  async function handleResetPassword() {
    if (!newPassword.trim()) {
      error = "Vui lòng nhập mật khẩu mới";
      return;
    }
    if (newPassword.length < 8) {
      error = "Mật khẩu phải có ít nhất 8 ký tự";
      return;
    }

    // Validate password strength
    const passwordRegex =
      /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
    if (!passwordRegex.test(newPassword)) {
      error =
        "Mật khẩu phải chứa ít nhất 1 chữ hoa, 1 chữ thường, 1 số và 1 ký tự đặc biệt (@$!%*?&)";
      return;
    }

    if (newPassword !== confirmPassword) {
      error = "Mật khẩu xác nhận không khớp";
      return;
    }

    try {
      isLoading = true;
      error = "";
      await resetPassword(resetToken, newPassword);
      successMessage =
        "Đổi mật khẩu thành công! Đang chuyển đến trang đăng nhập...";
      setTimeout(() => {
        push("/login");
      }, 2000);
    } catch (err: any) {
      error = err.message || "Không thể đổi mật khẩu";
    } finally {
      isLoading = false;
    }
  }

  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }

  function toggleConfirmPasswordVisibility() {
    showConfirmPassword = !showConfirmPassword;
  }

  function backToLogin() {
    push("/login");
  }
</script>

<div class="forgot-password-container">
  <div class="forgot-password-box">
    <h1>Quên mật khẩu</h1>

    {#if step === "email"}
      <p class="description">
        Nhập email của bạn và chúng tôi sẽ gửi mã OTP để đặt lại mật khẩu
      </p>

      <form
        onsubmit={(e) => {
          e.preventDefault();
          handleSendOTP();
        }}
      >
        <div class="input-group">
          <label for="email">Email</label>
          <input
            id="email"
            type="email"
            bind:value={email}
            placeholder="Nhập email của bạn"
            disabled={isLoading}
          />
        </div>

        {#if error}
          <div class="error" role="alert">{error}</div>
        {/if}

        {#if successMessage}
          <div class="success" role="alert">{successMessage}</div>
        {/if}

        <Button
          type="submit"
          label={isLoading ? "Đang gửi..." : "Gửi mã OTP"}
          variant="primary"
          disabled={isLoading}
        />
      </form>
    {/if}

    {#if step === "otp"}
      <p class="description">
        Nhập mã OTP đã được gửi đến <strong>{email}</strong>
      </p>

      <form
        onsubmit={(e) => {
          e.preventDefault();
          handleVerifyOTP();
        }}
      >
        <div class="input-group">
          <label for="otp">Mã OTP</label>
          <input
            id="otp"
            type="text"
            bind:value={otp}
            placeholder="Nhập mã 6 số"
            maxlength="6"
            disabled={isLoading}
          />
        </div>

        {#if error}
          <div class="error" role="alert">{error}</div>
        {/if}

        {#if successMessage}
          <div class="success" role="alert">{successMessage}</div>
        {/if}

        <Button
          type="submit"
          label={isLoading ? "Đang xác thực..." : "Xác thực OTP"}
          variant="primary"
          disabled={isLoading}
        />

        <button
          type="button"
          class="text-button"
          onclick={handleSendOTP}
          disabled={isLoading}
        >
          Gửi lại mã OTP
        </button>
      </form>
    {/if}

    {#if step === "newPassword"}
      <p class="description">Nhập mật khẩu mới của bạn</p>

      <form
        onsubmit={(e) => {
          e.preventDefault();
          handleResetPassword();
        }}
      >
        <div class="input-group password-group">
          <label for="newPassword">Mật khẩu mới</label>
          <input
            id="newPassword"
            type={showPassword ? "text" : "password"}
            bind:value={newPassword}
            placeholder="Nhập mật khẩu mới (tối thiểu 8 ký tự)"
            disabled={isLoading}
          />
          <button
            type="button"
            class="password-toggle-icon"
            onclick={togglePasswordVisibility}
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

        <div class="input-group password-group">
          <label for="confirmPassword">Xác nhận mật khẩu</label>
          <input
            id="confirmPassword"
            type={showConfirmPassword ? "text" : "password"}
            bind:value={confirmPassword}
            placeholder="Nhập lại mật khẩu mới"
            disabled={isLoading}
          />
          <button
            type="button"
            class="password-toggle-icon"
            onclick={toggleConfirmPasswordVisibility}
            aria-label={showConfirmPassword ? "Ẩn mật khẩu" : "Hiện mật khẩu"}
          >
            {#if showConfirmPassword}
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

        {#if successMessage}
          <div class="success" role="alert">{successMessage}</div>
        {/if}

        <Button
          type="submit"
          label={isLoading ? "Đang xử lý..." : "Đổi mật khẩu"}
          variant="primary"
          disabled={isLoading}
        />
      </form>
    {/if}

    <div class="back-to-login">
      <button type="button" class="text-button" onclick={backToLogin}>
        ← Quay lại đăng nhập
      </button>
    </div>
  </div>
</div>

<style>
  .forgot-password-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: #dae0e6;
    padding: 2rem;
  }

  .forgot-password-box {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    max-width: 400px;
    width: 100%;
  }

  h1 {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: #1a1a1b;
  }

  .description {
    color: #7c7c7c;
    font-size: 14px;
    margin-bottom: 1.5rem;
    line-height: 1.5;
  }

  .description strong {
    color: #1a1a1b;
  }

  form {
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

  .input-group input:disabled {
    background: #f6f7f8;
    cursor: not-allowed;
  }

  .password-group {
    position: relative;
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
    margin-top: 10px;
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

  .error {
    color: #ff4500;
    font-size: 12px;
    padding: 8px 12px;
    background: rgba(255, 69, 0, 0.1);
    border-radius: 4px;
  }

  .success {
    color: #46d160;
    font-size: 12px;
    padding: 8px 12px;
    background: rgba(70, 209, 96, 0.1);
    border-radius: 4px;
  }

  .text-button {
    background: none;
    border: none;
    padding: 0.5rem 0;
    color: var(--primary-color);
    font-weight: 600;
    cursor: pointer;
    font-size: 14px;
    text-align: center;
    width: 100%;
  }

  .text-button:hover {
    text-decoration: underline;
  }

  .text-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .back-to-login {
    margin-top: 1.5rem;
    text-align: center;
  }

  .back-to-login .text-button {
    color: #7c7c7c;
  }
</style>
