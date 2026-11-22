<template>
  <div>
    <div class="flex w-full flex-col tabs">
      <div class="flex w-full flex-col items-center">
        <div class="w-full">
          <div v-if="surveys.length === 0" class="app-muted py-8 text-center text-sm">
            {{ t('app.surveys.empty') }}
          </div>
          <div v-else class="overflow-x-auto">
            <table class="w-full text-left text-sm text-gray-500 dark:text-gray-400">
              <thead
                class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-white/5 dark:text-gray-300"
              >
                <tr>
                  <th scope="col" class="px-6 py-3">{{ t('app.surveys.nameTitle') }}</th>
                  <th scope="col" class="px-6 py-3">{{ t('app.surveys.build') }}</th>
                  <th scope="col" class="px-6 py-3">{{ t('app.surveys.delivery') }}</th>
                  <th scope="col" class="px-6 py-3">{{ t('app.surveys.share') }}</th>
                  <th scope="col" class="px-6 py-3">{{ t('app.surveys.responses') }}</th>
                  <th scope="col" class="px-6 py-3">{{ t('app.surveys.completion') }}</th>
                  <th scope="col" class="px-6 py-3">{{ t('surveyTabs.analytics') }}</th>
                </tr>
              </thead>
              <tbody>
                <SurveyRow v-for="survey in surveys" :key="survey.uuid" :survey="survey" />
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Survey } from '@/lib/types'
import SurveyRow from './SurveyRow.vue'

const { t } = useI18n()

interface Props {
  surveys: Survey[]
}

defineProps<Props>()
</script>
