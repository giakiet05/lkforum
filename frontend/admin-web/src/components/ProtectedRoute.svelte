<script lang="ts">
  import { onMount } from "svelte";
  import { isAuthenticated } from "../stores/auth-store";
  import { push } from "svelte-spa-router";

  let { children } = $props();
  let canShow = $state(false);

  onMount(() => {
    const unsubscribe = isAuthenticated.subscribe((auth) => {
      if (!auth) {
        push("/login");
        canShow = false;
      } else {
        canShow = true;
      }
    });
    return unsubscribe;
  });
</script>

{#if canShow}
  {@render children()}
{/if}
