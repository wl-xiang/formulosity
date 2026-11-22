import { use } from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import {
  GridComponent,
  LegendComponent,
  TooltipComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([
  BarChart,
  LineChart,
  PieChart,
  GridComponent,
  LegendComponent,
  TooltipComponent,
  CanvasRenderer,
])

export const chartPalette = [
  '#3b82f6',
  '#10b981',
  '#f59e0b',
  '#ef4444',
  '#8b5cf6',
  '#06b6d4',
  '#ec4899',
  '#84cc16',
]

export type ChartTheme = 'light' | 'dark'

export function chartColors(theme: ChartTheme) {
  if (theme === 'dark') {
    return {
      text: '#d1d5db',
      muted: '#9ca3af',
      axis: '#4b5563',
      grid: '#1f2937',
      tooltipBackground: '#111827',
      tooltipBorder: '#374151',
      palette: chartPalette,
    }
  }

  return {
    text: '#374151',
    muted: '#6b7280',
    axis: '#9ca3af',
    grid: '#e5e7eb',
    tooltipBackground: '#ffffff',
    tooltipBorder: '#e5e7eb',
    palette: chartPalette,
  }
}
