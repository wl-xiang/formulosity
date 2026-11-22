import { i18n } from '@/locales'

function t(key: string, named?: Record<string, string | number>): string {
  return i18n.global.t(key, named ?? {})
}

export function translateApiError(message?: string): string {
  if (!message) {
    return t('common.error')
  }

  const normalized = message.trim()

  const exactMessages: Record<string, string> = {
    'invalid option selected': 'errors.invalidOption',
    'duplicate option selected': 'errors.duplicateOption',
    'invalid date format, expected YYYY-MM-DD': 'errors.invalidDate',
    'session is not in progress': 'errors.sessionNotInProgress',
    'duplicate session for ip address': 'errors.duplicateSessionIp',
    'question is hidden': 'errors.questionHidden',
    'invalid request format': 'errors.invalidRequest',
    'file is required for this question type': 'errors.fileRequired',
    'unable to save file': 'errors.unableToSaveFile',
    'invalid answer': 'errors.invalidAnswer',
  }

  if (exactMessages[normalized]) {
    return t(exactMessages[normalized])
  }

  const patterns: Array<{
    pattern: RegExp
    key: string
    params: (matches: RegExpMatchArray) => Record<string, string | number>
  }> = [
    {
      pattern: /^select at least (\d+) options$/,
      key: 'errors.selectAtLeast',
      params: (matches) => ({ n: Number(matches[1]) }),
    },
    {
      pattern: /^select at most (\d+) options$/,
      key: 'errors.selectAtMost',
      params: (matches) => ({ n: Number(matches[1]) }),
    },
    {
      pattern: /^please write at least (\d+) characters$/,
      key: 'errors.writeAtLeast',
      params: (matches) => ({ n: Number(matches[1]) }),
    },
    {
      pattern: /^please write at most (\d+) characters$/,
      key: 'errors.writeAtMost',
      params: (matches) => ({ n: Number(matches[1]) }),
    },
    {
      pattern: /^minimum: (\d+)$/,
      key: 'errors.minimum',
      params: (matches) => ({ n: Number(matches[1]) }),
    },
    {
      pattern: /^maximum: (\d+)$/,
      key: 'errors.maximum',
      params: (matches) => ({ n: Number(matches[1]) }),
    },
    {
      pattern: /^file size exceeds the maximum size of (.+)$/,
      key: 'errors.fileSize',
      params: (matches) => ({ size: matches[1] }),
    },
    {
      pattern: /^file format is invalid: (.+)\. Allowed formats: (.+)$/,
      key: 'errors.fileFormat',
      params: (matches) => ({ format: matches[1], allowed: matches[2] }),
    },
  ]

  for (const { pattern, key, params } of patterns) {
    const matches = normalized.match(pattern)
    if (matches) {
      return t(key, params(matches))
    }
  }

  return normalized
}
