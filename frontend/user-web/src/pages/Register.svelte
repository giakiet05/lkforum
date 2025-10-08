<script lang="ts">
  import Button from "../components/Button.svelte";
  import { push } from "svelte-spa-router";

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

  // Các biến cho form đăng ký
  let email = "";
  let username = "";
  let password = "";
  let confirmPassword = "";
  let showPassword = false;

  // UI state
  let loading = false;
  let error: string | null = null;

  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }

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

      // Nếu backend trả token -> lưu và điều hướng về trang chính
      if (data.access_token) {
        localStorage.setItem("access_token", data.access_token);
      }
      if (data.refresh_token) {
        localStorage.setItem("refresh_token", data.refresh_token);
      }
      if (data.user) {
        localStorage.setItem("user", JSON.stringify(data.user));
      }

      if (data.access_token) {
        push("/");
      } else {
        alert(data.message || "Đăng ký thành công. Vui lòng đăng nhập.");
        push("/login");
      }
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

      <div class="signin-link">
        Đã có tài khoản? <a href="/">Đăng nhập</a>
      </div>
    </div>
  </div>

  <div class="decorative-section"></div>
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
    width: 500px;
    height: auto;
    border-radius: 12px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    object-fit: cover;
  }
  .form-section {
    /* Đổi tên từ login-form-section */
    position: relative;
    flex: 0 0 50%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    background-color: white;
    padding: 2rem;
  }
  .form-wrapper {
    width: 100%;
    max-width: 450px;
    padding-right: 5rem;
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
    position: absolute;
    top: 2rem;
    left: 2rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    text-decoration: none;
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
