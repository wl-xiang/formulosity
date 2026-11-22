<template>
  <tr class="group border-b border-gray-200 hover:bg-gray-50 dark:border-white/10 dark:hover:bg-white/5">
    <td class="p-4">
      <div class="space-y-1">
        <div class="font-medium text-gray-900 dark:text-white">{{ survey.name }}</div>
        <div v-if="survey.config" class="text-sm text-gray-600 dark:text-gray-400">
          {{ survey.config.title }}
        </div>
        <div class="text-xs text-gray-500 dark:text-gray-400">
          {{ formatDate(survey.created_at) }}
        </div>
      </div>
    </td>
    <td class="p-4">
      <span
        :class="[
          'inline-flex cursor-pointer items-center gap-1 rounded-full px-2 py-1 text-xs font-medium',
          getParseStatusColor(survey.parse_status),
        ]"
        @click="showErrorLog = !showErrorLog"
      >
        {{ t(`app.parseStatus.${survey.parse_status}`) }}
        <Icon
          v-if="survey.parse_status === SurveyParseStatus.Error"
          :icon="showErrorLog ? 'heroicons:chevron-up' : 'heroicons:chevron-down'"
          class="h-3 w-3"
        />
      </span>
      <div v-if="showErrorLog && survey.parse_status === SurveyParseStatus.Error" class="mt-3">
        <div class="rounded-lg bg-gray-100 p-3 text-xs dark:bg-white/10">
          <div class="mb-1 font-medium text-gray-900 dark:text-gray-100">
            {{ t('app.surveys.errorDetails') }}
          </div>
          <code class="text-gray-700 dark:text-gray-300">{{ survey.error_log }}</code>
        </div>
      </div>
    </td>
    <td class="p-4">
      <button
        v-if="isLaunched || canStartSurvey"
        :class="[
          'inline-flex items-center gap-1 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors',
          isLaunched
            ? 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-white/10 dark:text-gray-200 dark:hover:bg-white/20'
            : 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-500/15 dark:text-blue-300 dark:hover:bg-blue-500/25',
        ]"
        @click="updateSurveyStatus(survey.uuid, isLaunched ? 'stopped' : 'launched')"
      >
        <Icon :icon="isLaunched ? 'heroicons:pause' : 'heroicons:play'" class="h-4 w-4" />
        {{ isLaunched ? t('app.surveys.stop') : t('app.surveys.start') }}
      </button>
      <ErrCode v-if="errorMsg" :message="errorMsg" class="mt-2" />
    </td>
    <td class="p-4">
      <a
        v-if="survey.delivery_status === SurveyDeliveryStatus.Launched"
        :href="survey.url"
        target="_blank"
        class="inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
      >
        {{ t('app.surveys.publicLink') }}
        <Icon icon="heroicons:arrow-top-right-on-square" class="h-4 w-4" />
      </a>
    </td>
    <td class="p-4">
      <router-link
        :to="`/app/surveys/${survey.uuid}/responses`"
        class="font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
      >
        {{ survey.stats.sessions_count_completed }}
      </router-link>
    </td>
    <td class="p-4 text-gray-900 dark:text-gray-100">
      {{ formatPercent(survey.stats.completion_rate, 0) }}
    </td>
    <td class="p-4">
      <router-link
        :to="`/app/surveys/${survey.uuid}/analytics`"
        class="inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
      >
        <Icon icon="heroicons:chart-bar" class="h-4 w-4" />
        {{ t('surveyTabs.analytics') }}
      </router-link>
    </td>
  </tr>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import type { Survey } from '@/lib/types'
import { SurveyDeliveryStatus, SurveyParseStatus } from '@/lib/types'
import { updateSurvey } from '@/lib/api'
import { formatDate, formatPercent } from '@/lib/format'
import ErrCode from '@/components/ui/ErrCode.vue'

const { t } = useI18n()

interface Props {
  survey: Survey
}

const props = defineProps<Props>()

const errorMsg = ref<string>('')
const showErrorLog = ref<boolean>(false)

const isLaunched = computed(
  () => props.survey.delivery_status === SurveyDeliveryStatus.Launched,
)
const canStartSurvey = computed(
  () => props.survey.parse_status === SurveyParseStatus.Success && !isLaunched.value,
)

function getParseStatusColor(status: string): string {
  const colors: Record<string, string> = {
    [SurveyParseStatus.Success]:
      'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-300',
    [SurveyParseStatus.Error]:
      'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-300',
    [SurveyParseStatus.Deleted]:
      'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300',
  }
  return colors[status] || 'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300'
}

async function updateSurveyStatus(surveyUUID: string, status: string) {
  const res = await updateSurvey(surveyUUID, {
    delivery_status: status,
  })

  if (res.error) {
    errorMsg.value = res.error
  } else {
    window.location.href = '/app'
  }
}
</script>
