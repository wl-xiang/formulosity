<template>
  <article :class="['app-card', question.chart_kind === 'ranking' ? 'xl:col-span-2' : '']">
    <header class="mb-4 flex flex-col gap-2">
      <h3 class="text-base font-semibold text-gray-900 dark:text-white">
        {{ question.label }}
      </h3>
      <div class="flex flex-wrap items-center gap-2 text-xs">
        <span
          class="rounded-full bg-blue-100 px-2 py-0.5 font-medium text-blue-700 dark:bg-blue-500/15 dark:text-blue-300"
        >
          {{ t(`questionTypes.${question.type}`) }}
        </span>
        <span class="app-muted">
          {{ t('analytics.questions.answered', { count: question.answered_count }) }}
        </span>
        <span
          v-if="question.conditional"
          class="rounded-full bg-purple-100 px-2 py-0.5 font-medium text-purple-700 dark:bg-purple-500/15 dark:text-purple-300"
        >
          {{ t('analytics.questions.conditional') }}
        </span>
      </div>
    </header>

    <div
      v-if="question.answered_count === 0"
      class="app-muted py-10 text-center text-sm"
    >
      {{ t('analytics.questions.noAnswers') }}
    </div>

    <template v-else>
      <VChart v-if="chartOption" :option="chartOption" autoresize class="h-72 w-full" />

      <div
        v-if="question.numeric_summary"
        class="mt-3 flex flex-wrap gap-4 text-xs app-muted"
      >
        <span>
          {{ t('analytics.summary.average') }}:
          {{ formatOptional(question.numeric_summary.average) }}
        </span>
        <span>
          {{ t('analytics.summary.median') }}:
          {{ formatOptional(question.numeric_summary.median) }}
        </span>
        <span>
          {{ t('analytics.summary.mode') }}:
          {{ formatOptional(question.numeric_summary.mode) }}
        </span>
      </div>

      <div v-if="question.buckets?.length" class="mt-3 max-h-48 overflow-y-auto">
        <table class="w-full text-left text-xs">
          <thead class="text-gray-500 dark:text-gray-400">
            <tr>
              <th class="py-1 pr-2 font-medium">{{ t('analytics.chart.value') }}</th>
              <th class="px-2 py-1 text-right font-medium">
                {{ t('analytics.chart.count') }}
              </th>
              <th class="py-1 pl-2 text-right font-medium">
                {{ t('analytics.chart.percentage') }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="bucket in question.buckets"
              :key="bucket.key"
              class="border-t border-gray-200 dark:border-white/10"
            >
              <td class="py-1 pr-2 text-gray-900 dark:text-gray-100">
                {{ bucketLabel(bucket) }}
              </td>
              <td class="px-2 py-1 text-right">{{ formatNumber(bucket.count) }}</td>
              <td class="py-1 pl-2 text-right">{{ formatPercent(bucket.percentage) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="question.ranking?.length" class="mt-3 max-h-56 overflow-y-auto">
        <table class="w-full text-left text-xs">
          <thead class="text-gray-500 dark:text-gray-400">
            <tr>
              <th class="py-1 pr-2 font-medium">{{ t('analytics.ranking.option') }}</th>
              <th class="py-1 pl-2 text-right font-medium">
                {{ t('analytics.summary.averageRank') }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="option in question.ranking"
              :key="option.key"
              class="border-t border-gray-200 dark:border-white/10"
            >
              <td class="py-1 pr-2 text-gray-900 dark:text-gray-100">{{ option.label }}</td>
              <td class="py-1 pl-2 text-right">{{ formatOptional(option.average_rank) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <p
        v-if="question.chart_kind === 'multi_choice'"
        class="app-muted mt-3 text-xs"
      >
        {{ t('analytics.notes.multiSelect') }}
      </p>
      <p v-if="question.chart_kind === 'ranking'" class="app-muted mt-3 text-xs">
        {{ t('analytics.notes.ranking') }}
      </p>
    </template>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'
import '@/lib/echarts'
import {
  bucketDisplayLabel,
  buildQuestionChartOption,
} from '@/lib/charts'
import type { ChartLabels } from '@/lib/charts'
import { formatNumber, formatPercent } from '@/lib/format'
import { useTheme } from '@/composables/useTheme'
import type { AnalyticsBucket, QuestionAnalytics } from '@/lib/types'

interface Props {
  question: QuestionAnalytics
}

const props = defineProps<Props>()
const { t } = useI18n()
const { resolvedTheme } = useTheme()

const labels = computed<ChartLabels>(() => ({
  started: '',
  completed: '',
  yes: t('common.yes'),
  no: t('common.no'),
  rank: (n: number) => t('analytics.ranking.rank', { n }),
}))

const chartOption = computed(() =>
  buildQuestionChartOption(props.question, resolvedTheme.value, labels.value),
)

function bucketLabel(bucket: AnalyticsBucket): string {
  return bucketDisplayLabel(props.question, bucket, labels.value)
}

function formatOptional(value?: number): string {
  if (value === undefined) {
    return t('common.notAvailable')
  }

  return formatNumber(value, { maximumFractionDigits: 2 })
}
</script>
