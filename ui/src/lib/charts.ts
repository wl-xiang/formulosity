import type { EChartsOption } from 'echarts'
import { formatDate } from '@/lib/format'
import { chartColors } from '@/lib/echarts'
import type { ChartTheme } from '@/lib/echarts'
import type {
  AnalyticsBucket,
  AnalyticsTimelinePoint,
  QuestionAnalytics,
} from '@/lib/types'

export type ChartLabels = {
  started: string
  completed: string
  yes: string
  no: string
  rank: (n: number) => string
}

function tooltipStyle(theme: ChartTheme) {
  const colors = chartColors(theme)

  return {
    backgroundColor: colors.tooltipBackground,
    borderColor: colors.tooltipBorder,
    textStyle: {
      color: colors.text,
    },
  }
}

function truncate(value: string, maxLength = 20): string {
  return value.length > maxLength ? `${value.slice(0, maxLength)}…` : value
}

function dateBucketLabel(key: string, granularity?: string): string {
  if (granularity === 'year') {
    return formatDate(key, { year: 'numeric' })
  }

  if (granularity === 'month') {
    return formatDate(key, { year: 'numeric', month: 'short' })
  }

  return formatDate(key, { month: 'short', day: 'numeric' })
}

export function bucketDisplayLabel(
  question: QuestionAnalytics,
  bucket: AnalyticsBucket,
  labels: ChartLabels,
): string {
  if (question.chart_kind === 'boolean') {
    return bucket.key === 'true' ? labels.yes : labels.no
  }

  if (question.chart_kind === 'date') {
    return dateBucketLabel(bucket.key, question.granularity)
  }

  return bucket.label
}

export function buildTimelineOption(
  points: AnalyticsTimelinePoint[],
  theme: ChartTheme,
  labels: Pick<ChartLabels, 'started' | 'completed'>,
): EChartsOption {
  const colors = chartColors(theme)

  return {
    color: [colors.palette[0], colors.palette[1]],
    tooltip: {
      ...tooltipStyle(theme),
      trigger: 'axis',
    },
    legend: {
      top: 0,
      textStyle: { color: colors.muted },
    },
    grid: {
      left: 8,
      right: 16,
      top: 32,
      bottom: 8,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: points.map((point) => formatDate(point.date, { month: 'short', day: 'numeric' })),
      axisLabel: { color: colors.muted, hideOverlap: true },
      axisLine: { lineStyle: { color: colors.axis } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: colors.muted },
      splitLine: { lineStyle: { color: colors.grid } },
    },
    series: [
      {
        name: labels.started,
        type: 'line',
        symbolSize: 4,
        data: points.map((point) => point.started),
        lineStyle: { width: 2 },
        areaStyle: { opacity: 0.08 },
      },
      {
        name: labels.completed,
        type: 'line',
        symbolSize: 4,
        data: points.map((point) => point.completed),
        lineStyle: { width: 2 },
        areaStyle: { opacity: 0.08 },
      },
    ],
  }
}

export function buildQuestionChartOption(
  question: QuestionAnalytics,
  theme: ChartTheme,
  labels: ChartLabels,
): EChartsOption | null {
  if (question.chart_kind === 'ranking') {
    return buildRankingOption(question, theme, labels)
  }

  const buckets = question.buckets ?? []
  if (buckets.length === 0) {
    return null
  }

  switch (question.chart_kind) {
    case 'single_choice':
      return buckets.length <= 6
        ? buildDonutOption(question, buckets, theme, labels)
        : buildHorizontalBarOption(question, buckets, theme, labels)
    case 'boolean':
      return buildDonutOption(question, buckets, theme, labels)
    case 'multi_choice':
      return buildHorizontalBarOption(question, buckets, theme, labels)
    case 'rating':
      return buildVerticalBarOption(
        buckets.map((bucket) => bucket.label),
        buckets.map((bucket) => bucket.count),
        theme,
      )
    case 'date':
      return buildVerticalBarOption(
        buckets.map((bucket) => dateBucketLabel(bucket.key, question.granularity)),
        buckets.map((bucket) => bucket.count),
        theme,
      )
    default:
      return null
  }
}

