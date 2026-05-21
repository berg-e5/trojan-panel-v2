<template>
  <div ref="chartRef" class="echarts" :style="{ width: width, height: height }" />
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: 'ECharts',
  props: {
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '300px'
    },
    option: {
      type: Object,
      required: true
    },
    theme: {
      type: String,
      default: 'default'
    }
  },
  data() {
    return {
      chart: null
    }
  },
  watch: {
    option: {
      deep: true,
      handler(val) {
        if (this.chart) {
          this.chart.setOption(val, true)
        }
      }
    }
  },
  mounted() {
    this.initChart()
    window.addEventListener('resize', this.handleResize)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.handleResize)
    if (this.chart) {
      this.chart.dispose()
      this.chart = null
    }
  },
  methods: {
    initChart() {
      this.chart = echarts.init(this.$refs.chartRef, this.theme)
      this.chart.setOption(this.option, true)
    },
    handleResize() {
      if (this.chart) {
        this.chart.resize()
      }
    }
  }
}
</script>

<style scoped>
.echarts {
  width: 100%;
}
</style>
