<script lang="ts">
  import Post from "../components/Post.svelte";
  import type { PostResponse } from "../dtos/post-dto";

  // TODO: Replace with API call to fetch posts (keeping mock data for now)
  const posts: PostResponse[] = [
    {
      id: "675e1a2b3c4d5e6f7a8b9c0d",
      author: {
        id: "675e1a2b3c4d5e6f7a8b9c0e",
        username: "user123",
        avatar: { public_id: "avatar1", url: "/avatar.jpg" },
      },
      community: {
        id: "675e1a2b3c4d5e6f7a8b9c0f",
        name: "sveltejs",
      },
      title: "Svelte 5 is amazing!",
      type: "text",
      content: {
        text: "I just tried out the new Svelte 5 features and they are mind-blowing. The new runes system is so intuitive!",
      },
      votes_count: {
        up: 0,
        down: 0,
        score: 0,
      },
      comments_count: 42,
      created_at: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(), // 4 hours ago
    },
    {
      id: "675e1a2b3c4d5e6f7a8b9c10",
      author: {
        id: "675e1a2b3c4d5e6f7a8b9c11",
        username: "photographer",
        avatar: { public_id: "avatar2", url: "/avatar.jpg" },
      },
      community: {
        id: "675e1a2b3c4d5e6f7a8b9c12",
        name: "pics",
      },
      title: "Girl on wayhome, who is she?",
      type: "text",
      content: {
        images: [
          {
            public_id: "img1",
            url: "/GirlFromNowhere.jpg",
            width: 800,
            height: 600,
          },
        ],
      },
      votes_count: {
        up: 0,
        down: 0,
        score: 0,
      },
      comments_count: 89,
      created_at: new Date(Date.now() - 8 * 60 * 60 * 1000).toISOString(), // 8 hours ago
    },
    {
      id: "675e1a2b3c4d5e6f7a8b9c13",
      author: {
        id: "675e1a2b3c4d5e6f7a8b9c14",
        username: "pollmaster",
        avatar: { public_id: "avatar3", url: "/avatar.jpg" },
      },
      community: {
        id: "675e1a2b3c4d5e6f7a8b9c15",
        name: "polls",
      },
      title: "What is your favorite frontend framework?",
      type: "poll",
      content: {
        poll: {
          question: "What is your favorite frontend framework?",
          options: [
            { id: "opt1", text: "Svelte", votes: 450, percentage: 57 },
            { id: "opt2", text: "React", votes: 200, percentage: 25.3 },
            { id: "opt3", text: "Vue", votes: 139, percentage: 17.6 },
          ],
          total_votes: 789,
          allow_multiple: false,
        },
      },
      votes_count: {
        up: 0,
        down: 0,
        score: 0,
      },
      comments_count: 231,
      created_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(), // 1 day ago
    },
    {
      id: "675e1a2b3c4d5e6f7a8b9c16",
      author: {
        id: "675e1a2b3c4d5e6f7a8b9c17",
        username: "videographer",
        avatar: { public_id: "avatar4", url: "/avatar.jpg" },
      },
      community: {
        id: "675e1a2b3c4d5e6f7a8b9c18",
        name: "videos",
      },
      title: "Flashback AMV",
      type: "text",
      content: {
        videos: [
          {
            public_id: "video1",
            url: "/video.mp4",
            thumbnail_url:
              "https://i1.sndcdn.com/artworks-000307576689-fkq1mv-t500x500.jpg",
          },
        ],
      },
      votes_count: {
        up: 0,
        down: 0,
        score: 0,
      },
      comments_count: 60,
      created_at: new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString(), // 3 hours ago
    },
  ];

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
