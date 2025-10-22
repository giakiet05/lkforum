<script lang="ts">
  import Button from "../components/Button.svelte";
  import { push } from "svelte-spa-router";
  import { setAuth } from "../stores/auth-store";

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

  // Các biến cho form đăng ký
  let email = "";
  let username = "";
  let password = "";
  let confirmPassword = "";
  let showPassword = false;
  let otp = "";

  // UI state
  let loading = false;
  let error: string | null = null;
  let step: "register" | "verify" = "register"; // Step 1: register → Step 2: verify OTP

  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }

  // Step 1: Đăng ký - nhập email, username, password → backend gửi OTP
  async function handleRegisterSubmit() {
    // validate cơ bản
    if (!email || !username || !password || !confirmPassword) {
      error = "Vui lòng điền đầy đủ thông tin";
      return;
    }
    if (password !== confirmPassword) {
      error = "Mật khẩu xác nhận không khớp!";
      return;
    }

    loading = true;
    error = null;

    try {
      const res = await fetch(`${API_BASE_URL}/api/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, username, password }),
      });

      if (!res.ok) {
        let errObj: any = {};
        try {
          errObj = await res.json();
        } catch {
          try {
            const text = await res.text();
            errObj = { error: text || `HTTP ${res.status}` };
          } catch {
            errObj = { error: `HTTP ${res.status}` };
          }
        }
        throw errObj.error || errObj.message || "Lỗi đăng ký";
      }

      const data = await res.json();
      console.log("Register response:", data);

      // Chuyển sang step verify OTP
      step = "verify";
      error = null;
    } catch (err: any) {
      console.error("Register error:", err);
      if (typeof err === "string") error = err;
      else if (err && (err.message || err.error))
        error = err.message || err.error;
      else error = "Lỗi khi đăng ký. Vui lòng thử lại.";
    } finally {
      loading = false;
    }
  }

  // Step 2: Verify OTP
  async function handleVerifyOTP() {
    if (!otp || otp.length !== 6) {
      error = "Vui lòng nhập mã OTP 6 chữ số";
      return;
    }

    loading = true;
    error = null;

    try {
      const res = await fetch(`${API_BASE_URL}/api/auth/verify-email`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, otp }),
      });

      if (!res.ok) {
        let errObj: any = {};
        try {
          errObj = await res.json();
        } catch {
          try {
            const text = await res.text();
            errObj = { error: text || `HTTP ${res.status}` };
          } catch {
            errObj = { error: `HTTP ${res.status}` };
          }
        }
        throw errObj.error || errObj.message || "Mã OTP không đúng";
      }

      const data = await res.json();
      console.log("Verify OTP response:", data);

      // Lưu tokens và user vào localStorage
      if (data.access_token) {
        localStorage.setItem("access_token", data.access_token);
      }
      if (data.refresh_token) {
        localStorage.setItem("refresh_token", data.refresh_token);
      }
      if (data.user) {
        localStorage.setItem("user", JSON.stringify(data.user));
      }

      // Update authStore để UI hiển thị avatar ngay lập tức
      if (data.user && data.access_token) {
        setAuth(data.user, data.access_token);
      }

      // Redirect về trang chính
      alert("Đăng ký thành công! Chào mừng bạn đến với LKForum.");
      push("/");
    } catch (err: any) {
      console.error("Register error:", err);
      if (typeof err === "string") error = err;
      else if (err && (err.message || err.error))
        error = err.message || err.error;
      else error = "Lỗi khi đăng ký. Vui lòng thử lại.";
    } finally {
      loading = false;
    }
  }

  async function handleResendOTP() {
    loading = true;
    error = null;

    try {
      const res = await fetch(
        `${API_BASE_URL}/api/auth/resend-verification-email`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email }),
        }
      );

      if (!res.ok) {
        throw new Error("Không thể gửi lại mã OTP");
      }

      alert("Mã OTP mới đã được gửi đến email của bạn!");
      error = null;
    } catch (err: any) {
      console.error("Resend OTP error:", err);
      error = "Không thể gửi lại mã OTP. Vui lòng thử lại sau.";
    } finally {
      loading = false;
    }
  }

  function handleBackToRegister() {
    step = "register";
    otp = "";
    error = null;
  }
</script>

<div class="register-page">
  <div class="center-image-container">
    <img src="/discuss.jpg" alt="Brand Logo" class="center-image" />
  </div>

  <div class="form-section">
    <a href="/" class="brand-logo">
      <img src="/LKlogo.jpg" alt="LKForum Logo" />
      <span>LKForum</span>
    </a>
    <div class="form-wrapper">
      {#if step === "register"}
        <h2 style="color:black;">Tạo tài khoản mới</h2>
        <p>Tham gia cộng đồng LKForum ngay hôm nay.</p>

        <form
          on:submit|preventDefault={handleRegisterSubmit}
          class="register-form"
        >
          <div class="input-group">
            <label for="email">Email</label>
            <input
              type="email"
              id="email"
              bind:value={email}
              placeholder="Nhập email của bạn"
            />
          </div>

          <div class="input-group">
            <label for="username">Tên đăng nhập</label>
            <input
              type="text"
              id="username"
              bind:value={username}
              placeholder="Chọn một tên đăng nhập"
            />
          </div>

          <div class="input-group password-group">
            <label for="password">Mật khẩu</label>
            <input
              type={showPassword ? "text" : "password"}
              id="password"
              bind:value={password}
              placeholder="Tạo mật khẩu"
            />
            <span
              class="password-toggle-icon"
              on:click={togglePasswordVisibility}
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
                  ><path
                    d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"
                  /><circle cx="12" cy="12" r="3" /></svg
                >
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
                  ><path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" /><path
                    d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
                  /><path
                    d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"
                  /><line x1="2" x2="22" y1="2" y2="22" /></svg
                >
              {/if}
            </span>
          </div>

          <div class="input-group">
            <label for="confirmPassword">Xác nhận mật khẩu</label>
            <input
              type="password"
              id="confirmPassword"
              bind:value={confirmPassword}
              placeholder="Nhập lại mật khẩu"
            />
          </div>

          {#if error}
            <div class="error" role="alert">{error}</div>
          {/if}

          <Button
            type="submit"
            label={loading ? "Đang đăng ký..." : "Đăng Ký"}
            variant="primary"
            disabled={loading}
          />
        </form>

        <p class="login-link">
          Đã có tài khoản? <a href="/#/login">Đăng nhập</a>
        </p>
      {:else}
        <!-- Step 2: Verify OTP -->
        <h2 style="color:black;">Xác thực Email</h2>
        <p>Chúng tôi đã gửi mã OTP đến <strong>{email}</strong></p>
        <p class="otp-hint">
          Vui lòng kiểm tra hộp thư và nhập mã gồm 6 chữ số
        </p>

        <form on:submit|preventDefault={handleVerifyOTP} class="verify-form">
          <div class="input-group">
            <label for="otp">Mã OTP</label>
            <input
              type="text"
              id="otp"
              bind:value={otp}
              placeholder="Nhập 6 chữ số"
              maxlength="6"
              class="otp-input"
            />
          </div>

          {#if error}
            <div class="error" role="alert">{error}</div>
          {/if}

          <Button
            type="submit"
            label={loading ? "Đang xác thực..." : "Xác Nhận"}
            variant="primary"
            disabled={loading}
          />
        </form>

        <div class="otp-actions">
          <button
            type="button"
            class="link-btn"
            on:click={handleResendOTP}
            disabled={loading}
          >
            Gửi lại mã OTP
          </button>
          <button
            type="button"
            class="link-btn"
            on:click={handleBackToRegister}
            disabled={loading}
          >
            Quay lại
          </button>
        </div>

        <p class="login-link">
          Đã có tài khoản? <a href="/#/login">Đăng nhập</a>
        </p>
      {/if}
    </div>
  </div>
</div>

<style>
  /* Gần như toàn bộ style được sao chép từ trang Login để đồng nhất */
  /* Đổi tên class để tránh xung đột nếu cần, nhưng ở đây chúng ta giữ nguyên */
  .register-page {
    display: flex;
    width: 100vw;
    height: 100vh;
    font-family: var(--font-primary);
    position: relative;
    overflow: hidden;
  }
  .center-image-container {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    z-index: 10;
  }
  .center-image {
    display: block;
    /* Sửa lại: Dùng vw (viewport width) để ảnh co dãn theo màn hình */
    width: 25vw; /* Chiếm khoảng 25% chiều rộng màn hình */
    max-width: 450px; /* Nhưng không bao giờ to quá 450px */
    min-width: 250px; /* Và không bao giờ nhỏ hơn 250px */

    height: auto;
    border-radius: 12px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    object-fit: cover;
  }
  .form-section {
    /* Đổi tên từ login-form-section */
    flex: 0 0 50%;
    display: flex;
    flex-direction: column; /* Xếp dọc */
    justify-content: center; /* Căn giữa theo chiều dọc */
    align-items: flex-start; /* Căn trái */
    background-color: white;
    padding: 2rem 4rem; /* Tăng padding để đẹp hơn */
    box-sizing: border-box; /* Thêm vào để padding không làm vỡ layout */
  }
  .form-wrapper {
    width: 100%;
    max-width: 450px;
    padding-right: 12%; /* Sửa lại: Dùng % để nó co dãn theo */
  }
  .form-wrapper h2 {
    font-family: var(--font-secondary);
    font-size: 2.5em;
    font-weight: 700;
    color: var(--text-color);
    margin-bottom: 0.5rem;
  }
  .form-wrapper p {
    color: #666;
    margin-bottom: 2.5rem;
  }
  .input-group {
    margin-bottom: 1.5rem;
  }
  .input-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
  }
  .input-group input {
    width: 100%;
    padding: 1rem 0.2rem;
    border: none;
    border-radius: 0;
    border-bottom: 2px solid var(--border-color);
    font-size: 1em;
    box-sizing: border-box;
    background-color: transparent;
    transition: border-color 0.3s ease;
  }
  .input-group input:focus {
    outline: none;
    border-bottom-color: var(--primary-color);
  }
  .signin-link {
    /* Đổi tên từ signup-link */
    text-align: center;
    margin-top: 2rem;
    color: #555;
  }
  .signin-link a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 600;
  }

  /* OTP Step styles */
  .otp-hint {
    font-size: 0.9em;
    color: #888;
    margin-bottom: 1.5rem;
    margin-top: -1rem;
  }

  .otp-input {
    text-align: center;
    font-size: 1.5em;
    letter-spacing: 0.5em;
    font-weight: 600;
  }

  .otp-actions {
    display: flex;
    justify-content: space-between;
    margin-top: 1.5rem;
  }

  .link-btn {
    background: none;
    border: none;
    color: var(--primary-color);
    font-weight: 600;
    cursor: pointer;
    text-decoration: underline;
    font-size: 0.95em;
  }

  .link-btn:hover {
    color: var(--primary-color-hover);
  }

  .link-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .decorative-section {
    flex: 0 0 50%;
    background-image: url("/background.png");
    background-size: cover;
    background-position: center;
  }
  .password-group {
    position: relative;
  }
  .password-toggle-icon {
    position: absolute;
    top: 55%;
    right: 10px;
    transform: translateY(-50%);
    cursor: pointer;
    color: #888;
  }
  .brand-logo {
    /* Bỏ position: absolute */
    display: flex;
    align-items: center;
    gap: 0.75rem;
    text-decoration: none;
    margin-bottom: 5rem; /* Tạo khoảng cách với form */
  }
  .brand-logo img {
    width: 80px;
    height: 80px;
  }
  .brand-logo span {
    font-size: 1.5em;
    font-weight: 700;
    color: #213547;
    font-family: var(--font-secondary);
  }

  @media (max-width: 900px) {
    /* Ẩn tấm ảnh ở giữa và phần nền trang trí */
    .center-image-container,
    .decorative-section {
      display: none;
    }

    /* Cho form chiếm toàn bộ chiều rộng màn hình */
    .form-section {
      flex: 1; /* Hoặc flex: 0 0 100%; */
      justify-content: center;
      padding: 2rem;
    }

    /* Điều chỉnh lại padding cho form để cân đối hơn */
    .form-wrapper {
      padding-right: 0;
      max-width: 100%;
    }

    /* Căn giữa lại link đăng nhập */
    .signin-link {
      text-align: center;
    }
  }
</style>
