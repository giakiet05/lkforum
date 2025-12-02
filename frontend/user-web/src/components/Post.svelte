<script lang="ts">
  import { fade } from "svelte/transition";
  import { push } from "svelte-spa-router";
  import type { PostResponse } from "../dtos/post-dto";
  import {
    voteOnPost,
    savePost,
    unsavePost,
    hidePost,
    voteOnPoll,
    removePollVote,
    deletePost,
    reportPost,
  } from "../services/post-service";
  import { authStore } from "../stores/auth-store";

  type PostProps = {
    post: PostResponse;
    onUpdate?: () => void; // Callback to refresh post data
  };

  let { post, onUpdate }: PostProps = $props();

  let selectedOptions = $state<string[]>([]);
  let currentImageIndex = $state(0);

  // Check if user has voted based on backend data
  const hasVoted = $derived(
    post.content.poll?.user_vote_ids &&
      post.content.poll.user_vote_ids.length > 0
  );

  // Vote state
  let userVote = $state<"up" | "down" | "">(
    (post.user_vote as "up" | "down" | "") || ""
  );
  let votesCount = $state(post.votes_count?.score || 0);
  let isVoting = $state(false);

  // Save state
  let isSaved = $state(false);
  let isSaving = $state(false);

  // Menu state
  let showMenu = $state(false);
  let showReportModal = $state(false);
  let reportReason = $state("");
  let reportDetails = $state("");
  let isReporting = $state(false);

  const currentUser = $derived($authStore.user);
  const isOwnPost = $derived(currentUser && post.author.id === currentUser.id);

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
    const votedOptions = post.content.poll?.user_vote_ids || [];
    const alreadyVotedThisOption = votedOptions.includes(optionId);

    if (post.content.poll?.allow_multiple) {
      // Multiple choice: toggle selection
      const index = selectedOptions.indexOf(optionId);
      if (index > -1) {
        selectedOptions.splice(index, 1);
      } else {
        selectedOptions.push(optionId);
      }
      selectedOptions = selectedOptions;
    } else {
      // Single choice: can change vote to different option
      selectedOptions = [optionId];
    }
  }

  async function submitVote() {
    if (selectedOptions.length === 0) return;
    if (!currentUser) {
      alert("Please login to vote on polls");
      return;
    }

    try {
      const optionId = selectedOptions[0];
      const hasVotedBefore =
        post.content.poll?.user_vote_ids &&
        post.content.poll.user_vote_ids.length > 0;

      // If single choice and already voted, remove old vote first
      if (!post.content.poll?.allow_multiple && hasVotedBefore) {
        await removePollVote(post.id);
        // Backend returns updated poll, but we'll fetch fresh data
      }

      // Submit new vote
      const updatedPoll = await voteOnPoll(post.id, optionId);

      // Update poll with fresh data from backend
      if (post.content.poll && updatedPoll) {
        post.content.poll.options = updatedPoll.options;
        post.content.poll.total_votes = updatedPoll.total_votes;
        post.content.poll.user_vote_ids = updatedPoll.user_vote_ids || [];
      }

      // Clear selection
      selectedOptions = [];
    } catch (error) {
      console.error("Failed to vote on poll:", error);
      alert("Failed to submit vote. Please try again.");
    }
  }

  function handlePollOptionClick(e: MouseEvent, optionId: string) {
    e.stopPropagation();
    handleVote(optionId);
  }

  function handleVoteSubmit(e: MouseEvent) {
    e.stopPropagation();
    submitVote();
  }

  async function handleUnvote(e: MouseEvent) {
    e.stopPropagation();
    if (!currentUser) {
      alert("Please login to manage votes");
      return;
    }

    try {
      const updatedPoll = await removePollVote(post.id);

      // Update poll with fresh data from backend
      if (post.content.poll && updatedPoll) {
        post.content.poll.options = updatedPoll.options;
        post.content.poll.total_votes = updatedPoll.total_votes;
        post.content.poll.user_vote_ids = updatedPoll.user_vote_ids || [];
      }
    } catch (error) {
      console.error("Failed to remove poll vote:", error);
      alert("Failed to remove vote. Please try again.");
    }
  }

  const getVotePercentage = (votes: number, total: number) => {
    if (total === 0) return 0;
    return (votes / total) * 100;
  };

  async function handleUpvote(e: MouseEvent) {
    e.stopPropagation();
    if (!currentUser) {
      alert("Please login to vote");
      return;
    }
    if (isOwnPost) {
      alert("You cannot vote on your own post");
      return;
    }
    if (isVoting) return;

    try {
      isVoting = true;
      const newVote = userVote === "up" ? null : true; // Toggle or set upvote

      await voteOnPost(post.id, newVote!);

      // Update local state
      if (userVote === "up") {
        // Remove upvote
        userVote = "";
        votesCount--;
      } else if (userVote === "down") {
        // Change from downvote to upvote
        userVote = "up";
        votesCount += 2; // Remove -1 and add +1
      } else {
        // Add upvote
        userVote = "up";
        votesCount++;
      }
    } catch (error) {
      console.error("Failed to vote:", error);
    } finally {
      isVoting = false;
    }
  }

  async function handleDownvote(e: MouseEvent) {
    e.stopPropagation();
    if (!currentUser) {
      alert("Please login to vote");
      return;
    }
    if (isOwnPost) {
      alert("You cannot vote on your own post");
      return;
    }
    if (isVoting) return;

    try {
      isVoting = true;
      const newVote = userVote === "down" ? null : false; // Toggle or set downvote

      await voteOnPost(post.id, newVote!);

      // Update local state
      if (userVote === "down") {
        // Remove downvote
        userVote = "";
        votesCount++;
      } else if (userVote === "up") {
        // Change from upvote to downvote
        userVote = "down";
        votesCount -= 2; // Remove +1 and add -1
      } else {
        // Add downvote
        userVote = "down";
        votesCount--;
      }
    } catch (error) {
      console.error("Failed to vote:", error);
    } finally {
      isVoting = false;
    }
  }

  async function handleSave(e: MouseEvent) {
    e.stopPropagation();
    if (!currentUser) {
      alert("Please login to save posts");
      return;
    }
    if (isSaving) return;

    try {
      isSaving = true;
      if (isSaved) {
        await unsavePost(post.id);
        isSaved = false;
      } else {
        await savePost(post.id);
        isSaved = true;
      }
    } catch (error) {
      console.error("Failed to save post:", error);
    } finally {
      isSaving = false;
    }
  }

  async function handleHide(e: MouseEvent) {
    e.stopPropagation();
    if (!currentUser) {
      alert("Please login to hide posts");
      return;
    }

    try {
      await hidePost(post.id);
      if (onUpdate) onUpdate(); // Refresh parent to remove hidden post
    } catch (error) {
      console.error("Failed to hide post:", error);
    }
  }

  function handleShare(e: MouseEvent) {
    e.stopPropagation();
    const url = `${window.location.origin}/#/post/${post.id}`;

    if (navigator.share) {
      navigator
        .share({
          title: post.title,
          url: url,
        })
        .catch(() => {
          // Fallback to copy
          copyToClipboard(url);
        });
    } else {
      copyToClipboard(url);
    }
  }

  function copyToClipboard(text: string) {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        alert("Link copied to clipboard!");
      })
      .catch(() => {
        alert("Failed to copy link");
      });
  }

  function toggleMenu(e: MouseEvent) {
    e.stopPropagation();
    showMenu = !showMenu;
  }

  function handleEdit(e: MouseEvent) {
    e.stopPropagation();
    showMenu = false;
    push(`/post/${post.id}/edit`);
  }

  async function handleDelete(e: MouseEvent) {
    e.stopPropagation();
    showMenu = false;

    if (!confirm("Are you sure you want to delete this post?")) return;

    try {
      await deletePost(post.id);
      if (onUpdate) onUpdate();
    } catch (error) {
      console.error("Failed to delete post:", error);
      alert("Failed to delete post. Please try again.");
    }
  }

  function openReportModal(e: MouseEvent) {
    e.stopPropagation();
    showMenu = false;
    showReportModal = true;
  }

  function closeReportModal() {
    showReportModal = false;
    reportReason = "";
    reportDetails = "";
  }

  async function handleReport(e: Event) {
    e.preventDefault();
    if (!reportReason.trim()) {
      alert("Please select a reason");
      return;
    }

    try {
      isReporting = true;
      await reportPost(post.id, {
        reason: reportReason,
        details: reportDetails,
      });
      alert("Report submitted successfully");
      closeReportModal();
    } catch (error) {
      console.error("Failed to report post:", error);
      alert("Failed to submit report. Please try again.");
    } finally {
      isReporting = false;
    }
  }
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
        <div class="menu-container">
          <button class="more-btn" onclick={toggleMenu} title="More options">
            <img src="/dot.png" alt="" width="20" height="20" />
          </button>
          {#if showMenu}
            <div class="dropdown-menu" transition:fade={{ duration: 150 }}>
              {#if isOwnPost}
                <button class="menu-item" onclick={handleEdit}>
                  <img
                    src="/write_icon.svg"
                    alt="Edit"
                    width="16"
                    height="16"
                  />
                  <span>Edit Post</span>
                </button>
                <button class="menu-item delete" onclick={handleDelete}>
                  <img
                    src="/delete_icon.svg"
                    alt="Delete"
                    width="16"
                    height="16"
                  />
                  <span>Delete Post</span>
                </button>
              {:else}
                <button class="menu-item" onclick={openReportModal}>
                  <span>⚠️ Report</span>
                </button>
                <button class="menu-item" onclick={handleHide}>
                  <span>🚫 Hide</span>
                </button>
              {/if}
            </div>
          {/if}
        </div>
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
      {#if post.content.videos && post.content.videos.length > 0}
        <video
          controls
          poster={post.content.videos[0].thumbnail}
          class="post-video"
        >
          <source src={post.content.videos[0].url} type="video/mp4" />
          <track kind="captions" />
          Your browser does not support the video tag.
        </video>
      {/if}
      {#if post.type === "poll" && post.content.poll}
        <div class="poll-container">
          <h3 class="poll-question">{post.content.poll.question}</h3>
          <div class="poll-options">
            {#each post.content.poll.options as option}
              {@const isVotedOption = post.content.poll.user_vote_ids?.includes(
                option.id
              )}
              <button
                class="poll-option"
                class:selected={selectedOptions.includes(option.id)}
                class:voted={isVotedOption}
                onclick={(e) => handlePollOptionClick(e, option.id)}
                disabled={hasVoted && !post.content.poll.allow_multiple}
              >
                <div
                  class="poll-result-bar"
                  style="width: {option.percentage}%;"
                ></div>
                <div class="poll-option-content">
                  {#if !hasVoted}
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
                  {/if}
                  <span class="poll-option-text">{option.text}</span>
                  <span class="poll-option-percentage"
                    >{option.percentage.toFixed(1)}%</span
                  >
                </div>
              </button>
            {/each}
          </div>
          <div class="poll-actions">
            {#if !hasVoted}
              <button
                class="vote-submit-btn"
                onclick={handleVoteSubmit}
                disabled={selectedOptions.length === 0}
              >
                Vote
              </button>
            {:else}
              <button class="unvote-btn" onclick={handleUnvote}>
                Remove Vote
              </button>
            {/if}
          </div>
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
          class:voted={userVote === "up"}
          aria-label="Upvote"
          disabled={isVoting || isOwnPost}
          onclick={handleUpvote}
        >
          ▲
        </button>
        <span class="vote-count">{votesCount}</span>
        <button
          class="footer-btn vote-btn"
          class:voted={userVote === "down"}
          aria-label="Downvote"
          disabled={isVoting || isOwnPost}
          onclick={handleDownvote}
        >
          ▼
        </button>
      </div>
      <button class="footer-btn" onclick={handlePostClick}>
        <img src="/CommentIcon.svg" alt="Comments" width="20" height="20" />
        <span>{post.comments_count} Comments</span>
      </button>
      <button class="footer-btn" onclick={handleShare}>
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
      <button
        class="footer-btn"
        class:saved={isSaved}
        disabled={isSaving}
        onclick={handleSave}
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill={isSaved ? "currentColor" : "none"}
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          ><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"
          ></path></svg
        >
        <span>{isSaved ? "Saved" : "Save"}</span>
      </button>
    </div>
  </div>
</article>

<!-- Report Modal -->
{#if showReportModal}
  <div
    class="modal-overlay"
    onclick={closeReportModal}
    transition:fade={{ duration: 200 }}
  >
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Report Post</h3>
        <button class="close-btn" onclick={closeReportModal}>×</button>
      </div>
      <form onsubmit={handleReport}>
        <div class="form-group">
          <label for="report-reason">Reason *</label>
          <select id="report-reason" bind:value={reportReason} required>
            <option value="">Select a reason</option>
            <option value="spam">Spam</option>
            <option value="harassment">Harassment or Bullying</option>
            <option value="hate">Hate Speech</option>
            <option value="violence">Violence or Threat</option>
            <option value="misinformation">Misinformation</option>
            <option value="nsfw">NSFW Content</option>
            <option value="copyright">Copyright Violation</option>
            <option value="other">Other</option>
          </select>
        </div>
        <div class="form-group">
          <label for="report-details">Additional Details (Optional)</label>
          <textarea
            id="report-details"
            bind:value={reportDetails}
            placeholder="Provide more context about why you're reporting this post..."
            rows="4"
          ></textarea>
        </div>
        <div class="modal-actions">
          <button type="button" class="btn-cancel" onclick={closeReportModal}>
            Cancel
          </button>
          <button type="submit" class="btn-submit" disabled={isReporting}>
            {isReporting ? "Submitting..." : "Submit Report"}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

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
    font-weight: 600;
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
    padding: 0;
    border: 1px solid #ccc;
    border-radius: 4px;
    background-color: white;
    color: #000000;
    cursor: pointer;
    text-align: left;
    overflow: hidden;
    min-height: 40px;
  }
  .poll-option:not(:disabled):hover {
    border-color: #878a8c;
  }
  .poll-option.selected {
    border-color: var(--primary-color);
    background-color: #f0f8ff;
  }
  .poll-option.voted {
    border-color: var(--primary-color);
    border-width: 2px;
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
    opacity: 0.15;
    transition: width 0.5s ease-in-out;
  }

  .poll-option-content {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: center;
    width: 100%;
    padding: 8px 12px;
    gap: 8px;
  }

  .poll-option-text {
    flex-grow: 1;
  }

  .poll-option-stats {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
  }

  .poll-option-percentage {
    font-weight: 600;
    color: var(--primary-color);
  }

  .poll-option-votes {
    color: #666;
  }

  .poll-actions {
    margin-top: 12px;
    display: flex;
    gap: 8px;
  }

  .vote-submit-btn,
  .unvote-btn {
    padding: 8px 24px;
    border: none;
    border-radius: 20px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .vote-submit-btn {
    background-color: var(--primary-color);
    color: white;
  }

  .vote-submit-btn:hover:not(:disabled) {
    background-color: var(--primary-color-dark);
  }

  .vote-submit-btn:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }

  .unvote-btn {
    background-color: #f0f0f0;
    color: #333;
    border: 1px solid #ddd;
  }

  .unvote-btn:hover {
    background-color: #e0e0e0;
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

  .poll-footer {
    font-size: 12px;
    color: #878a8c;
    margin-top: 10px;
  }

  /* Vote button states */
  .vote-btn.voted {
    color: var(--blue--);
    font-weight: bold;
  }

  .vote-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Save button state */
  .footer-btn.saved {
    color: var(--blue--);
    font-weight: 600;
  }

  .footer-btn.saved svg {
    fill: var(--blue--);
  }

  /* Menu Dropdown */
  .menu-container {
    position: relative;
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 4px;
    background: white;
    border: 1px solid #ccc;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    min-width: 150px;
    z-index: 1000;
  }

  .menu-item {
    width: 100%;
    padding: 10px 16px;
    border: none;
    background: white;
    text-align: left;
    cursor: pointer;
    font-size: 14px;
    transition: background-color 0.2s;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .menu-item:hover {
    background-color: #f6f7f8;
  }

  .menu-item.delete {
    color: #d93025;
  }

  .menu-item.delete:hover {
    background-color: #fef1f0;
  }

  /* Report Modal */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
  }

  .modal-content {
    background: white;
    border-radius: 8px;
    width: 90%;
    max-width: 500px;
    max-height: 80vh;
    overflow-y: auto;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid #eee;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 20px;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 28px;
    cursor: pointer;
    color: #878a8c;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .close-btn:hover {
    color: #000;
  }

  .modal-content form {
    padding: 20px;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 14px;
    font-family: inherit;
  }

  .form-group textarea {
    resize: vertical;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }

  .btn-cancel,
  .btn-submit {
    padding: 10px 20px;
    border-radius: 20px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    border: none;
  }

  .btn-cancel {
    background: #f6f7f8;
    color: #1c1c1c;
  }

  .btn-cancel:hover {
    background: #e9ebed;
  }

  .btn-submit {
    background: var(--blue--);
    color: white;
  }

  .btn-submit:hover:not(:disabled) {
    background: var(--darkblue--);
  }

  .btn-submit:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
