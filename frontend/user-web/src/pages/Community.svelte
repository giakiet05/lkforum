<script lang="ts">
  import { onMount } from "svelte";
  import Post from "../components/Post.svelte";
  import type { PostData } from "../types/post";

  type CommunityProps = {
    params?: { name: string };
  };

  let { params = { name: "sveltejs" } }: CommunityProps = $props();

  let activeSort = $state("hot");
  let isJoined = $state(false);

  // Mock community data
  const community = {
    name: params.name,
    displayName: `lk/${params.name}`,
    description: "A community for all things related to " + params.name,
    members: 125000,
    online: 3200,
    createdAt: "Jan 1, 2020",
    banner: "/banner_sample1.jpg",
    icon: "/LKlogo.jpg",
  };

  // Mock posts for this community
  const posts: PostData[] = [
    {
      id: "1",
      type: "text",
      community: params.name,
      author: "user123",
      time: "4 hours ago",
      title: "Welcome to " + params.name + "!",
      upvotes: 123,
      downvotes: 5,
      commentsCount: 42,
      content:
        "This is the community page. Join us to see more content and participate in discussions!",
    },
    {
      id: "2",
      type: "image",
      community: params.name,
      author: "photographer",
      time: "8 hours ago",
      title: "Check out this amazing photo!",
      upvotes: 456,
      downvotes: 12,
      commentsCount: 89,
      images: ["/GirlFromNowhere.jpg"],
    },
  ];

  function toggleJoin() {
    isJoined = !isJoined;
  }

  onMount(() => {
    window.scrollTo(0, 0);
  });
</script>

