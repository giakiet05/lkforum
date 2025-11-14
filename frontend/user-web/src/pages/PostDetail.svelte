<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import CommentSection from "../components/CommentSection.svelte";
  import type { PostResponse } from "../dtos/post-dto";

  type PostDetailProps = {
    params?: { id: string };
  };

  let { params = { id: "1" } }: PostDetailProps = $props();

  let post = $state<PostResponse | null>(null);
  let selectedOptions = $state<string[]>([]);
  let hasVoted = $state(false);

  onMount(() => {
    // Scroll to top when component mounts
    window.scrollTo(0, 0);

    // TODO: Replace with API call to fetch post by ID
    // Mock data for now
    // post = {
    //   id: params.id,
    //   author: {
    //     id: "user1",
    //     username: "nguyenvana",
    //     avatar: {
    //       public_id: "avatar1",
    //       url: "https://i.pravatar.cc/150?img=1",
    //     },
    //   },
    //   community: {
    //     id: "comm1",
    //     name: "programming",
    //   },
    //   title: "What's your favorite programming language and why?",
    //   type: "text",
    //   content: {
    //     text: "I've been coding for 5 years now and I've tried many languages. Currently loving TypeScript for its type safety and great tooling. What about you guys?\n\nI'm curious to hear different perspectives, especially from people working in different domains (web, mobile, systems programming, etc.).",
    //   },
    //   votes_count: {
    //     up: 156,
    //     down: 23,
    //     score: 133,
    //   },
    //   user_vote: "",
    //   comments_count: 47,
    //   created_at: new Date(Date.now() - 1000 * 60 * 60 * 3).toISOString(), // 3 hours ago
    // };

    // Uncomment below for Image post example
    post = {
      id: params.id,
      author: {
        id: "user2",
        username: "photographer_pro",
        avatar: {
          public_id: "avatar2",
          url: "https://i.pravatar.cc/150?img=2",
        },
      },
      community: {
        id: "comm2",
        name: "photography",
      },
      title: "Golden hour at the beach 🌅",
      type: "text",
      content: {
        text: "Captured this beautiful sunset yesterday!",
        images: [
          {
            public_id: "img1",
            url: "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800",
            width: 800,
            height: 600,
          },
          {
            public_id: "img2",
            url: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=800",
            width: 800,
            height: 600,
          },
        ],
      },
      votes_count: {
        up: 892,
        down: 12,
        score: 880,
      },
      user_vote: "up",
      comments_count: 156,
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 8).toISOString(),
    };

    // Uncomment below for Poll post example
    // post = {
    //   id: params.id,
    //   author: {
    //     id: "user3",
    //     username: "poll_master",
    //   },
    //   community: {
    //     id: "comm3",
    //     name: "polls"
    //   },
    //   title: "Which framework do you prefer for frontend development?",
    //   type: "poll",
    //   content: {
    //     poll: {
    //       question: "Select your favorite frontend framework:",
    //       options: [
    //         { id: "opt1", text: "React", votes: 450, percentage: 45 },
    //         { id: "opt2", text: "Vue.js", votes: 250, percentage: 25 },
    //         { id: "opt3", text: "Angular", votes: 150, percentage: 15 },
    //         { id: "opt4", text: "Svelte", votes: 150, percentage: 15 }
    //       ],
    //       total_votes: 1000,
    //       allow_multiple: false,
    //       expires_at: new Date(Date.now() + 1000 * 60 * 60 * 24 * 7).toISOString() // 7 days from now
    //     }
    //   },
    //   votes_count: {
    //     up: 234,
    //     down: 8,
    //     score: 226
    //   },
    //   comments_count: 89,
    //   created_at: new Date(Date.now() - 1000 * 60 * 60 * 12).toISOString(),
    // };
  });

  function handleVote(optionId: string) {
    if (!post || hasVoted) return;

    if (post.content.poll?.allow_multiple) {
      const index = selectedOptions.indexOf(optionId);
      if (index > -1) {
        selectedOptions.splice(index, 1);
      } else {
        selectedOptions.push(optionId);
      }
      selectedOptions = selectedOptions;
    } else {
      selectedOptions = [optionId];
    }
  }

  function submitVote() {
    if (!post || selectedOptions.length === 0) return;

    if (post.content.poll) {
      for (const option of post.content.poll.options) {
        if (selectedOptions.includes(option.id)) {
          option.votes++;
        }
      }
      post.content.poll.total_votes += selectedOptions.length;
    }
    hasVoted = true;
  }

  const getVotePercentage = (votes: number, total: number) => {
    if (total === 0) return 0;
    return (votes / total) * 100;
  };

  function goBack() {
    window.history.back();
  }
