<script lang="ts">
  import type { Comment } from "../mocks/comments.mock";
  import CommentComponent from "./Comment.svelte";

  type CommentProps = {
    comment: Comment;
    depth?: number;
  };

  let { comment, depth = 0 }: CommentProps = $props();

  let isCollapsed = $state(false);
  let isEditing = $state(false);
  let editContent = $state(comment.content);
  let showReplyBox = $state(false);
  let replyContent = $state("");
  let userVote = $state<"up" | "down" | null>(null); // Track user's vote

  const toggleCollapse = () => {
    isCollapsed = !isCollapsed;
  };

  const handleEdit = () => {
    isEditing = true;
    editContent = comment.content;
  };

  const saveEdit = () => {
    comment.content = editContent;
    isEditing = false;
  };

  const cancelEdit = () => {
    isEditing = false;
    editContent = comment.content;
  };

  const handleDelete = () => {
    if (confirm("Are you sure you want to delete this comment?")) {
      // Mock delete - in real app, would update parent state
      console.log("Deleting comment:", comment.id);
    }
  };

  const handleVote = (voteType: "up" | "down") => {
    if (userVote === voteType) {
      // Remove vote
      if (voteType === "up") {
        comment.votesCount.up--;
      } else {
        comment.votesCount.down--;
      }
      userVote = null;
    } else {
      // Change or add vote
      if (userVote === "up") {
        comment.votesCount.up--;
      } else if (userVote === "down") {
        comment.votesCount.down--;
      }

      if (voteType === "up") {
        comment.votesCount.up++;
      } else {
        comment.votesCount.down++;
      }
      userVote = voteType;
    }
  };

  const submitReply = () => {
    if (replyContent.trim()) {
      // Mock reply - in real app, would add to parent state
      console.log("Replying to comment:", comment.id, replyContent);
      replyContent = "";
      showReplyBox = false;
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    return `${diffDays}d ago`;
  };

  const getScore = () => {
    return comment.votesCount.up - comment.votesCount.down;
  };
</script>

<div class="comment" style="margin-left: {depth * 28}px">
  <div class="comment-main">
    <!-- Vote Section -->
    <div class="vote-section">
      <button
        class="vote-btn"
        class:voted={userVote === "up"}
        onclick={() => handleVote("up")}
        aria-label="Upvote"
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 3l5 7H3l5-7z" />
        </svg>
      </button>
      <span
        class="vote-count"
        class:positive={getScore() > 0}
        class:negative={getScore() < 0}
      >
        {getScore()}
      </span>
      <button
        class="vote-btn"
        class:voted={userVote === "down"}
        onclick={() => handleVote("down")}
        aria-label="Downvote"
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 13L3 6h10l-5 7z" />
        </svg>
      </button>
    </div>

    <!-- Comment Content -->
    <div class="comment-content">
      <!-- Header -->
      <div class="comment-header">
        <img
          src={comment.authorAvatar || "https://i.pravatar.cc/150?img=1"}
          alt={comment.authorUsername}
          class="author-avatar"
        />
        <span class="author-name">u/{comment.authorUsername}</span>
        <span class="comment-time">{formatTime(comment.createdAt)}</span>
        {#if comment.replies.length > 0}
          <button class="collapse-btn" onclick={toggleCollapse}>
            {isCollapsed ? "[+]" : "[-]"}
          </button>
        {/if}
      </div>

      {#if !isCollapsed}
        <!-- Body -->
        {#if isEditing}
          <div class="edit-section">
            <textarea bind:value={editContent} class="edit-textarea" rows="4"
            ></textarea>
            <div class="edit-actions">
              <button class="save-btn" onclick={saveEdit}>Save</button>
              <button class="cancel-btn" onclick={cancelEdit}>Cancel</button>
            </div>
          </div>
        {:else}
          <div class="comment-text">
            {comment.content}
          </div>
        {/if}

        <!-- Actions -->
        {#if !isEditing}
          <div class="comment-actions">
            <button
              class="action-btn"
              onclick={() => (showReplyBox = !showReplyBox)}
            >
              Reply
            </button>
            <button class="action-btn" onclick={handleEdit}> Edit </button>
            <button class="action-btn delete" onclick={handleDelete}>
              Delete
            </button>
          </div>
        {/if}

        <!-- Reply Box -->
        {#if showReplyBox}
          <div class="reply-box">
            <textarea
              bind:value={replyContent}
              placeholder="What are your thoughts?"
              class="reply-textarea"
              rows="3"
            ></textarea>
            <div class="reply-actions">
              <button class="submit-btn" onclick={submitReply}>Comment</button>
              <button
                class="cancel-btn"
                onclick={() => {
                  showReplyBox = false;
                  replyContent = "";
                }}
              >
                Cancel
              </button>
            </div>
          </div>
        {/if}

        <!-- Replies -->
        {#if comment.replies.length > 0}
          <div class="replies">
            {#each comment.replies as reply}
              <CommentComponent comment={reply} depth={depth + 1} />
            {/each}
          </div>
        {/if}
      {/if}
    </div>
  </div>
</div>

<style>
  .comment {
    margin-bottom: 8px;
  }

  .comment-main {
    display: flex;
    gap: 8px;
  }

  /* Vote Section */
  .vote-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding-top: 4px;
  }

  .vote-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    color: #878a8c;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 2px;
    transition: all 0.2s;
  }

  .vote-btn:hover {
    background: rgba(135, 138, 140, 0.1);
  }

  .vote-btn.voted {
    color: #ff4500;
  }

  .vote-count {
    font-size: 12px;
    font-weight: 700;
    color: #1c1c1c;
    min-width: 24px;
    text-align: center;
  }

  .vote-count.positive {
    color: #ff4500;
  }

  .vote-count.negative {
    color: #7193ff;
  }

  /* Comment Content */
  .comment-content {
    flex: 1;
    min-width: 0;
  }

  .comment-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
    font-size: 12px;
  }

  .author-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    object-fit: cover;
  }

  .author-name {
    font-weight: 700;
    color: #1c1c1c;
  }

  .author-name:hover {
    text-decoration: underline;
    cursor: pointer;
  }

  .comment-time {
    color: #7c7c7c;
  }

  .collapse-btn {
    background: none;
    border: none;
    cursor: pointer;
    color: #878a8c;
    font-size: 12px;
    font-weight: 700;
    padding: 2px 6px;
    margin-left: auto;
  }

  .collapse-btn:hover {
    background: rgba(135, 138, 140, 0.1);
    border-radius: 2px;
  }

  .comment-text {
    color: #1c1c1c;
    font-size: 14px;
    line-height: 21px;
    margin-bottom: 8px;
    word-wrap: break-word;
  }

  /* Actions */
  .comment-actions {
    display: flex;
    gap: 12px;
    margin-bottom: 8px;
  }

  .action-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 12px;
    font-weight: 700;
    color: #878a8c;
    padding: 4px 8px;
    border-radius: 2px;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s;
  }

  .action-btn:hover {
    background: rgba(135, 138, 140, 0.1);
  }

  .action-btn.delete {
    color: #ea0027;
  }

  /* Edit Section */
  .edit-section {
    margin-bottom: 8px;
  }

  .edit-textarea {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
    font-family: "Roboto", sans-serif;
    resize: vertical;
    box-sizing: border-box;
    margin-bottom: 8px;
  }

  .edit-textarea:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .edit-actions {
    display: flex;
    gap: 8px;
  }

  .save-btn {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 6px 16px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .save-btn:hover {
    background: var(--darkblue--);
  }

  .cancel-btn {
    background: none;
    color: var(--blue--);
    border: 1px solid var(--blue--);
    padding: 6px 16px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s;
  }

  .cancel-btn:hover {
    background: rgba(21, 48, 96, 0.1);
  }

  /* Reply Box */
  .reply-box {
    margin-top: 8px;
    margin-bottom: 12px;
  }

  .reply-textarea {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
    font-family: "Roboto", sans-serif;
    resize: vertical;
    box-sizing: border-box;
    margin-bottom: 8px;
  }

  .reply-textarea:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .reply-actions {
    display: flex;
    gap: 8px;
  }

  .submit-btn {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 6px 16px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .submit-btn:hover {
    background: var(--darkblue--);
  }

  /* Replies */
  .replies {
    margin-top: 12px;
  }
</style>
