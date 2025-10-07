<script lang="ts">
  // Đường dẫn import đã chính xác
  import Button from "../components/Button.svelte";

  let username = "";
  let password = "";

  let showPassword = false; // Mặc định là ẩn mật khẩu

  function togglePasswordVisibility() {
    showPassword = !showPassword; // Đảo ngược trạng thái true/false
  }
  function handleLoginSubmit() {
    console.log("Login submitted with:", { username, password });
    alert(`Đăng nhập với:\nUsername: ${username}\nPassword: ${password}`);
  }

  function handleGoogleLogin() {
    console.log("Attempting Google Login...");
    alert("Bắt đầu đăng nhập với Google!");
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
        <div class="input-group">
          <label for="username">Tên đăng nhập</label>
          <input
            type="text"
            id="username"
            bind:value={username}
            placeholder="Nhập tên đăng nhập của bạn"
          />
        </div>

        <div class="input-group password-group">
          <label for="password">Mật khẩu</label>
          <input
            type={showPassword ? "text" : "password"}
            id="password"
            bind:value={password}
            placeholder="Mật khẩu"
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

        <Button type="submit" label="Đăng Nhập" variant="primary" />
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
    width: 500px; /* Tăng kích thước ảnh cho ấn tượng hơn */
    height: auto; /* Giữ đúng tỷ lệ ảnh */
    border-radius: 12px; /* Bo góc nhẹ cho đẹp */
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2); /* Đổ bóng trực tiếp cho ảnh */
    object-fit: cover;
  }

  .login-form-section {
    flex: 0 0 50%; /* Chỉnh lại cho bố cục 50/50 */
    display: flex;
    position: relative;
    justify-content: flex-start;
    align-items: center;
    background-color: var(white);
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
    background-image: url("https://images.unsplash.com/photo-1557683316-973673baf926?q=80&w=2029&auto=format&fit=crop");
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
    color: var(--darkblue--color);
    font-family: var(--font-secondary);
  }
</style>
