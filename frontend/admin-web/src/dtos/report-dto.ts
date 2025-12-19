export interface Report {
  id: string;
  reporter_id: string;
  target_id: string;
  target_type: "user" | "post" | "comment";
  reason: string;
  description?: string;
  is_deleted: boolean;
  deleted_at?: string;
  created_at: string;
}

export interface GetReportsQuery {
  target_type?: "user" | "post" | "comment";
  is_deleted?: boolean;
  limit?: number;
  offset?: number;
}

export interface ReportsResponse {
  reports: Report[];
  total: number;
}
