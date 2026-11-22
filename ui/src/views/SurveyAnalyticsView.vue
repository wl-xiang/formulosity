<template>
  <AppLayout>
    <div v-if="loading">{{ t('common.loading') }}</div>
    <ErrCode v-else-if="errorDisplay" :message="errorDisplay" />
    <div v-else-if="analytics">
      <SurveyTabs :survey-uuid="analytics.survey_uuid" class="mb-6" />

      <div class="mb-6 flex flex-col gap-1">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('analytics.title') }}
        </h1>
        <p class="text-lg text-gray-600 dark:text-gray-400">{{ analytics.survey_name }}</p>
        <p class="app-muted text-sm">{{ t('analytics.subtitle') }}</p>
      </div>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <StatCard
          :label="t('analytics.cards.total')"
          :value="formatNumber(analytics.totals.sessions_total)"
        />
        <StatCard
          :label="t('analytics.cards.completed')"
          :value="formatNumber(analytics.totals.sessions_completed)"
        />
        <StatCard
          :label="t('analytics.cards.inProgress')"
          :value="formatNumber(analytics.totals.sessions_in_progress)"
        />
        <StatCard
          :label="t('analytics.cards.completionRate')"
          :value="formatPercent(analytics.totals.completion_rate)"
          :hint="`${analytics.totals.sessions_completed} / ${analytics.totals.sessions_total}`"
        />
        <StatCard
          :label="t('analytics.cards.averageCompletionTime')"
          :value="averageDuration"
        />
      </div>

      <div class="mt-6">
        <ResponseTimelineChart :points="analytics.timeline" />
      </div>

      <section class="mt-6">
        <div class="mb-3 flex flex-wrap items-end justify-between gap-2">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('analytics.questions.title') }}
          </h2>
          <span class="app-muted text-xs">
            {{ t('analytics.cards.total') }}: {{ formatNumber(analytics.totals.sessions_total) }}
          </span>
        </div>

        <div
          v-if="analytics.totals.sessions_total === 0"
          class="app-card py-10 text-center"
        >
          <p class="text-base font-medium text-gray-900 dark:text-white">
            {{ t('analytics.emptyTitle') }}
          </p>
          <p class="app-muted mt-2 text-sm">{{ t('analytics.emptyDescription') }}</p>
        </div>

        <div
          v-else-if="analytics.questions.length === 0"
          class="app-card app-muted py-8 text-center text-sm"
        >
          {{ t('analytics.questions.noChartableQuestions') }}
        </div>

        <div v-else class="grid grid-cols-1 gap-4 xl:grid-cols-2">
          <QuestionChartCard
            v-for="question in analytics.questions"
            :key="question.question_uuid"
            :question="question"
          />
        </div>

        <p class="app-muted mt-4 text-xs">
          {{ t('analytics.questions.textQuestionsNote') }}
        </p>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/app/AppLayout.vue'
import SurveyTabs from '@/components/app/SurveyTabs.vue'
import ErrCode from '@/components/ui/ErrCode.vue'
import StatCard from '@/components/app/analytics/StatCard.vue'
import ResponseTimelineChart from '@/components/app/analytics/ResponseTimelineChart.vue'
import QuestionChartCard from '@/components/app/analytics/QuestionChartCard.vue'
import { getSurveyAnalytics } from '@/lib/api'
import {
  formatDuration,
  formatNumber,
  formatPercent,
} from '@/lib/format'
import type { SurveyAnalytics } from '@/lib/types'

const route = useRoute()
const { t } = useI18n()

const analytics = ref<SurveyAnalytics | null>(null)
const loading = ref<boolean>(true)
const errorKey = ref<string>('')

const errorDisplay = computed(() => (errorKey.value ? t(errorKey.value) : ''))

const averageDuration = computed(() => {
  const seconds = analytics.value?.totals.average_completion_seconds
  if (seconds === undefined) {
    return t('common.notAvailable')
  }

  return formatDuration(seconds)
})

onMounted(async () => {
  const surveyUuid = route.params.surveyUuid as string
  if (!surveyUuid) {
    errorKey.value = 'analytics.surveyNotFound'
    loading.value = false
    return
  }

  const response = await getSurveyAnalytics(surveyUuid)
  if (response.error) {
    errorKey.value =
      response.status === 404 ? 'analytics.surveyNotFound' : 'analytics.failedToLoad'
  } else if (response.data?.data) {
    analytics.value = response.data.data as SurveyAnalytics
  } else {
    errorKey.value = 'analytics.failedToLoad'
  }

  loading.value = false
})
</script>
