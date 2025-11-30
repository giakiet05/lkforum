// --- Request DTOs ---

export type ReportTargetType = "post" | "comment" | "user";

export interface GetReportsFilterQuery {
    reporter_id?: string;
    target_id?: string;
    target_type?: ReportTargetType;
    reason?: string;
    start_date?: string; // ISO date string
    end_date?: string;   // ISO date string
    page?: number;
    page_size?: number;
}

export interface DeleteReportsRequest {
    report_ids: string[];
}

// --- Response DTOs ---

export interface ReportResponse {
    id: string;
    reporter_id: string;
    target_id: string;
    target_type: ReportTargetType;
    reason: string;
    description?: string;
    created_at: string; // ISO date string
}

export interface PaginatedReportResponse {
    reports: ReportResponse[];
    pagination: {
        page: number;
        page_size: number;
        total: number;
    };
}
