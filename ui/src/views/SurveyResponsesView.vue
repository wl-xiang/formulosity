<template>
  <AppLayout>
    <div v-if="loading">{{ t('common.loading') }}</div>
    <ErrCode v-else-if="errorDisplay" :message="errorDisplay" />
    <div v-else-if="survey" class="survey-responses">
      <SurveyTabs :survey-uuid="survey.uuid" class="mb-6" />
      <h1 class="mb-2 text-2xl font-bold">{{ t('responses.title') }}</h1>
      <h2 class="mb-4 text-lg text-gray-600 dark:text-gray-400">{{ survey.name }}</h2>
      <p class="mb-6 text-gray-500 dark:text-gray-400">
        {{ t('responses.responseCount', { count: survey.stats.sessions_count_completed }) }}
      </p>

      <!-- Export button -->
      <div class="mb-6">
        <button
          :disabled="downloading"
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
          @click="exportResponses"
        >
          <Icon icon="heroicons:arrow-down-tray" class="h-4 w-4" />
          {{ downloading ? t('responses.exporting') : t('responses.exportJson') }}
        </button>
      </div>

      <div class="mt-8">
        <div v-if="sessions.length === 0" class="py-8 text-center">
          <p>{{ t('responses.noResponses') }}</p>
        </div>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm text-gray-500 dark:text-gray-400">
            <thead
              class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-white/5 dark:text-gray-300"
            >
              <tr>
                <th
                  v-for="col in columns"
                  :key="col.key"
                  scope="col"
                  class="cursor-pointer p-4 hover:bg-gray-100 dark:hover:bg-white/10"
                  @click="sortBy(col.key)"
                >
                  <div class="flex items-center gap-1">
                    <Icon
                      v-if="sortByField === col.key"
                      :icon="sortOrder === 'asc' ? 'heroicons:chevron-up' : 'heroicons:chevron-down'"
                      class="h-3 w-3"
                    />
                    {{ col.label }}
                  </div>
                </th>
                <th scope="col" class="p-4">{{ t('responses.webhookStatus') }}</th>
                <th scope="col" class="p-4">{{ t('responses.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="session in sessions"
                :key="session.uuid"
                class="group border-b border-gray-200 hover:bg-gray-50 dark:border-white/10 dark:hover:bg-white/5"
              >
                <td class="p-4 text-gray-900 dark:text-gray-100">{{ session.uuid }}</td>
                <td class="p-4">
                  <span
                    :class="[
                      'inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs font-medium',
                      session.status === SurveySessionStatus.Completed
                        ? 'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-300'
                        : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/15 dark:text-yellow-300',
                    ]"
                  >
                    <Icon
                      :icon="
                        session.status === SurveySessionStatus.Completed
                          ? 'heroicons:check-circle'
                          : 'heroicons:clock'
                      "
                      class="h-3 w-3"
                    />
                    {{
                      session.status === SurveySessionStatus.Completed
                        ? t('responses.completed')
                        : t('responses.inProgress')
                    }}
                  </span>
                </td>
                <td class="p-4">{{ formatDateTime(session.created_at) }}</td>
                <td class="p-4">
                  {{ session.completed_at ? formatDateTime(session.completed_at) : '-' }}
                </td>
                <td class="p-4">{{ session.webhookData.statusCode || '-' }}</td>
                <td class="p-4">
                  <div class="flex items-center gap-2">
                    <button
                      class="inline-flex items-center gap-1 rounded bg-blue-100 px-2 py-1 text-xs text-blue-700 hover:bg-blue-200 dark:bg-blue-500/15 dark:text-blue-300 dark:hover:bg-blue-500/25"
                      @click="viewSession = session"
                    >
                      <Icon icon="heroicons:eye" class="h-3 w-3" />
                      {{ t('responses.view') }}
                    </button>
                    <button
                      class="inline-flex items-center gap-1 rounded bg-gray-100 px-2 py-1 text-xs text-gray-700 hover:bg-gray-200 dark:bg-white/10 dark:text-gray-200 dark:hover:bg-white/20"
                      @click="deleteSession(session)"
                    >
                      <Icon icon="heroicons:trash" class="h-3 w-3" />
                      {{ t('responses.delete') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="mt-6 flex justify-center">
          <div class="flex items-center gap-1">
            <button
              v-for="page in totalPages"
              :key="page"
              :class="[
                'rounded px-3 py-1 text-sm',
                page === currentPage
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-white/10 dark:text-gray-200 dark:hover:bg-white/20',
              ]"
              @click="changePage(page)"
            >
              {{ page }}
            </button>
          </div>
        </div>
      </div>

      <!-- Modal for viewing session details -->
      <div
        v-if="viewSession"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="viewSession = undefined"
      >
        <div
          class="max-h-96 max-w-4xl overflow-y-auto rounded-lg bg-white p-6 text-gray-900 dark:bg-gray-900 dark:text-gray-100"
        >
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-semibold">
              {{ t('responses.modalTitle', { uuid: viewSession.uuid }) }}
            </h3>
            <button
              class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
              @click="viewSession = undefined"
            >
              <Icon icon="heroicons:x-mark" class="h-6 w-6" />
            </button>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead
                class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-white/5 dark:text-gray-300"
              >
                <tr>
                  <th class="p-3 text-left">{{ t('responses.questionId') }}</th>
                  <th class="p-3 text-left">{{ t('responses.question') }}</th>
                  <th class="p-3 text-left">{{ t('responses.response') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="answer in viewSession.question_answers"
                  :key="answer.question_uuid"
                  class="border-b border-gray-200 dark:border-white/10"
                >
                  <td class="p-3">{{ answer.question_id }}</td>
                  <td class="p-3">{{ getQuestionLabel(answer.question_uuid) }}</td>
                  <td class="p-3">
                    <div v-if="isFileAnswer(answer.question_uuid)">
                      <button
                        class="inline-flex items-center gap-1 rounded bg-blue-100 px-2 py-1 text-xs text-blue-700 hover:bg-blue-200 dark:bg-blue-500/15 dark:text-blue-300 dark:hover:bg-blue-500/25"
                        @click="downloadFile(answer.answer.value as string)"
                      >
                        <Icon icon="heroicons:arrow-down-tray" class="h-3 w-3" />
                        {{ t('common.download') }}
                      </button>
                    </div>
                    <div v-else>{{ formatAnswer(answer.answer.value) }}</div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <button
            class="mt-4 flex items-center gap-1 rounded bg-gray-500 px-4 py-2 text-white hover:bg-gray-600"
            @click="viewSession = undefined"
          >
            <Icon icon="heroicons:x-mark" class="h-4 w-4" />
            {{ t('responses.close') }}
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import AppLayout from '@/components/app/AppLayout.vue'
import SurveyTabs from '@/components/app/SurveyTabs.vue'
import ErrCode from '@/components/ui/ErrCode.vue'
import { deleteSurveySession, download, getSurveySessions, getSurveys } from '@/lib/api'
import { formatDateTime } from '@/lib/format'
import { translateApiError } from '@/lib/errors'
import type { Survey, SurveySession } from '@/lib/types'
import { SurveyQuestionType, SurveySessionStatus, SurveySessionsLimit } from '@/lib/types'

const { t } = useI18n()
const route = useRoute()

const survey = ref<Survey | null>(null)
const sessions = ref<SurveySession[]>([])
const errorKey = ref<string>('')
const errorMsg = ref<string>('')
const loading = ref<boolean>(true)
const downloading = ref<boolean>(false)
const viewSession = ref<SurveySession | undefined>(undefined)

const errorDisplay = computed(() => (errorKey.value ? t(errorKey.value) : errorMsg.value))

// Pagination and sorting
const currentPage = ref<number>(1)
const sortByField = ref<string>('created_at')
const sortOrder = ref<string>('desc')

const columns = computed(() => [
  { label: t('responses.sessionId'), key: 'uuid' },
  { label: t('responses.status'), key: 'status' },
  { label: t('responses.startedAt'), key: 'created_at' },
  { label: t('responses.completedAt'), key: 'completed_at' },
])

const totalPages = computed(() => {
  if (!survey.value) return 1
  return survey.value.pages_count || 1
})

onMounted(async () => {
  const surveyUuid = route.params.surveyUuid as string
  if (!surveyUuid) return

  // First get the survey details
  const surveysResp = await getSurveys()
  if (surveysResp.error) {
    errorKey.value = 'responses.failedToFetch'
    loading.value = false
    return
  }

  const foundSurvey = surveysResp.data.data.find((s: Survey) => s.uuid === surveyUuid)
  if (!foundSurvey) {
    errorKey.value = 'responses.surveyNotFound'
    loading.value = false
    return
  }

  survey.value = foundSurvey
  await fetchResponses(currentPage.value, sortByField.value, sortOrder.value)
  loading.value = false
})

async function fetchResponses(page: number, sortBy: string, order: string) {
  if (!survey.value) return

  errorKey.value = ''
  errorMsg.value = ''
  const limit = SurveySessionsLimit
  const offset = (page - 1) * limit

  const sessionsResp = await getSurveySessions(
    survey.value.uuid,
    `limit=${limit}&offset=${offset}&sort_by=${sortBy}&order=${order}`,
  )

  if (sessionsResp.error) {
    errorKey.value = 'responses.unableToLoad'
  } else {
    sessions.value = sessionsResp.data.data.sessions || []
    // Update pages count from the response
    if (sessionsResp.data.data.pages_count) {
      survey.value.pages_count = sessionsResp.data.data.pages_count
    }
  }
}

async function changePage(page: number) {
  currentPage.value = page
  await fetchResponses(page, sortByField.value, sortOrder.value)
}

async function sortBy(field: string) {
  let newOrder = 'asc'
  if (sortByField.value === field) {
    newOrder = sortOrder.value === 'asc' ? 'desc' : 'asc'
  }

  sortByField.value = field
  sortOrder.value = newOrder
  currentPage.value = 1

  await fetchResponses(currentPage.value, field, newOrder)
}

async function exportResponses() {
  if (!survey.value) return

  downloading.value = true
  errorKey.value = ''
  errorMsg.value = ''

  const allSessionsResp = await getSurveySessions(
    survey.value.uuid,
    'limit=1000000&offset=0&sort_by=created_at&order=desc',
  )

  downloading.value = false

  if (allSessionsResp.error) {
    errorKey.value = 'responses.unableToExport'
    return
  }

  const element = document.createElement('a')
  const file = new Blob([JSON.stringify(allSessionsResp.data.data.sessions, null, 2)], {
    type: 'application/json',
  })
  element.href = URL.createObjectURL(file)
  element.download = t('responses.exportFileName')
  document.body.appendChild(element)
  element.click()
  document.body.removeChild(element)
  URL.revokeObjectURL(element.href)
}

async function deleteSession(session: SurveySession) {
  if (!survey.value) return

  errorKey.value = ''
  errorMsg.value = ''

  const deleteResp = await deleteSurveySession(survey.value.uuid, session.uuid)

  if (deleteResp.error) {
    errorMsg.value = translateApiError(deleteResp.error)
  } else {
    await fetchResponses(1, sortByField.value, sortOrder.value)
    currentPage.value = 1
  }
}

function getQuestionLabel(questionUuid: string): string {
  if (!survey.value?.config?.questions?.questions) return t('common.unknownQuestion')
  const question = survey.value.config.questions.questions.find((q) => q.uuid === questionUuid)
  return question ? question.label : t('common.unknownQuestion')
}

function isFileAnswer(questionUuid: string): boolean {
  if (!survey.value?.config?.questions?.questions) return false
  const question = survey.value.config.questions.questions.find((q) => q.uuid === questionUuid)
  return question?.type === SurveyQuestionType.File
}

function formatAnswer(value: string | string[] | number | boolean): string {
  if (Array.isArray(value)) {
    return value.join(', ')
  }
  if (typeof value === 'boolean') {
    return value ? t('common.yes') : t('common.no')
  }
  return String(value)
}

async function downloadFile(filePath: string) {
  if (!survey.value) return
  const fileName = filePath.substring(filePath.lastIndexOf('/') + 1)
  await download(survey.value.uuid, fileName)
}
</script>
