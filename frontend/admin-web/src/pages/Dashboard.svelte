<script lang="ts">
  import { onMount } from "svelte";
  import { toastStore } from "../stores/toast-store";
  import { isAuthenticated } from "../stores/auth-store";

  import StatsCard from "../components/StatsCard.svelte";

  import {
    getPlatformOverview,
    getUserStats,
    getContentStats,
  } from "../services/stats-service";

  import type {
    PlatformOverview,
    UserStats,
    ContentStats,
  } from "../dtos/stats-dto";

  let overview = $state<PlatformOverview | null>(null);
  let userStats = $state<UserStats | null>(null);
  let contentStats = $state<ContentStats | null>(null);

  let loading = $state(false);

  async function loadOverview() {
    // Only load if authenticated
    let authenticated = false;
    const unsubscribe = isAuthenticated.subscribe((val) => {
      authenticated = val;
    });
    unsubscribe();

    if (!authenticated) {
      return;
    }

    loading = true;
    try {
      overview = await getPlatformOverview();
      userStats = await getUserStats({ period: "week" });
      contentStats = await getContentStats({ period: "week" });
    } catch (error) {
      console.error("Failed to load overview:", error);
      toastStore.error("Không thể tải thống kê");
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    const timer = setTimeout(() => {
      loadOverview();
    }, 100);

    return () => clearTimeout(timer);
  });
</script>

<div class="dashboard">
  <div class="content">
    <div class="overview">
      {#if loading}
        <div class="loading">Đang tải...</div>
      {:else if overview}
        <section class="stats-section">
          <h2>Thống kê người dùng</h2>
          <div class="stats-grid">
            <StatsCard
              title="Tổng người dùng"
              value={overview.total_users}
              icon="👥"
              color="#4a70a9"
              subtitle={userStats
                ? `+${userStats.new_users} người dùng mới`
                : ""}
            />
            <StatsCard
              title="Người dùng hoạt động"
              value={overview.active_users}
              icon="✅"
              color="#28a745"
            />
            <StatsCard
              title="Đã bị cấm"
              value={overview.banned_users}
              icon="🚫"
              color="#dc3545"
            />
          </div>
        </section>

        <section class="stats-section">
          <h2>Thống kê cộng đồng</h2>
          <div class="stats-grid">
            <StatsCard
              title="Tổng cộng đồng"
              value={overview.total_communities}
              icon="🏘️"
              color="#4a70a9"
            />
            <StatsCard
              title="Cộng đồng hoạt động"
              value={overview.active_communities}
              icon="✅"
              color="#28a745"
            />
            <StatsCard
              title="Đã bị cấm"
              value={overview.banned_communities}
              icon="🚫"
              color="#dc3545"
            />
          </div>
        </section>

        <section class="stats-section">
          <h2>Thống kê nội dung</h2>
          <div class="stats-grid">
            <StatsCard
              title="Tổng bài viết"
              value={overview.total_posts}
              icon="📝"
              color="#4a70a9"
              subtitle={contentStats
                ? `+${contentStats.new_posts} bài viết mới`
                : ""}
            />
            <StatsCard
              title="Tổng bình luận"
              value={overview.total_comments}
              icon="💬"
              color="#17a2b8"
              subtitle={contentStats
                ? `+${contentStats.new_comments} bình luận mới`
                : ""}
            />
            <StatsCard
              title="Báo cáo chờ xử lý"
              value={overview.pending_reports}
              icon="⚠️"
              color="#ffc107"
            />
          </div>
        </section>
      {/if}
    </div>
  </div>
</div>

<style>
  .dashboard {
    min-height: 100%;
  }

  .content {
    width: 100%;
  }

  .overview {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .stats-section h2 {
    font-size: 1.25rem;
    margin-bottom: 1rem;
    color: #333;
    font-weight: 600;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
  }

  .loading {
    padding: 3rem;
    text-align: center;
    color: #6c757d;
    font-size: 1rem;
  }
</style>
