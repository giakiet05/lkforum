<script lang="ts">
  import Post from "../components/Post.svelte";
  import type { PostResponse } from "../dtos/post-dto";
  import { getPosts } from "../services/post-service";
  import { onMount } from "svelte";

  type SortType = "best" | "hot" | "new" | "top" | "rising";

  let posts = $state<PostResponse[]>([]);
  let loading = $state(false);
  let loadingMore = $state(false);
  let error = $state<string | null>(null);
  let sortBy = $state<SortType>("best");
  let currentPage = $state(1);
  let hasMore = $state(true);
  let sentinelElement: HTMLDivElement | null = null;

  async function loadPosts(reset = true) {
    console.log(`🔄 [All/loadPosts] reset=${reset}, currentPage=${currentPage}`);
    if (reset) {
      loading = true;
      currentPage = 1;
      posts = [];
      hasMore = true;
    } else {
      loadingMore = true;
    }

    error = null;

    try {
      console.log(`📡 [All/API Call] page=${currentPage}, limit=20`);
      const newPosts = await getPosts({
        feed_type: "all",
        sort: sortBy || undefined,
        page: currentPage,
        limit: 20,
      });

      console.log(`✅ [All/API Response] received ${newPosts.length} posts`);

      if (reset) {
        posts = newPosts;
      } else {
        posts = [...posts, ...newPosts];
      }

      hasMore = newPosts.length === 20;
      console.log(`📊 [All/State] Total posts: ${posts.length}, hasMore: ${hasMore}`);
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load posts";
      console.error("❌ [All/Error] loading posts:", err);
    } finally {
      loading = false;
      loadingMore = false;
    }
  }

  async function loadMorePosts() {
    console.log(`🔽 [All/loadMorePosts] called - loadingMore=${loadingMore}, hasMore=${hasMore}`);
    if (loadingMore || !hasMore) {
      console.log(`⛔ [All/loadMorePosts] blocked`);
      return;
    }
    currentPage++;
    console.log(`➕ [All/loadMorePosts] incrementing page to ${currentPage}`);
    await loadPosts(false);
  }

  // Reload posts when sortBy changes (not on initial)
  let previousSortBy: string | null = null;
  $effect(() => {
    if (previousSortBy !== null && sortBy !== previousSortBy) {
      loadPosts(true);
    }
    previousSortBy = sortBy;
  });

  // Intersection Observer for infinite scroll
  $effect(() => {
    if (!sentinelElement || loading) {
      console.log(`⏸️ [All/Observer] not initialized - sentinelElement=${!!sentinelElement}, loading=${loading}`);
      return;
    }

    console.log(`👁️ [All/Observer] initializing...`);

    const observer = new IntersectionObserver(
      (entries) => {
        console.log(`🎯 [All/Observer] callback triggered - isIntersecting=${entries[0].isIntersecting}, hasMore=${hasMore}, loadingMore=${loadingMore}`);
        if (entries[0].isIntersecting && hasMore && !loadingMore) {
          console.log(`✅ [All/Observer] calling loadMorePosts()`);
          loadMorePosts();
        }
      },
      { threshold: 0.1 },
    );

    observer.observe(sentinelElement);
    console.log(`✅ [All/Observer] observing sentinel element`);

    return () => {
      console.log(`🛑 [All/Observer] disconnecting`);
      observer.disconnect();
    };
  });

  onMount(() => {
    loadPosts(true);
  });
</script>

<div class="page-container">
  <h1 class="page-title">Tất cả</h1>

  <div class="sort-options">
    <select bind:value={sortBy}>
      <option value="" disabled selected hidden>Sắp xếp theo</option>
      <option value="best">Tốt nhất</option>
      <option value="hot">Nổi bật</option>
      <option value="new">Mới nhất</option>
      <option value="top">Top</option>
      <option value="rising">Đang lên</option>
    </select>
  </div>

  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Đang tải bài viết...</p>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
      <button onclick={() => loadPosts()}>Thử lại</button>
    </div>
  {:else if posts.length === 0}
    <div class="empty">
      <p>Không tìm thấy bài viết nào</p>
    </div>
  {:else}
    <div class="post-list">
      {#each posts as post}
        <Post {post} />
      {/each}
    </div>

    <!-- Sentinel for infinite scroll -->
    {#if hasMore}
      <div class="sentinel" bind:this={sentinelElement}>
        {#if loadingMore}
          <div class="loading-more">
            <div class="spinner"></div>
            <p>Đang tải thêm...</p>
          </div>
        {/if}
      </div>
    {:else if posts.length > 0}
      <div class="end-message">
        <p>Bạn đã xem hết bài viết</p>
      </div>
    {/if}
  {/if}
</div>

<style>
  .page-container {
    padding: 16px 24px;
  }

  @media (max-width: 768px) {
    .page-container {
      padding: 0;
    }

    .page-title {
      padding: 0 12px;
    }
  }

  .page-title {
    font-size: 24px;
    font-weight: 700;
    color: var(--blue--);
    margin-bottom: 20px;
  }

  .sort-options {
    margin-bottom: 16px;
    display: inline-block;
  }

  .sort-options select {
    padding: 8px 32px 8px 12px;
    border: none;
    border-radius: 4px;
    background-color: transparent;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23153060' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 8px center;
    background-size: 16px;
    color: var(--blue--);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    transition: all 0.2s ease;
  }

  .sort-options select:hover {
    background-color: rgba(21, 48, 96, 0.08);
  }

  .sort-options select:focus {
    outline: none;
    background-color: rgba(21, 48, 96, 0.12);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    border-radius: 8px;
  }

  .loading,
  .error,
  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    text-align: center;
  }

  .loading p,
  .error p,
  .empty p {
    color: #5a5a5a;
    font-size: 16px;
    margin: 16px 0;
  }

  .spinner {
    border: 3px solid #f3f3f3;
    border-top: 3px solid var(--blue--);
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .error button {
    padding: 10px 20px;
    background-color: var(--blue--);
    color: white;
    border: none;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }

  .error button:hover {
    background-color: #0d2849;
  }

  .post-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .sentinel {
    min-height: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .loading-more {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 24px;
    text-align: center;
  }

  .loading-more p {
    color: #5a5a5a;
    font-size: 14px;
    margin: 12px 0 0 0;
  }

  .end-message {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32px 24px;
    text-align: center;
  }

  .end-message p {
    color: #878a8c;
    font-size: 14px;
    font-weight: 500;
  }
</style>
