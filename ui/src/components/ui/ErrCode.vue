<template>
  <div
    class="mb-4 flex items-center rounded-lg border border-orange-300 bg-orange-50 p-4 text-sm text-orange-800 dark:border-orange-500/30 dark:bg-orange-500/10 dark:text-orange-200"
  >
    <Icon icon="heroicons:information-circle" class="me-3 inline h-4 w-4 flex-shrink-0" />
    <span class="sr-only">{{ t('common.error') }}</span>
    <div>
      <p>{{ displayMessage }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'

interface Props {
  message?: string
}

const props = withDefaults(defineProps<Props>(), {
  message: undefined,
})

const { t } = useI18n()
const displayMessage = computed(() => props.message ?? t('common.error'))
</script>
