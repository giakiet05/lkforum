<script lang="ts">
  interface Props {
    title: string;
    value: number | string;
    icon?: string;
    color?: string;
    subtitle?: string;
  }

  let { title, value, icon, color = "#4a70a9", subtitle }: Props = $props();
</script>

<div class="stats-card" style="border-left: 4px solid {color};">
  <div class="stats-header">
    {#if icon}
      <span class="stats-icon">{icon}</span>
    {/if}
    <h3>{title}</h3>
  </div>
  <div class="stats-value">{value}</div>
  {#if subtitle}
    <div class="stats-subtitle">{subtitle}</div>
  {/if}
</div>

<style>
  .stats-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition:
      transform 0.2s,
      box-shadow 0.2s;
  }

  .stats-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }

  .stats-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .stats-icon {
    font-size: 1.5rem;
  }

  .stats-header h3 {
    font-size: 0.875rem;
    font-weight: 500;
    color: #666;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .stats-value {
    font-size: 2rem;
    font-weight: 700;
    color: #333;
    margin-bottom: 0.5rem;
  }

  .stats-subtitle {
    font-size: 0.875rem;
    color: #888;
  }
</style>
