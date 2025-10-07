<script lang="ts">
  import Router from "svelte-spa-router";
  import routes from "./routes";
</script>

<main>
  <!-- <nav>
    <a href="/#/">Home</a> |
    <a href="/#/login">Login</a>
  </nav> -->
  <Router {routes} />
</main>

<style>
</style>
