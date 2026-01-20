import type {
    GetCommentsFilterQuery,
    GetCommentByPostIDQuery,
    CreateCommentRequest,
    CommentResponse,
} from "../dtos/comment-dto";
import { publicFetch, authenticatedFetch, handleApiResponse } from "./api";

/**
 * Get comments by post ID (main endpoint for PostDetail page)
 */
export async function getCommentsByPostId(query: GetCommentByPostIDQuery): Promise<{ comments: CommentResponse[], pagination: any }> {
    const params = new URLSearchParams();
    
    params.append("post_id", query.post_id);
    if (query.depth !== undefined) params.append("depth", query.depth.toString());
    if (query.children_page_size) params.append("children_page_size", query.children_page_size.toString());
    if (query.page) params.append("page", query.page.toString());
    if (query.page_size) params.append("page_size", query.page_size.toString());

    const res = await publicFetch(`/api/comments/post?${params.toString()}`, {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Create a new comment or reply
 */
export async function createComment(data: CreateCommentRequest): Promise<CommentResponse> {
    const res = await authenticatedFetch("/api/comments", {
        method: "POST",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Delete a comment by ID
 */
export async function deleteComment(commentId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/comments/${commentId}`, {
        method: "DELETE",
    });

    return await handleApiResponse(res);
}

/**
 * Get comments with filters
 */
export async function getCommentsFilter(query?: GetCommentsFilterQuery): Promise<CommentResponse[]> {
    const params = new URLSearchParams();
    
    if (query?.post_id) params.append("post_id", query.post_id);
    if (query?.parent_id) params.append("parent_id", query.parent_id);
    if (query?.user_id) params.append("user_id", query.user_id);
    if (query?.content) params.append("content", query.content);
    if (query?.page) params.append("page", query.page.toString());
    if (query?.page_size) params.append("page_size", query.page_size.toString());

    const queryString = params.toString();
    const url = queryString ? `/api/comments/filter?${queryString}` : "/api/comments/filter";

    const res = await publicFetch(url, {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Get comments by a specific user (for profile page)
 */
export async function getCommentsByUserId(userId: string, page = 1, pageSize = 10): Promise<CommentResponse[]> {
    return getCommentsFilter({
        user_id: userId,
        page,
        page_size: pageSize,
    });
}
