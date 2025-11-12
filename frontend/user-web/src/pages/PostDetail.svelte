<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import CommentSection from "../components/CommentSection.svelte";
  import type { PostResponse } from "../dtos/post-dto";

  type PostDetailProps = {
    params?: { id: string };
  };

  let { params = { id: "1" } }: PostDetailProps = $props();

  let post: PostResponse | undefined = $state(undefined);
  let selectedOptions = $state<string[]>([]);
  let hasVoted = $state(false);

  onMount(() => {
    // Scroll to top when component mounts
    window.scrollTo(0, 0);

    // TODO: Replace with API call to fetch post by ID
    // For now, post will be undefined (no mock data)
    // const response = await fetch(`/api/posts/${params.id}`);
    // post = await response.json();
  });

  function handleVote(optionId: number) {
    if (!post || hasVoted) return;

    if (post.poll?.multipleChoice) {
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

    if (post.poll) {
      for (const option of post.poll.options) {
        if (selectedOptions.includes(option.id)) {
          option.votes++;
        }
      }
      post.poll.totalVotes += selectedOptions.length;
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
          <span class="community-name">lk/{post.community}</span>
          <span class="meta-divider">•</span>
          <span class="author">Posted by u/{post.author}</span>
          <span class="time">{post.time}</span>
        </div>

        <h1 class="post-title">{post.title}</h1>

        <div class="post-content">
          {#if post.type === "text"}
            <p class="text-content">{post.content}</p>
          {:else if post.type === "image" && post.images}
            <div class="image-gallery">
              {#each post.images as src, i}
                <img {src} alt="Post content {i + 1}" class="post-image" />
              {/each}
            </div>
          {:else if post.type === "video" && post.videoUrl}
            <video controls poster={post.thumbnailUrl} class="post-video">
              <source src={post.videoUrl} type="video/mp4" />
              Your browser does not support the video tag.
            </video>
          {:else if post.type === "poll" && post.poll}
            <div class="poll-container">
              <h3 class="poll-question">{post.poll.question}</h3>
              <div class="poll-options">
                {#each post.poll.options as option}
                  <button
                    class="poll-option"
                    class:selected={selectedOptions.includes(option.id)}
                    onclick={() => handleVote(option.id)}
                    disabled={hasVoted}
                  >
                    {#if hasVoted}
                      <div
                        class="poll-result-bar"
                        style="width: {getVotePercentage(
                          option.votes,
                          post.poll.totalVotes
                        )}%;"
                      ></div>
                      <span class="poll-option-text">{option.text}</span>
                      <span class="poll-option-votes">{option.votes} votes</span
                      >
                    {:else}
                      <div class="radio-check">
                        {#if post.poll.multipleChoice}
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
                {post.poll.totalVotes} votes • {post.poll.multipleChoice
                  ? "Multiple choices allowed"
                  : "Single choice"}
              </p>
            </div>
          {/if}
        </div>

        <div class="post-footer">
          <div class="vote-actions">
            <button class="footer-btn vote-btn" aria-label="Upvote">▲</button>
            <span class="vote-count">{post.upvotes - post.downvotes}</span>
            <button class="footer-btn vote-btn" aria-label="Downvote">▼</button>
          </div>
          <button class="footer-btn">
            <img src="/CommentIcon.svg" alt="Comments" width="20" height="20" />
            <span>{post.commentsCount} Comments</span>
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
