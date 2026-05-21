<template>
  <el-tooltip :content="isDark ? $t('navbar.lightMode') : $t('navbar.darkMode')" effect="dark" placement="bottom">
    <div class="dark-mode-toggle right-menu-item hover-effect" @click="toggleDark">
      <svg-icon :icon-class="isDark ? 'sun' : 'moon'" class-name="dark-icon" />
    </div>
  </el-tooltip>
</template>

<script>
export default {
  name: 'DarkModeToggle',
  computed: {
    isDark() {
      return this.$store.getters.sideTheme === 'dark'
    }
  },
  methods: {
    toggleDark() {
      const theme = this.isDark ? 'light' : 'dark'
      this.$store.dispatch('settings/changeSideTheme', theme)
      this.$nextTick(() => {
        if (theme === 'dark') {
          document.body.classList.add('dark')
        } else {
          document.body.classList.remove('dark')
        }
      })
    }
  }
}
</script>

<style scoped>
.dark-mode-toggle {
  display: inline-block;
  padding: 0 8px;
  height: 100%;
  font-size: 18px;
  color: #5a5e66;
  vertical-align: text-bottom;
  cursor: pointer;
}

.dark-mode-toggle.hover-effect:hover {
  background: rgba(0, 0, 0, 0.025);
}

.dark-icon {
  font-size: 18px;
}

body.dark .dark-mode-toggle {
  color: #e6e6e6;
}

body.dark .dark-mode-toggle.hover-effect:hover {
  background: rgba(255, 255, 255, 0.05);
}
</style>
