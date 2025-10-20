<script lang="ts">
  import type { SvelteComponent } from "svelte";
  import MainLayout from "./MainLayout.svelte";

  export let component: typeof SvelteComponent;
  export let layout: typeof MainLayout | null = null;
</script>

{#if layout}
  <svelte:component this={layout}>
    <svelte:component this={component} />
  </svelte:component>
{:else}
  <svelte:component this={component} />
{/if}
