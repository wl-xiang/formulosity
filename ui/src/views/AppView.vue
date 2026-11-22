<template>
  <AppLayout>
    <div v-if="loading">{{ t('common.loading') }}</div>
    <ErrCode v-else-if="errorKey" :message="t(errorKey)" />
    <SurveysPage v-else :surveys="surveys" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/app/AppLayout.vue'
import SurveysPage from '@/components/app/SurveysPage.vue'
import ErrCode from '@/components/ui/ErrCode.vue'
import { getSurveys } from '@/lib/api'
import type { Survey } from '@/lib/types'

const { t } = useI18n()

const surveys = ref<Survey[]>([])
const errorKey = ref<string>('')
const loading = ref<boolean>(true)

onMounted(async () => {
  const surveysResp = await getSurveys()
  if (surveysResp.error) {
    errorKey.value = 'app.surveys.failedToFetch'
  } else {
    surveys.value = surveysResp.data.data
  }
  loading.value = false
})
</script>
