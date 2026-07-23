<template>
  <section class="panel result-panel">
    <div class="panel-head">
      <h2 class="panel-title">Results</h2>
      <div class="result-actions">
        <span v-if="result" class="muted">{{ result.rowCount }} rows / {{ result.durationMs }} ms</span>
        <el-button size="small" :disabled="!canExport" data-testid="export-csv" @click="exportResult('csv')">Export CSV</el-button>
        <el-button size="small" type="primary" :disabled="!canExport" data-testid="export-json" @click="exportResult('json')">Export JSON</el-button>
      </div>
    </div>
    <el-empty v-if="!result" description="Run a query" />
    <el-empty v-else-if="result.empty" description="No rows" />
    <template v-else>
      <div class="result-table-scroll">
        <el-table :data="pagedRows" height="340" :style="{ minWidth: tableMinWidth }">
          <el-table-column
            v-for="column in result.columns"
            :key="column.name"
            :prop="column.name"
            :label="column.name"
            min-width="170"
            show-overflow-tooltip
          />
        </el-table>
      </div>
      <div class="pagination-row">
        <span class="page-summary" data-testid="page-summary">Showing {{ visibleStart }}-{{ visibleEnd }} of {{ totalRows }}</span>
        <el-pagination
          v-if="showPagination"
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="pageSizes"
          :total="totalRows"
          layout="sizes, prev, pager, next"
          small
          background
        />
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { ExportFormat, QueryResult } from '../api/types'
import { buildExportFile, canExportResult, downloadExportFile } from '../utils/exportResults'

const props = defineProps<{ result: QueryResult | null; loading: boolean; dbName: string | null }>()

const pageSizes = [25, 50, 100, 200]
const currentPage = ref(1)
const pageSize = ref(50)
const canExport = computed(() => canExportResult(props.result, props.loading))
const totalRows = computed(() => props.result?.rows.length ?? 0)
const showPagination = computed(() => totalRows.value > pageSize.value)
const visibleStart = computed(() => (totalRows.value === 0 ? 0 : (currentPage.value - 1) * pageSize.value + 1))
const visibleEnd = computed(() => Math.min(currentPage.value * pageSize.value, totalRows.value))
const tableMinWidth = computed(() => `${Math.max(760, (props.result?.columns.length ?? 0) * 170)}px`)
const pagedRows = computed(() => {
  if (!props.result) return []
  const start = (currentPage.value - 1) * pageSize.value
  return props.result.rows.slice(start, start + pageSize.value)
})

watch(() => props.result, () => {
  currentPage.value = 1
})

watch([totalRows, pageSize], () => {
  const lastPage = Math.max(1, Math.ceil(totalRows.value / pageSize.value))
  if (currentPage.value > lastPage) currentPage.value = lastPage
})

function exportResult(format: ExportFormat) {
  if (!props.result || !canExport.value) return
  downloadExportFile(buildExportFile(props.result, format, props.dbName))
}
</script>

<style scoped>
.result-panel { padding: 14px; min-height: 220px; min-width: 0; overflow: hidden; }
.panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.result-actions { display: flex; align-items: center; justify-content: flex-end; gap: 8px; flex-wrap: wrap; }
.result-table-scroll {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  border: var(--dq-border);
}
.result-table-scroll :deep(.el-table) {
  border: 0 !important;
}
.pagination-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-top: 10px;
}
.page-summary {
  color: var(--dq-slate);
  font-size: 12px;
  font-weight: 800;
  white-space: nowrap;
}
@media (max-width: 720px) {
  .panel-head,
  .pagination-row {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
