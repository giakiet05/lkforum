<script lang="ts">
  import { mockComments, type Comment } from "../mocks/comments.mock";
  import CommentComponent from "./Comment.svelte";

  type CommentSectionProps = {
    postId: string;
  };

  let { postId }: CommentSectionProps = $props();

  type SortType = "top" | "newest" | "oldest" | "controversial";

  let sortBy = $state<SortType>("top");
  let showSortDropdown = $state(false);
  let newCommentContent = $state("");
  let selectedImage = $state<File | null>(null);
  let imagePreview = $state<string | null>(null);

  // Filter comments for this post
  const postComments = $derived(
    mockComments.filter((c) => c.postId === postId)
  );

  // Sort comments
  const sortedComments = $derived(
    (() => {
      const comments = [...postComments];

      switch (sortBy) {
        case "top":
          return comments.sort((a, b) => {
            const scoreA = a.votesCount.up - a.votesCount.down;
            const scoreB = b.votesCount.up - b.votesCount.down;
            return scoreB - scoreA;
          });
        case "newest":
          return comments.sort(
            (a, b) =>
              new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
          );
        case "oldest":
          return comments.sort(
            (a, b) =>
              new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
          );
        case "controversial":
          return comments.sort((a, b) => {
            // Comments with similar up/down votes are more controversial
            const controversyA = Math.min(a.votesCount.up, a.votesCount.down);
            const controversyB = Math.min(b.votesCount.up, b.votesCount.down);
            return controversyB - controversyA;
          });
        default:
          return comments;
      }
    })()
  );

  const handleSortChange = (sort: SortType) => {
    sortBy = sort;
    showSortDropdown = false;
  };

  const submitComment = () => {
    if (newCommentContent.trim()) {
      // Mock submit - in real app, would send to backend
      console.log("Submitting comment:", newCommentContent);
      if (selectedImage) {
        console.log("With image:", selectedImage.name);
      }
      newCommentContent = "";
      selectedImage = null;
      imagePreview = null;
    }
  };

  const handleImageSelect = (event: Event) => {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      selectedImage = file;

      // Create preview
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
    const countReplies = (comment: Comment): number => {
      return (
        1 + comment.replies.reduce((sum, reply) => sum + countReplies(reply), 0)
      );
    };
    return postComments.reduce(
      (sum, comment) => sum + countReplies(comment),
      0
    );
  };
</script>

<div class="comment-section">
  <!-- Comment Input Box -->
  <div class="add-comment">
    <textarea
      bind:value={newCommentContent}
      placeholder="What are your thoughts?"
      class="comment-textarea"
      rows="4"
    ></textarea>

    {#if imagePreview}
      <div class="image-preview">
        <img src={imagePreview} alt="Preview" />
        <button class="remove-image" onclick={removeImage} title="Remove image">
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
          title="Add image"
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
          title="Add picture"
        >
          <img src="/comment_picture.png" alt="" width="20" height="20" />
        </button>
      </div>
      <button class="submit-btn" onclick={submitComment}>Comment</button>
    </div>
  </div>

  <!-- Sort Bar -->
  <div class="sort-bar">
    <span class="comment-count">{getTotalComments()} Comments</span>
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
          >{sortBy.charAt(0).toUpperCase() + sortBy.slice(1)}</span
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
            Top
          </button>
          <button
            class="dropdown-item"
            class:active={sortBy === "newest"}
            onclick={() => handleSortChange("newest")}
          >
            Newest
          </button>
          <button
            class="dropdown-item"
            class:active={sortBy === "oldest"}
            onclick={() => handleSortChange("oldest")}
          >
            Oldest
          </button>
          <button
            class="dropdown-item"
            class:active={sortBy === "controversial"}
            onclick={() => handleSortChange("controversial")}
          >
            Controversial
          </button>
        </div>
      {/if}
    </div>
  </div>

  <!-- Comments List -->
  <div class="comments-list">
    {#each sortedComments as comment}
      <CommentComponent {comment} depth={0} />
    {/each}
    {#if sortedComments.length === 0}
      <div class="no-comments">
        <p>No comments yet</p>
        <p class="no-comments-subtitle">
          Be the first to share what you think!
        </p>
      </div>
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

  .comment-submit {
    display: flex;
    justify-content: flex-end;
    padding-top: 8px;
    border-top: 1px solid #edeff1;
    margin-top: 8px;
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
</style>
