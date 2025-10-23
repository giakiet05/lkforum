<script lang="ts">
  import { push } from "svelte-spa-router";
  import { mockQueuePosts } from "../mocks/mod-queue.mock";
  import type { QueuePost } from "../mocks/mod-queue.mock";

  export interface Props {
    params?: { name?: string };
  }

  let { params }: Props = $props();

  const communityName = params?.name || "";

  type SidebarItem = "queue" | "restricted" | "members" | "rules";
  type QueueTab = "unmoderated" | "edited" | "removed" | "reported";
  type SortOption = "newest" | "oldest" | "most-reported";

  let activeSidebarItem = $state<SidebarItem>("queue");
  let activeQueueTab = $state<QueueTab>("unmoderated");
  let sortBy = $state<SortOption>("newest");

  // Rules form state
  let ruleName = $state("");
  let ruleDescription = $state("");
  let reportReason = $state("");

  const isRuleFormValid = $derived(
    ruleName.trim().length > 0 && ruleDescription.trim().length > 0
  );

  const filteredPosts = $derived(() => {
    let posts = mockQueuePosts.filter((p) => p.queueType === activeQueueTab);

    // Sort posts
    if (sortBy === "newest") {
      posts.sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      );
    } else if (sortBy === "oldest") {
      posts.sort(
        (a, b) =>
          new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      );
    } else if (sortBy === "most-reported") {
      posts.sort((a, b) => (b.reportCount || 0) - (a.reportCount || 0));
    }

    return posts;
  });

  function handleExitModTools() {
    push(`/lk/${communityName}`);
  }

  function handleSidebarClick(item: SidebarItem) {
    activeSidebarItem = item;
  }

  function formatTime(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffHours / 24);

    if (diffHours < 1) return "just now";
    if (diffHours < 24) return `${diffHours} hours ago`;
    if (diffDays < 7) return `${diffDays} days ago`;
    return date.toLocaleDateString();
  }

  function handleApprove(postId: string) {
    console.log("Approve post:", postId);
    // TODO: Implement approve functionality
  }

  function handleRemove(postId: string) {
    console.log("Remove post:", postId);
    // TODO: Implement remove functionality
  }

  function handleSaveRule() {
    if (!isRuleFormValid) return;
    console.log("Save rule:", { ruleName, ruleDescription, reportReason });
    alert("Rule saved!");
    // Reset form
    ruleName = "";
    ruleDescription = "";
    reportReason = "";
  }
</script>

