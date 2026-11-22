import { computed, ref, watch } from 'vue'

export type ThemePreference = 'system' | 'light' | 'dark'
export type ResolvedTheme = 'light' | 'dark'

const STORAGE_KEY = 'formulosity_theme'
const mediaQuery =
  typeof window !== 'undefined'
    ? window.matchMedia('(prefers-color-scheme: dark)')
    : null

const themePreference = ref<ThemePreference>(readStoredPreference())
const systemTheme = ref<ResolvedTheme>(mediaQuery?.matches ? 'dark' : 'light')
const resolvedTheme = computed<ResolvedTheme>(() =>
  themePreference.value === 'system' ? systemTheme.value : themePreference.value,
)

function readStoredPreference(): ThemePreference {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'system' || stored === 'light' || stored === 'dark') {
      return stored
    }
  } catch {
    // Ignore storage access errors and follow the system theme.
  }

  return 'system'
}

function applyTheme(theme: ResolvedTheme) {
  if (typeof document === 'undefined') {
    return
  }

  document.documentElement.classList.toggle('dark', theme === 'dark')
  document.documentElement.style.colorScheme = theme
}

watch(resolvedTheme, applyTheme, { immediate: true })

function handleSystemThemeChange(event: MediaQueryListEvent) {
  systemTheme.value = event.matches ? 'dark' : 'light'
}

if (mediaQuery) {
  if (typeof mediaQuery.addEventListener === 'function') {
    mediaQuery.addEventListener('change', handleSystemThemeChange)
  } else {
    mediaQuery.addListener(handleSystemThemeChange)
  }
}

export function setThemePreference(preference: ThemePreference) {
  themePreference.value = preference

  try {
    localStorage.setItem(STORAGE_KEY, preference)
  } catch {
    // Persisting the preference is best-effort.
  }
}

export function initTheme() {
  applyTheme(resolvedTheme.value)
}

export function useTheme() {
  return {
    themePreference,
    resolvedTheme,
    isDark: computed(() => resolvedTheme.value === 'dark'),
    setThemePreference,
  }
}
