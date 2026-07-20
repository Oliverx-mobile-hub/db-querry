<template>
  <section class="panel result-panel">
    <div class="panel-head">
      <h2 class="panel-title">Results</h2>
      <div class="result-actions">
        <span v-if="result" class="muted">{{ result.rowCount }} rows · {{ result.durationMs }} ms</span>
        <el-button size="small" :disabled="!canExport" data-testid="export-csv" @click="exportResult('csv')">Export CSV</el-button>
        <el-button size="small" type="primary" :disabled="!canExport" data-testid="export-json" @click="exportResult('json')">Export JSON</el-button>
      </div>
    </div>
    <el-empty v-if="!result" description="Run a query" />
    <el-empty v-else-if="result.empty" description="No rows" />
    <el-table v-else :data="result.rows" height="320">
      <el-table-column v-for="column in result.columns" :key="column.name" :prop="column.name" :label="column.name" min-width="160" />
    </el-table>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ExportFormat, QueryResult } from '../api/types'
import { buildExportFile, canExportResult, downloadExportFile } from '../utils/exportResults'

const props = defineProps<{ result: QueryResult | null; loading: boolean; dbName: string | null }>()

const canExport = computed(() => canExportResult(props.result, props.loading))

function exportResult(format: ExportFormat) {
  if (!props.result || !canExport.value) return
  downloadExportFile(buildExportFile(props.result, format, props.dbName))
}
</script>

<style scoped>
.result-panel { padding: 14px; min-height: 220px; }
.panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.result-actions { display: flex; align-items: center; justify-content: flex-end; gap: 8px; flex-wrap: wrap; }
</style>
