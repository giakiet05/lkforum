<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import CommentSection from "../components/CommentSection.svelte";
  import type { PostResponse } from "../dtos/post-dto";
  import {
    voteOnPost,
    savePost,
    unsavePost,
    getPostById,
    voteOnPoll,
    removePollVote,
    deletePost,
    reportPost,
  } from "../services/post-service";
  import { authStore } from "../stores/auth-store";

  type PostDetailProps = {
    params?: { id: string };
  };

  let { params = { id: "1" } }: PostDetailProps = $props();

  let post = $state<PostResponse | null>(null);
  let selectedOptions = $state<string[]>([]);
  let hasVoted = $state(false);

  // Vote state
  let userVote = $state<"up" | "down" | "">("");
  let votesCount = $state(0);
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
  const isOwnPost = $derived(
    currentUser && post && post.author.id === currentUser.id
  );

  async function handleUpvote() {
    if (!currentUser) {
      alert("Please login to vote");
      return;
    }
    if (isOwnPost) {
      alert("You cannot vote on your own post");
      return;
    }
    if (isVoting || !post) return;

    try {
      isVoting = true;
      const newVote = userVote === "up" ? null : true;
      await voteOnPost(post.id, newVote!);

      if (userVote === "up") {
        userVote = "";
        votesCount--;
      } else if (userVote === "down") {
        userVote = "up";
        votesCount += 2;
      } else {
        userVote = "up";
        votesCount++;
      }
    } catch (error) {
      console.error("Failed to vote:", error);
    } finally {
      isVoting = false;
    }
  }

  async function handleDownvote() {
    if (!currentUser) {
      alert("Please login to vote");
      return;
    }
    if (isOwnPost) {
      alert("You cannot vote on your own post");
      return;
    }
    if (isVoting || !post) return;

    try {
      isVoting = true;
      const newVote = userVote === "down" ? null : false;
      await voteOnPost(post.id, newVote!);

      if (userVote === "down") {
        userVote = "";
        votesCount++;
      } else if (userVote === "up") {
        userVote = "down";
        votesCount -= 2;
      } else {
        userVote = "down";
        votesCount--;
      }
    } catch (error) {
      console.error("Failed to vote:", error);
    } finally {
      isVoting = false;
    }
  }

  async function handleSave() {
    if (!currentUser) {
      alert("Please login to save posts");
      return;
    }
    if (isSaving || !post) return;

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

  onMount(async () => {
    // Scroll to top when component mounts
    window.scrollTo(0, 0);

    // Fetch post from API
    try {
      const fetchedPost = await getPostById(params.id);
      post = fetchedPost;

      // Initialize vote state from post data
      userVote = (post.user_vote as "up" | "down" | "") || "";
      votesCount = post.votes_count?.score || 0;
    } catch (error) {
      console.error("Failed to fetch post:", error);
    }
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

  async function submitVote() {
    if (!post || selectedOptions.length === 0) return;
    if (!currentUser) {
      alert("Please login to vote on polls");
      return;
    }

    try {
      const optionId = selectedOptions[0];
      await voteOnPoll(post.id, optionId);

      // Update local state
      if (post.content.poll) {
        for (const option of post.content.poll.options) {
          if (option.id === optionId) {
            option.votes++;
            option.percentage =
              (option.votes / (post.content.poll.total_votes + 1)) * 100;
          } else {
            option.percentage =
              (option.votes / (post.content.poll.total_votes + 1)) * 100;
          }
        }
        post.content.poll.total_votes++;
        if (!post.content.poll.user_vote_ids) {
          post.content.poll.user_vote_ids = [];
        }
        post.content.poll.user_vote_ids.push(optionId);
      }
      hasVoted = true;
    } catch (error) {
      console.error("Failed to vote on poll:", error);
      alert("Failed to submit vote. Please try again.");
    }
  }

  const getVotePercentage = (votes: number, total: number) => {
    if (total === 0) return 0;
    return (votes / total) * 100;
  };

  function handleShare() {
    if (!post) return;
    const url = `${window.location.origin}/#/post/${post.id}`;

    if (navigator.share) {
      navigator
        .share({
          title: post.title,
          url: url,
        })
        .catch(() => {
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

  function toggleMenu() {
    showMenu = !showMenu;
  }

  function handleEdit() {
    showMenu = false;
    push(`/post/${post!.id}/edit`);
  }

  async function handleDelete() {
    showMenu = false;
    if (!post || !confirm("Are you sure you want to delete this post?")) return;

    try {
      await deletePost(post.id);
      push("/"); // Redirect to home after delete
    } catch (error) {
      console.error("Failed to delete post:", error);
      alert("Failed to delete post. Please try again.");
    }
  }

  function openReportModal() {
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
    if (!post || !reportReason.trim()) {
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
          <div class="post-header-left">
            <span class="community-name">lk/{post.community.name}</span>
            <span class="meta-divider">•</span>
            <span class="author">Posted by u/{post.author.username}</span>
            <span class="time"
              >{new Date(post.created_at).toLocaleDateString()}</span
            >
          </div>
          <div class="post-header-right">
            <div class="menu-container">
              <button
                class="more-btn"
                onclick={toggleMenu}
                title="More options"
              >
                •••
              </button>
              {#if showMenu}
                <div class="dropdown-menu">
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
                  {/if}
                </div>
              {/if}
            </div>
          </div>
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
          {:else if post.content.videos && post.content.videos.length > 0}
            <video
              controls
              poster={post.content.videos[0].thumbnail_url}
              class="post-video"
            >
              <source src={post.content.videos[0].url} type="video/mp4" />
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
          <button class="footer-btn">
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
            >
              <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path>
              <polyline points="16 6 12 2 8 6"></polyline>
              <line x1="12" y1="2" x2="12" y2="15"></line>
            </svg>
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
            >
              <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"
              ></path>
            </svg>
            <span>{isSaved ? "Saved" : "Save"}</span>
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

<!-- Report Modal -->
{#if showReportModal}
  <div class="modal-overlay" onclick={closeReportModal}>
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
    justify-content: space-between;
    font-size: 12px;
    margin-bottom: 12px;
  }

  .post-header-left {
    display: flex;
    align-items: center;
  }

  .post-header-right {
    display: flex;
    align-items: center;
    gap: 8px;
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

  /* Menu Dropdown */
  .menu-container {
    position: relative;
  }

  .more-btn {
    background: none;
    border: none;
    font-size: 20px;
    font-weight: bold;
    color: #878a8c;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
  }

  .more-btn:hover {
    background: #f6f7f8;
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

  .vote-btn.voted {
    color: var(--blue--);
    font-weight: bold;
  }

  .footer-btn.saved {
    color: var(--blue--);
    font-weight: 600;
  }

  .footer-btn.saved svg {
    fill: var(--blue--);
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
