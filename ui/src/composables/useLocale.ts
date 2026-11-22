import { ref } from 'vue'
import { i18n } from '@/locales'

export type LocalePreference = 'system' | 'en' | 'zh-CN'
export type ResolvedLocale = 'en' | 'zh-CN'

const STORAGE_KEY = 'formulosity_locale'
const supportedLocales: ResolvedLocale[] = ['en', 'zh-CN']

const localePreference = ref<LocalePreference>(readStoredPreference())
const resolvedLocale = ref<ResolvedLocale>(resolveLocale(localePreference.value))

function readStoredPreference(): LocalePreference {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'system' || stored === 'en' || stored === 'zh-CN') {
      return stored
    }
  } catch {
    // Ignore storage access errors and fall back to the system locale.
  }

  return 'system'
}

function detectSystemLocale(): ResolvedLocale {
  if (typeof navigator === 'undefined') {
    return 'en'
  }

  const languages = navigator.languages?.length
    ? navigator.languages
    : [navigator.language]

  for (const language of languages) {
    const normalized = String(language).toLowerCase()
    if (normalized.startsWith('zh')) {
      return 'zh-CN'
    }
    if (normalized.startsWith('en')) {
      return 'en'
    }
  }

  return 'en'
}

function resolveLocale(preference: LocalePreference): ResolvedLocale {
  if (preference === 'system') {
    return detectSystemLocale()
  }

  return supportedLocales.includes(preference) ? preference : 'en'
}

function applyLocale(locale: ResolvedLocale) {
  i18n.global.locale.value = locale

  if (typeof document !== 'undefined') {
    document.documentElement.lang = locale
  }
}

export function setLocalePreference(preference: LocalePreference) {
  localePreference.value = preference
  resolvedLocale.value = resolveLocale(preference)
  applyLocale(resolvedLocale.value)

  try {
    localStorage.setItem(STORAGE_KEY, preference)
  } catch {
    // Persisting the preference is best-effort.
  }
}

export function initLocale() {
  resolvedLocale.value = resolveLocale(localePreference.value)
  applyLocale(resolvedLocale.value)
}

export function useLocale() {
  return {
    localePreference,
    resolvedLocale,
    setLocalePreference,
  }
}
