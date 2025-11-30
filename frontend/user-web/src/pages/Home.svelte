<script lang="ts">
  import Post from "../components/Post.svelte";
  import type { PostResponse } from "../dtos/post-dto";

  // TODO: Replace with API call to fetch posts
  const posts: PostResponse[] = [];

  /* Old mock data - needs to be replaced with real API data
  const posts: PostResponse[] = [
    {
      id: "1",
      type: "text",
      report: "sveltejs",
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
      report: "pics",
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
      report: "polls",
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
      report: "videos",
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
  */

  let sortBy: "best" | "hot" | "new" | "top" | "rising" | "" = "";
</script>

<div>
  <div class="sort-options">
    <select bind:value={sortBy}>
      <option value="" disabled selected hidden>Sort by</option>
      <option value="best">Best</option>
      <option value="hot">Hot</option>
      <option value="new">New</option>
      <option value="top">Top</option>
      <option value="rising">Rising</option>
    </select>
  </div>
  <div class="post-list">
    {#each posts as post}
      <Post {post} />
    {/each}
  </div>
</div>

<style>
  div {
    padding: 16px 24px;
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

  .sort-options select option {
    padding: 12px 16px;
    background-color: white;
    color: #1a1a1b;
    font-size: 14px;
    font-weight: 400;
    border: none;
  }

  .sort-options select option:first-child {
    color: #878a8c;
    font-weight: 500;
  }

  .sort-options select option:hover {
    background-color: #f8f9fa;
  }

  .sort-options select option:checked {
    background-color: #f0f1f2;
    font-weight: 500;
  }

  /* Custom dropdown appearance */
  @media screen and (-webkit-min-device-pixel-ratio: 0) {
    .sort-options select option {
      padding: 12px 16px;
    }

    .sort-options select option:checked {
      background: #f0f1f2;
    }
  }
</style>