<div class="community-page">
  <!-- Banner -->
  <div class="community-banner">
    <img src={community.banner} alt="Community banner" />
  </div>

  <!-- Community Header -->
  <div class="community-header">
    <div class="community-header-content">
      <div class="community-info">
        <img src={community.icon} alt="Community icon" class="community-icon" />
        <div class="community-title">
          <h1>{community.displayName}</h1>
          <p class="community-name">lk/{community.name}</p>
        </div>
      </div>
      <button class="join-btn" class:joined={isJoined} onclick={toggleJoin}>
        {isJoined ? "Joined" : "Join"}
      </button>
    </div>
  </div>

  <div class="community-container">
    <!-- Main Content -->
    <div class="community-main">
      <!-- Sort Bar -->
      <div class="sorting-bar">
        <button
          class="sort-btn"
          class:active={activeSort === "hot"}
          onclick={() => (activeSort = "hot")}
        >
          Hot
        </button>
        <button
          class="sort-btn"
          class:active={activeSort === "new"}
          onclick={() => (activeSort = "new")}
        >
          New
        </button>
        <button
          class="sort-btn"
          class:active={activeSort === "top"}
          onclick={() => (activeSort = "top")}
        >
          Top
        </button>
      </div>

      <!-- Posts -->
      <div class="post-list">
        {#each posts as post}
          <Post {post} />
        {/each}
      </div>
    </div>

    <!-- Sidebar -->
    <div class="community-sidebar">
      <div class="about-card">
        <h3>About Community</h3>
        <p class="about-description">{community.description}</p>

        <div class="community-stats">
          <div class="stat">
            <div class="stat-value">{community.members.toLocaleString()}</div>
            <div class="stat-label">Members</div>
          </div>
          <div class="stat">
            <div class="stat-value">{community.online.toLocaleString()}</div>
            <div class="stat-label">Online</div>
          </div>
        </div>

        <div class="community-created">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
            <path
              d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
            />
          </svg>
          <span>Created {community.createdAt}</span>
        </div>

        <button class="create-post-btn">Create Post</button>
      </div>

      <!-- Rules Card -->
      <div class="rules-card">
        <h3>Community Rules</h3>
        <ol class="rules-list">
          <li>Be respectful and civil</li>
          <li>No spam or self-promotion</li>
          <li>Stay on topic</li>
          <li>No personal attacks</li>
          <li>Follow Reddit's content policy</li>
        </ol>
      </div>

      <!-- Moderators Card -->
      <div class="moderators-card">
        <h3>Moderators</h3>
        <div class="moderator-list">
          <div class="moderator">
            <img src="/avatar.jpg" alt="Moderator" class="mod-avatar" />
            <span class="mod-name">u/moderator1</span>
          </div>
          <div class="moderator">
            <img src="/avatar.jpg" alt="Moderator" class="mod-avatar" />
            <span class="mod-name">u/moderator2</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .community-page {
    min-height: 100vh;
    background-color: white;
  }

  /* Banner */
  .community-banner {
    width: calc(100% - 48px);
    height: 200px;
    overflow: hidden;
    background: linear-gradient(to bottom, #33a8ff, #0079d3);
    border-radius: 8px;
    margin: 8px 24px;
  }

  .community-banner img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  /* Community Header */
  .community-header {
    background-color: white;
    border-bottom: 1px solid #edeff1;
    padding: 16px 0;
  }

  .community-header-content {
    max-width: 100%;
    padding: 0 24px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .community-info {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .community-icon {
    width: 72px;
    height: 72px;
    border-radius: 50%;
    border: 4px solid white;
    background: white;
    margin-top: -36px;
    object-fit: cover;
  }

  .community-title h1 {
    font-size: 28px;
    font-weight: 700;
    margin: 0;
    color: #1c1c1c;
  }

  .community-name {
    font-size: 14px;
    color: #7c7c7c;
    margin: 4px 0 0 0;
  }

  .join-btn {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 8px 24px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .join-btn:hover {
    background: var(--darkblue--);
  }

  .join-btn.joined {
    background: #edeff1;
    color: #1c1c1c;
    border: 1px solid #edeff1;
  }

  .join-btn.joined:hover {
    background: #d7dadc;
  }

  /* Container */
  .community-container {
    max-width: 100%;
    margin: 0 auto;
    padding: 24px;
    display: grid;
    grid-template-columns: 1fr 312px;
    gap: 24px;
  }

  /* Main Content */
  .community-main {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .sorting-bar {
    background-color: #f6f7f8;
    border-radius: 4px;
    padding: 4px;
    display: flex;
    gap: 0px;
  }

  .sort-btn {
    flex: 1;
    padding: 8px 12px;
    border: none;
    background-color: transparent;
    color: #878a8c;
    font-weight: bold;
    border-radius: 4px;
    cursor: pointer;
    transition:
      background-color 0.2s,
      color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .sort-btn.active {
    background-color: white;
    color: var(--blue--);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .sort-btn:not(.active):hover {
    background-color: #e9ebee;
  }

  .post-list {
    display: flex;
    flex-direction: column;
  }

  /* Sidebar */
  .community-sidebar {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .about-card,
  .rules-card,
  .moderators-card {
    background: white;
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 12px;
  }

  .about-card h3,
  .rules-card h3,
  .moderators-card h3 {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #1c1c1c;
    margin: 0 0 8px 0;
    padding-bottom: 12px;
    border-bottom: 1px solid #edeff1;
  }

  .about-description {
    font-size: 14px;
    line-height: 21px;
    color: #1c1c1c;
    margin: 12px 0;
  }

  .community-stats {
    display: flex;
    gap: 24px;
    padding: 12px 0;
    border-top: 1px solid #edeff1;
    border-bottom: 1px solid #edeff1;
    margin: 12px 0;
  }

  .stat {
    display: flex;
    flex-direction: column;
  }

  .stat-value {
    font-size: 16px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .stat-label {
    font-size: 12px;
    color: #7c7c7c;
  }

  .community-created {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #7c7c7c;
    font-size: 14px;
    margin: 12px 0;
  }

  .create-post-btn {
    width: 100%;
    background: var(--blue--);
    color: white;
    border: none;
    padding: 8px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
    margin-top: 12px;
  }

  .create-post-btn:hover {
    background: var(--darkblue--);
  }

  /* Rules */
  .rules-list {
    margin: 0;
    padding-left: 20px;
  }

  .rules-list li {
    font-size: 14px;
    color: #1c1c1c;
    margin: 8px 0;
    line-height: 18px;
  }

  /* Moderators */
  .moderator-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .moderator {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .mod-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    object-fit: cover;
  }

  .mod-name {
    font-size: 12px;
    color: #1c1c1c;
    font-weight: 500;
  }

  .mod-name:hover {
    text-decoration: underline;
    cursor: pointer;
  }

  /* Responsive */
  @media (max-width: 960px) {
    .community-container {
      grid-template-columns: 1fr;
    }

    .community-sidebar {
      order: -1;
    }
  }
</style>