<div class="mod-tools-page">
  <!-- Sidebar -->
  <aside class="mod-sidebar">
    <button class="exit-mod-tools" onclick={handleExitModTools}>
      <!-- Use /arrowback_icon.svg (place file into public folder as arrowback_icon.svg) -->
      <img src="/arrowback_icon.svg" alt="" width="20" height="20" />
      Exit mod tools
    </button>

    <nav class="mod-nav">
      <button
        class="nav-item"
        class:active={activeSidebarItem === "queue"}
        onclick={() => handleSidebarClick("queue")}
      >
        <img src="/queue_icon.svg" alt="" width="20" height="20" />
        <span>Queue</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "restricted"}
        onclick={() => handleSidebarClick("restricted")}
      >
        <img src="/restricted_icon.svg" alt="" width="20" height="20" />
        <span>Restricted Users</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "members"}
        onclick={() => handleSidebarClick("members")}
      >
        <img src="/member_icon.svg" alt="" width="20" height="20" />
        <span>Mods & Members</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "rules"}
        onclick={() => handleSidebarClick("rules")}
      >
        <img src="/rule_icon.svg" alt="" width="20" height="20" />
        <span>Rules</span>
      </button>
    </nav>
  </aside>

  <!-- Main Content -->
  <main class="mod-content">
    {#if activeSidebarItem === "queue"}
      <div class="queue-section">
        <div class="queue-header">
          <h1>Queue</h1>
          <div class="sort-options">
            <select bind:value={sortBy}>
              <option value="newest">Newest</option>
              <option value="oldest">Oldest</option>
              <option value="most-reported">Most Reported</option>
            </select>
          </div>
        </div>

        <!-- Queue Tabs -->
        <div class="queue-tabs">
          <button
            class="tab-btn"
            class:active={activeQueueTab === "unmoderated"}
            onclick={() => (activeQueueTab = "unmoderated")}
          >
            Unmoderated
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "edited"}
            onclick={() => (activeQueueTab = "edited")}
          >
            Edited
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "removed"}
            onclick={() => (activeQueueTab = "removed")}
          >
            Removed
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "reported"}
            onclick={() => (activeQueueTab = "reported")}
          >
            Reported
          </button>
        </div>

        <!-- Posts List -->
        <div class="posts-list">
          {#each filteredPosts() as post (post.id)}
            <div class="post-card">
              <div class="post-header">
                <div class="post-author">
                  {#if post.authorAvatar}
                    <img
                      src={post.authorAvatar}
                      alt={post.author}
                      class="author-avatar"
                    />
                  {:else}
                    <div class="author-avatar-placeholder">
                      {post.author[0].toUpperCase()}
                    </div>
                  {/if}
                  <div class="author-info">
                    <span class="author-name">u/{post.author}</span>
                    <span class="post-time">{formatTime(post.createdAt)}</span>
                  </div>
                </div>
                {#if post.reportCount}
                  <span class="report-badge">{post.reportCount} reports</span>
                {/if}
              </div>

              <h3 class="post-title">{post.title}</h3>
              <p class="post-content">{post.content}</p>

              {#if post.reportReason}
                <div class="report-info">
                  <strong>Report reason:</strong>
                  {post.reportReason}
                </div>
              {/if}

              {#if post.removedReason}
                <div class="removed-info">
                  <strong>Removed by:</strong>
                  {post.removedBy} - {post.removedReason}
                </div>
              {/if}

              <div class="post-actions">
                <button
                  class="action-btn approve"
                  onclick={() => handleApprove(post.id)}
                >
                  Approve
                </button>
                <button
                  class="action-btn remove"
                  onclick={() => handleRemove(post.id)}
                >
                  Remove
                </button>
              </div>
            </div>
          {:else}
            <div class="empty-state">
              <p>No posts in this queue</p>
            </div>
          {/each}
        </div>
      </div>
    {:else if activeSidebarItem === "restricted"}
      <div class="placeholder-section">
        <h1>Restricted Users</h1>
        <p>Coming soon...</p>
      </div>
    {:else if activeSidebarItem === "members"}
      <div class="placeholder-section">
        <h1>Mods & Members</h1>
        <p>Coming soon...</p>
      </div>
    {:else if activeSidebarItem === "rules"}
      <div class="rules-section">
        <div class="rules-header">
          <div class="rules-title">
            <h2>Name and describe your rule</h2>
            <p class="sub">
              Rules set the expectations for members and redditors visiting your
              community
            </p>
          </div>
        </div>

        <form class="rule-form">
          <div class="form-row">
            <label>Rule name<span class="required">*</span></label>
            <div class="pill-input">
              <input
                type="text"
                placeholder="Rule name"
                maxlength="100"
                bind:value={ruleName}
              />
              <span class="char-count">{100 - ruleName.length}</span>
            </div>
          </div>

          <div class="form-row">
            <label>Description<span class="required">*</span></label>
            <div class="pill-input textarea">
              <textarea
                placeholder="Description"
                maxlength="500"
                bind:value={ruleDescription}
              ></textarea>
              <span class="char-count">{500 - ruleDescription.length}</span>
            </div>
          </div>

          <h3 class="section-heading">Reporting</h3>
          <p class="small">
            Users or mods can select a report reason when reporting content
          </p>

          <div class="form-row">
            <label>Report reason</label>
            <div class="pill-input">
              <input
                type="text"
                placeholder="Report reason"
                maxlength="100"
                bind:value={reportReason}
              />
              <span class="char-count">{100 - reportReason.length}</span>
            </div>
            <p class="hint">
              By default, this is the same as your rule name. Max characters 100
            </p>
          </div>

          <div class="form-row save-row">
            <button
              class="save-btn"
              class:enabled={isRuleFormValid}
              disabled={!isRuleFormValid}
              onclick={handleSaveRule}
              type="button"
            >
              Save
            </button>
          </div>
        </form>
      </div>
    {/if}
  </main>
</div>

<style>
  .mod-tools-page {
    display: flex;
    min-height: 100vh;
    background: #f6f7f8;
  }

  /* Sidebar */
  .mod-sidebar {
    width: 240px;
    background: white;
    border-right: 1px solid #edeff1;
    padding: 20px 0;
    position: sticky;
    top: 0;
    height: 100vh;
    overflow-y: auto;
  }

  .exit-mod-tools {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    width: 100%;
    background: transparent;
    border: none;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    transition: background 0.2s;
  }

  .exit-mod-tools:hover {
    background: #f6f7f8;
  }

  .mod-nav {
    margin-top: 20px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 20px;
    width: 100%;
    background: transparent;
    border: none;
    border-left: 3px solid transparent;
    font-size: 14px;
    font-weight: 500;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
  }

  .nav-item:hover {
    background: #f6f7f8;
  }

  .nav-item.active {
    background: #f6f7f8;
    border-left-color: var(--blue--);
    color: var(--blue--);
    font-weight: 600;
  }

  .nav-item img {
    flex-shrink: 0;
  }

  /* Main Content */
  .mod-content {
    flex: 1;
    padding: 24px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
  }

  .queue-section {
    background: white;
    border-radius: 8px;
    overflow: hidden;
  }

  .queue-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 24px;
    border-bottom: 1px solid #edeff1;
  }

  .queue-header h1 {
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0;
  }

  /* Sort Options - Same style as Profile page */
  .sort-options {
    padding: 0 1rem;
    position: relative;
  }

  .sort-options::after {
    content: "";
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    width: 20px;
    height: 20px;
    background-image: url("/Sort.jpg");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    pointer-events: none;
    opacity: 0.8;
  }

  .sort-options select {
    padding: 0.5rem 2rem 0.5rem 2.75rem;
    border: none;
    border-radius: 4px;
    background-color: #f8f9fa;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 1em;
    color: #1a1a1b;
    font-size: 0.9rem;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    font-weight: 400;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    transition: all 0.2s ease;
  }

  .sort-options select:not(:focus) {
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .sort-options select:hover {
    background-color: #f0f1f2;
  }

  .sort-options select:focus {
    outline: none;
    background-color: #fff;
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.06);
  }

  .sort-options select option {
    padding: 0.75rem 1rem;
    background-color: white;
    color: #1a1a1b;
    font-size: 0.9rem;
    font-weight: 400;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .sort-options select option:hover {
    background-color: #f8f9fa;
    color: #00008b;
  }

  .sort-options select option:checked {
    background-color: #f0f1f2;
    font-weight: 500;
  }

  .sort-options select:focus {
    border-radius: 4px;
  }

  @media screen and (-webkit-min-device-pixel-ratio: 0) {
    .sort-options select {
      border-radius: 4px !important;
    }

    .sort-options select:focus {
      border: none;
    }

    .sort-options select option:checked {
      background: #f0f1f2 linear-gradient(0deg, #f0f1f2 0%, #f0f1f2 100%);
      font-weight: 500;
    }

    .sort-options select option:hover {
      background: #e8f0fe linear-gradient(0deg, #e8f0fe 0%, #e8f0fe 100%);
      color: #00008b;
    }
  }

  .sort-select {
    padding: 8px 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
    background: white;
    cursor: pointer;
  }

  /* Tabs */
  .queue-tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid #edeff1;
    padding: 0 24px;
    background: white;
  }

  .tab-btn {
    padding: 12px 16px;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 14px;
    font-weight: 600;
    color: #7c7c7c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tab-btn.active {
    color: #1c1c1c;
    border-bottom-color: var(--blue--);
  }

  .tab-btn:hover {
    color: #1c1c1c;
  }

  /* Posts List */
  .posts-list {
    padding: 24px;
  }

  .post-card {
    background: #f6f7f8;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
  }

  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .post-author {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .author-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }

  .author-avatar-placeholder {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--blue--);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
  }

  .author-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .author-name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .post-time {
    font-size: 12px;
    color: var(--grayfont);
  }

  .report-badge {
    padding: 4px 8px;
    background: #ff4444;
    color: white;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
  }

  .post-title {
    font-size: 16px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0 0 8px 0;
  }

  .post-content {
    font-size: 14px;
    color: #1c1c1c;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }

  .report-info,
  .removed-info {
    padding: 12px;
    background: #fff3cd;
    border-left: 3px solid #ffc107;
    border-radius: 4px;
    font-size: 13px;
    margin-bottom: 12px;
  }

  .removed-info {
    background: #f8d7da;
    border-left-color: #dc3545;
  }

  .post-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .action-btn {
    padding: 8px 16px;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 13px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn:hover {
    background: #f6f7f8;
  }

  .action-btn.approve {
    background: var(--blue--);
    color: white;
    border: 1px solid var(--blue--);
    border-radius: 16px;
  }

  .action-btn.approve:hover {
    background: #0000cd;
    border-color: #0000cd;
  }

  .action-btn.remove {
    background: var(--button-secondary-bg);
    color: #1c1c1c;
    border: 1px solid transparent;
    border-radius: 16px;
  }

  .action-btn.remove:hover {
    background: rgba(214, 216, 222, 0.6);
  }

  .empty-state {
    text-align: center;
    padding: 48px 24px;
    color: var(--grayfont);
  }

  .placeholder-section {
    background: white;
    border-radius: 8px;
    padding: 48px 24px;
    text-align: center;
  }

  /* Rules section styles */
  .rules-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .rules-header {
    margin-bottom: 24px;
  }

  .rules-title h2 {
    margin: 0;
    font-size: 22px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .rules-title .sub {
    margin: 6px 0 0 0;
    color: var(--grayfont);
    font-size: 14px;
  }

  .rule-form {
    margin-top: 12px;
  }

  .form-row {
    margin-bottom: 18px;
  }

  .form-row label {
    display: block;
    font-size: 14px;
    color: #1c1c1c;
    margin-bottom: 8px;
    font-weight: 600;
  }

  .required {
    color: #d9534f;
    margin-left: 6px;
  }

  .pill-input {
    display: flex;
    align-items: center;
    background: #eef2f4;
    padding: 12px 16px;
    border-radius: 16px;
    position: relative;
  }

  .pill-input input {
    border: none;
    background: transparent;
    width: 100%;
    font-size: 14px;
    color: #1c1c1c;
    outline: none;
  }

  .pill-input.textarea {
    align-items: flex-start;
    padding-bottom: 36px;
  }

  .pill-input.textarea textarea {
    width: 100%;
    border: none;
    background: transparent;
    min-height: 120px;
    resize: vertical;
    font-size: 14px;
    color: #1c1c1c;
    outline: none;
    font-family: inherit;
  }

  .char-count {
    position: absolute;
    right: 16px;
    bottom: 12px;
    font-size: 12px;
    color: var(--grayfont);
  }

  .section-heading {
    margin-top: 16px;
    margin-bottom: 8px;
    font-size: 16px;
    font-weight: 700;
  }

  .small {
    color: var(--grayfont);
    font-size: 13px;
  }

  .hint {
    font-size: 12px;
    color: var(--grayfont);
    margin-top: 6px;
  }

  .save-row {
    display: flex;
    justify-content: flex-end;
    margin-top: 24px;
  }

  .save-btn {
    padding: 10px 24px;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    border: 1px solid transparent;
  }

  .save-btn:disabled {
    background: var(--button-secondary-bg);
    color: #1c1c1c;
    cursor: not-allowed;
  }

  .save-btn.enabled {
    background: var(--blue--);
    color: white;
    border-color: var(--blue--);
    cursor: pointer;
  }

  .save-btn.enabled:hover {
    background: #0000cd;
    border-color: #0000cd;
  }

  .placeholder-section h1 {
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0 0 16px 0;
  }

  .placeholder-section p {
    font-size: 16px;
    color: var(--grayfont);
    margin: 0;
  }
</style>
