<template>
  <section class="panel result-panel">
    <div class="panel-head">
      <h2 class="panel-title">Results</h2>
      <span v-if="result" class="muted">{{ result.rowCount }} rows · {{ result.durationMs }} ms</span>
    </div>
    <el-empty v-if="!result" description="Run a query" />
    <el-empty v-else-if="result.empty" description="No rows" />
    <el-table v-else :data="result.rows" height="320">
      <el-table-column v-for="column in result.columns" :key="column.name" :prop="column.name" :label="column.name" min-width="160" />
    </el-table>
  </section>
</template>

<script setup lang="ts">
import type { QueryResult } from '../api/types'
defineProps<{ result: QueryResult | null }>()
</script>

<style scoped>
.result-panel { padding: 14px; min-height: 220px; }
.panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
