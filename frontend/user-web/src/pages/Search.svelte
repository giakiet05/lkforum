<script lang="ts">
  import { push } from "svelte-spa-router";
  import { getCommunities } from "../services/community-service";
  import { searchUsers } from "../services/user-service";
  import type { CommunityResponse } from "../dtos/community-dto";
  import type { UserResponse } from "../dtos/user-dto";

  let searchQuery = $state("");
  let communityResults = $state<CommunityResponse[]>([]);
  let userResults = $state<UserResponse[]>([]);
  let isSearching = $state(false);
  let searchTimeout: number | null = null;

  async function performSearch(query: string) {
    if (!query.trim()) {
      communityResults = [];
      userResults = [];
      return;
    }

    try {
      isSearching = true;
      const cleanQuery = query.replace(/^(c\/|lk\/|u\/)/i, "").trim();
      if (!cleanQuery) {
        communityResults = [];
        userResults = [];
        return;
      }

      const [communities, users] = await Promise.all([
        getCommunities({ search: cleanQuery, limit: 5 }),
        searchUsers(cleanQuery, 5),
      ]);

      communityResults = communities.data;
      userResults = users;
    } catch (err) {
      console.error("Search error:", err);
    } finally {
      isSearching = false;
    }
  }

  function handleSearchInput() {
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    searchTimeout = setTimeout(() => {
      performSearch(searchQuery);
    }, 300);
  }

  function handleCommunityClick(communityName: string) {
    push(`/lk/${communityName}`);
  }

  function handleUserClick(username: string) {
    push(`/profile/${username}`);
  }

  function handleBack() {
    window.history.back();
  }
</script>

<div class="search-page">
  <div class="search-header">
    <button class="back-button" onclick={handleBack} aria-label="Quay lại">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
        <path
          d="M15 18L9 12L15 6"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>

    <div class="search-input-wrapper">
      <svg
        class="search-icon"
        width="20"
        height="20"
        viewBox="0 0 20 20"
        fill="none"
      >
        <circle
          cx="8.5"
          cy="8.5"
          r="5.5"
          stroke="currentColor"
          stroke-width="2"
        />
        <path
          d="M12.5 12.5L17 17"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        />
      </svg>
      <input
        type="text"
        class="search-input"
        placeholder="Tìm cộng đồng và người dùng"
        bind:value={searchQuery}
        oninput={handleSearchInput}
        autofocus
      />
    </div>
  </div>

  <div class="search-results">
    {#if isSearching}
      <div class="loading">Đang tìm kiếm...</div>
    {:else if searchQuery.trim() && communityResults.length === 0 && userResults.length === 0}
      <div class="no-results">Không tìm thấy kết quả</div>
    {:else}
      {#if communityResults.length > 0}
        <div class="results-section">
          <h3 class="section-title">Cộng đồng</h3>
          {#each communityResults as community}
            <button
              class="result-item"
              onclick={() => handleCommunityClick(community.name)}
            >
              <img
                src={community.avatar || "/default-community.png"}
                alt=""
                class="result-avatar"
              />
              <div class="result-info">
                <div class="result-name">lk/{community.name}</div>
                <div class="result-meta">
                  {community.member_count} thành viên
                </div>
              </div>
            </button>
          {/each}
        </div>
      {/if}

      {#if userResults.length > 0}
        <div class="results-section">
          <h3 class="section-title">Người dùng</h3>
          {#each userResults as user}
            <button
              class="result-item"
              onclick={() => handleUserClick(user.username)}
            >
              {#if user.profile?.avatar?.url}
                <img
                  src={user.profile.avatar.url}
                  alt=""
                  class="result-avatar"
                />
              {:else}
                <div class="result-avatar fallback">
                  {user.username.charAt(0).toUpperCase()}
                </div>
              {/if}
              <div class="result-info">
                <div class="result-name">u/{user.username}</div>
                <div class="result-meta">{user.reputation || 0} karma</div>
              </div>
            </button>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .search-page {
    min-height: 100vh;
    background: var(--background);
  }

  .search-header {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: var(--topbar-background, #ffffff);
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .back-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border: none;
    background: transparent;
    color: var(--foreground);
    cursor: pointer;
    border-radius: 50%;
    transition: background 0.2s;
    flex-shrink: 0;
  }

  .back-button:hover {
    background: var(--hover-background, rgba(0, 0, 0, 0.05));
  }

  .search-input-wrapper {
    flex: 1;
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-icon {
    position: absolute;
    left: 12px;
    color: var(--muted-foreground);
    pointer-events: none;
  }

  .search-input {
    width: 100%;
    padding: 10px 16px 10px 40px;
    background: var(--input-background, #f6f7f8);
    border: 1px solid var(--border);
    border-radius: 20px;
    font-size: 16px;
    color: var(--foreground);
    outline: none;
    transition: all 0.2s;
  }

  .search-input:focus {
    background: var(--background);
    border-color: var(--primary);
  }

  .search-results {
    padding: 16px;
  }

  .loading,
  .no-results {
    text-align: center;
    padding: 32px;
    color: var(--muted-foreground);
  }

  .results-section {
    margin-bottom: 24px;
  }

  .section-title {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--muted-foreground);
    margin-bottom: 8px;
    padding: 0 4px;
  }

  .result-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;
    text-align: left;
  }

  .result-item:hover {
    background: var(--hover-background, rgba(0, 0, 0, 0.05));
  }

  .result-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
  }

  .result-avatar.fallback {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--primary);
    color: white;
    font-weight: 600;
  }

  .result-info {
    flex: 1;
    min-width: 0;
  }

  .result-name {
    font-weight: 600;
    color: var(--foreground);
    margin-bottom: 2px;
  }

  .result-meta {
    font-size: 12px;
    color: var(--muted-foreground);
  }
</style>
