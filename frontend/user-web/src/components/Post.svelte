<script lang="ts">
  import { fade } from "svelte/transition";
  import { push } from "svelte-spa-router";
  import type { PostResponse } from "../dtos/post-dto";

  type PostProps = {
    post: PostResponse;
  };

  let { post }: PostProps = $props();

  let selectedOptions = $state<string[]>([]);
  let hasVoted = $state(false);
  let currentImageIndex = $state(0);

  function handlePostClick() {
    push(`/post/${post.id}`);
  }

  function handleButtonClick(e: MouseEvent) {
    // Prevent navigation when clicking buttons
    e.stopPropagation();
  }

  function handleCommunityClick(e: MouseEvent) {
    e.stopPropagation();
    push(`/lk/${post.community.name}`);
  }

  function nextImage(e: MouseEvent) {
    e.stopPropagation();
    if (
      post.content.images &&
      currentImageIndex < post.content.images.length - 1
    ) {
      currentImageIndex++;
    }
  }

  function prevImage(e: MouseEvent) {
    e.stopPropagation();
    if (currentImageIndex > 0) {
      currentImageIndex--;
    }
  }

  function handleVote(optionId: string) {
    if (hasVoted) return;

    if (post.content.poll?.allow_multiple) {
      const index = selectedOptions.indexOf(optionId);
      if (index > -1) {
        selectedOptions.splice(index, 1);
      } else {
        selectedOptions.push(optionId);
      }
      selectedOptions = selectedOptions; // Trigger reactivity
    } else {
      selectedOptions = [optionId];
    }
  }

  function submitVote() {
    if (selectedOptions.length === 0) return;
    // Here you would typically send the vote to a server
    // For this example, we'll just simulate it
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

  function handlePollOptionClick(e: MouseEvent, optionId: string) {
    e.stopPropagation();
    handleVote(optionId);
  }

  function handleVoteSubmit(e: MouseEvent) {
    e.stopPropagation();
    submitVote();
  }

  const getVotePercentage = (votes: number, total: number) => {
    if (total === 0) return 0;
    return (votes / total) * 100;
  };
</script>

<article class="post-container" transition:fade onclick={handlePostClick}>
  <div class="post-main">
    <div class="post-header">
      <div class="post-header-left">
        <img
          src={post.author.avatar?.url || "/avatar.jpg"}
          alt="User avatar"
          class="author-avatar"
        />
        <span class="community-name" onclick={handleCommunityClick}
          >lk/{post.community.name}</span
        >
        <span class="meta-divider">•</span>
        <span class="author">Posted by u/{post.author.username}</span>
        <span class="time"
          >{new Date(post.created_at).toLocaleDateString()}</span
        >
      </div>
      <div class="post-header-right">
        <button class="join-btn" onclick={handleButtonClick}>Join</button>
        <button
          class="more-btn"
          onclick={handleButtonClick}
          title="More options"
        >
          <img src="/dot.png" alt="" width="20" height="20" />
        </button>
      </div>
    </div>

    <h2 class="post-title">{post.title}</h2>

    <div class="post-content">
      {#if post.type === "text" && post.content.text}
        <p class="text-content">{post.content.text}</p>
      {/if}
      {#if post.content.images && post.content.images.length > 0}
        <div class="image-carousel">
          <img
            src={post.content.images[currentImageIndex].url}
            alt="Post content {currentImageIndex + 1}"
            class="post-image"
          />

          {#if post.content.images.length > 1}
            <button
              class="carousel-btn prev-btn"
              onclick={prevImage}
              disabled={currentImageIndex === 0}
              aria-label="Previous image"
            >
              ‹
            </button>
            <button
              class="carousel-btn next-btn"
              onclick={nextImage}
              disabled={currentImageIndex === post.content.images.length - 1}
              aria-label="Next image"
            >
              ›
            </button>
            <div class="image-counter">
              {currentImageIndex + 1} / {post.content.images.length}
            </div>
          {/if}
        </div>
      {/if}
      {#if post.content.video}
        <video
          controls
          poster={post.content.video.thumbnail}
          class="post-video"
        >
          <source src={post.content.video.url} type="video/mp4" />
          Your browser does not support the video tag.
        </video>
      {/if}
      {#if post.type === "poll" && post.content.poll}
        <div class="poll-container">
          <h3 class="poll-question">{post.content.poll.question}</h3>
          <div class="poll-options">
            {#each post.content.poll.options as option}
              <button
                class="poll-option"
                class:selected={selectedOptions.includes(option.id)}
                onclick={(e) => handlePollOptionClick(e, option.id)}
                disabled={hasVoted}
              >
                {#if hasVoted}
                  <div
                    class="poll-result-bar"
                    style="width: {option.percentage}%;"
                  ></div>
                  <span class="poll-option-text">{option.text}</span>
                  <span class="poll-option-votes">{option.votes} votes</span>
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
              onclick={handleVoteSubmit}
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
        <button
          class="footer-btn vote-btn"
          aria-label="Upvote"
          onclick={handleButtonClick}>▲</button
        >
        <span class="vote-count">{post.votes_count?.score || 0}</span>
        <button
          class="footer-btn vote-btn"
          aria-label="Downvote"
          onclick={handleButtonClick}>▼</button
        >
      </div>
      <button class="footer-btn" onclick={handlePostClick}>
        <img src="/CommentIcon.svg" alt="Comments" width="20" height="20" />
        <span>{post.comments_count} Comments</span>
      </button>
      <button class="footer-btn" onclick={handleButtonClick}>
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          ><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path><polyline
            points="16 6 12 2 8 6"
          ></polyline><line x1="12" y1="2" x2="12" y2="15"></line></svg
        >
        <span>Share</span>
      </button>
      <button class="footer-btn" onclick={handleButtonClick}>
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          ><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"
          ></path></svg
        >
        <span>Save</span>
      </button>
    </div>
  </div>
</article>

<style>
  .post-container {
    background-color: white;
    border: 1px solid #eaebef;
    border-radius: 4px;
    margin-bottom: 10px;
    color: #000000;
    font-family: Arial, sans-serif;
    cursor: pointer;
    transition: border-color 0.2s;
  }
  .post-container:hover {
    border-color: var(--button-secondary-background);
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
    padding: 8px 16px;
    flex-grow: 1;
    overflow: hidden;
  }

  .post-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 12px;
    margin-bottom: 8px;
  }

  .post-header-left {
    display: flex;
    align-items: center;
    flex: 1;
  }

  .post-header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .author-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    object-fit: cover;
    margin-right: 8px;
  }

  .community-name {
    font-weight: bold;
    color: #000000;
    cursor: pointer;
  }

  .community-name:hover {
    text-decoration: underline;
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

  .join-btn {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 4px 12px;
    border-radius: 9999px;
    font-size: 12px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .join-btn:hover {
    background: var(--darkblue--);
  }

  .more-btn {
    background: transparent;
    border: none;
    padding: 4px;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
  }

  .more-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .more-btn img {
    display: block;
  }

  .post-title {
    font-size: 18px;
    font-weight: 500;
    color: #000000;
    margin: 0 0 8px 0;
  }

  .post-content {
    margin-bottom: 8px;
  }
  .text-content {
    font-size: 14px;
    line-height: 21px;
    white-space: pre-wrap;
    color: rgba(0, 0, 0, 0.6);
  }

  /* Image Carousel */
  .image-carousel {
    position: relative;
    width: 100%;
    max-height: 500px;
    border-radius: 4px;
    overflow: hidden;
    background: #000;
    margin-top: 8px;
  }

  .post-image {
    width: 100%;
    height: 500px;
    object-fit: contain;
    display: block;
    background: #000;
  }

  .carousel-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(255, 255, 255, 0.9);
    border: none;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    font-size: 24px;
    font-weight: bold;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    color: #000;
    z-index: 2;
  }

  .carousel-btn:hover:not(:disabled) {
    background: rgba(255, 255, 255, 1);
    transform: translateY(-50%) scale(1.1);
  }

  .carousel-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .prev-btn {
    left: 12px;
  }

  .next-btn {
    right: 12px;
  }

  .image-counter {
    position: absolute;
    bottom: 12px;
    right: 12px;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
    z-index: 2;
  }

  .post-video {
    width: 100%;
    max-height: 500px;
    border-radius: 4px;
    background-color: #000;
  }

  .post-footer {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .vote-actions {
    display: flex;
    align-items: center;
    background-color: var(--button-secondary-background);
    border-radius: 20px;
  }

  .footer-btn {
    background: rgba(
      214,
      216,
      222,
      0.4
    ); /* --button-secondary-background at 40% */
    border: none;
    color: #000000;
    font-weight: bold;
    font-size: 12px;
    padding: 10px 14px;
    border-radius: 20px; /* Make it rounder */
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
    transition: background-color 0.2s;
  }

  /* Keep vote buttons inside the container transparent */
  .vote-actions .footer-btn {
    background: transparent;
    border-radius: 4px; /* Reset border-radius if needed */
    padding: 8px 10px;
  }

  .footer-btn:hover {
    background-color: var(--button-secondary-background-hover);
  }

  /* Poll Styles */
  .poll-container {
    margin-top: 12px;
  }
  .poll-question {
    font-size: 15px;
    margin: 0 0 10px;
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
    padding: 8px 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
    background-color: white;
    color: #000000;
    cursor: pointer;
    text-align: left;
    overflow: hidden;
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
    width: 16px;
    height: 16px;
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
    padding: 8px 16px;
    border: 1px solid var(--primary-color);
    background-color: var(--primary-color);
    color: white;
    border-radius: 20px;
    font-weight: bold;
    cursor: pointer;
  }
  .vote-submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .poll-footer {
    font-size: 12px;
    color: #878a8c;
    margin-top: 10px;
  }
</style>
