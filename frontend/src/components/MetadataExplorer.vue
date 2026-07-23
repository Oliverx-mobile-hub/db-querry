<template>
  <section class="panel metadata-panel">
    <div class="panel-head">
      <h2 class="panel-title">Metadata</h2>
      <span class="muted">{{ objectCount }} objects</span>
    </div>

    <el-empty v-if="!metadata || metadata.schemas.length === 0" description="No metadata" />
    <template v-else>
      <div class="database-toolbar">
        <div class="database-title">
          <span class="database-icon">DB</span>
          <strong>{{ displayDbName }}</strong>
          <span class="dialect-chip">{{ metadata.databaseType }}</span>
        </div>
        <el-button size="small" :disabled="!dbName" @click="$emit('refresh')">Refresh</el-button>
      </div>

      <el-input v-model="filterText" class="metadata-search" clearable placeholder="Search tables, views, columns" />

      <div v-if="filteredGroups.length === 0" class="empty-filter">No matching metadata</div>
      <div v-else class="metadata-browser">
        <section v-for="group in filteredGroups" :key="group.type" class="metadata-group">
          <div class="group-title">
            <span>{{ group.label }}</span>
            <span>{{ group.objects.length }}</span>
          </div>

          <div v-for="object in group.objects" :key="objectKey(object)" class="object-entry">
            <button class="object-row" type="button" @click="toggleObject(objectKey(object))">
              <span class="disclosure" :class="{ expanded: expandedKeys.has(objectKey(object)) }" />
              <span class="object-content">
                <span class="object-name-line">
                  <strong class="object-name">{{ object.name }}</strong>
                  <span class="kind-chip">{{ object.type }}</span>
                </span>
                <span class="object-meta-line">
                  <span>{{ object.schema }}</span>
                  <span>{{ object.columns.length }} fields</span>
                </span>
              </span>
            </button>

            <div v-if="expandedKeys.has(objectKey(object))" class="column-list">
              <div v-for="column in object.columns" :key="column.name" class="column-row">
                <div class="column-head">
                  <strong class="column-name">{{ column.name }}</strong>
                  <span v-if="column.primaryKey" class="pk-chip">PK</span>
                  <span v-if="!column.nullable" class="not-null-chip">NOT NULL</span>
                </div>
                <code class="data-type">{{ column.dataType }}</code>
              </div>
            </div>
          </div>
        </section>
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { MetadataDocument, MetadataObject } from '../api/types'

type MetadataGroup = {
  type: MetadataObject['type']
  label: string
  objects: MetadataObject[]
}

const props = defineProps<{ metadata: MetadataDocument | null; dbName: string | null }>()
defineEmits<{ refresh: [] }>()

const filterText = ref('')
const expandedKeys = ref(new Set<string>())
const objectCount = computed(() => props.metadata?.schemas.reduce((sum, schema) => sum + schema.objects.length, 0) ?? 0)
const displayDbName = computed(() => props.dbName || props.metadata?.schemas[0]?.name || 'Database')
const allObjects = computed(() => props.metadata?.schemas.flatMap(schema => schema.objects) ?? [])
const filteredObjects = computed(() => {
  const filter = filterText.value.trim().toLowerCase()
  if (!filter) return allObjects.value
  return allObjects.value.filter(object => {
    const objectMatches = [object.schema, object.name, object.type].some(value => value.toLowerCase().includes(filter))
    const columnMatches = object.columns.some(column => [column.name, column.dataType].some(value => value.toLowerCase().includes(filter)))
    return objectMatches || columnMatches
  })
})
const filteredGroups = computed<MetadataGroup[]>(() => {
  const tables = filteredObjects.value.filter(object => object.type === 'table')
  const views = filteredObjects.value.filter(object => object.type === 'view')
  const groups: MetadataGroup[] = [
    { type: 'table', label: 'Tables', objects: tables },
    { type: 'view', label: 'Views', objects: views },
  ]
  return groups.filter(group => group.objects.length > 0)
})

watch(() => props.metadata, () => {
  expandedKeys.value = new Set()
  filterText.value = ''
})

watch(filterText, value => {
  if (!value.trim()) return
  expandedKeys.value = new Set(filteredObjects.value.map(objectKey))
})

