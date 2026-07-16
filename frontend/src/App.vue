<template>
  <main class="app-shell">
    <aside class="database-zone">
      <ConnectionPanel :dbs="dbs" :selected-name="selectedDb" :loading="loading" @select="selectDb" @add="addDb" @delete="deleteDb" />
      <MetadataExplorer :metadata="metadata" />
    </aside>
    <section class="workspace">
      <QueryEditor v-model="sqlText" :db-name="selectedDb" :loading="queryLoading" @execute="runQuery" />
      <NaturalLanguagePanel :db-name="selectedDb" :draft="draft" :loading="generateLoading" @generate="generateSql" @use-sql="sqlText = $event" />
      <ResultTable :result="result" />
    </section>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessageBox } from 'element-plus'
import { api, ApiClientError } from './api/client'
import type { DbSummary, GeneratedSqlDraft, MetadataDocument, QueryResult } from './api/types'
import ConnectionPanel from './components/ConnectionPanel.vue'
import MetadataExplorer from './components/MetadataExplorer.vue'
import QueryEditor from './components/QueryEditor.vue'
import NaturalLanguagePanel from './components/NaturalLanguagePanel.vue'
import ResultTable from './components/ResultTable.vue'

const dbs = ref<DbSummary[]>([])
const selectedDb = ref<string | null>(null)
const metadata = ref<MetadataDocument | null>(null)
const result = ref<QueryResult | null>(null)
const draft = ref<GeneratedSqlDraft | null>(null)
const sqlText = ref('SELECT * FROM users')
const loading = ref(false)
const queryLoading = ref(false)
const generateLoading = ref(false)
const error = ref('')

onMounted(loadDbs)

async function loadDbs() {
  await withError(async () => {
    const data = await api.listDbs()
    dbs.value = data.dbs
    if (!selectedDb.value && data.dbs.length > 0) await selectDb(data.dbs[0].name)
  })
}

async function addDb(payload: { name: string; url: string }) {
  loading.value = true
  await withError(async () => {
    await api.putDb(payload.name, payload.url)
    await loadDbs()
    await selectDb(payload.name)
  })
  loading.value = false
}

async function deleteDb(name: string) {
  try {
    await ElMessageBox.confirm(`删除数据库连接 ${name}？这会移除本地保存的连接和 metadata。`, 'DELETE DATABASE', {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning',
      autofocus: false,
    })
  } catch {
    return
  }
  loading.value = true
  await withError(async () => {
    await api.deleteDb(name)
    if (selectedDb.value === name) {
      selectedDb.value = null
      metadata.value = null
      result.value = null
      draft.value = null
    }
    await loadDbs()
  })
  loading.value = false
}

async function selectDb(name: string) {
  selectedDb.value = name
  await withError(async () => {
    const data = await api.getMetadata(name)
    metadata.value = data.metadata
  })
}

async function runQuery(sql: string) {
  if (!selectedDb.value) return
  queryLoading.value = true
  await withError(async () => { result.value = await api.query(selectedDb.value as string, sql) })
  queryLoading.value = false
}

async function generateSql(prompt: string) {
  if (!selectedDb.value) return
  generateLoading.value = true
  await withError(async () => { draft.value = await api.generateSql(selectedDb.value as string, prompt) })
  generateLoading.value = false
}

async function withError(action: () => Promise<void>) {
  try {
    error.value = ''
    await action()
  } catch (err) {
    error.value = err instanceof ApiClientError ? err.message : '请求失败'
    await ElMessageBox.alert(error.value, 'REQUEST FAILED', {
      type: 'error',
      confirmButtonText: 'OK',
      autofocus: false,
    })
  }
}
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(560px, 640px) minmax(0, 1fr);
  gap: 14px;
  padding: 14px;
}
.database-zone {
  display: grid;
  grid-template-columns: 280px minmax(260px, 1fr);
  gap: 14px;
  min-width: 0;
}
.workspace { display: grid; grid-template-rows: auto auto minmax(0, 1fr) auto; gap: 14px; min-width: 0; }
@media (max-width: 1180px) {
  .app-shell { grid-template-columns: 1fr; }
  .database-zone { grid-template-columns: 1fr; }
}
</style>
