import { i18n } from '@/locales'

function currentLocale(): string {
  return i18n.global.locale.value
}

function parseDateValue(value: string): Date {
  if (/^\d{4}$/.test(value)) {
    return new Date(Number(value), 0, 1)
  }

  if (/^\d{4}-\d{2}$/.test(value)) {
    const [year, month] = value.split('-').map(Number)
    return new Date(year, month - 1, 1)
  }

  if (/^\d{4}-\d{2}-\d{2}$/.test(value)) {
    const [year, month, day] = value.split('-').map(Number)
    return new Date(year, month - 1, day)
  }

  return new Date(value)
}

export function formatDate(
  value: string | Date,
  options: Intl.DateTimeFormatOptions = { dateStyle: 'medium' },
): string {
  const date = value instanceof Date ? value : parseDateValue(value)
  if (Number.isNaN(date.getTime())) {
    return String(value)
  }

  return new Intl.DateTimeFormat(currentLocale(), options).format(date)
}

export function formatDateTime(value: string | Date): string {
  return formatDate(value, { dateStyle: 'medium', timeStyle: 'short' })
}

export function formatNumber(
  value: number,
  options: Intl.NumberFormatOptions = {},
): string {
  return new Intl.NumberFormat(currentLocale(), options).format(value)
}

export function formatPercent(value: number, maximumFractionDigits = 1): string {
  return new Intl.NumberFormat(currentLocale(), {
    style: 'percent',
    maximumFractionDigits,
  }).format(value / 100)
}

export function formatDuration(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds <= 0) {
    return formatNumber(0, { style: 'unit', unit: 'second', unitDisplay: 'short' })
  }

  if (seconds < 60) {
    return new Intl.NumberFormat(currentLocale(), {
      style: 'unit',
      unit: 'second',
      unitDisplay: 'short',
      maximumFractionDigits: 0,
    }).format(seconds)
  }

  if (seconds < 3600) {
    return new Intl.NumberFormat(currentLocale(), {
      style: 'unit',
      unit: 'minute',
      unitDisplay: 'short',
      maximumFractionDigits: 1,
    }).format(seconds / 60)
  }

  return new Intl.NumberFormat(currentLocale(), {
    style: 'unit',
    unit: 'hour',
    unitDisplay: 'short',
    maximumFractionDigits: 1,
  }).format(seconds / 3600)
}
