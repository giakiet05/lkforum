<script lang="ts">
  import { push } from "svelte-spa-router";

  type Community = {
    id: string;
    name: string;
    description: string;
    icon?: string;
    isFavorite: boolean;
    isJoined: boolean;
  };

  // Mock data - sau này sẽ fetch từ API
  let allCommunities = $state<Community[]>([
    {
      id: "1",
      name: "3amjokes",
      description:
        "lk/3amjokes - for all the stupid humor of sleep deprivation. Have you been up for longer than a normal human being can operate? Good. Have you just laughed at a joke that wouldn't be funny otherwise? submit your insomniac dad jokes today",
      icon: "🌙",
      isFavorite: true,
      isJoined: true,
    },
    {
      id: "2",
      name: "Animesuggest",
      description:
        "Subreddit for anime and manga fans which allows suggestions and requests for anything related to anime and manga subculture.",
      icon: "💭",
      isFavorite: true,
      isJoined: true,
    },
    {
      id: "3",
      name: "acne",
      description: "A subreddit for discussing acne and how to best treat it.",
      icon: "🩺",
      isFavorite: false,
      isJoined: true,
    },
    {
      id: "4",
      name: "AdviceAnimals",
      description: "Reddit's Gold Mine",
      icon: "🦄",
      isFavorite: false,
      isJoined: true,
    },
    {
      id: "5",
      name: "AmItheAsshole",
      description:
        "A catharsis for the frustrated moral philosopher in all of us, and a place to finally find out if you were wrong in an argument that's been bothering you. Tell us about any non-violent conflict you have experienced; give us both sides of the story, and find out if you're right, or you're the asshole. See our ~~*Best Of*~~ \"Most Controversial\" at /r/AITAFiltered!",
      icon: "🤔",
      isFavorite: false,
      isJoined: true,
    },
    {
      id: "6",
      name: "announcements",
      description: "Official announcements from Reddit, Inc.",
      icon: "📢",
      isFavorite: false,
      isJoined: true,
    },
    {
      id: "7",
      name: "apple",
      description:
        "An unofficial report about Apple and all of its devices and software.",
      icon: "🍎",
      isFavorite: false,
      isJoined: true,
    },
  ]);

  let searchQuery = $state("");
  let activeTab = $state<"all" | "favorites">("all");

  let filteredCommunities = $derived(() => {
    let filtered = allCommunities;

    // Filter by tab
    if (activeTab === "favorites") {
      filtered = filtered.filter((c) => c.isFavorite);
    }

    // Filter by search
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (c) =>
          c.name.toLowerCase().includes(query) ||
          c.description.toLowerCase().includes(query)
      );
    }

    return filtered;
  });

  function toggleFavorite(communityId: string) {
    const report = allCommunities.find((c) => c.id === communityId);
    if (report) {
      report.isFavorite = !report.isFavorite;
    }
  }

  function leaveCommunity(communityId: string) {
    const report = allCommunities.find((c) => c.id === communityId);
    if (report && confirm(`Leave lk/${report.name}?`)) {
      // TODO: Call API to leave report
      allCommunities = allCommunities.filter((c) => c.id !== communityId);
    }
  }

  function navigateToCommunity(communityName: string) {
    push(`/lk/${communityName}`);
  }
</script>

