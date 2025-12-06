<script lang="ts">
  import Router from "svelte-spa-router";
  import { routes } from "./routes";
  import { isAuthenticated } from "./stores/auth-store";
  import { push } from "svelte-spa-router";

  // Redirect to login if not authenticated
  $effect(() => {
    if (!$isAuthenticated && window.location.hash !== "#/login") {
      push("/login");
    } else if ($isAuthenticated && window.location.hash === "#/login") {
      push("/dashboard");
    }
  });
</script>

<Router {routes} />

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      "Helvetica Neue", Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
</style>
