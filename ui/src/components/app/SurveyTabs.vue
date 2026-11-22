<template>
  <nav
    class="flex items-center gap-1 border-b border-gray-200 dark:border-white/10"
    :aria-label="t('surveyTabs.responses')"
  >
    <router-link
      :to="`/app/surveys/${surveyUuid}/responses`"
      class="border-b-2 border-transparent px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
      active-class="border-blue-600 text-blue-600 dark:border-blue-400 dark:text-blue-400"
    >
      {{ t('surveyTabs.responses') }}
    </router-link>
    <router-link
      :to="`/app/surveys/${surveyUuid}/analytics`"
      class="border-b-2 border-transparent px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
      active-class="border-blue-600 text-blue-600 dark:border-blue-400 dark:text-blue-400"
    >
      {{ t('surveyTabs.analytics') }}
    </router-link>
  </nav>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface Props {
  surveyUuid: string
}

defineProps<Props>()
</script>
