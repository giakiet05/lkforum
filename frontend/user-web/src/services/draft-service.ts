import type {
    CreateDraftRequest,
    UpdateDraftRequest,
    DraftResponse,
    PaginatedDraftsResponse,
} from "../dtos/draft-dto";
import { authenticatedFetch, handleApiResponse } from "./api";

/**
 * Create a new draft
 */
export async function createDraft(data: CreateDraftRequest): Promise<DraftResponse> {
    const res = await authenticatedFetch("/api/drafts", {
        method: "POST",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Get all drafts for current user
 */
export async function getDrafts(page = 1, pageSize = 10): Promise<PaginatedDraftsResponse> {
    const params = new URLSearchParams({
        page: page.toString(),
        pageSize: pageSize.toString(),
    });

    const res = await authenticatedFetch(`/api/drafts?${params.toString()}`, {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Get draft by ID
 */
export async function getDraftById(draftId: string): Promise<DraftResponse> {
    const res = await authenticatedFetch(`/api/drafts/${draftId}`, {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Update a draft
 */
export async function updateDraft(draftId: string, data: UpdateDraftRequest): Promise<DraftResponse> {
    const res = await authenticatedFetch(`/api/drafts/${draftId}`, {
        method: "PUT",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Delete a draft
 */
export async function deleteDraft(draftId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/drafts/${draftId}`, {
        method: "DELETE",
    });

    await handleApiResponse(res);
}

/**
 * Publish a draft as a post
 */
export async function publishDraft(draftId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/drafts/${draftId}/publish`, {
        method: "POST",
    });

    await handleApiResponse(res);
}