function buildDonutOption(
  question: QuestionAnalytics,
  buckets: AnalyticsBucket[],
  theme: ChartTheme,
  labels: ChartLabels,
): EChartsOption {
  const colors = chartColors(theme)

  return {
    color: colors.palette,
    tooltip: {
      ...tooltipStyle(theme),
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)',
    },
    legend: {
      bottom: 0,
      type: 'scroll',
      textStyle: { color: colors.muted },
    },
    series: [
      {
        type: 'pie',
        radius: ['45%', '68%'],
        center: ['50%', '44%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderColor: colors.tooltipBackground,
          borderWidth: 2,
        },
        label: {
          color: colors.text,
          formatter: '{b}\n{d}%',
        },
        labelLine: {
          lineStyle: { color: colors.axis },
        },
        data: buckets.map((bucket) => ({
          name: bucketDisplayLabel(question, bucket, labels),
          value: bucket.count,
        })),
      },
    ],
  }
}

function buildHorizontalBarOption(
  question: QuestionAnalytics,
  buckets: AnalyticsBucket[],
  theme: ChartTheme,
  labels: ChartLabels,
): EChartsOption {
  const colors = chartColors(theme)
  const rows = [...buckets].reverse()

  return {
    color: [colors.palette[0]],
    tooltip: {
      ...tooltipStyle(theme),
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
    },
    grid: {
      left: 8,
      right: 40,
      top: 8,
      bottom: 8,
      containLabel: true,
    },
    xAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: colors.muted },
      splitLine: { lineStyle: { color: colors.grid } },
    },
    yAxis: {
      type: 'category',
      data: rows.map((bucket) => bucketDisplayLabel(question, bucket, labels)),
      axisLabel: {
        color: colors.muted,
        formatter: (value: string) => truncate(value),
      },
      axisLine: { lineStyle: { color: colors.axis } },
      axisTick: { show: false },
    },
    series: [
      {
        type: 'bar',
        data: rows.map((bucket) => bucket.count),
        barMaxWidth: 24,
        itemStyle: { borderRadius: [0, 4, 4, 0] },
        label: {
          show: true,
          position: 'right',
          color: colors.text,
        },
      },
    ],
  }
}

function buildVerticalBarOption(
  categories: string[],
  values: number[],
  theme: ChartTheme,
): EChartsOption {
  const colors = chartColors(theme)

  return {
    color: [colors.palette[0]],
    tooltip: {
      ...tooltipStyle(theme),
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
    },
    grid: {
      left: 8,
      right: 16,
      top: 24,
      bottom: 8,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: categories,
      axisLabel: {
        color: colors.muted,
        rotate: categories.length > 10 ? 45 : 0,
        hideOverlap: true,
      },
      axisLine: { lineStyle: { color: colors.axis } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: colors.muted },
      splitLine: { lineStyle: { color: colors.grid } },
    },
    series: [
      {
        type: 'bar',
        data: values,
        barMaxWidth: 32,
        itemStyle: { borderRadius: [4, 4, 0, 0] },
        label: {
          show: true,
          position: 'top',
          color: colors.text,
        },
      },
    ],
  }
}

function buildRankingOption(
  question: QuestionAnalytics,
  theme: ChartTheme,
  labels: ChartLabels,
): EChartsOption | null {
  const ranking = question.ranking ?? []
  if (ranking.length === 0) {
    return null
  }

  const colors = chartColors(theme)
  const rows = [...ranking].reverse()
  const rankCount = Math.max(...ranking.map((option) => option.rank_counts.length), 0)
  if (rankCount === 0) {
    return null
  }

  return {
    color: colors.palette,
    tooltip: {
      ...tooltipStyle(theme),
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
    },
    legend: {
      bottom: 0,
      type: 'scroll',
      textStyle: { color: colors.muted },
    },
    grid: {
      left: 8,
      right: 24,
      top: 8,
      bottom: 36,
      containLabel: true,
    },
    xAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: colors.muted },
      splitLine: { lineStyle: { color: colors.grid } },
    },
    yAxis: {
      type: 'category',
      data: rows.map((option) => truncate(option.label)),
      axisLabel: { color: colors.muted },
      axisLine: { lineStyle: { color: colors.axis } },
      axisTick: { show: false },
    },
    series: Array.from({ length: rankCount }, (_, index) => ({
      name: labels.rank(index + 1),
      type: 'bar' as const,
      stack: 'total',
      barMaxWidth: 24,
      itemStyle: { color: colors.palette[index % colors.palette.length] },
      data: rows.map((option) => option.rank_counts[index] ?? 0),
    })),
  }
}
