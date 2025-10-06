<script lang="ts">
  import { login } from "../services/auth-service";
  import type { LoginRequest } from "../dtos/auth-dto";
  import { push } from "svelte-spa-router";

  let identifier = "";
  let password = "";
  let error = "";
  let success = "";

  async function handleLogin() {
    error = "";
    if (!identifier || !password) {
      error = "Please enter both email and password.";
      return;
    }
    const credentials: LoginRequest = { identifier, password };
    try {
      await login(credentials);
      success = "Login successful!";
      push("/");
    } catch (e) {
      error = e instanceof Error ? e.message : "Login failed";
      success = "";
    }
  }
</script>

<main>
  <h1>Login</h1>
  <form on:submit|preventDefault={handleLogin}>
    <div>
      <label for="identifier">Email:</label>
      <input id="identifier" type="email" bind:value={identifier} />
    </div>
    <div>
      <label for="password">Password:</label>
      <input id="password" type="password" bind:value={password} />
    </div>
    {#if error}
      <p style="color: red">{error}</p>
    {/if}
    {#if success}
      <p style="color: green">{success}</p>
    {/if}
    <button type="submit">Login</button>
  </form>
</main>
