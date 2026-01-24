<script lang="ts">
  import { push } from "svelte-spa-router";
  import type { CommentResponse } from "../dtos/comment-dto";
  import CommentComponent from "./Comment.svelte";
  import { deleteComment, createComment } from "../services/comment-service";
  import { authStore } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
  import ConfirmModal from "./ConfirmModal.svelte";

  type CommentProps = {
    comment: CommentResponse;
    depth?: number;
    onUpdate?: () => void;
  };

  let { comment, depth = 0, onUpdate }: CommentProps = $props();

  const currentUser = $derived($authStore.user);
  const isOwnComment = $derived(
    currentUser && comment.author.id === currentUser.id,
  );

  let isCollapsed = $state(false);
  let isEditing = $state(false);
  let editContent = $state(comment.content);
  let showReplyBox = $state(false);
  let replyContent = $state("");
  let userVote = $state<"up" | "down" | null>(null); // Track user's vote
  let replyImage = $state<File | null>(null);
  let replyImagePreview = $state<string | null>(null);
  let replyErrorMessage = $state<string | null>(null);
  let showDeleteConfirm = $state(false);

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
    showDeleteConfirm = true;
  };

  const confirmDelete = async () => {
    showDeleteConfirm = false;
    try {
      await deleteComment(comment.id);
      if (onUpdate) onUpdate();
    } catch (error) {
      console.error("Failed to delete comment:", error);
      toastStore.error("Không thể xóa bình luận. Vui lòng thử lại.");
    }
  };

  const handleVote = (voteType: "up" | "down") => {
    // TODO: Implement vote API when backend adds comment voting
    console.log("Vote not implemented yet");
  };

  const submitReply = async () => {
    if (!replyContent.trim()) return;
    if (!currentUser) {
      replyErrorMessage = "Vui lòng đăng nhập để trả lời";
      return;
    }

    try {
      replyErrorMessage = null;
      await createComment({
        post_id: comment.post_id,
        parent_id: comment.id,
        content: replyContent,
      });

      // Clear reply form
      replyContent = "";
      replyImage = null;
      replyImagePreview = null;
      showReplyBox = false;

      // Auto-expand to show the new reply
      isCollapsed = false;

      // Reload comments to show the new reply
      if (onUpdate) {
        onUpdate();
      }
    } catch (error: any) {
      console.error("Failed to submit reply:", error);
      // Show specific error message from backend
      replyErrorMessage =
        error?.message || "Không thể gửi trả lời. Vui lòng thử lại.";
    }
  };

  const handleReplyImageSelect = (event: Event) => {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      replyImage = file;

      const reader = new FileReader();
      reader.onload = (e) => {
        replyImagePreview = e.target?.result as string;
      };
      reader.readAsDataURL(file);
    }
  };

  const removeReplyImage = () => {
    replyImage = null;
    replyImagePreview = null;
  };

  const formatTime = (dateString: string) => {
    // Backend sends: "2025-12-07 10:54:03" (local server time without timezone)
    // Parse it as UTC to avoid timezone conversion issues
    const date = new Date(dateString.replace(" ", "T") + "Z");
    const now = new Date();

    const diffMs = now.getTime() - date.getTime();
    const diffSecs = Math.floor(diffMs / 1000);
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    // Handle future dates or just posted (< 1 second)
    if (diffSecs < 1) return `vừa xong`;
    if (diffSecs < 60) return `${diffSecs} giây trước`;
    if (diffMins < 60) return `${diffMins} phút trước`;
    if (diffHours < 24) return `${diffHours} giờ trước`;
    if (diffDays < 30) return `${diffDays} ngày trước`;

    // For older comments, show the actual date
    return date.toLocaleDateString();
  };
</script>

