<script lang="ts">
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

  let activeSort = "hot";
</script>

<div>
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
</div>

<style>
  .sorting-bar {
    background-color: #f6f7f8;
    border-radius: 4px;
    padding: 4px;
    margin-bottom: 16px;
    margin-top: -16px;
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
    color: #0079d3;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .sort-btn:not(.active):hover {
    background-color: #e9ebee;
  }
</style>