<div class="manage-communities-page">
  <div class="page-header">
    <h1>Manage communities</h1>
  </div>

  <div class="page-content">
    <!-- Search Bar -->
    <div class="search-bar">
      <svg
        width="20"
        height="20"
        viewBox="0 0 20 20"
        fill="none"
        class="search-icon"
      >
        <path
          d="M9 17A8 8 0 1 0 9 1a8 8 0 0 0 0 16zM18 18l-4-4"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        />
      </svg>
      <input
        type="text"
        placeholder="Filter your communities"
        bind:value={searchQuery}
      />
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        class="tab"
        class:active={activeTab === "all"}
        onclick={() => (activeTab = "all")}
      >
        All Communities
      </button>
      <button
        class="tab"
        class:active={activeTab === "favorites"}
        onclick={() => (activeTab = "favorites")}
      >
        Favorites
      </button>
    </div>

    <!-- Communities List -->
    <div class="communities-list">
      {#each filteredCommunities() as report (report.id)}
        <div class="report-card">
          <div class="report-main">
            <button
              class="report-info"
              onclick={() => navigateToCommunity(report.name)}
            >
              <div class="report-avatar">{report.icon || "📁"}</div>
              <div class="report-details">
                <h3 class="report-name">lk/{report.name}</h3>
                <p class="report-description">{report.description}</p>
              </div>
            </button>

            <div class="report-actions">
              <button
                class="favorite-btn"
                class:active={report.isFavorite}
                onclick={() => toggleFavorite(report.id)}
                title={report.isFavorite
                  ? "Remove from favorites"
                  : "Add to favorites"}
              >
                {#if report.isFavorite}
                  <svg
                    width="20"
                    height="20"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      d="M10 2l2 6h6l-4.5 3.5 2 6L10 14l-5.5 3.5 2-6L2 8h6l2-6z"
                    />
                  </svg>
                {:else}
                  <svg
                    width="20"
                    height="20"
                    viewBox="0 0 20 20"
                    fill="none"
                    stroke="currentColor"
                  >
                    <path
                      d="M10 2l2 6h6l-4.5 3.5 2 6L10 14l-5.5 3.5 2-6L2 8h6l2-6z"
                      stroke-width="1.5"
                    />
                  </svg>
                {/if}
              </button>

              <button class="joined-btn">Joined</button>
            </div>
          </div>

          <!-- Leave button (appears on hover) -->
          <button
            class="leave-btn"
            onclick={() => leaveCommunity(report.id)}
          >
            Leave lk/{report.name}
          </button>
        </div>
      {/each}

      {#if filteredCommunities().length === 0}
        <div class="empty-state">
          <p>No communities found</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .manage-communities-page {
    max-width: 1000px;
    margin: 0 auto;
    padding: 20px;
  }

  .page-header {
    margin-bottom: 24px;
  }

  .page-header h1 {
    font-size: 28px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0;
  }

  .page-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  /* Search Bar */
  .search-bar {
    position: relative;
    width: 100%;
  }

  .search-icon {
    position: absolute;
    left: 16px;
    top: 50%;
    transform: translateY(-50%);
    color: #878a8c;
  }

  .search-bar input {
    width: 100%;
    padding: 12px 16px 12px 48px;
    border: 1px solid #ccc;
    border-radius: 24px;
    font-size: 14px;
    background: white;
    transition: border-color 0.2s;
  }

  .search-bar input:focus {
    outline: none;
    border-color: #0079d3;
  }

  .search-bar input::placeholder {
    color: #878a8c;
  }

  /* Tabs */
  .tabs {
    display: flex;
    gap: 0;
    background: #f6f7f8;
    border-radius: 8px;
    padding: 4px;
    width: fit-content;
  }

  .tab {
    padding: 10px 24px;
    background: transparent;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tab:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .tab.active {
    background: white;
    color: #1055c9;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  /* Communities List */
  .communities-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .report-card {
    background: white;
    border: 1px solid #ccc;
    border-radius: 8px;
    padding: 16px;
    transition: all 0.2s;
    position: relative;
  }

  .report-card:hover {
    border-color: #878a8c;
  }

  .report-card:hover .leave-btn {
    opacity: 1;
    pointer-events: auto;
  }

  .report-main {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  .report-info {
    flex: 1;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    text-align: left;
    min-width: 0;
  }

  .report-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
    flex-shrink: 0;
    background: #f6f7f8;
  }

  .report-details {
    flex: 1;
    min-width: 0;
  }

  .report-name {
    font-size: 16px;
    font-weight: 600;
    color: #1c1c1c;
    margin: 0 0 4px 0;
  }

  .report-description {
    font-size: 14px;
    color: #7c7c7c;
    margin: 0;
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .report-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
  }

  .favorite-btn {
    width: 36px;
    height: 36px;
    padding: 0;
    background: transparent;
    border: 1px solid #edeff1;
    border-radius: 50%;
    cursor: pointer;
    color: #d3d3d3;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .favorite-btn:hover {
    background: #f6f7f8;
    color: #ffd700;
  }

  .favorite-btn.active {
    color: #ffd700;
    border-color: #ffd700;
  }

  .joined-btn {
    padding: 8px 24px;
    background: transparent;
    border: 1px solid var(--blue--);
    border-radius: 24px;
    color: var(--blue--);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .joined-btn:hover {
    background: var(--blue--);
    color: white;
  }

  .leave-btn {
    position: absolute;
    bottom: 16px;
    left: 50%;
    transform: translateX(-50%);
    padding: 8px 16px;
    background: #1c1c1c;
    color: white;
    border: none;
    border-radius: 24px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    opacity: 0;
    pointer-events: none;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .leave-btn:hover {
    background: #000;
  }

  .empty-state {
    text-align: center;
    padding: 48px 20px;
    color: #878a8c;
  }

  .empty-state p {
    font-size: 16px;
    margin: 0;
  }
</style>
