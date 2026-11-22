<template>
  <div class="app-card">
    <p class="text-xs font-medium uppercase tracking-wide app-muted">{{ label }}</p>
    <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">{{ value }}</p>
    <p v-if="hint" class="mt-1 text-xs app-muted">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  label: string
  value: string
  hint?: string
}

defineProps<Props>()
</script>
