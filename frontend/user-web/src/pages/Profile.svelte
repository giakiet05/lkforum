<script lang="ts">
  import Post from '../components/Post.svelte';
  import type { PostData } from '../types/post';

  const user = {
    name: 'Long',
    avatar: '/avatar.jpg',
    karma: 12345,
    cakeDay: 'October 18, 2025',
    coverImage: '/background.png',
  };

  const posts: PostData[] = [
    {
      id: '1',
      type: 'text',
      community: 'sveltejs',
      author: 'Long',
      time: '4 hours ago',
      title: 'Svelte 5 is amazing!',
      upvotes: 123,
      downvotes: 5,
      commentsCount: 42,
      content: 'I just tried out the new Svelte 5 features and they are mind-blowing. The new runes system is so intuitive!',
    },
    {
      id: '2',
      type: 'image',
      community: 'pics',
      author: 'Long',
      time: '8 hours ago',
      title: 'A picture I took',
      upvotes: 456,
      downvotes: 12,
      commentsCount: 89,
      images: ['/discuss.jpg'],
    },
  ];

  let activeTab: 'posts' | 'comments' | 'about' = 'posts';
</script>

<div class="profile-page">
  <div class="profile-header">
    <div class="cover-image-wrapper">
      <img src={user.coverImage} alt="Cover" class="cover-image" />
    </div>
    <div class="profile-info-bar">
        <div class="profile-details">
            <img src={user.avatar} alt={user.name} class="profile-avatar" />
            <div class="profile-text">
              <h1 class="username">{user.name}</h1>
              <p class="user-meta">{user.karma} karma</p>
            </div>
        </div>
      <button class="edit-profile-btn">Edit Profile</button>
    </div>
  </div>

  <div class="profile-content">
    <div class="profile-main-content">
        <div class="profile-tabs">
            <button class="tab-btn" class:active={activeTab === 'posts'} on:click={() => activeTab = 'posts'}>Posts</button>
            <button class="tab-btn" class:active={activeTab === 'comments'} on:click={() => activeTab = 'comments'}>Comments</button>
            <button class="tab-btn" class:active={activeTab === 'about'} on:click={() => activeTab = 'about'}>About</button>
          </div>
      
          <div class="tab-content">
            {#if activeTab === 'posts'}
              <div class="post-list">
                {#each posts as post}
                  <Post {post} />
                {/each}
              </div>
            {:else if activeTab === 'comments'}
              <p>Comments will be shown here.</p>
            {:else if activeTab === 'about'}
              <p>About section will be shown here.</p>
            {/if}
          </div>
    </div>
    <div class="profile-sidebar">
        <div class="user-card">
            <div class="user-card-header">
                <h2 class="username">{user.name}</h2>
                <p class="user-meta">{user.karma} karma</p>
            </div>
            <div class="user-card-body">
                <div class="user-stats">
                    <div class="stat">
                        <span class="stat-value">{user.karma}</span>
                        <span class="stat-label">Karma</span>
                    </div>
                    <div class="stat">
                        <span class="stat-value">{user.cakeDay}</span>
                        <span class="stat-label">Cake day</span>
                    </div>
                </div>
                <button class="new-post-btn">New post</button>
            </div>
        </div>
    </div>
  </div>
</div>

<style>
  .profile-page {
    background-color: #dae0e6;
    min-height: 100vh;
  }

  .profile-header {
    background-color: white;
  }

  .cover-image-wrapper {
    height: 200px;
    overflow: hidden;
  }

  .cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .profile-info-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 2rem;
    max-width: 1000px;
    margin: 0 auto;
  }

  .profile-details {
    display: flex;
    align-items: flex-end;
    transform: translateY(-50px);
  }

  .profile-avatar {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    border: 4px solid white;
  }

  .profile-text {
    margin-left: 1rem;
    margin-bottom: 1rem;
  }

  .username {
    font-size: 2rem;
    font-weight: bold;
    margin: 0;
  }

  .user-meta {
    color: #878a8c;
    margin: 0;
  }

  .edit-profile-btn {
    background-color: transparent;
    border: 1px solid #0079d3;
    color: #0079d3;
    padding: 0.5rem 1rem;
    border-radius: 9999px;
    font-weight: bold;
    cursor: pointer;
  }

  .profile-content {
    display: flex;
    gap: 1rem;
    padding: 1rem 2rem;
    max-width: 1000px;
    margin: 0 auto;
  }

  .profile-main-content {
    flex-grow: 1;
  }

  .profile-tabs {
    background-color: white;
    border-bottom: 1px solid #ccc;
    display: flex;
  }

  .tab-btn {
    background: none;
    border: none;
    padding: 1rem;
    font-weight: bold;
    cursor: pointer;
    color: #878a8c;
  }

  .tab-btn.active {
    color: #0079d3;
    border-bottom: 2px solid #0079d3;
  }

  .tab-content {
    padding-top: 1rem;
  }

  .profile-sidebar {
    width: 300px;
  }

  .user-card {
      background-color: white;
      border-radius: 4px;
      border: 1px solid #ccc;
  }

  .user-card-header {
      background-color: #0079d3;
      color: white;
      padding: 1rem;
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
  }

  .user-card-body {
      padding: 1rem;
  }

  .user-stats {
      display: flex;
      justify-content: space-between;
      margin-bottom: 1rem;
  }

  .stat {
      text-align: center;
  }

  .stat-value {
      font-weight: bold;
      display: block;
  }

  .stat-label {
      color: #878a8c;
      font-size: 0.8rem;
  }

  .new-post-btn {
      width: 100%;
      background-color: #0079d3;
      color: white;
      border: none;
      padding: 0.75rem;
      border-radius: 9999px;
      font-weight: bold;
      cursor: pointer;
  }
</style>