<div class="comment" style="margin-left: {depth * 28}px">
  <div class="comment-main">
    <!-- Vote Section - Hidden until backend supports comment voting -->
    <!-- <div class="vote-section">
      <button
        class="vote-btn"
        class:voted={userVote === "up"}
        onclick={() => handleVote("up")}
        aria-label="Bình chọn tăng"
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
        aria-label="Bình chọn giảm"
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 13L3 6h10l-5 7z" />
        </svg>
      </button>
    </div> -->

    <!-- Comment Content -->
    <div class="comment-content">
      <!-- Header -->
      <div class="comment-header">
        <img
          src={comment.author.avatar?.url || "/user.jpg"}
          alt={comment.author.username}
          class="author-avatar clickable"
          onclick={() => push(`/profile/${comment.author.username}`)}
          role="button"
          tabindex="0"
        />
        <span
          class="author-name clickable"
          onclick={() => push(`/profile/${comment.author.username}`)}
          role="button"
          tabindex="0"
        >
          u/{comment.author.username}
        </span>
        <span class="comment-time">{formatTime(comment.created_at)}</span>
        {#if comment.children && comment.children.length > 0}
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
              <button class="save-btn" onclick={saveEdit}>Lưu</button>
              <button class="cancel-btn" onclick={cancelEdit}>Hủy</button>
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
            <!-- Hide Reply button at max depth (backend limit is depth 0-2) -->
            {#if depth < 2}
              <button
                class="action-btn"
                onclick={() => (showReplyBox = !showReplyBox)}
              >
                Trả lời
              </button>
            {/if}
            {#if isOwnComment}
              <button class="action-btn" onclick={handleEdit}> Sửa </button>
              <button class="action-btn delete" onclick={handleDelete}>
                Xóa
              </button>
            {/if}
          </div>
        {/if}

        <!-- Reply Box -->
        {#if showReplyBox}
          <div class="reply-box">
            {#if replyErrorMessage}
              <div class="reply-error-banner">
                <img src="/error.svg" alt="Error" width="16" height="16" />
                <span>{replyErrorMessage}</span>
                <button
                  class="error-close-btn"
                  onclick={() => (replyErrorMessage = null)}>×</button
                >
              </div>
            {/if}
            <textarea
              bind:value={replyContent}
              placeholder="Bạn nghĩ gì?"
              class="reply-textarea"
              rows="3"
              onkeydown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  submitReply();
                }
              }}
            ></textarea>

            {#if replyImagePreview}
              <div class="reply-image-preview">
                <img src={replyImagePreview} alt="Preview" />
                <button class="remove-reply-image" onclick={removeReplyImage}>
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                  >
                    <path
                      d="M4 4l8 8M12 4l-8 8"
                      stroke="currentColor"
                      stroke-width="2"
                    />
                  </svg>
                </button>
              </div>
            {/if}

            <div class="reply-actions">
              <div class="reply-attachment-buttons">
                <input
                  type="file"
                  id="reply-image-upload-{comment.id}"
                  accept="image/*"
                  onchange={handleReplyImageSelect}
                  style="display: none;"
                />
                <button
                  class="attachment-btn"
                  onclick={() =>
                    document
                      .getElementById(`reply-image-upload-{comment.id}`)
                      ?.click()}
                  title="Thêm ảnh"
                >
                  <img
                    src="/icon_comment.png"
                    alt="Add comment icon"
                    width="20"
                    height="20"
                  />
                </button>
                <button
                  class="attachment-btn"
                  onclick={() =>
                    document
                      .getElementById(`reply-image-upload-{comment.id}`)
                      ?.click()}
                  title="Thêm ảnh"
                >
                  <img
                    src="/comment_picture.png"
                    alt="Thêm ảnh"
                    width="20"
                    height="20"
                  />
                </button>
              </div>
              <div class="reply-action-buttons">
                <button class="submit-btn" onclick={submitReply}
                  >Bình luận</button
                >
                <button
                  class="cancel-btn"
                  onclick={() => {
                    showReplyBox = false;
                    replyContent = "";
                    replyImage = null;
                    replyImagePreview = null;
                  }}
                >
                  Hủy
                </button>
              </div>
            </div>
          </div>
        {/if}

        <!-- Replies/Children -->
        {#if comment.children && comment.children.length > 0}
          <div class="replies">
            {#each comment.children as reply}
              <CommentComponent comment={reply} depth={depth + 1} {onUpdate} />
            {/each}
          </div>
        {/if}
      {/if}
    </div>
  </div>
</div>

<ConfirmModal
  show={showDeleteConfirm}
  title="Xác nhận xóa"
  message="Bạn có chắc muốn xóa bình luận này? Hành động này không thể hoàn tác."
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDelete}
  onCancel={() => (showDeleteConfirm = false)}
/>

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
    color: var(--darkblue--);
  }

  .vote-count {
    font-size: 12px;
    font-weight: 700;
    color: var(--darkblue--);
    min-width: 24px;
    text-align: center;
  }

  .vote-count.positive {
    color: var(--darkblue--);
  }

  .vote-count.negative {
    color: var(--darkblue--);
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

  .author-avatar.clickable {
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .author-avatar.clickable:hover {
    opacity: 0.8;
  }

  .author-name {
    font-weight: 700;
    color: #1c1c1c;
  }

  .author-name.clickable {
    cursor: pointer;
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

  .reply-error-banner {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    margin-bottom: 8px;
    background-color: #fee;
    border: 1px solid #fcc;
    border-radius: 4px;
    color: #c00;
    font-size: 13px;
    font-weight: 500;
  }

  .reply-error-banner svg {
    flex-shrink: 0;
  }

  .reply-error-banner span {
    flex: 1;
  }

  .error-close-btn {
    background: none;
    border: none;
    color: #c00;
    font-size: 20px;
    font-weight: 700;
    line-height: 1;
    cursor: pointer;
    padding: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 2px;
    transition: background 0.2s;
  }

  .error-close-btn:hover {
    background: rgba(204, 0, 0, 0.1);
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

  .reply-image-preview {
    position: relative;
    margin-bottom: 8px;
  }

  .reply-image-preview img {
    max-width: 200px;
    max-height: 200px;
    border-radius: 4px;
    display: block;
  }

  .remove-reply-image {
    position: absolute;
    top: 4px;
    right: 4px;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    border: none;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
  }

  .remove-reply-image:hover {
    background: rgba(0, 0, 0, 0.9);
  }

  .reply-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .reply-attachment-buttons {
    display: flex;
    gap: 8px;
  }

  .reply-action-buttons {
    display: flex;
    gap: 8px;
  }

  .attachment-btn {
    background: transparent;
    border: none;
    padding: 6px;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
    color: #7c7c7c;
  }

  .attachment-btn:hover {
    background: #f6f7f8;
  }

  .attachment-btn img {
    display: block;
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