</script>

<div class="post-detail-page">
  {#if post}
    <!-- Back button -->
    <button class="back-btn" onclick={goBack}>
      <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path d="M19 12H5M12 19l-7-7 7-7" />
      </svg>
      Back
    </button>

    <!-- Post Content -->
    <article class="post-detail-container">
      <div class="post-main">
        <div class="post-header">
          <span class="community-name">lk/{post.community.name}</span>
          <span class="meta-divider">•</span>
          <span class="author">Posted by u/{post.author.username}</span>
          <span class="time"
            >{new Date(post.created_at).toLocaleDateString()}</span
          >
        </div>

        <h1 class="post-title">{post.title}</h1>

        <div class="post-content">
          {#if post.type === "text" && post.content.text}
            <p class="text-content">{post.content.text}</p>
          {:else if post.content.images && post.content.images.length > 0}
            <div class="image-gallery">
              {#each post.content.images as image, i}
                <img
                  src={image.url}
                  alt="Post content {i + 1}"
                  class="post-image"
                />
              {/each}
            </div>
          {:else if post.content.video}
            <video
              controls
              poster={post.content.video.thumbnail}
              class="post-video"
            >
              <source src={post.content.video.url} type="video/mp4" />
              <track kind="captions" />
              Your browser does not support the video tag.
            </video>
          {:else if post.type === "poll" && post.content.poll}
            <div class="poll-container">
              <h3 class="poll-question">{post.content.poll.question}</h3>
              <div class="poll-options">
                {#each post.content.poll.options as option}
                  <button
                    class="poll-option"
                    class:selected={selectedOptions.includes(option.id)}
                    onclick={() => handleVote(option.id)}
                    disabled={hasVoted}
                  >
                    {#if hasVoted}
                      <div
                        class="poll-result-bar"
                        style="width: {option.percentage}%;"
                      ></div>
                      <span class="poll-option-text">{option.text}</span>
                      <span class="poll-option-votes">{option.votes} votes</span
                      >
                    {:else}
                      <div class="radio-check">
                        {#if post.content.poll.allow_multiple}
                          <div
                            class="checkbox"
                            class:checked={selectedOptions.includes(option.id)}
                          ></div>
                        {:else}
                          <div
                            class="radio"
                            class:checked={selectedOptions.includes(option.id)}
                          ></div>
                        {/if}
                      </div>
                      <span class="poll-option-text">{option.text}</span>
                    {/if}
                  </button>
                {/each}
              </div>
              {#if !hasVoted}
                <button
                  class="vote-submit-btn"
                  onclick={submitVote}
                  disabled={selectedOptions.length === 0}
                >
                  Vote
                </button>
              {/if}
              <p class="poll-footer">
                {post.content.poll.total_votes} votes • {post.content.poll
                  .allow_multiple
                  ? "Multiple choices allowed"
                  : "Single choice"}
              </p>
            </div>
          {/if}
        </div>

        <div class="post-footer">
          <div class="vote-actions">
            <button class="footer-btn vote-btn" aria-label="Upvote">▲</button>
            <span class="vote-count"
              >{post.votes_count ? post.votes_count.score : 0}</span
            >
            <button class="footer-btn vote-btn" aria-label="Downvote">▼</button>
          </div>
          <button class="footer-btn">
            <img src="/CommentIcon.svg" alt="Comments" width="20" height="20" />
            <span>{post.comments_count} Comments</span>
          </button>
          <button class="footer-btn">
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path>
              <polyline points="16 6 12 2 8 6"></polyline>
              <line x1="12" y1="2" x2="12" y2="15"></line>
            </svg>
            <span>Share</span>
          </button>
          <button class="footer-btn">
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"
              ></path>
            </svg>
            <span>Save</span>
          </button>
        </div>
      </div>
    </article>

    <!-- Comment Section -->
    <CommentSection postId={post.id} />
  {:else}
    <div class="loading">Loading...</div>
  {/if}
</div>

<style>
  .post-detail-page {
    max-width: 900px;
    margin: 0 auto;
    padding: 16px;
  }

  .back-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    background: none;
    border: none;
    color: #1c1c1c;
    font-size: 14px;
    font-weight: 400;
    cursor: pointer;
    padding: 8px 12px;
    margin-bottom: 16px;
    border-radius: 4px;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .back-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .post-detail-container {
    background-color: white;
    border: 1px solid #eaebef;
    border-radius: 4px;
    color: #000000;
    font-family: "Roboto", Arial, sans-serif;
  }

  .vote-btn {
    color: #878a8c;
  }

  .vote-btn:hover {
    color: var(--darkblue--);
  }

  .vote-count {
    font-weight: bold;
    font-size: 12px;
    margin: 0 4px;
    color: var(--darkblue--);
  }

  .post-main {
    padding: 16px 24px;
  }

  .post-header {
    display: flex;
    align-items: center;
    font-size: 12px;
    margin-bottom: 12px;
  }

  .community-name {
    font-weight: bold;
    color: #000000;
  }

  .meta-divider {
    margin: 0 4px;
    color: #878a8c;
  }

  .author,
  .time {
    color: #878a8c;
  }

  .author {
    margin-right: 4px;
  }

  .post-title {
    font-size: 24px;
    font-weight: 600;
    color: #000000;
    margin: 0 0 16px 0;
    line-height: 1.3;
  }

  .post-content {
    margin-bottom: 16px;
  }

  .text-content {
    font-size: 16px;
    line-height: 24px;
    white-space: pre-wrap;
    color: #1c1c1c;
  }

  .image-gallery {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .post-image {
    max-width: 100%;
    max-height: 600px;
    border-radius: 4px;
    display: block;
    margin: auto;
  }

  .post-video {
    width: 100%;
    max-height: 600px;
    border-radius: 4px;
    background-color: #000;
  }

  .post-footer {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-top: 12px;
    border-top: 1px solid #eaebef;
  }

  .vote-actions {
    display: flex;
    align-items: center;
    background-color: var(--button-secondary-background);
    border-radius: 20px;
  }

  .footer-btn {
    background: rgba(214, 216, 222, 0.4);
    border: none;
    color: #000000;
    font-weight: bold;
    font-size: 12px;
    padding: 8px 12px;
    border-radius: 20px;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
    transition: background-color 0.2s;
  }

  .vote-actions .footer-btn {
    background: transparent;
    border-radius: 4px;
  }

  .footer-btn:hover {
    background-color: var(--button-secondary-background-hover);
  }

  /* Poll Styles */
  .poll-container {
    margin-top: 12px;
  }

  .poll-question {
    font-size: 18px;
    font-weight: 600;
    margin: 0 0 16px;
  }

  .poll-options {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .poll-option {
    position: relative;
    display: flex;
    align-items: center;
    width: 100%;
    padding: 12px 16px;
    border: 2px solid #ccc;
    border-radius: 4px;
    background-color: white;
    color: #000000;
    cursor: pointer;
    text-align: left;
    overflow: hidden;
    font-size: 15px;
  }

  .poll-option:not(:disabled):hover {
    border-color: #878a8c;
  }

  .poll-option.selected {
    border-color: var(--primary-color);
    background-color: #f0f8ff;
  }

  .poll-option:disabled {
    cursor: not-allowed;
  }

  .poll-result-bar {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    background-color: var(--primary-color);
    opacity: 0.2;
    transition: width 0.5s ease-in-out;
  }

  .poll-option-text {
    position: relative;
    z-index: 1;
    flex-grow: 1;
  }

  .poll-option-votes {
    position: relative;
    z-index: 1;
    font-weight: bold;
  }

  .radio-check {
    margin-right: 12px;
    flex-shrink: 0;
  }

  .radio,
  .checkbox {
    width: 18px;
    height: 18px;
    border: 2px solid #878a8c;
    display: inline-block;
  }

  .radio {
    border-radius: 50%;
  }

  .checkbox {
    border-radius: 3px;
  }

  .radio.checked,
  .checkbox.checked {
    border-color: var(--primary-color);
    background-color: var(--primary-color);
  }

  .vote-submit-btn {
    margin-top: 12px;
    padding: 10px 24px;
    border: 1px solid var(--primary-color);
    background-color: var(--primary-color);
    color: white;
    border-radius: 20px;
    font-weight: bold;
    font-size: 14px;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
  }

  .vote-submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .vote-submit-btn:not(:disabled):hover {
    background-color: var(--primary-color-hover);
  }

  .poll-footer {
    font-size: 13px;
    color: #878a8c;
    margin-top: 12px;
  }

  .loading {
    text-align: center;
    padding: 40px;
    font-size: 16px;
    color: #878a8c;
  }
</style>
