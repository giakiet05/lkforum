<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import type { PostResponse } from "../dtos/post-dto";
  import {
    getPostById,
    updatePost,
    removePostImages,
    removePostVideos,
    uploadPostImages,
    uploadPostVideos,
  } from "../services/post-service";
  import { authStore } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
  import { extractPostId } from "../utils/slug";
  import ConfirmModal from "../components/ConfirmModal.svelte";

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
  let showCancelConfirm = $state(false);

  // Media management state
  let existingImages = $state<Array<{ url: string; public_id: string }>>([]);
  let existingVideos = $state<Array<{ url: string; public_id: string }>>([]);
  let imagesToRemove = $state<string[]>([]);
  let videosToRemove = $state<string[]>([]);
  let newImageFiles = $state<File[]>([]);
  let newVideoFiles = $state<File[]>([]);
  let showRemoveMediaConfirm = $state(false);
  let mediaToRemove = $state<{
    type: "image" | "video";
    publicId: string;
    url: string;
  } | null>(null);

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

      // Initialize existing media
      if (post.content.images) {
        existingImages = post.content.images.map((img) => ({
          url: img.url,
          public_id: img.public_id,
        }));
      }
      if (post.content.videos) {
        existingVideos = post.content.videos.map((vid) => ({
          url: vid.url,
          public_id: vid.public_id,
        }));
      }
    } catch (err) {
      console.error("Failed to load post:", err);
      error = "Failed to load post";
    } finally {
      isLoading = false;
    }
  });

  // Filtered media (excluding ones marked for removal)
  const displayedImages = $derived(
    existingImages.filter((img) => !imagesToRemove.includes(img.public_id)),
  );
  const displayedVideos = $derived(
    existingVideos.filter((vid) => !videosToRemove.includes(vid.public_id)),
  );

  function handleRemoveExistingImage(publicId: string, url: string) {
    mediaToRemove = { type: "image", publicId, url };
    showRemoveMediaConfirm = true;
  }

  function handleRemoveExistingVideo(publicId: string, url: string) {
    mediaToRemove = { type: "video", publicId, url };
    showRemoveMediaConfirm = true;
  }

  function confirmRemoveMedia() {
    if (mediaToRemove) {
      if (mediaToRemove.type === "image") {
        imagesToRemove = [...imagesToRemove, mediaToRemove.publicId];
      } else {
        videosToRemove = [...videosToRemove, mediaToRemove.publicId];
      }
    }
    showRemoveMediaConfirm = false;
    mediaToRemove = null;
  }

  function handleNewImageSelect(e: Event) {
    const input = e.target as HTMLInputElement;
    const files = Array.from(input.files || []).filter((f) =>
      f.type.startsWith("image/"),
    );
    newImageFiles = [...newImageFiles, ...files];
    input.value = "";
  }

  function handleNewVideoSelect(e: Event) {
    const input = e.target as HTMLInputElement;
    const files = Array.from(input.files || []).filter((f) =>
      f.type.startsWith("video/"),
    );
    newVideoFiles = [...newVideoFiles, ...files];
    input.value = "";
  }

  function removeNewImage(index: number) {
    newImageFiles = newImageFiles.filter((_, i) => i !== index);
  }

  function removeNewVideo(index: number) {
    newVideoFiles = newVideoFiles.filter((_, i) => i !== index);
  }

  function getFilePreviewUrl(file: File): string {
    return URL.createObjectURL(file);
  }

  async function handleSave() {
    if (!title.trim()) {
      toastStore.warning("Tiêu đề là bắt buộc");
      return;
    }

    try {
      isSaving = true;

      // 1. Update post text content
      await updatePost(postId, {
        title: title.trim(),
        content: {
          text: content.trim(),
        },
      });

      // 2. Remove images if any marked
      if (imagesToRemove.length > 0) {
        await removePostImages(postId, imagesToRemove);
      }

      // 3. Remove videos if any marked
      if (videosToRemove.length > 0) {
        await removePostVideos(postId, videosToRemove);
      }

      // 4. Upload new images if any
      if (newImageFiles.length > 0) {
        await uploadPostImages(postId, newImageFiles);
      }

      // 5. Upload new videos if any
      if (newVideoFiles.length > 0) {
        await uploadPostVideos(postId, newVideoFiles);
      }

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
    showCancelConfirm = true;
  }

  function confirmCancel() {
    showCancelConfirm = false;
    push(`/post/${params.slugId}`);
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

        <!-- Existing Images Section -->
        {#if displayedImages.length > 0 || newImageFiles.length > 0}
          <div class="media-section">
            <label>Hình ảnh</label>
            <div class="media-grid">
              {#each displayedImages as img}
                <div class="media-item">
                  <img src={img.url} alt="Post image" />
                  <button
                    class="remove-btn"
                    onclick={() =>
                      handleRemoveExistingImage(img.public_id, img.url)}
                    title="Xóa ảnh này"
                  >
                    ✕
                  </button>
                </div>
              {/each}
              {#each newImageFiles as file, index}
                <div class="media-item new">
                  <img src={getFilePreviewUrl(file)} alt={file.name} />
                  <button
                    class="remove-btn"
                    onclick={() => removeNewImage(index)}
                    title="Xóa ảnh này"
                  >
                    ✕
                  </button>
                  <span class="new-badge">Mới</span>
                </div>
              {/each}
            </div>
            <label class="add-media-btn">
              <input
                type="file"
                accept="image/*"
                multiple
                onchange={handleNewImageSelect}
                hidden
              />
              + Thêm ảnh
            </label>
          </div>
        {:else if post.type !== "poll"}
          <div class="media-section">
            <label>Hình ảnh</label>
            <label class="add-media-btn">
              <input
                type="file"
                accept="image/*"
                multiple
                onchange={handleNewImageSelect}
                hidden
              />
              + Thêm ảnh
            </label>
            {#if newImageFiles.length > 0}
              <div class="media-grid">
                {#each newImageFiles as file, index}
                  <div class="media-item new">
                    <img src={getFilePreviewUrl(file)} alt={file.name} />
                    <button
                      class="remove-btn"
                      onclick={() => removeNewImage(index)}
                      title="Xóa ảnh này"
                    >
                      ✕
                    </button>
                    <span class="new-badge">Mới</span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}

        <!-- Existing Videos Section -->
        {#if displayedVideos.length > 0 || newVideoFiles.length > 0}
          <div class="media-section">
            <label>Video</label>
            <div class="media-grid">
              {#each displayedVideos as vid}
                <div class="media-item video">
                  <video src={vid.url} muted></video>
                  <div class="video-overlay">▶</div>
                  <button
                    class="remove-btn"
                    onclick={() =>
                      handleRemoveExistingVideo(vid.public_id, vid.url)}
                    title="Xóa video này"
                  >
                    ✕
                  </button>
                </div>
              {/each}
              {#each newVideoFiles as file, index}
                <div class="media-item video new">
                  <video src={getFilePreviewUrl(file)} muted></video>
                  <div class="video-overlay">▶</div>
                  <button
                    class="remove-btn"
                    onclick={() => removeNewVideo(index)}
                    title="Xóa video này"
                  >
                    ✕
                  </button>
                  <span class="new-badge">Mới</span>
                </div>
              {/each}
            </div>
            <label class="add-media-btn">
              <input
                type="file"
                accept="video/*"
                multiple
                onchange={handleNewVideoSelect}
                hidden
              />
              + Thêm video
            </label>
          </div>
        {:else if post.type !== "poll"}
          <div class="media-section">
            <label>Video</label>
            <label class="add-media-btn">
              <input
                type="file"
                accept="video/*"
                multiple
                onchange={handleNewVideoSelect}
                hidden
              />
              + Thêm video
            </label>
            {#if newVideoFiles.length > 0}
              <div class="media-grid">
                {#each newVideoFiles as file, index}
                  <div class="media-item video new">
                    <video src={getFilePreviewUrl(file)} muted></video>
                    <div class="video-overlay">▶</div>
                    <button
                      class="remove-btn"
                      onclick={() => removeNewVideo(index)}
                      title="Xóa video này"
                    >
                      ✕
                    </button>
                    <span class="new-badge">Mới</span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}

        {#if imagesToRemove.length > 0 || videosToRemove.length > 0}
          <div class="warning-message">
            <p>
              ⚠️ {imagesToRemove.length + videosToRemove.length} media sẽ bị xóa
              khi lưu bài viết.
            </p>
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

<ConfirmModal
  show={showCancelConfirm}
  title="Xác nhận hủy"
  message="Bạn có chắc muốn bỏ các thay đổi? Các thay đổi chưa lưu sẽ bị mất."
  confirmText="Bỏ thay đổi"
  cancelText="Tiếp tục chỉnh sửa"
  confirmVariant="danger"
  onConfirm={confirmCancel}
  onCancel={() => (showCancelConfirm = false)}
/>

<ConfirmModal
  show={showRemoveMediaConfirm}
  title="Xác nhận xóa"
  message="Bạn có chắc muốn xóa {mediaToRemove?.type === 'image'
    ? 'ảnh'
    : 'video'} này? Hành động này không thể hoàn tác sau khi lưu."
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmRemoveMedia}
  onCancel={() => {
    showRemoveMediaConfirm = false;
    mediaToRemove = null;
  }}
/>

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

  .warning-message {
    padding: 12px;
    background: #f8d7da;
    border: 1px solid #f5c6cb;
    border-radius: 4px;
    color: #721c24;
  }

  .warning-message p {
    margin: 0;
    font-size: 14px;
  }

  /* Media Section Styles */
  .media-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .media-section > label:first-child {
    font-weight: 600;
    font-size: 14px;
    color: #1c1c1c;
  }

  .media-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
  }

  .media-item {
    position: relative;
    width: 120px;
    height: 120px;
    border-radius: 8px;
    overflow: hidden;
    background: #f6f7f8;
    border: 1px solid #edeff1;
  }

  .media-item.new {
    border: 2px solid var(--blue--);
  }

  .media-item img,
  .media-item video {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .media-item .remove-btn {
    position: absolute;
    top: 4px;
    right: 4px;
    width: 24px;
    height: 24px;
    background: rgba(0, 0, 0, 0.7);
    border: none;
    border-radius: 50%;
    color: white;
    font-size: 14px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
    z-index: 1;
  }

  .media-item .remove-btn:hover {
    background: #c00;
  }

  .media-item .video-overlay {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 36px;
    height: 36px;
    background: rgba(0, 0, 0, 0.6);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 14px;
    pointer-events: none;
  }

  .media-item .new-badge {
    position: absolute;
    bottom: 4px;
    left: 4px;
    padding: 2px 6px;
    background: var(--blue--);
    color: white;
    font-size: 10px;
    font-weight: 600;
    border-radius: 4px;
  }

  .add-media-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 8px 16px;
    background: transparent;
    border: 1px dashed var(--blue--);
    border-radius: 8px;
    color: var(--blue--);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    width: fit-content;
  }

  .add-media-btn:hover {
    background: rgba(21, 48, 96, 0.05);
  }
</style>
