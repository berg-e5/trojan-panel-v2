<template>
  <div class="stats-container">
    <el-row :gutter="16">
      <!-- 用户统计卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card stat-users">
          <div class="stat-icon">
            <i class="el-icon-user" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ userStats.totalUsers || 0 }}</div>
            <div class="stat-label">{{ $t('stats.totalUsers') }}</div>
          </div>
          <div class="stat-detail">
            <span class="text-success">{{ userStats.activeUsers || 0 }} {{ $t('stats.active') }}</span>
            <span class="text-danger">{{ userStats.disabledUsers || 0 }} {{ $t('stats.disabled') }}</span>
          </div>
        </div>
      </el-col>

      <!-- 节点统计卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card stat-nodes">
          <div class="stat-icon">
            <i class="el-icon-s-grid" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ nodeStats.totalNodes || 0 }}</div>
            <div class="stat-label">{{ $t('stats.totalNodes') }}</div>
          </div>
          <div class="stat-detail">
            <span class="text-success">{{ nodeStats.onlineNodes || 0 }} {{ $t('stats.online') }}</span>
            <span class="text-danger">{{ nodeStats.offlineNodes || 0 }} {{ $t('stats.offline') }}</span>
          </div>
        </div>
      </el-col>

      <!-- 流量统计卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card stat-traffic">
          <div class="stat-icon">
            <i class="el-icon-arrow-up" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatBytes(trafficStats.usedBandwidth || 0) }}</div>
            <div class="stat-label">{{ $t('stats.usedTraffic') }}</div>
          </div>
          <div class="stat-detail">
            <span>{{ $t('stats.usage') }}: {{ (trafficStats.usagePercent || 0).toFixed(1) }}%</span>
          </div>
        </div>
      </el-col>

      <!-- 剩余流量卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card stat-remaining">
          <div class="stat-icon">
            <i class="el-icon-arrow-down" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatBytes(trafficStats.unusedBandwidth || 0) }}</div>
            <div class="stat-label">{{ $t('stats.remainingTraffic') }}</div>
          </div>
          <div class="stat-detail">
            <span>{{ $t('stats.total') }}: {{ formatBytes(trafficStats.totalBandwidth || 0) }}</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <!-- 用户活跃饼图 -->
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div slot="header">{{ $t('stats.userDistribution') }}</div>
          <e-charts
            v-if="userPieOption"
            :option="userPieOption"
            height="250px"
          />
        </el-card>
      </el-col>

      <!-- 协议分布饼图 -->
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div slot="header">{{ $t('stats.protocolDistribution') }}</div>
          <e-charts
            v-if="protocolPieOption"
            :option="protocolPieOption"
            height="250px"
          />
        </el-card>
      </el-col>

      <!-- 流量使用进度 -->
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div slot="header">{{ $t('stats.trafficUsageRate') }}</div>
          <e-charts
            v-if="trafficGaugeOption"
            :option="trafficGaugeOption"
            height="250px"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import ECharts from '@/components/ECharts/index.vue'
import { userStats, trafficStats, nodeStats, protocolStats } from '@/api/dashboard'

