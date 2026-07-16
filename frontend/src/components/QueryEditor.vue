<template>
  <section class="panel query-panel">
    <div class="panel-head">
      <h2 class="panel-title">SQL Editor</h2>
      <el-button type="primary" :disabled="!dbName" :loading="loading" @click="$emit('execute', sql)">Run</el-button>
    </div>
    <textarea v-model="sql" class="sql-editor" spellcheck="false" />
  </section>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{ dbName: string | null; loading: boolean; modelValue: string }>()
const emit = defineEmits<{ execute: [sql: string]; 'update:modelValue': [sql: string] }>()
const sql = ref(props.modelValue || 'SELECT * FROM users')

watch(() => props.modelValue, value => { if (value !== sql.value) sql.value = value })
watch(sql, value => emit('update:modelValue', value))
</script>

<style scoped>
.query-panel { padding: 14px; }
.panel-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.sql-editor {
  width: 100%;
  min-height: 220px;
  resize: vertical;
  background: #111;
  color: #f8f8f2;
  border: var(--dq-border);
  border-radius: var(--dq-radius);
  padding: 14px;
  font-family: "Cascadia Code", Consolas, monospace;
  font-size: 14px;
  line-height: 1.55;
}
</style>
