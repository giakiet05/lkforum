import type {
    GetPostsQuery,
    PostResponse,
    CreatePostRequest,
    PostVoteRequest,
    PaginatedPostsResponse,
} from "../dtos/post-dto";
import { publicFetch, authenticatedFetch, handleApiResponse } from "./api";

/**
 * Get posts with optional filters
 */
export async function getPosts(query?: GetPostsQuery): Promise<PostResponse[]> {
    const params = new URLSearchParams();
    
    if (query?.community_id) params.append("community_id", query.community_id);
    if (query?.author_id) params.append("author_id", query.author_id);
    if (query?.type) params.append("type", query.type);
    if (query?.sort) params.append("sort", query.sort);
    if (query?.time) params.append("time", query.time);
    if (query?.page) params.append("page", query.page.toString());
    if (query?.limit) params.append("limit", query.limit.toString());

    const queryString = params.toString();
    const url = queryString ? `/api/posts?${queryString}` : "/api/posts";

    const res = await publicFetch(url, {
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
    const res = await authenticatedFetch("/api/posts/", {
        method: "POST",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Vote on a post (upvote or downvote)
 */
export async function voteOnPost(postId: string, value: boolean): Promise<void> {
    const res = await authenticatedFetch(`/api/posts/${postId}/vote`, {
        method: "POST",
        body: JSON.stringify({ value }),
    });

    await handleApiResponse(res);
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
