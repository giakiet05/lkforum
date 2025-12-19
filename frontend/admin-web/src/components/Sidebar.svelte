<script lang="ts">
  import { push, location } from "svelte-spa-router";

  type SidebarItem = {
    id: string;
    label: string;
    to?: string;
    icon?: string;
  };

  type SidebarProps = {
    items: SidebarItem[];
    compact?: boolean;
    activeRoute?: string;
    onNavigate?: (item: SidebarItem) => void;
  };

  let {
    items,
    compact = $bindable(false),
    activeRoute = "",
    onNavigate,
  }: SidebarProps = $props();

  function handleNavigate(item: SidebarItem) {
    if (item.to) {
      push(item.to);
    }
    onNavigate?.(item);
  }

  function handleToggleCompact() {
    compact = !compact;
  }

  function isActive(item: SidebarItem): boolean {
    const loc = $location || "";
    const derived = loc.startsWith("#") ? loc.slice(1) : loc;
    const current = activeRoute || derived || "";
    return item.to === current;
  }
</script>

<aside class="sidebar" class:compact>
  <nav class="sidebar-nav">
    <ul class="nav-list">
      {#each items as item (item.id)}
        <li class="nav-item">
          <button
            class="nav-button"
            class:active={isActive(item)}
            onclick={() => handleNavigate(item)}
          >
            {#if item.icon}
              <span class="nav-icon" aria-hidden="true">{@html item.icon}</span>
            {/if}
            {#if !compact}
              <span class="nav-label">{item.label}</span>
            {/if}
          </button>
        </li>
      {/each}
    </ul>
  </nav>

  <div class="sidebar-footer">
    <button
      class="toggle-button"
      onclick={handleToggleCompact}
      title={compact ? "Mở rộng" : "Thu gọn"}
    >
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
        {#if compact}
          <path
            d="M7 4L13 10L7 16"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        {:else}
          <path
            d="M13 4L7 10L13 16"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        {/if}
      </svg>
    </button>
  </div>
</aside>

<style>
  .sidebar {
    position: fixed;
    top: var(--topbar-height);
    left: 0;
    height: calc(100vh - var(--topbar-height));
    width: var(--sidebar-width);
    background: var(--sidebar-background);
    border-right: 1px solid var(--sidebar-border);
    display: flex;
    flex-direction: column;
    transition: width 0.2s ease;
    z-index: 100;
    box-sizing: border-box;
    overflow-x: hidden;
  }

  .sidebar.compact {
    width: var(--sidebar-compact-width);
  }

  .sidebar-nav {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 12px 8px;
  }

  .nav-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .nav-item {
    list-style: none;
  }

  .nav-button {
    all: unset;
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 10px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.15s ease;
    color: var(--sidebar-foreground);
    font-size: 14px;
    font-weight: 500;
  }

  .nav-button:hover {
    background: var(--sidebar-accent);
  }

  .nav-button.active {
    background: var(--sidebar-accent);
    color: var(--primary-color);
    font-weight: 700;
  }

  .nav-icon {
    flex-shrink: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .nav-icon :global(svg) {
    width: 20px;
    height: 20px;
  }

  .compact .nav-button {
    justify-content: center;
    padding: 12px;
  }

  .nav-label {
    flex: 1;
    text-align: left;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .sidebar-footer {
    padding: 12px 8px;
    border-top: 1px solid var(--sidebar-border);
  }

  .toggle-button {
    all: unset;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 8px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.15s ease;
    color: var(--sidebar-foreground);
  }

  .toggle-button:hover {
    background: var(--sidebar-accent);
  }
</style>
