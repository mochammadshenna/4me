<template>
  <v-app>
    <!-- Mobile App Bar -->
    <v-app-bar
      v-if="isMobile"
      color="#1a1d29"
      density="compact"
      elevation="0"
      class="mobile-app-bar"
    >
      <v-app-bar-nav-icon @click="drawer = !drawer" color="#e4e7eb"></v-app-bar-nav-icon>
      <v-toolbar-title class="app-title">4me Todos</v-toolbar-title>
    </v-app-bar>

    <!-- Navigation Drawer / Sidebar -->
    <v-navigation-drawer
      v-model="drawer"
      :permanent="!isMobile"
      :temporary="isMobile"
      :rail="false"
      width="280"
      class="modern-sidebar"
    >
      <!-- Logo Header -->
      <div class="sidebar-logo-header">
        <div class="logo-container">
          <v-icon size="32" color="#6C5CE7">mdi-checkbox-multiple-marked</v-icon>
          <span class="logo-text">4me Todos</span>
        </div>
      </div>

      <!-- Navigation Menu -->
      <v-list density="compact" class="sidebar-menu">
        <v-list-item
          :class="{ 'menu-item': true, 'active': $route.name === 'Dashboard' }"
          @click="navigateTo('/')"
          ripple
        >
          <template v-slot:prepend>
            <v-icon>mdi-view-dashboard</v-icon>
          </template>
          <v-list-item-title>Dashboard</v-list-item-title>
        </v-list-item>

        <v-list-item
          :class="{ 'menu-item': true, 'active': $route.name === 'ProjectView' }"
          @click="navigateTo('/projects')"
          ripple
        >
          <template v-slot:prepend>
            <v-icon>mdi-folder</v-icon>
          </template>
          <v-list-item-title>Projects</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-calendar</v-icon>
          </template>
          <v-list-item-title>Calendar</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-message</v-icon>
          </template>
          <v-list-item-title>Comments</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-magnify</v-icon>
          </template>
          <v-list-item-title>Search</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-chart-line</v-icon>
          </template>
          <v-list-item-title>Analytics</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-cog</v-icon>
          </template>
          <v-list-item-title>Settings</v-list-item-title>
        </v-list-item>

        <!-- Divider -->
        <v-divider class="sidebar-divider"></v-divider>

        <!-- Toolbar Items -->
        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-share-variant</v-icon>
          </template>
          <v-list-item-title>Share</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-download</v-icon>
          </template>
          <v-list-item-title>Export</v-list-item-title>
        </v-list-item>

        <v-list-item class="menu-item logout-item" @click="handleLogout" ripple>
          <template v-slot:prepend>
            <v-icon>mdi-logout</v-icon>
          </template>
          <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>

      <!-- Footer -->
      <template v-slot:append>
        <div class="sidebar-footer">
          <div class="dark-mode-toggle" @click="toggleDarkMode">
            <v-icon size="20">mdi-weather-night</v-icon>
            <span>Dark Mode</span>
          </div>
        </div>
      </template>
    </v-navigation-drawer>

    <!-- Main Content -->
    <v-main class="main-container">
      <div class="content-wrapper">
        <slot />
      </div>
    </v-main>
  </v-app>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const drawer = ref(true)
const isMobile = ref(false)

function navigateTo(route) {
  router.push(route)
  if (isMobile.value) {
    drawer.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

function toggleDarkMode() {
  console.log('Toggle dark mode')
}

function checkScreenSize() {
  isMobile.value = window.innerWidth < 768
  drawer.value = !isMobile.value
}

onMounted(() => {
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize)
})
</script>

<style scoped>
/* Mobile App Bar */
.mobile-app-bar {
  border-bottom: 1px solid #2d3142 !important;
}

.app-title {
  font-size: 18px;
  font-weight: 600;
  color: #e4e7eb;
  letter-spacing: -0.02em;
}

/* Sidebar Styles */
.modern-sidebar {
  background: #1a1d29 !important;
  border-right: 1px solid #2d3142 !important;
}

/* Logo Header */
.sidebar-logo-header {
  padding: 24px 20px;
  border-bottom: 1px solid #2d3142;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #e4e7eb;
  letter-spacing: -0.02em;
}

/* Menu Styles */
.sidebar-menu {
  padding: 16px 12px;
  background: transparent !important;
}

.menu-item {
  margin-bottom: 4px;
  border-radius: 8px !important;
  min-height: 44px !important;
  transition: all 0.2s ease;
}

.menu-item:hover {
  background: #252836 !important;
}

.menu-item.active {
  background: #2d3142 !important;
}

.menu-item :deep(.v-list-item__prepend) {
  margin-right: 12px;
}

.menu-item :deep(.v-icon) {
  color: #9ca3af;
  font-size: 20px;
}

.menu-item.active :deep(.v-icon) {
  color: #6C5CE7;
}

.menu-item :deep(.v-list-item-title) {
  color: #9ca3af;
  font-weight: 500;
  font-size: 14px;
}

.menu-item.active :deep(.v-list-item-title) {
  color: #e4e7eb;
  font-weight: 600;
}

.sidebar-divider {
  margin: 16px 16px;
  border-color: #2d3142 !important;
  opacity: 1 !important;
}

.logout-item :deep(.v-icon) {
  color: #ef4444 !important;
}

.logout-item :deep(.v-list-item-title) {
  color: #ef4444 !important;
}

.logout-item:hover {
  background: rgba(239, 68, 68, 0.1) !important;
}

/* Sidebar Footer */
.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid #2d3142;
}

.dark-mode-toggle {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #252836;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #9ca3af;
  font-size: 14px;
  font-weight: 500;
}

.dark-mode-toggle:hover {
  background: #2d3142;
  color: #e4e7eb;
}

.dark-mode-toggle :deep(.v-icon) {
  color: #9ca3af;
}

.dark-mode-toggle:hover :deep(.v-icon) {
  color: #6C5CE7;
}

/* Main Content */
.main-container {
  background: #f5f7fa;
}

.content-wrapper {
  padding: 32px;
  min-height: 100vh;
  max-width: 100%;
  margin: 0 auto;
}

/* Responsive Design */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 16px;
  }
}

@media (min-width: 1200px) {
  .content-wrapper {
    padding: 40px;
  }
}

@media (min-width: 1600px) {
  .content-wrapper {
    padding: 48px;
  }
}

@media (min-width: 2000px) {
  .content-wrapper {
    padding: 56px;
    max-width: 1800px;
  }
}
</style>
