<script lang="ts">
  import { push, location } from "svelte-spa-router";
  import CommunitiesList from "./CommunitiesList.svelte";

  type SidebarItem = {
    id: string;
    label: string;
    to?: string;
    icon?: string;
    children?: SidebarItem[];
  };

  type SidebarProps = {
    items: SidebarItem[];
    compact?: boolean;
    activeRoute?: string;
    onNavigate?: (item: SidebarItem) => void;
    onToggleCompact?: () => void;
    isOpen?: boolean;
    onClose?: () => void;
  };

  let {
    items,
    compact = $bindable(false),
    activeRoute = "",
    onNavigate,
    onToggleCompact,
    isOpen = false,
    onClose,
  }: SidebarProps = $props();

  let expandedGroups = $state(new Set<string>());
  function handleNavigate(item: SidebarItem) {
    if (item.to) {
      push(item.to);
    }
    onNavigate?.(item);
  }

  function toggleGroup(id: string) {
    if (expandedGroups.has(id)) {
      expandedGroups.delete(id);
    } else {
      expandedGroups.add(id);
    }
    // bắt buộc gán lại để kích hoạt reactivity
    expandedGroups = new Set(expandedGroups);
  }

  function handleToggleCompact() {
    console.log("Toggle compact clicked! Current:", compact);
    compact = !compact;
    console.log("After toggle:", compact);
    onToggleCompact?.();
  }

  function isActive(item: SidebarItem): boolean {
    const loc = $location || "";
    const derived = loc.startsWith("#") ? loc.slice(1) : loc;
    const current = activeRoute || derived || "";
    return item.to === current;
  }

  function handleOverlayClick() {
    onClose?.();
  }
</script>

<!-- Mobile overlay -->
{#if isOpen}
  <div
    class="sidebar-overlay"
    onclick={handleOverlayClick}
    role="presentation"
  ></div>
{/if}

<aside class="sidebar" class:compact class:open={isOpen}>
  <!-- header removed: brand/logo handled by topbar now -->

  <nav class="sidebar-nav">
    <ul class="nav-list">
      {#each items as item (item.id)}
        <li class="nav-item">
          {#if item.children && item.children.length > 0}
            <button
              class="nav-button group-button"
              onclick={() => toggleGroup(item.id)}
              aria-expanded={expandedGroups.has(item.id)}
            >
              {#if item.icon}
                <span class="nav-icon" aria-hidden="true"
                  >{@html item.icon}</span
                >
              {/if}
              {#if !compact}
                <span class="nav-label">{item.label}</span>
                <span
                  class="expand-icon"
                  class:expanded={expandedGroups.has(item.id)}
                >
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path
                      d="M4 6L8 10L12 6"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </span>
              {/if}
            </button>

            {#if expandedGroups.has(item.id) && !compact}
              <ul class="nav-sublist">
                {#each item.children as child (child.id)}
                  <li class="nav-subitem">
                    <button
                      class="nav-button sub-button"
                      class:active={isActive(child)}
                      onclick={() => handleNavigate(child)}
                    >
                      {#if child.icon}
                        <span class="nav-icon" aria-hidden="true"
                          >{@html child.icon}</span
                        >
                      {/if}
                      <span class="nav-label">{child.label}</span>
                    </button>
                  </li>
                {/each}
              </ul>
            {/if}
          {:else}
            <button
              class="nav-button"
              class:active={isActive(item)}
              onclick={() => handleNavigate(item)}
            >
              {#if item.icon}
                <span class="nav-icon" aria-hidden="true"
                  >{@html item.icon}</span
                >
              {/if}
              {#if !compact}
                <span class="nav-label">{item.label}</span>
              {/if}
            </button>
          {/if}
        </li>
      {/each}
    </ul>

    <!-- Communities Section -->
    <CommunitiesList {compact} />
  </nav>

  <div class="sidebar-footer">
    <button
      class="toggle-button"
      onclick={handleToggleCompact}
      title={compact ? "Expand" : "Collapse"}
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
    transition:
      width 0.2s ease,
      transform 0.2s ease;
    z-index: 100;
    box-sizing: border-box;
  }

  .sidebar.compact {
    width: var(--sidebar-compact-width);
  }

  /* give a little top padding so items don't touch top edge */
  .sidebar-nav {
    flex: 1;
    overflow-y: auto;
    padding: 12px 8px;
  }

  .nav-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .nav-item {
    width: 100%;
  }

  .nav-button {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    color: var(--sidebar-foreground);
    font-size: 14px;
    font-weight: 500;
    transition: all 0.15s ease;
    text-align: left;
  }

  .sidebar.compact .nav-button {
    justify-content: center;
    padding: 10px;
  }

  .nav-button:hover {
    background: rgba(0, 0, 0, 0.08);
  }

  .nav-button.active {
    background: var(--sidebar-active-background);
    color: #00008b;
    border-left: 3px solid var(--sidebar-active-border);
    font-weight: 600;
  }

  .nav-button.active .nav-icon {
    color: hsl(var(--sidebar-active-foreground));
  }

  .nav-button:active {
    transform: scale(0.98);
  }

  .nav-icon {
    width: 20px;
    height: 20px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: hsl(var(--sidebar-foreground) / 0.7);
  }

  .nav-button:hover .nav-icon {
    color: hsl(var(--sidebar-foreground));
  }

  .nav-label {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .expand-icon {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    transition: transform 0.2s ease;
    color: hsl(var(--sidebar-foreground) / 0.5);
  }

  .expand-icon.expanded {
    transform: rotate(0deg);
  }

  .expand-icon:not(.expanded) {
    transform: rotate(-90deg);
  }

  .nav-sublist {
    list-style: none;
    padding: 4px 0 4px 32px;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .nav-subitem {
    width: 100%;
  }

  .sub-button {
    font-weight: 400;
    font-size: 13px;
    padding: 8px 12px;
  }

  .sidebar-footer {
    padding: 12px;
    border-top: 2px solid hsl(var(--sidebar-border));
  }

  .toggle-button {
    width: 100%;
    padding: 10px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    color: hsl(var(--sidebar-foreground) / 0.6);
    transition: all 0.15s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .toggle-button:hover {
    background: hsl(var(--sidebar-accent));
    color: hsl(var(--sidebar-foreground));
  }

  .sidebar-overlay {
    position: fixed;
    top: var(--topbar-height);
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 189;
    display: none;
  }

  @media (max-width: 1024px) {
    .sidebar-overlay {
      display: block;
    }

    .sidebar {
      transform: translateX(-100%);
      top: var(--topbar-height);
      width: 280px;
      max-width: 80vw;
      transition: transform 0.3s ease;
      z-index: 190;
      box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
    }

    .sidebar.open {
      transform: translateX(0);
    }

    .sidebar.compact {
      transform: translateX(-100%);
    }
  }

  @media (max-width: 480px) {
    .sidebar {
      width: 260px;
    }
  }
</style>