function objectKey(object: MetadataObject) {
  return `${object.schema}.${object.name}`
}

function toggleObject(key: string) {
  const next = new Set(expandedKeys.value)
  if (next.has(key)) {
    next.delete(key)
  } else {
    next.add(key)
  }
  expandedKeys.value = next
}
</script>

<style scoped>
.metadata-panel {
  min-height: 100%;
  overflow: hidden;
}
.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 14px 14px 10px;
}
.database-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 12px 14px;
  border-top: var(--dq-border);
  border-bottom: var(--dq-border);
  background: var(--dq-sunbeam);
}
.database-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.database-title strong {
  min-width: 0;
  font-size: 14px;
  font-weight: 900;
  text-transform: uppercase;
  overflow-wrap: anywhere;
}
.database-icon,
.dialect-chip,
.kind-chip,
.pk-chip,
.not-null-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 20px;
  padding: 2px 6px;
  border: 1px solid var(--dq-graphite);
  background: var(--dq-cloud);
  color: var(--dq-ink);
  font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
  font-size: 11px;
  font-weight: 900;
  line-height: 1.2;
  overflow-wrap: anywhere;
  white-space: normal;
}
.metadata-search {
  display: block;
  padding: 12px 14px;
}
.metadata-browser {
  max-height: calc(100vh - 190px);
  overflow: auto;
  padding: 0 14px 14px;
}
.metadata-group {
  display: grid;
  gap: 8px;
  margin-top: 8px;
}
.group-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  color: var(--dq-slate);
  font-size: 12px;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.object-entry {
  border-left: 2px solid #d8d8d4;
  padding-left: 10px;
}
.object-row {
  width: 100%;
  display: grid;
  grid-template-columns: 14px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
  border: 0;
  background: transparent;
  color: var(--dq-ink);
  padding: 4px 0;
  text-align: left;
  cursor: pointer;
}
.object-row:hover .object-name {
  color: #005f95;
  text-decoration: underline;
  text-underline-offset: 3px;
}
.disclosure {
  width: 0;
  height: 0;
  margin-top: 6px;
  border-top: 5px solid transparent;
  border-bottom: 5px solid transparent;
  border-left: 6px solid #9ca3af;
  transition: transform 120ms ease;
}
.disclosure.expanded {
  transform: rotate(90deg);
}
.object-content {
  display: grid;
  gap: 4px;
  min-width: 0;
}
.object-name-line {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  min-width: 0;
  flex-wrap: wrap;
}
.object-name {
  min-width: 0;
  color: #00639b;
  font-size: 13px;
  line-height: 1.35;
  overflow-wrap: anywhere;
}
.object-meta-line {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  color: var(--dq-slate);
  font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
  font-size: 11px;
  line-height: 1.3;
}
.object-meta-line span {
  max-width: 100%;
  padding: 2px 5px;
  border: 1px solid #cfcfca;
  background: var(--dq-fog);
  overflow-wrap: anywhere;
}
.column-list {
  display: grid;
  gap: 8px;
  margin: 6px 0 10px 18px;
  padding-left: 14px;
  border-left: 1px solid #d8d8d4;
}
.column-row {
  display: grid;
  gap: 4px;
  min-width: 0;
  padding: 4px 0;
}
.column-head {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  flex-wrap: wrap;
  min-width: 0;
}
.column-name {
  min-width: 0;
  font-size: 12px;
  line-height: 1.35;
  overflow-wrap: anywhere;
}
.data-type {
  justify-self: start;
  max-width: 100%;
  padding: 3px 6px;
  border: 1px solid var(--dq-graphite);
  background: var(--dq-cloud);
  color: var(--dq-ink);
  font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
  font-size: 11px;
  line-height: 1.35;
  overflow-wrap: anywhere;
  white-space: normal;
  word-break: break-word;
}
.pk-chip {
  color: #b01919;
  border-color: #b01919;
}
.not-null-chip {
  color: #5b3aa4;
  border-color: #9c88d8;
  background: #f4efff;
}
.empty-filter {
  margin: 0 14px 14px;
  padding: 18px 10px;
  border: 1px dashed #b8b8b3;
  color: var(--dq-slate);
  font-size: 12px;
  font-weight: 800;
  text-align: center;
}
</style>
