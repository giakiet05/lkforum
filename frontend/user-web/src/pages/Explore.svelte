<script lang="ts">
  import Post from "../components/Post.svelte";
  import type { PostResponse } from "../dtos/post-dto";
  import { getPosts } from "../services/post-service";
  import { onMount } from "svelte";

  type SortType = "best" | "hot" | "new" | "top" | "rising";

  let posts = $state<PostResponse[]>([]);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let sortBy = $state<SortType | "">("");

  async function loadPosts() {
    loading = true;
    error = null;
    try {
      posts = await getPosts({
        feed_type: "explore",
        sort: sortBy || undefined,
        limit: 20,
      });
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load posts";
      console.error("Error loading posts:", err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (sortBy !== "") {
      loadPosts();
    }
  });

  onMount(() => {
    loadPosts();
  });
</script>

<div class="page-container">
  <h1 class="page-title">Khám phá</h1>

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
  {/if}
</div>

<style>
  .page-container {
    padding: 16px 24px;
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
</style>
