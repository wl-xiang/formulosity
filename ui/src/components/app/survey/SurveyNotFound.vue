<template>
  <div class="m-auto flex items-center justify-center gap-4">
    <h1 class="h1">{{ t('survey.notFound') }}</h1>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
