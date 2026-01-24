<script lang="ts">
  import { onMount } from "svelte";
  import type { CommentResponse } from "../dtos/comment-dto";
  import CommentComponent from "./Comment.svelte";
  import {
    getCommentsByPostId,
    createComment,
  } from "../services/comment-service";
  import { authStore } from "../stores/auth-store";

  type CommentSectionProps = {
    postId: string;
    onCommentAdded?: () => void;
    onTotalCommentsChange?: (total: number) => void;
  };

  let { postId, onCommentAdded, onTotalCommentsChange }: CommentSectionProps =
    $props();

  type SortType = "top" | "newest" | "oldest" | "controversial";

  let sortBy = $state<SortType>("top");
  let showSortDropdown = $state(false);
  let newCommentContent = $state("");
  let selectedImage = $state<File | null>(null);
  let imagePreview = $state<string | null>(null);
  let comments = $state<CommentResponse[]>([]);
  let isLoading = $state(true);
  let isSubmitting = $state(false);
  let totalComments = $state(0);
  let errorMessage = $state<string | null>(null);

  const currentUser = $derived($authStore.user);

  // Fetch comments on mount
  onMount(async () => {
    await loadComments();
  });

  async function loadComments() {
    try {
      isLoading = true;
      const response = await getCommentsByPostId({
        post_id: postId,
        page: 1,
        page_size: 50,
        depth: 2, // Backend allows max depth of 2
      });
      console.log("📥 loadComments response:", response);
      console.log(
        "📅 First comment created_at:",
        response.comments?.[0]?.created_at,
      );
      // Force reactivity by creating new array reference
      comments = [...(response.comments || [])];
      console.log("📝 comments array:", comments);
      totalComments = response.pagination?.total || comments.length;

      // Notify parent of total comments count
      if (onTotalCommentsChange) onTotalCommentsChange(totalComments);
    } catch (error) {
      console.error("Failed to load comments:", error);
    } finally {
      isLoading = false;
    }
  }

  // Sort comments
  const sortedComments = $derived(
    (() => {
      const commentsToSort = [...comments];
      console.log(
        "🔄 sortedComments recalculating, comments.length:",
        comments.length,
      );
      console.log("🔄 sortedComments array:", commentsToSort);

      switch (sortBy) {
        case "top":
          return commentsToSort;
        case "newest":
          return commentsToSort.sort(
            (a, b) =>
              new Date(b.created_at).getTime() -
              new Date(a.created_at).getTime(),
          );
        case "oldest":
          return commentsToSort.sort(
            (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          );
        case "controversial":
          return commentsToSort;
        default:
          return commentsToSort;
      }
    })(),
  );

  const handleSortChange = (sort: SortType) => {
    sortBy = sort;
    showSortDropdown = false;
  };

  const submitComment = async () => {
    if (!newCommentContent.trim()) return;
    if (!currentUser) {
      errorMessage = "Vui lòng đăng nhập để bình luận";
      return;
    }

    try {
      isSubmitting = true;
      errorMessage = null;
      const newComment = await createComment({
        post_id: postId,
        content: newCommentContent,
      });

      console.log("✅ Comment created:", newComment);
      console.log("📅 New comment created_at:", newComment.created_at);

      // Optimistic UI update: add the new comment immediately
      comments = [newComment, ...comments];
      totalComments += 1;
      console.log("✅ Comment added to UI, total:", comments.length);

      // Notify parent component
      if (onCommentAdded) onCommentAdded();

      newCommentContent = "";
      selectedImage = null;
      imagePreview = null;
    } catch (error: any) {
      console.error("❌ Failed to submit comment:", error);
      // Show specific error message from backend
      errorMessage =
        error?.message || "Không thể đăng bình luận. Vui lòng thử lại.";
    } finally {
      isSubmitting = false;
    }
  };

  const handleImageSelect = (event: Event) => {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      selectedImage = file;

      const reader = new FileReader();
      reader.onload = (e) => {
        imagePreview = e.target?.result as string;
      };
      reader.readAsDataURL(file);
    }
  };

  const removeImage = () => {
    selectedImage = null;
    imagePreview = null;
  };

  const getTotalComments = () => {
    return totalComments;
  };