export default {
  name: 'StatsCharts',
  components: { ECharts },
  data() {
    return {
      userStats: {},
      nodeStats: {},
      trafficStats: {},
      protocolStats: [],
      userPieOption: null,
      protocolPieOption: null,
      trafficGaugeOption: null
    }
  },
  created() {
    this.fetchAll()
  },
  methods: {
    fetchAll() {
      this.fetchUserStats()
      this.fetchTrafficStats()
      this.fetchNodeStats()
      this.fetchProtocolStats()
    },
    async fetchUserStats() {
      try {
        const res = await userStats()
        if (res.data.code === 0) {
          this.userStats = res.data.data
          this.buildUserPieChart()
        }
      } catch (e) {
        console.error('fetchUserStats error', e)
      }
    },
    async fetchTrafficStats() {
      try {
        const res = await trafficStats()
        if (res.data.code === 0) {
          this.trafficStats = res.data.data
          this.buildTrafficGauge()
        }
      } catch (e) {
        console.error('fetchTrafficStats error', e)
      }
    },
    async fetchNodeStats() {
      try {
        const res = await nodeStats()
        if (res.data.code === 0) {
          this.nodeStats = res.data.data
        }
      } catch (e) {
        console.error('fetchNodeStats error', e)
      }
    },
    async fetchProtocolStats() {
      try {
        const res = await protocolStats()
        if (res.data.code === 0) {
          this.protocolStats = res.data.data || []
          this.buildProtocolPieChart()
        }
      } catch (e) {
        console.error('fetchProtocolStats error', e)
      }
    },
    buildUserPieChart() {
      const u = this.userStats
      this.userPieOption = {
        tooltip: {
          trigger: 'item',
          formatter: '{b}: {c} ({d}%)'
        },
        legend: {
          bottom: 0,
          left: 'center'
        },
        series: [
          {
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 6,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: true,
              formatter: '{b}: {c}'
            },
            data: [
              { value: u.activeUsers || 0, name: this.$t('stats.active'), itemStyle: { color: '#13ce66' } },
              { value: u.inactiveUsers || 0, name: this.$t('stats.inactive'), itemStyle: { color: '#ffba00' } },
              { value: u.disabledUsers || 0, name: this.$t('stats.disabled'), itemStyle: { color: '#ff4949' } }
            ]
          }
        ]
      }
    },
    buildProtocolPieChart() {
      const colors = ['#409eff', '#36a3f7', '#f4516c', '#ffba00', '#13ce66', '#4ab7bd']
      const data = (this.protocolStats || []).map((p, i) => ({
        name: p.name,
        value: p.count,
        itemStyle: { color: colors[i % colors.length] }
      }))
      this.protocolPieOption = {
        tooltip: {
          trigger: 'item',
          formatter: '{b}: {c} ({d}%)'
        },
        legend: {
          bottom: 0,
          left: 'center'
        },
        series: [
          {
            type: 'pie',
            radius: ['40%', '70%'],
            itemStyle: {
              borderRadius: 6,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: true,
              formatter: '{b}: {c}'
            },
            data
          }
        ]
      }
    },
    buildTrafficGauge() {
      const percent = Math.min((this.trafficStats.usagePercent || 0), 100)
      this.trafficGaugeOption = {
        series: [
          {
            type: 'gauge',
            startAngle: 200,
            endAngle: -20,
            min: 0,
            max: 100,
            splitNumber: 5,
            radius: '90%',
            axisLine: {
              lineStyle: {
                width: 15,
                color: [
                  [0.6, '#13ce66'],
                  [0.8, '#ffba00'],
                  [1, '#ff4949']
                ]
              }
            },
            pointer: {
              itemStyle: {
                color: '#666'
              }
            },
            axisTick: {
              distance: -20,
              length: 5
            },
            splitLine: {
              distance: -20,
              length: 10
            },
            axisLabel: {
              distance: 20,
              color: '#999',
              fontSize: 10
            },
            detail: {
              valueAnimation: true,
              formatter: '{value}%',
              fontSize: 24,
              offsetCenter: [0, '70%']
            },
            data: [{ value: parseFloat(percent.toFixed(1)), name: this.$t('stats.usageRate') }]
          }
        ]
      }
    },
    formatBytes(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }
  }
}
</script>

<style scoped>
.stats-container {
  margin-top: 16px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-icon {
  font-size: 36px;
  opacity: 0.8;
}

.stat-users .stat-icon { color: #409eff; }
.stat-nodes .stat-icon { color: #f4516c; }
.stat-traffic .stat-icon { color: #36a3f7; }
.stat-remaining .stat-icon { color: #13ce66; }

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.stat-detail {
  font-size: 11px;
  color: #666;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.text-success { color: #13ce66; }
.text-danger { color: #ff4949; }

body.dark .stat-card {
  background: #1a1a2e;
  color: #e6e6e6;
}

body.dark .stat-value {
  color: #e6e6e6;
}

body.dark .stat-label,
body.dark .stat-detail {
  color: #a0a0a0;
}
</style>
