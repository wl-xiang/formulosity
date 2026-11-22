<template>
  <section class="app-card">
    <h2 class="text-base font-semibold text-gray-900 dark:text-white">
      {{ t('analytics.timeline.title') }}
    </h2>

    <p v-if="points.length === 0" class="app-muted py-10 text-center text-sm">
      {{ t('analytics.timeline.empty') }}
    </p>
    <VChart v-else :option="option" autoresize class="mt-2 h-64 w-full" />
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'
import '@/lib/echarts'
import { buildTimelineOption } from '@/lib/charts'
import { useTheme } from '@/composables/useTheme'
import type { AnalyticsTimelinePoint } from '@/lib/types'

interface Props {
  points: AnalyticsTimelinePoint[]
}

const props = defineProps<Props>()
const { t } = useI18n()
const { resolvedTheme } = useTheme()

const option = computed(() =>
  buildTimelineOption(props.points, resolvedTheme.value, {
    started: t('analytics.timeline.started'),
    completed: t('analytics.timeline.completed'),
  }),
)
</script>
