<script lang="ts">
  import Button from "../components/Button.svelte";
  import { login } from "../services/auth-service";
  import { push } from "svelte-spa-router";

  // form fields
  let identifier = ""; // username hoặc email tuỳ backend
  let password = "";

  // UI state
  let loading = false;
  let error: string | null = null;
  let showPassword = false;

  // Validator đơn giản
  function validate() {
    if (!identifier || !password) {
      error = "Vui lòng nhập tên đăng nhập và mật khẩu";
      return false;
    }
    error = null;
    return true;
  }

  // Chuyển đổi hiển thị mật khẩu
  function toggleShowPasswordVisibility() {
    showPassword = !showPassword;
  }

  // Xử lý submit form
  async function handleLoginSubmit() {
    if (!validate()) return;
    loading = true;
    error = null;

    try {
      await login({ identifier, password }); // service đã lưu token và set store
      // redirect sau khi login thành công
      push("/");
    } catch (err: any) {
      console.error("Login error:", err);
      if (typeof err === "string") error = err;
      else if (err && (err.message || err.error))
        error = err.message || err.error;
      else error = "Lỗi khi đăng nhập. Vui lòng thử lại.";
    } finally {
      loading = false;
    }
  }

  function handleGoogleLogin() {
    // Google để sau
    alert("Google login sẽ làm sau");
  }
</script>

<div class="login-page">
  <div class="center-image-container">
    <img src="/discuss.jpg" alt="Brand Logo" class="center-image" />
  </div>

  <div class="login-form-section">
    <a href="/" class="brand-logo">
      <img src="/LKlogo.jpg" alt="LKForum Logo" />
      <span>LKForum</span>
    </a>
    <div class="form-wrapper">
      <h2 style="color:black;">Đăng nhập với LKForum</h2>
      <p>Đăng nhập để tiếp tục khám phá.</p>

      <form on:submit|preventDefault={handleLoginSubmit} class="login-form">
        <!-- input identifier -->
        <div class="input-group">
          <label for="identifier">Tên đăng nhập</label>
          <input
            id="identifier"
            type="text"
            bind:value={identifier}
            placeholder="Nhập tên đăng nhập của bạn"
          />
        </div>

        <!-- password-group -->
        <div class="input-group password-group">
          <label for="password">Mật khẩu</label>
          <input
            id="password"
            type={showPassword ? "text" : "password"}
            bind:value={password}
            placeholder="Mật khẩu"
          />
          <span
            class="password-toggle-icon"
            on:click={toggleShowPasswordVisibility}
            role="button"
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
          </span>
        </div>

        {#if error}
          <div class="error" role="alert">{error}</div>
        {/if}

        <Button
          type="submit"
          label={loading ? "Đang đăng nhập..." : "Đăng Nhập"}
          variant="primary"
          disabled={loading}
        />
      </form>

      <div class="separator">
        <span>HOẶC</span>
      </div>

      <Button
        label="Đăng nhập với Google"
        variant="google"
        on:click={handleGoogleLogin}
      />

      <div class="signup-link">
        Chưa có tài khoản? <a href="/register">Đăng ký ngay</a>
      </div>
    </div>
  </div>

  <div class="decorative-section"></div>
</div>

<style>
  .login-page {
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
    /* Không cần background, padding hay bo tròn cho khung nữa */
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

  .login-form-section {
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
    padding-right: 12%;
  }
  .form-wrapper h2 {
    font-family: var(--font-secondary);
    font-size: 2.5em;
    font-weight: 700;
    color: var(--text-color);
    margin-bottom: 0.5rem;
  }
  .form-wrapper {
    width: 100%;
    max-width: 450px;
    /* Bỏ padding-right */
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
    transition: border-color 0.3s ease; /* Thêm hiệu ứng chuyển động */
  }

  /* Style khi người dùng nhấn (focus) vào ô input */
  .input-group input:focus {
    outline: none; /* Bỏ viền outline mặc định */
    border-bottom-color: var(--primary-color); /* Đổi màu gạch chân */
  }
  .separator {
    display: flex;
    align-items: center;
    text-align: center;
    color: #aaa;
    margin: 2rem 0;
  }
  .separator::before,
  .separator::after {
    content: "";
    flex: 1;
    border-bottom: 1px solid var(--border-color);
  }
  .separator:not(:empty)::before {
    margin-right: 0.5em;
  }
  .separator:not(:empty)::after {
    margin-left: 0.5em;
  }
  .signup-link {
    text-align: center;
    margin-top: 2rem;
    color: #555;
  }
  .signup-link a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 600;
  }

  .decorative-section {
    flex: 0 0 50%; /* Chỉnh lại cho bố cục 50/50 */
    background-image: url("/background.png");
    background-size: cover;
    background-position: center;
  }

  @media (max-width: 900px) {
    .decorative-section {
      display: none;
    }
    .login-form-section {
      flex-basis: 100%;
    }
    .center-image-container {
      display: none;
    }
  }
  .password-group {
    position: relative;
  }

  .password-toggle-icon {
    position: absolute;
    top: 55%; /* Điều chỉnh để căn giữa với ô input */
    right: 10px;
    transform: translateY(-50%);
    cursor: pointer;
    color: #888;
  }
  /* CSS CHO LOGO */
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
    color: var(--darkblue--color);
    font-family: var(--font-secondary);
  }
</style>
