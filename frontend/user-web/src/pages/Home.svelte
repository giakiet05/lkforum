<script lang="ts">
  import { run } from "svelte/legacy";

  // Import các thành phần cần thiết
  import Sidebar from "../components/Sidebar.svelte"; // Updated import path
  import Topbar from "../components/Topbar.svelte";
  import { authStore } from "../stores/auth-store";
  import { logout } from "../services/auth-service";
  import { push } from "svelte-spa-router";
  import { get } from "svelte/store";
  import type { User } from "../models/user";

  import Post from "../components/Post.svelte";
  import type { PostData } from "../types/post.ts";

  const posts: PostData[] = [
    {
      id: "1",
      type: "text",
      community: "sveltejs",
      author: "user123",
      time: "4 hours ago",
      title: "Svelte 5 is amazing!",
      upvotes: 123,
      downvotes: 5,
      commentsCount: 42,
      content:
        "I just tried out the new Svelte 5 features and they are mind-blowing. The new runes system is so intuitive!",
    },
    {
      id: "2",
      type: "image",
      community: "pics",
      author: "photographer",
      time: "8 hours ago",
      title: "Girl on wayhome, who is she?",
      upvotes: 456,
      downvotes: 12,
      commentsCount: 89,
      images: ["/GirlFromNowhere.jpg"],
    },
    {
      id: "3",
      type: "poll",
      community: "polls",
      author: "pollmaster",
      time: "1 day ago",
      title: "What is your favorite frontend framework?",
      upvotes: 789,
      downvotes: 50,
      commentsCount: 231,
      poll: {
        question: "What is your favorite frontend framework?",
        options: [
          { id: 1, text: "Svelte", votes: 450 },
          { id: 2, text: "React", votes: 200 },
          { id: 3, text: "Vue", votes: 139 },
        ],
        multipleChoice: false,
        totalVotes: 789,
      },
    },
    {
      id: "4",
      type: "video",
      community: "videos",
      author: "videographer",
      time: "3 hours ago",
      title: "Flashback AMV",
      upvotes: 250,
      downvotes: 15,
      commentsCount: 60,
      videoUrl: "./video.mp4",
      thumbnailUrl:
        "https://i1.sndcdn.com/artworks-000307576689-fkq1mv-t500x500.jpg",
    },
  ];

  const sidebarItems = [
    {
      id: "home",
      label: "Home",
      to: "/",
      icon: `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M3 11.5L12 4l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1v-8.5z" fill="currentColor"/>
    </svg>`,
    },
    {
      id: "popular",
      label: "Popular",
      to: "/popular",
      icon: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g clip-path="url(#clip0_20_157)"><path d="M23.2499 12.751L12.7769 23.25" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><path d="M17.25 12.751H23.25V18.75" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><path d="M18.75 0.75V5.25H12.75V11.25H6.75V17.25H0.75V23.25" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></g><defs><clipPath id="clip0_20_157"><rect width="24" height="24" fill="white"/></clipPath></defs></svg>`,
    },
    {
      id: "community",
      label: "Community",
      icon: `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 2a5 5 0 100 10 5 5 0 000-10zm0 12c-4 0-8 2-8 4v2h16v-2c0-2-4-4-8-4z" fill="currentColor"/>
    </svg>`,
      children: [
        { id: "lk-all", label: "lk/all", to: "/lk/all", icon: "🌐" },
        { id: "lk-news", label: "lk/news", to: "/lk/news", icon: "📰" },
        { id: "lk-dev", label: "lk/dev", to: "/lk/dev", icon: "💻" },
        // thêm sau: { id: 'r-your', label: 'r/your', to: '/r/your', icon: '...' }
      ],
    },
    { id: "explore", label: "Explore", to: "/explore", icon: "🧭" },
    {
      id: "all",
      label: "All",
      to: "/all",
      icon: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g clip-path="url(#clip0_20_165)"><path d="M23.25 21.25C23.25 22.35 22.35 23.25 21.25 23.25H2.75C1.65 23.25 0.75 22.35 0.75 21.25V2.75C0.75 1.65 1.65 0.75 2.75 0.75H21.25C22.35 0.75 23.25 1.65 23.25 2.75V21.25Z" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-miterlimit="10" stroke-linecap="round" stroke-linejoin="round"/><path d="M7.44971 18.4499C7.44971 19.0499 7.04971 19.4499 6.44971 19.4499H5.44971C4.84971 19.4499 4.44971 19.0499 4.44971 18.4499V15.6499C4.44971 15.0499 4.84971 14.6499 5.44971 14.6499H6.44971C7.04971 14.6499 7.44971 15.0499 7.44971 15.6499V18.4499Z" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-miterlimit="10" stroke-linecap="round" stroke-linejoin="round"/><path d="M13.4502 18.4499C13.4502 19.0499 13.0502 19.4499 12.4502 19.4499H11.4502C10.8502 19.4499 10.4502 19.0499 10.4502 18.4499V6.6499C10.4502 6.0499 10.8502 5.6499 11.4502 5.6499H12.4502C13.0502 5.6499 13.4502 6.0499 13.4502 6.6499V18.4499Z" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-miterlimit="10" stroke-linecap="round" stroke-linejoin="round"/><path d="M19.4502 18.4499C19.4502 19.0499 19.0502 19.4499 18.4502 19.4499H17.4502C16.8502 19.4499 16.4502 19.0499 16.4502 18.4499V11.6499C16.4502 11.0499 16.8502 10.6499 17.4502 10.6499H18.4502C19.0502 10.6499 19.4502 11.0499 19.4502 11.6499V18.4499Z" stroke="currentColor" stroke-opacity="0.7" stroke-width="1.5" stroke-miterlimit="10" stroke-linecap="round" stroke-linejoin="round"/></g><defs><clipPath id="clip0_20_165"><rect width="24" height="24" fill="white"/></clipPath></defs></svg>`,
    },
  ];

  let activeSort = "hot";

  // Logic của trang Home giữ nguyên
  let user: User | null = null;
  // read authStore once on mount
  run(() => {
    const state = get(authStore);
    user = state.user;
  });

  // track whether sidebar is compact (bind from Sidebar)
  let isSidebarCompact = false;

  // pass a simple user object to Topbar (read from store on mount)
  let topbarUser:
    | { name: string; avatar?: string; karma?: number }
    | undefined = undefined;
  run(() => {
    const state = get(authStore);
    if (state.user) {
      topbarUser = {
        name: state.user.username || state.user.email || "User",
        avatar: "../avatar.jpg",
        karma: 20,
      };
    }
  });

  function handleLogout() {
    alert("Đăng xuất!");
    logout();

    // push("/login");
  }

  function handleNavigate(item: any) {
    push(item.to);
  }
