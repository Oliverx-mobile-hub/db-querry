<template>
  <section class="panel nl-panel">
    <div class="panel-head">
      <h2 class="panel-title">Natural SQL</h2>
      <el-button :disabled="!dbName" :loading="loading" @click="$emit('generate', prompt)">Generate</el-button>
    </div>
    <el-input v-model="prompt" type="textarea" :rows="3" placeholder="查询用户表的所有信息" />
    <div v-if="draft" class="draft">
      <strong>Generated SQL</strong>
      <pre>{{ draft.sql }}</pre>
      <p>{{ draft.explanation }}</p>
      <el-tag :type="draft.validation.executable ? 'success' : 'danger'">{{ draft.validation.executable ? 'executable' : 'blocked' }}</el-tag>
      <el-button type="primary" @click="$emit('use-sql', draft.sql)">Use SQL</el-button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { GeneratedSqlDraft } from '../api/types'

defineProps<{ dbName: string | null; draft: GeneratedSqlDraft | null; loading: boolean }>()
defineEmits<{ generate: [prompt: string]; 'use-sql': [sql: string] }>()
const prompt = ref('')
</script>

<style scoped>
.nl-panel { padding: 14px; background: var(--dq-sunbeam); }
.panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.draft { margin-top: 12px; padding: 12px; background: var(--dq-cloud); border: var(--dq-border); }
pre { white-space: pre-wrap; overflow-wrap: anywhere; margin: 8px 0; font-family: "Cascadia Code", Consolas, monospace; }
p { margin: 8px 0; }
</style>
