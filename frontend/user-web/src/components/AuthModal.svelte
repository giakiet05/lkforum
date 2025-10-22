<script lang="ts">
  import Modal from "./Modal.svelte";
  import Login from "../pages/Login.svelte";
  import { login } from "../services/auth-service";
  import type { LoginRequest } from "../dtos/auth-dto";
  import { push } from "svelte-spa-router";

  let { show = false, onClose }: { show: boolean; onClose: () => void } =
    $props();

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

  function handleSwitchToRegister() {
    // Đóng modal và chuyển sang trang Register
    onClose();
    push("/register");
  }
</script>

<Modal {show} {onClose} title="Đăng nhập vào LKForum">
  <Login
    mode="login"
    onSubmit={(data) => handleLogin(data as LoginRequest)}
    {isLoading}
    {error}
    on:switchMode={handleSwitchToRegister}
  />
</Modal>

<style>
  :global(.modal-container) {
    width: 400px !important;
    padding: 24px !important;
  }
</style>
