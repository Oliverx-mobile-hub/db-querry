<template>
  <section class="panel connection-panel">
    <div class="panel-head">
      <h2 class="panel-title">Databases</h2>
      <el-button size="small" type="primary" @click="dialogVisible = true">Add</el-button>
    </div>

    <div class="db-list">
      <button
        v-for="db in dbs"
        :key="db.name"
        class="db-item"
        :class="{ active: db.name === selectedName }"
        @click="$emit('select', db.name)"
      >
        <span class="db-row">
          <strong>{{ db.name }}</strong>
          <el-button class="delete-button" size="small" plain @click.stop="$emit('delete', db.name)">Delete</el-button>
        </span>
        <span>{{ db.displayDsn }}</span>
        <span class="status-pill" :class="`status-${db.connectionStatus}`">
          <span class="status-dot" />
          {{ statusLabel(db.connectionStatus) }}
        </span>
        <small class="metadata-note">metadata: {{ db.metadataStatus }}</small>
      </button>
    </div>

    <el-dialog v-model="dialogVisible" title="ADD DATABASE" width="520px">
      <el-form label-position="top" @submit.prevent>
        <el-form-item label="Name">
          <el-input v-model="name" placeholder="local" />
        </el-form-item>
        <el-form-item label="Postgres URL">
          <el-input v-model="url" placeholder="postgres://postgres:postgre@localhost:5432/postgres" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="loading" @click="submit">Save</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { DbSummary } from '../api/types'

defineProps<{ dbs: DbSummary[]; selectedName: string | null; loading: boolean }>()
const emit = defineEmits<{ select: [name: string]; add: [payload: { name: string; url: string }]; delete: [name: string] }>()

const dialogVisible = ref(false)
const name = ref('local')
const url = ref('')

function submit() {
  emit('add', { name: name.value.trim(), url: url.value.trim() })
  dialogVisible.value = false
}

function statusLabel(status: DbSummary['connectionStatus']) {
  if (status === 'online') return 'ONLINE'
  if (status === 'offline') return 'OFFLINE'
  return 'CHECKING'
}
</script>

<style scoped>
.connection-panel { padding: 14px; min-height: 100%; }
.panel-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 14px; }
.db-list { display: grid; gap: 10px; }
.db-item {
  width: 100%;
  display: grid;
  gap: 6px;
  padding: 12px;
  background: var(--dq-fog);
  border: var(--dq-border);
  border-radius: var(--dq-radius);
  text-align: left;
  cursor: pointer;
}
.db-item.active { background: var(--dq-soft-blue); box-shadow: 4px 4px 0 var(--dq-graphite); }
.db-item strong { font-size: 14px; text-transform: uppercase; }
.db-item span { color: var(--dq-slate); font-size: 12px; overflow-wrap: anywhere; }
.db-row { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.delete-button {
  padding: 4px 8px !important;
  background: #fff4f4 !important;
  border-width: 1px !important;
  color: var(--dq-error) !important;
}
.status-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  min-height: 28px;
  padding: 5px 10px;
  border: var(--dq-border);
  border-radius: var(--dq-radius);
  color: var(--dq-ink) !important;
  font-size: 12px !important;
  font-weight: 900;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  box-shadow: 3px 3px 0 var(--dq-graphite);
}
.status-online {
  background: var(--dq-success);
}
.status-offline {
  background: var(--dq-error);
}
.status-unknown {
  background: var(--dq-sunbeam);
}
.status-dot {
  width: 8px;
  height: 8px;
  border: 1px solid var(--dq-graphite);
  border-radius: 50%;
  background: var(--dq-cloud);
}
.metadata-note {
  color: var(--dq-slate);
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
</style>
