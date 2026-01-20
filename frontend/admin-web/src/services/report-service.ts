import { authenticatedFetch, handleApiResponse } from "./api";
import type {
  Report,
  GetReportsQuery,
  ReportsResponse,
} from "../dtos/report-dto";

export async function getReports(
  query: GetReportsQuery = {}
): Promise<ReportsResponse> {
  const params = new URLSearchParams();
  if (query.target_type) params.append("target_type", query.target_type);
  if (query.is_deleted !== undefined)
    params.append("is_deleted", query.is_deleted.toString());
  if (query.limit) params.append("limit", query.limit.toString());
  if (query.offset) params.append("offset", query.offset.toString());

  const url = `/api/admin/reports${params.toString() ? `?${params}` : ""}`;
  const response = await authenticatedFetch(url);
  return handleApiResponse<ReportsResponse>(response);
}

export async function getReportById(reportId: string): Promise<Report> {
  const response = await authenticatedFetch(`/api/admin/reports/${reportId}`);
  return handleApiResponse<{ data: Report }>(response).then((res) => res.data);
}

export async function deleteReport(reportId: string): Promise<void> {
  const response = await authenticatedFetch(`/api/admin/reports/${reportId}`, {
    method: "DELETE",
  });
  return handleApiResponse<void>(response);
}

export async function deleteReports(reportIds: string[]): Promise<void> {
  const response = await authenticatedFetch("/api/admin/reports/batch", {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ report_ids: reportIds }),
  });
  return handleApiResponse<void>(response);
}
