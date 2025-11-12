import type {
    GetCommentsFilterQuery,
    CommentResponse,
} from "../dtos/comment-dto";
import { publicFetch, handleApiResponse } from "./api";

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
