import type {
    GetPostsQuery,
    PostResponse,
    CreatePostRequest,
    PostVoteRequest,
    PaginatedPostsResponse,
    UpdatePostRequest,
    PollVoteRequest,
    ReportPostRequest,
    RemoveImagesRequest,
    PollResponse,
} from "../dtos/post-dto";
import { publicFetch, authenticatedFetch, handleApiResponse, optionalAuthFetch } from "./api";

/**
 * Get posts with optional filters
 * Uses optional auth - if logged in, sends token for personalized feed
 */
export async function getPosts(query?: GetPostsQuery): Promise<PostResponse[]> {
    const params = new URLSearchParams();
    
    if (query?.community_id) params.append("community_id", query.community_id);
    if (query?.author_id) params.append("author_id", query.author_id);
    if (query?.type) params.append("type", query.type);
    if (query?.sort) params.append("sort", query.sort);
    if (query?.time) params.append("time", query.time);
    if (query?.feed_type) params.append("feed_type", query.feed_type);
    if (query?.search) params.append("search", query.search);
    if (query?.page) params.append("page", query.page.toString());
    if (query?.limit) params.append("limit", query.limit.toString());

    const queryString = params.toString();
    const url = queryString ? `/api/posts?${queryString}` : "/api/posts";

    // Use optionalAuthFetch - sends token if available for personalized home feed
    const res = await optionalAuthFetch(url, {
        method: "GET",
    });

    const response: PaginatedPostsResponse = await handleApiResponse(res);
    console.log("🔍 getPosts response:", response);
    return response.posts;
}

/**
 * Get current user's posts (requires authentication)
 */
export async function getMyPosts(query?: GetPostsQuery): Promise<PostResponse[]> {
    const params = new URLSearchParams();
    
    if (query?.community_id) params.append("community_id", query.community_id);
    if (query?.type) params.append("type", query.type);
    if (query?.sort) params.append("sort", query.sort);
    if (query?.time) params.append("time", query.time);
    if (query?.page) params.append("page", query.page.toString());
    if (query?.limit) params.append("limit", query.limit.toString());

    const queryString = params.toString();
    const url = queryString ? `/api/posts/me?${queryString}` : "/api/posts/me";

    const res = await authenticatedFetch(url, {
        method: "GET",
    });

    const response: PaginatedPostsResponse = await handleApiResponse(res);
    return response.posts;
}

/**
 * Get posts by a specific user (for profile page)
 */
export async function getPostsByUserId(userId: string, page = 1, limit = 10): Promise<PostResponse[]> {
    return getPosts({
        author_id: userId,
        page,
        limit,
    });
}

/**
 * Create a new post (text or poll)
 */
export async function createPost(data: CreatePostRequest): Promise<PostResponse> {
    console.log("📡 Calling POST /api/posts/ with data:", data);
    const res = await authenticatedFetch("/api/posts/", {
        method: "POST",
        body: JSON.stringify(data),
    });

    console.log("📥 Response status:", res.status);
    return await handleApiResponse(res);
}

/**
 * Vote on a post (upvote or downvote)
 * Backend tự động toggle: POST cùng giá trị sẽ remove vote
 */
export async function voteOnPost(postId: string, value: boolean): Promise<void> {
    console.log("🗳️ Voting on post:", postId, "value:", value);
    const res = await authenticatedFetch(`/api/votes/post/${postId}`, {
        method: "POST",
        body: JSON.stringify({ value }),
    });

    console.log("🗳️ Vote response status:", res.status);
    const result = await handleApiResponse(res);
    console.log("🗳️ Vote result:", result);
    return result;
}

/**
 * Upload images to an existing post
 * NOTE: This requires FormData with images
 */
export async function uploadPostImages(postId: string, images: File[]): Promise<PostResponse> {
    const formData = new FormData();
    images.forEach((image) => {
        formData.append("images", image);
    });

    const res = await authenticatedFetch(`/api/posts/${postId}/images`, {
        method: "POST",
        body: formData,
    });

    return await handleApiResponse(res);
}

/**
 * Remove images from an existing post
 */
