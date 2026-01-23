<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import type { PostResponse } from "../dtos/post-dto";
  import { getPostById, updatePost } from "../services/post-service";
  import { authStore } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
  import { extractPostId } from "../utils/slug";

  type EditPostProps = {
    params?: { slugId: string };
  };

  let { params = { slugId: "" } }: EditPostProps = $props();

  let post = $state<PostResponse | null>(null);
  let title = $state("");
  let content = $state("");
  let isLoading = $state(true);
  let isSaving = $state(false);
  let error = $state<string | null>(null);

  const currentUser = $derived($authStore.user);
  const postId = $derived(extractPostId(params.slugId));

  onMount(async () => {
    if (!currentUser) {
      push("/login");
      return;
    }

    try {
      isLoading = true;
      post = await getPostById(postId);

      // Check if user owns the post
      if (post.author.id !== currentUser.id) {
        toastStore.error("Bạn không có quyền chỉnh sửa bài viết này");
        push(`/post/${params.slugId}`);
        return;
      }

      // Initialize form
      title = post.title;
      content = post.content.text || "";
    } catch (err) {
      console.error("Failed to load post:", err);
      error = "Failed to load post";
    } finally {
      isLoading = false;
    }
  });

  async function handleSave() {
    if (!title.trim()) {
      toastStore.warning("Tiêu đề là bắt buộc");
      return;
    }

    try {
      isSaving = true;
      await updatePost(postId, {
        title: title.trim(),
        content: {
          text: content.trim(),
        },
      });

      toastStore.success("Cập nhật bài viết thành công");
      push(`/post/${params.slugId}`);
    } catch (err) {
      console.error("Failed to update post:", err);
      toastStore.error("Không thể cập nhật bài viết. Vui lòng thử lại.");
    } finally {
      isSaving = false;
    }
  }

  function handleCancel() {
    if (confirm("Bỏ thay đổi?")) {
      push(`/post/${params.slugId}`);
    }
  }
</script>

<div class="edit-post-page">
  {#if isLoading}
    <div class="loading">Đang tải...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else if post}
    <div class="edit-container">
      <div class="header">
        <h1>Chỉnh sửa bài viết</h1>
      </div>

      <div class="form-content">
        <div class="community-info">
          <span>Đăng tại lk/{post.community.name}</span>
        </div>

        <div class="form-group">
          <label for="title">Tiêu đề *</label>
          <input
            id="title"
            type="text"
            bind:value={title}
            placeholder="Tiêu đề"
            maxlength="300"
            required
          />
          <span class="char-count">{title.length}/300</span>
        </div>

        <div class="form-group">
          <label for="content">Nội dung</label>
          <textarea
            id="content"
            bind:value={content}
            placeholder="Nội dung (tuỳ chọn)"
            rows="15"
          ></textarea>
        </div>

        {#if post.type === "poll"}
          <div class="info-message">
            <p>
              ⚠️ Chưa hỗ trợ chỉnh sửa khảo sát. Chỉ có thể chỉnh sửa tiêu đề và
              nội dung.
            </p>
          </div>
        {/if}

        {#if post.content.images && post.content.images.length > 0}
          <div class="info-message">
            <p>⚠️ Chưa hỗ trợ chỉnh sửa hình ảnh. Hình ảnh sẽ giữ nguyên.</p>
          </div>
        {/if}

        {#if post.content.videos && post.content.videos.length > 0}
          <div class="info-message">
            <p>⚠️ Chưa hỗ trợ chỉnh sửa video. Video sẽ giữ nguyên.</p>
          </div>
        {/if}

        <div class="footer-actions">
          <button class="btn-cancel" onclick={handleCancel} disabled={isSaving}>
            Hủy
          </button>
          <button class="btn-save" onclick={handleSave} disabled={isSaving}>
            {isSaving ? "Đang lưu..." : "Lưu"}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .edit-post-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  @media (max-width: 768px) {
    .edit-post-page {
      padding: 12px;
    }
  }

  .loading,
  .error {
    text-align: center;
    padding: 40px;
    font-size: 18px;
  }

  .error {
    color: #d93025;
  }

  .edit-container {
    background: white;
    border: 1px solid #eaebef;
    border-radius: 4px;
    padding: 24px;
  }

  .header {
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #eaebef;
  }

  .header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
  }

  .footer-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid #eaebef;
  }

  .btn-cancel,
  .btn-save {
    padding: 10px 24px;
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

  .btn-cancel:hover:not(:disabled) {
    background: #e9ebed;
  }

  .btn-save {
    background: var(--blue--);
    color: white;
  }

  .btn-save:hover:not(:disabled) {
    background: var(--darkblue--);
  }

  .btn-cancel:disabled,
  .btn-save:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .form-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .community-info {
    padding: 12px;
    background: #f6f7f8;
    border-radius: 4px;
    font-size: 14px;
    font-weight: 500;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    position: relative;
  }

  .form-group label {
    font-weight: 600;
    font-size: 14px;
    color: #1c1c1c;
  }

  .form-group input,
  .form-group textarea {
    width: 100%;
    padding: 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 14px;
    font-family: inherit;
    box-sizing: border-box;
  }

  .form-group input:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: var(--blue--);
  }

  .form-group textarea {
    resize: vertical;
    min-height: 150px;
  }

  .char-count {
    position: absolute;
    right: 12px;
    bottom: -20px;
    font-size: 12px;
    color: #878a8c;
  }

  .info-message {
    padding: 12px;
    background: #fff3cd;
    border: 1px solid #ffc107;
    border-radius: 4px;
    color: #856404;
  }

  .info-message p {
    margin: 0;
    font-size: 14px;
  }
</style>
