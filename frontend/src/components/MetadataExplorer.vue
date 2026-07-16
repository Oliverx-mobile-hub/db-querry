<template>
  <section class="panel metadata-panel">
    <div class="panel-head">
      <h2 class="panel-title">Metadata</h2>
      <span class="muted">{{ objectCount }} objects</span>
    </div>
    <el-empty v-if="!metadata || metadata.schemas.length === 0" description="No metadata" />
    <el-collapse v-else>
      <el-collapse-item v-for="schema in metadata.schemas" :key="schema.name" :title="schema.name">
        <div v-for="object in schema.objects" :key="`${object.schema}.${object.name}`" class="object-block">
          <div class="object-title">
            <strong>{{ object.name }}</strong>
            <el-tag size="small">{{ object.type }}</el-tag>
          </div>
          <div v-for="column in object.columns" :key="column.name" class="column-row">
            <span>{{ column.name }}</span>
            <code>{{ column.dataType }}</code>
            <b v-if="column.primaryKey">PK</b>
          </div>
        </div>
      </el-collapse-item>
    </el-collapse>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { MetadataDocument } from '../api/types'

const props = defineProps<{ metadata: MetadataDocument | null }>()
const objectCount = computed(() => props.metadata?.schemas.reduce((sum, schema) => sum + schema.objects.length, 0) ?? 0)
</script>

<style scoped>
.metadata-panel { padding: 14px; min-height: 100%; overflow: auto; }
.panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.object-block { border: var(--dq-border); padding: 10px; margin-bottom: 10px; background: var(--dq-cloud); }
.object-title { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
.column-row { display: grid; grid-template-columns: minmax(0, 1fr) auto auto; gap: 8px; align-items: center; font-size: 12px; padding: 4px 0; border-top: 1px solid #ddd; }
code { background: var(--dq-fog); padding: 2px 4px; border: 1px solid var(--dq-graphite); }
b { color: var(--dq-success); }
</style>
