<script lang="ts">
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import Topbar from "./components/Topbar.svelte";
  import TestToggle from "./components/TestToggle.svelte";
</script>

<main>
  <!-- <nav>
    <a href="/#/">Home</a> |
    <a href="/#/login">Login</a> |
    <a href="/#/register">register</a>
  </nav> -->
  <Router {routes} />
</main>

<style>
</style>
