<script lang="ts">
  import { authStore } from "../stores/auth-store";
  import { logout } from "../services/auth-service";
  import { push } from "svelte-spa-router";
  import { get } from "svelte/store";
  import type { User } from "../models/user";

  let user: User | null = null;
  $: {
    const state = get(authStore);
    user = state.user;
  }

  function handleLogout() {
    logout();
    push("/login");
  }
</script>

<main>
  <h1>Home Page</h1>
  {#if user}
    <p><strong>Username:</strong> {user.username}</p>
    <p><strong>Email:</strong> {user.email}</p>
    <button on:click={handleLogout}>Logout</button>
  {:else}
    <p>You are not logged in.</p>
  {/if}
</main>