</script>

<div class="comment-section">
  <!-- Error Message -->
  {#if errorMessage}
    <div class="error-banner">
      <img src="/error.svg" alt="Error" width="20" height="20" />
      <span>{errorMessage}</span>
      <button class="close-btn" onclick={() => (errorMessage = null)}>×</button>
    </div>
  {/if}

  <!-- Comment Input Box -->
  <div class="add-comment">
    <textarea
      bind:value={newCommentContent}
      placeholder="Bạn nghĩ gì?"
      class="comment-textarea"
      rows="4"
      onkeydown={(e) => {
        if (e.key === "Enter" && !e.shiftKey) {
          e.preventDefault();
          submitComment();
        }
      }}
    ></textarea>

    {#if imagePreview}
      <div class="image-preview">
        <img src={imagePreview} alt="Preview" />
        <button class="remove-image" onclick={removeImage} title="Xóa ảnh">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path
              d="M4 4l8 8M12 4l-8 8"
              stroke="currentColor"
              stroke-width="2"
            />
          </svg>
        </button>
      </div>
    {/if}

    <div class="comment-actions">
      <div class="attachment-buttons">
        <input
          type="file"
          id="comment-image-upload"
          accept="image/*"
          onchange={handleImageSelect}
          style="display: none;"
        />
        <button
          class="attachment-btn"
          onclick={() =>
            document.getElementById("comment-image-upload")?.click()}
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
            document.getElementById("comment-image-upload")?.click()}
          title="Thêm ảnh"
        >
          <img src="/comment_picture.png" alt="" width="20" height="20" />
        </button>
      </div>
      <button
        class="submit-btn"
        onclick={submitComment}
        disabled={isSubmitting}
      >
        {isSubmitting ? "Đang gửi..." : "Bình luận"}
      </button>
    </div>
  </div>

  <!-- Sort Bar -->
  <div class="sort-bar">
    <span class="comment-count">{getTotalComments()} Bình luận</span>
    <div class="sort-dropdown">
      <button
        class="sort-btn"
        onclick={() => (showSortDropdown = !showSortDropdown)}
      >
        <svg
          width="16"
          height="16"
          viewBox="0 0 16 16"
          fill="currentColor"
          class="sort-icon"
        >
          <path d="M3 4h10M3 8h7M3 12h4" />
        </svg>
        Sort by:
        <span class="sort-label"
          >{sortBy === "top"
            ? "Nổi bật"
            : sortBy === "newest"
              ? "Mới nhất"
              : sortBy === "oldest"
                ? "Cũ nhất"
                : "Tranh cãi"}</span
        >
        <svg
          width="12"
          height="12"
          viewBox="0 0 12 12"
          fill="currentColor"
          class="chevron"
        >
          <path d="M2 4l4 4 4-4" />
        </svg>
      </button>
      {#if showSortDropdown}
        <div class="dropdown-menu">
          <button
            class="dropdown-item"
            class:active={sortBy === "top"}
            onclick={() => handleSortChange("top")}
          >
            Nổi bật
          </button>
          <button
            class="dropdown-item"
            class:active={sortBy === "newest"}
            onclick={() => handleSortChange("newest")}
          >
            Mới nhất
          </button>
          <button
            class="dropdown-item"
            class:active={sortBy === "oldest"}
            onclick={() => handleSortChange("oldest")}
          >
            Cũ nhất
          </button>
          <button
            class="dropdown-item"
            class:active={sortBy === "controversial"}
            onclick={() => handleSortChange("controversial")}
          >
            Tranh cãi
          </button>
        </div>
      {/if}
    </div>
  </div>

  <!-- Comments List -->
  <div class="comments-list">
    {#if isLoading}
      <div class="loading">Đang tải bình luận...</div>
    {:else}
      {#each sortedComments as comment (comment.id)}
        <CommentComponent {comment} depth={0} onUpdate={loadComments} />
      {/each}
      {#if sortedComments.length === 0}
        <div class="no-comments">
          <p>Chưa có bình luận</p>
          <p class="no-comments-subtitle">
            Hãy là người đầu tiên chia sẻ suy nghĩ!
          </p>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .comment-section {
    background: white;
    border-radius: 4px;
    padding: 16px;
    margin-top: 16px;
  }

  .error-banner {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    margin-bottom: 16px;
    background-color: #fee;
    border: 1px solid #fcc;
    border-radius: 4px;
    color: #c00;
    font-size: 14px;
    font-weight: 500;
  }

  .error-banner svg {
    flex-shrink: 0;
  }

  .error-banner span {
    flex: 1;
  }

  .close-btn {
    background: none;
    border: none;
    color: #c00;
    font-size: 24px;
    font-weight: 700;
    line-height: 1;
    cursor: pointer;
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: background 0.2s;
  }

  .close-btn:hover {
    background: rgba(204, 0, 0, 0.1);
  }

  /* Add Comment */
  .add-comment {
    margin-bottom: 16px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    padding: 8px;
  }

  .comment-textarea {
    width: 100%;
    padding: 8px 12px;
    border: none;
    font-size: 14px;
    font-family: "Roboto", sans-serif;
    resize: vertical;
    box-sizing: border-box;
    min-height: 100px;
  }

  .comment-textarea:focus {
    outline: none;
  }

  .image-preview {
    position: relative;
    margin-top: 8px;
    margin-bottom: 8px;
  }

  .image-preview img {
    max-width: 200px;
    max-height: 200px;
    border-radius: 4px;
    display: block;
  }

  .remove-image {
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

  .remove-image:hover {
    background: rgba(0, 0, 0, 0.9);
  }

  .comment-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 8px;
    border-top: 1px solid #edeff1;
    margin-top: 8px;
  }

  .attachment-buttons {
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
    padding: 8px 24px;
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

  .submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Sort Bar */
  .sort-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-bottom: 12px;
    border-bottom: 1px solid #edeff1;
    margin-bottom: 16px;
  }

  .comment-count {
    font-size: 14px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .sort-dropdown {
    position: relative;
  }

  .sort-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: none;
    border: 1px solid #edeff1;
    padding: 6px 12px;
    border-radius: 9999px;
    font-size: 13px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s;
  }

  .sort-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .sort-icon {
    color: #878a8c;
  }

  .sort-label {
    font-weight: 700;
  }

  .chevron {
    color: #878a8c;
  }

  .dropdown-menu {
    position: absolute;
    top: calc(100% + 4px);
    right: 0;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    min-width: 140px;
    z-index: 100;
    overflow: hidden;
  }

  .dropdown-item {
    width: 100%;
    padding: 10px 16px;
    background: none;
    border: none;
    text-align: left;
    font-size: 14px;
    font-weight: 500;
    color: #1c1c1c;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    transition: background 0.2s;
  }

  .dropdown-item:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .dropdown-item.active {
    background: rgba(21, 48, 96, 0.1);
    color: var(--blue--);
    font-weight: 700;
  }

  /* Comments List */
  .comments-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .no-comments {
    text-align: center;
    padding: 40px 20px;
    color: #7c7c7c;
  }

  .no-comments p {
    margin: 0;
    font-size: 16px;
    font-weight: 500;
  }

  .no-comments-subtitle {
    font-size: 14px;
    font-weight: 400;
    margin-top: 8px !important;
  }

  .loading {
    text-align: center;
    padding: 40px 20px;
    color: #7c7c7c;
    font-size: 14px;
  }
</style>
