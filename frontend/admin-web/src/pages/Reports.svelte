<script lang="ts">
  import { onMount } from "svelte";
  import ReportTable from "../components/ReportTable.svelte";
  import ConfirmModal from "../components/ConfirmModal.svelte";
  import { getReports, deleteReport } from "../services/report-service";
  import { toastStore } from "../stores/toast-store";
  import type { Report } from "../dtos/report-dto";

  let loading = $state(false);
  let reports = $state<Report[]>([]);
  let showDeleteConfirm = $state(false);
  let reportToDelete = $state<string | null>(null);

  async function loadReports() {
    loading = true;
    try {
      const response = await getReports({ limit: 100 });
      reports = response.reports;
    } catch (error) {
      console.error("Failed to load reports:", error);
      toastStore.error("Không thể tải danh sách báo cáo");
    } finally {
      loading = false;
    }
  }

  async function handleDeleteReport(reportId: string) {
    reportToDelete = reportId;
    showDeleteConfirm = true;
  }

  async function confirmDeleteReport() {
    if (!reportToDelete) return;
    showDeleteConfirm = false;

    try {
      await deleteReport(reportToDelete);
      await loadReports();
      toastStore.success("Đã xóa báo cáo");
    } catch (error) {
      toastStore.error("Không thể xóa báo cáo");
    } finally {
      reportToDelete = null;
    }
  }

  onMount(() => {
    loadReports();
  });
</script>

<div class="reports-page">
  <div class="page-header">
    <h1>Quản lý báo cáo</h1>
  </div>

  <div class="page-content">
    {#if loading}
      <div class="loading">Đang tải...</div>
    {:else}
      <ReportTable {reports} onDelete={handleDeleteReport} />
    {/if}
  </div>
</div>

<ConfirmModal
  show={showDeleteConfirm}
  title="Xác nhận xóa"
  message="Bạn có chắc muốn xóa báo cáo này?"
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDeleteReport}
  onCancel={() => {
    showDeleteConfirm = false;
    reportToDelete = null;
  }}
/>

<style>
  .reports-page {
    height: 100%;
  }

  .page-header {
    margin-bottom: 2rem;
  }

  .page-header h1 {
    font-size: 1.75rem;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0;
  }

  .page-content {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .loading {
    padding: 3rem;
    text-align: center;
    color: #6c757d;
    font-size: 1rem;
  }
</style>
