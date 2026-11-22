<template>
  <div class="flex flex-wrap items-center justify-end gap-2">
    <label class="sr-only" for="locale-select">{{ t('preferences.language') }}</label>
    <select
      id="locale-select"
      :value="localePreference"
      :title="t('preferences.language')"
      class="rounded-lg border border-gray-300 bg-white px-2 py-1.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100"
      @change="onLocaleChange"
    >
      <option value="system">{{ t('preferences.followSystem') }}</option>
      <option value="en">{{ t('preferences.english') }}</option>
      <option value="zh-CN">{{ t('preferences.chineseSimplified') }}</option>
    </select>

    <div
      class="flex items-center gap-1 rounded-lg border border-gray-300 bg-white p-1 dark:border-gray-700 dark:bg-gray-900"
      role="group"
      :aria-label="t('preferences.theme')"
    >
      <button
        v-for="option in themeOptions"
        :key="option.value"
        type="button"
        :aria-label="t(option.label)"
        :title="t(option.label)"
        :aria-pressed="themePreference === option.value"
        :class="[
          'rounded-md p-1.5 transition-colors',
          themePreference === option.value
            ? 'bg-blue-600 text-white'
            : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-white/10',
        ]"
        @click="setThemePreference(option.value)"
      >
        <Icon :icon="option.icon" class="h-4 w-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import { useLocale } from '@/composables/useLocale'
import type { LocalePreference } from '@/composables/useLocale'
import { useTheme } from '@/composables/useTheme'
import type { ThemePreference } from '@/composables/useTheme'

const { t } = useI18n()
const { localePreference, setLocalePreference } = useLocale()
const { themePreference, setThemePreference } = useTheme()

const themeOptions: Array<{
  value: ThemePreference
  icon: string
  label: string
}> = [
  {
    value: 'system',
    icon: 'heroicons:computer-desktop',
    label: 'preferences.switchToSystem',
  },
  { value: 'light', icon: 'heroicons:sun', label: 'preferences.switchToLight' },
  { value: 'dark', icon: 'heroicons:moon', label: 'preferences.switchToDark' },
]

function onLocaleChange(event: Event) {
  const value = (event.target as HTMLSelectElement).value as LocalePreference
  setLocalePreference(value)
}
</script>
