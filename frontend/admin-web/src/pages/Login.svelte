<script lang="ts">
  import { onMount } from "svelte";
  import { login } from "../services/auth-service";
  import { setAuthenticated, isAuthenticated } from "../stores/auth-store";
  import { replace } from "svelte-spa-router";

  let username = $state("");
  let password = $state("");
  let error = $state("");
  let loading = $state(false);

  // Redirect if already authenticated - only check once on mount
  onMount(() => {
    if ($isAuthenticated) {
      replace("/dashboard");
    }
  });

  async function handleLogin() {
    if (!username || !password) {
      error = "Please enter username and password";
      return;
    }

    loading = true;
    error = "";

    try {
      await login({ identifier: username, password });
      setAuthenticated(true);
      replace("/dashboard");
    } catch (err: any) {
      error = err.message || "Login failed";
    } finally {
      loading = false;
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <h1>Admin Login</h1>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <form
      onsubmit={(e) => {
        e.preventDefault();
        handleLogin();
      }}
    >
      <div class="form-group">
        <label for="username">Username</label>
        <input
          id="username"
          type="text"
          bind:value={username}
          placeholder="Enter username"
          disabled={loading}
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          placeholder="Enter password"
          disabled={loading}
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? "Logging in..." : "Login"}
      </button>
    </form>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f5f5f5;
  }

  .login-card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }

  h1 {
    margin: 0 0 1.5rem;
    text-align: center;
    color: #333;
  }

  .error {
    background: #fee;
    color: #c33;
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-size: 14px;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #555;
  }

  input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
    box-sizing: border-box;
  }

  input:focus {
    outline: none;
    border-color: #4caf50;
  }

  button {
    width: 100%;
    padding: 0.75rem;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    margin-top: 0.5rem;
  }

  button:hover:not(:disabled) {
    background: #45a049;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
