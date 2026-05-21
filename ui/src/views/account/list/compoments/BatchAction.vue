<template>
  <el-dialog
    :title="$t('table.batchAction').toString()"
    :visible="dialogVisible"
    width="500px"
    @close="handleClose"
  >
    <el-form ref="dataForm" :model="temp" label-position="left">
      <!-- 操作类型 -->
      <el-form-item :label="$t('table.batchActionType').toString()" prop="action">
        <el-select v-model="temp.action" :placeholder="$t('table.selectAction').toString()">
          <el-option :label="$t('table.enable').toString()" value="enable" />
          <el-option :label="$t('table.disable').toString()" value="disable" />
          <el-option :label="$t('table.extend').toString()" value="extend" />
          <el-option :label="$t('table.resetTraffic').toString()" value="reset" />
          <el-option :label="$t('table.delete').toString()" value="delete" style="color: #f56c6c" />
        </el-select>
      </el-form-item>

      <!-- 续期天数（仅续期操作显示） -->
      <el-form-item
        v-if="temp.action === 'extend'"
        :label="$t('table.extendDays').toString()"
        prop="days"
      >
        <el-input-number
          v-model="temp.days"
          :min="1"
          :max="365"
          controls-position="right"
        />
        <span style="margin-left: 8px; color: #909399">{{ $t('table.daysUnit').toString() }}</span>
      </el-form-item>

      <!-- 确认提示 -->
      <el-alert
        v-if="temp.action"
        :title="confirmText"
        :type="alertType"
        :closable="false"
        show-icon
        style="margin-top: 16px"
      />
    </el-form>

    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose">{{ $t('table.cancel') }}</el-button>
      <el-button type="primary" @click="handleConfirm" :loading="confirmLoading">
        {{ $t('table.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
import { batchAccount } from '@/api/account'

export default {
  name: 'BatchAction',
  props: {
    dialogVisible: {
      type: Boolean,
      default: false
    },
    selectedIds: {
      type: Array,
      default: () => []
    },
    getList: {
      type: Function,
      required: true
    }
  },
  data() {
    return {
      temp: {
        action: '',
        days: 30
      },
      confirmLoading: false
    }
  },
  computed: {
    confirmText() {
      const count = this.selectedIds.length
      const actionMap = {
        enable: this.$t('table.enableConfirm'),
        disable: this.$t('table.disableConfirm'),
        extend: this.$t('table.extendConfirm'),
        reset: this.$t('table.resetConfirm'),
        delete: this.$t('table.deleteConfirm')
      }
      return `${actionMap[this.temp.action] || ''} (${count} ${this.$t('table.accounts').toString()})`
    },
    alertType() {
      const map = {
        enable: 'success',
        disable: 'warning',
        extend: 'success',
        reset: 'info',
        delete: 'error'
      }
      return map[this.temp.action] || 'info'
    }
  },
  methods: {
    handleClose() {
      this.temp.action = ''
      this.temp.days = 30
      this.$emit('update:dialogVisible', false)
    },
    handleConfirm() {
      if (!this.temp.action) {
        this.$message.warning(this.$t('table.selectActionFirst'))
        return
      }

      this.$confirm(this.confirmText, this.$t('table.confirmTitle'), {
        confirmButtonText: this.$t('table.confirm'),
        cancelButtonText: this.$t('table.cancel'),
        type: this.alertType
      }).then(() => {
        this.doBatchAction()
      }).catch(() => {})
    },
    doBatchAction() {
      this.confirmLoading = true
      const data = {
        ids: this.selectedIds,
        action: this.temp.action
      }
      if (this.temp.action === 'extend' && this.temp.days) {
        data.days = this.temp.days
      }

      batchAccount(data).then((res) => {
        this.confirmLoading = false
        if (res.data.code === 0) {
          this.$notify({
            title: this.$t('confirm.success'),
            message: res.data.data.message || this.$t('confirm.success'),
            type: 'success',
            duration: 3000
          })
          this.getList()
          this.handleClose()
        } else {
          this.$notify({
            title: this.$t('confirm.fail'),
            message: res.data.msg,
            type: 'error',
            duration: 3000
          })
        }
      }).catch((err) => {
        this.confirmLoading = false
        this.$notify({
          title: this.$t('confirm.fail'),
          message: err.message || this.$t('confirm.operationFailed'),
          type: 'error',
          duration: 3000
        })
      })
    }
  }
}
</script>

<style scoped>
.el-select {
  width: 100%;
}
</style>