export async function removePostImages(postId: string, publicIds: string[]): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/images`, {
        method: "DELETE",
        body: JSON.stringify({ public_ids: publicIds }),
    });

    await handleApiResponse(res);
}

/**
 * Upload videos to an existing post
 */
export async function uploadPostVideos(postId: string, videos: File[]): Promise<PostResponse> {
    const formData = new FormData();
    videos.forEach((video) => {
        formData.append("videos", video);
    });

    const res = await authenticatedFetch(`/api/posts/${postId}/videos`, {
        method: "POST",
        body: formData,
    });

    return await handleApiResponse(res);
}

/**
 * Remove videos from an existing post
 */
export async function removePostVideos(postId: string, publicIds: string[]): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/videos`, {
        method: "DELETE",
        body: JSON.stringify({ public_ids: publicIds }),
    });

    await handleApiResponse(res);
}

/**
 * Get a single post by ID
 */
export async function getPostById(postId: string): Promise<PostResponse> {
    const res = await publicFetch(`/api/posts/${postId}`, {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Update a post
 */
export async function updatePost(postId: string, data: UpdatePostRequest): Promise<PostResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}`, {
        method: "PUT",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Delete a post
 */
export async function deletePost(postId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}`, {
        method: "DELETE",
    });

    await handleApiResponse(res);
}

/**
 * Save a post
 */
export async function savePost(postId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/save`, {
        method: "POST",
    });

    await handleApiResponse(res);
}

/**
 * Unsave a post
 */
export async function unsavePost(postId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/save`, {
        method: "DELETE",
    });

    await handleApiResponse(res);
}

/**
 * Get saved posts
 */
export async function getSavedPosts(page = 1, limit = 10): Promise<PostResponse[]> {
    const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
    });

    const res = await authenticatedFetch(`/api/posts/saved?${params.toString()}`, {
        method: "GET",
    });

    const response: PaginatedPostsResponse = await handleApiResponse(res);
    return response.posts;
}

/**
 * Hide a post
 */
export async function hidePost(postId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/hide`, {
        method: "POST",
    });

    await handleApiResponse(res);
}

/**
 * Unhide a post
 */
export async function unhidePost(postId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/unhide`, {
        method: "POST",
    });

    await handleApiResponse(res);
}

/**
 * Get hidden posts
 */
export async function getHiddenPosts(page = 1, limit = 10): Promise<PostResponse[]> {
    const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
    });

    const res = await authenticatedFetch(`/api/posts/hidden?${params.toString()}`, {
        method: "GET",
    });

    const response: PaginatedPostsResponse = await handleApiResponse(res);
    return response.posts;
}

/**
 * Report a post
 */
export async function reportPost(postId: string, data: ReportPostRequest): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/report`, {
        method: "POST",
        body: JSON.stringify(data),
    });

    await handleApiResponse(res);
}

/**
 * Vote on a poll
 */
export async function voteOnPoll(postId: string, optionId: string): Promise<PollResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}/poll/vote`, {
        method: "POST",
        body: JSON.stringify({ option_id: optionId }),
    });

    return await handleApiResponse(res);
}

/**
 * Remove poll vote
 */
export async function removePollVote(postId: string): Promise<PollResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}/poll/vote`, {
        method: "DELETE",
    });

    return await handleApiResponse(res);
}

/**
 * Update a poll (question, expires_at, allow_multiple)
 */
export async function updatePoll(postId: string, data: { question?: string; expires_at?: string; allow_multiple?: boolean }): Promise<PostResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}/poll`, {
        method: "PUT",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Add options to a poll
 */
export async function addPollOptions(postId: string, options: string[]): Promise<PostResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}/poll/options`, {
        method: "POST",
        body: JSON.stringify({ options }),
    });

    return await handleApiResponse(res);
}

/**
 * Remove options from a poll
 */
export async function removePollOptions(postId: string, optionIds: string[]): Promise<PostResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}/poll/options`, {
        method: "DELETE",
        body: JSON.stringify({ option_ids: optionIds }),
    });

    return await handleApiResponse(res);
}

/**
 * Update a specific poll option text
 */
export async function updatePollOption(postId: string, optionId: string, text: string): Promise<PostResponse> {
    const res = await authenticatedFetch(`/api/posts/${postId}/poll/options/${optionId}`, {
        method: "PUT",
        body: JSON.stringify({ text }),
    });

    return await handleApiResponse(res);
}
