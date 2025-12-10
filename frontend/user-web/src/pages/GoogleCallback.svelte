<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { setAuth } from "../stores/auth-store";
  import {
    ACCESS_TOKEN_KEY,
    REFRESH_TOKEN_KEY,
    USER_KEY,
  } from "../constants/auth-constants";

  let status = $state<"loading" | "success" | "error">("loading");
  let errorMessage = $state("");

  /**
   * Decode JWT token để lấy user info
   * JWT format: header.payload.signature
   */
  function decodeJWT(token: string): any {
    try {
      const base64Url = token.split(".")[1];
      const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split("")
          .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
          .join("")
      );
      return JSON.parse(jsonPayload);
    } catch (err) {
      console.error("Error decoding JWT:", err);
      return null;
    }
  }

  onMount(async () => {
    // Đọc tokens từ query params trong hash
    // URL format: /#/auth/callback?access_token=...&refresh_token=...
    // window.location.hash = "#/auth/callback?access_token=...&refresh_token=..."
    const hash = window.location.hash; // "#/auth/callback?access_token=..."

    // Tìm vị trí dấu ?
    const queryIndex = hash.indexOf("?");
    if (queryIndex === -1) {
      status = "error";
      errorMessage = "Thiếu thông tin xác thực";
      return;
    }

    // Lấy phần query string sau dấu ?
    const queryString = hash.substring(queryIndex + 1); // "access_token=...&refresh_token=..."
    const urlParams = new URLSearchParams(queryString);
    const accessToken = urlParams.get("access_token");
    const refreshToken = urlParams.get("refresh_token");

    console.log("Full URL:", window.location.href);
    console.log("Hash:", hash);
    console.log("Query string:", queryString);
    console.log("Access token:", accessToken ? "Found" : "Not found");
    console.log("Refresh token:", refreshToken ? "Found" : "Not found");

    if (accessToken && refreshToken) {
      try {
        // Decode access token để lấy user ID
        const payload = decodeJWT(accessToken);
        if (!payload || !payload.sub) {
          throw new Error("Invalid token payload");
        }

        // Lưu tokens vào localStorage trước
        localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
        localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);

        // Fetch user info từ backend
        const response = await fetch("http://localhost:8080/api/users/me", {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch user info");
        }

        const result = await response.json();
        console.log("API response from /api/users/me:", result);

        // Extract user from response (API returns {data: user})
        const user = result.data || result;
        console.log("User data to save:", user);

        // Lưu user vào localStorage
        localStorage.setItem(USER_KEY, JSON.stringify(user));

        // Update auth store
        setAuth(user);

        status = "success";

        // Redirect về home sau 1 giây
        setTimeout(() => {
          push("/");
        }, 1000);
      } catch (err) {
        console.error("Error processing tokens:", err);
        status = "error";
        errorMessage = "Không thể xử lý thông tin đăng nhập";
      }
    } else {
      status = "error";
      errorMessage = "Thiếu thông tin xác thực";
    }
  });
</script>

<div class="callback-container">
  {#if status === "loading"}
    <div class="spinner"></div>
    <h2>Đang xử lý đăng nhập...</h2>
  {:else if status === "success"}
    <div class="success-icon">✓</div>
    <h2>Đăng nhập thành công!</h2>
    <p>Đang chuyển hướng...</p>
  {:else}
    <div class="error-icon">✗</div>
    <h2>Đăng nhập thất bại</h2>
    <p>{errorMessage}</p>
    <button class="retry-btn" onclick={() => push("/")}>
      Quay về trang chủ
    </button>
  {/if}
</div>

<style>
  .callback-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 20px;
    text-align: center;
  }

  .spinner {
    width: 48px;
    height: 48px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid var(--blue--);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 24px;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .success-icon {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background-color: #4caf50;
    color: white;
    font-size: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 24px;
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
    margin-bottom: 24px;
  }

  h2 {
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 12px 0;
  }

  p {
    font-size: 16px;
    color: #7c7c7c;
    margin: 0 0 24px 0;
  }

  .retry-btn {
    padding: 10px 24px;
    background-color: var(--blue--);
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
