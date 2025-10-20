<script lang="ts">
  import Modal from "./Modal.svelte";
  import Login from "../pages/Login.svelte";
  import Register from "../pages/Register.svelte";
  import { register, login } from "../services/auth-service";
  import type { LoginRequest } from "../dtos/auth-dto";
  import type { RegisterDto } from "../dtos/auth-dto";

  let { show = false, onClose }: { show: boolean; onClose: () => void } =
    $props();

  let activeTab = $state<"login" | "register">("login");
  let isLoading = $state(false);
  let error = $state("");

  const handleLogin = async (data: LoginRequest) => {
    try {
      isLoading = true;
      error = "";
      await login(data);
      onClose();
    } catch (err) {
      error = err instanceof Error ? err.message : "Login failed";
    } finally {
      isLoading = false;
    }
  };

  const handleRegister = async (data: RegisterDto) => {
    try {
      isLoading = true;
      error = "";
      await register(data);
      activeTab = "login"; // Switch to login after successful registration
    } catch (err) {
      error = err instanceof Error ? err.message : "Registration failed";
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
    error = "";
  }
</script>

<Modal
  {show}
  {onClose}
  title={activeTab === "login" ? "Đăng nhập vào LKForum" : "Tạo tài khoản mới"}
>
  {#if activeTab === "login"}
    <Login
      mode="login"
      onSubmit={(data) => handleLogin(data as LoginRequest)}
      {isLoading}
      {error}
      on:switchMode={handleSwitchToRegister}
    />
  {:else}
    <Login
      mode="register"
      onSubmit={(data) => handleRegister(data as RegisterDto)}
      {isLoading}
      {error}
      on:switchMode={handleSwitchToLogin}
    />
  {/if}
</Modal>

<style>
  :global(.modal-container) {
    width: 400px !important;
    padding: 24px !important;
  }
</style>