</script>

<div class="app-layout">
  <!-- Sidebar is fixed; main content reserves space using CSS variable -->
  <Topbar user={topbarUser} onLogout={handleLogout} />

  <Sidebar
    items={sidebarItems}
    onNavigate={handleNavigate}
    bind:compact={isSidebarCompact}
  />

  <main class="main-content" data-compact={isSidebarCompact}>
    <div class="sorting-bar">
      <button
        class="sort-btn"
        class:active={activeSort === "hot"}
        on:click={() => (activeSort = "hot")}
      >
        🔥 Hot
      </button>
      <button
        class="sort-btn"
        class:active={activeSort === "new"}
        on:click={() => (activeSort = "new")}
      >
        ✨ New
      </button>
      <button
        class="sort-btn"
        class:active={activeSort === "top"}
        on:click={() => (activeSort = "top")}
      >
        🏆 Top
      </button>
    </div>
    <div class="post-list">
      {#each posts as post}
        <Post {post} />
      {/each}
    </div>
  </main>
</div>

<style>
  :root {
    --sidebar-width: 256px;
    --sidebar-compact-width: 64px;
  }

  .app-layout {
    position: relative;
    min-height: 100vh;
    background-color: white;
  }

  /* main-content reserves space equal to the sidebar width to avoid overlap */
  .main-content {
    margin-left: var(--sidebar-width); /* reserves space for the sidebar */
    transition: margin-left 0.2s ease;
    padding-top: calc(var(--topbar-height) + 1rem); /* reserve topbar space */
    padding-right: 2rem;
    padding-left: 2rem;
    padding-bottom: 2rem;
    min-height: 100vh;
    box-sizing: border-box;
    overflow-y: auto;
  }

  /* compact sidebar state (bind:compact -> data-compact attribute) */
  .main-content[data-compact="true"] {
    margin-left: var(--sidebar-compact-width);
  }

  /* small screens: sidebar becomes overlaying, so remove left margin */
  @media (max-width: 768px) {
    .main-content {
      margin-left: 0;
    }
  }

  .sorting-bar {
    background-color: #f6f7f8;
    border-radius: 4px;
    padding: 4px;
    margin-bottom: 16px;
    display: flex;
    gap: 4px;
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
    color: #0079d3;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .sort-btn:not(.active):hover {
    background-color: #e9ebee;
  }

  .user-info {
    background-color: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    max-width: 400px;
    margin: 1rem 0;
  }

  .user-info p {
    margin: 0.5rem 0;
  }

  button {
    /* Style đơn giản cho nút logout */
    padding: 0.6rem 1.2rem;
    border: none;
    background-color: #dc3545;
    color: white;
    border-radius: 5px;
    cursor: pointer;
  }
</style>
