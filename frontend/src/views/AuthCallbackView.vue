<template>
  <v-layout class="auth-callback-layout">
    <!-- StudyBuddy Style Auth Callback Page -->
    <div class="auth-callback-container">
      <!-- Background Elements -->
      <div class="auth-callback-bg"></div>
      
      <!-- Auth Callback Card -->
      <v-card class="auth-callback-card" elevation="0">
        <!-- Loading Section -->
        <div class="loading-section">
          <div class="loading-icon">
            <v-progress-circular
              indeterminate
              color="#1976D2"
              size="64"
              width="4"
            />
          </div>
          <h2 class="loading-title">Completing Authentication</h2>
          <p class="loading-subtitle">Please wait while we set up your account...</p>
        </div>
      </v-card>
    </div>
  </v-layout>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(() => {
  const token = route.query.token
  const refreshToken = route.query.refresh_token
  
  if (token && refreshToken) {
    authStore.handleGoogleCallback(token, refreshToken)
    router.push('/')
  } else {
    router.push('/login')
  }
})
</script>

<style scoped>
/* StudyBuddy Auth Callback Page Styling */
.auth-callback-layout {
  background: #FFF8E1;
  min-height: 100vh;
}

.auth-callback-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
  position: relative;
}

.auth-callback-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 80%, rgba(227, 242, 253, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 243, 224, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(240, 248, 255, 0.2) 0%, transparent 50%);
  z-index: -1;
}

.auth-callback-card {
  background: white !important;
  border-radius: 24px !important;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2);
  width: 100%;
  max-width: 400px;
  overflow: hidden;
}

.loading-section {
  padding: 60px 40px;
  text-align: center;
}

.loading-icon {
  margin-bottom: 24px;
}

.loading-title {
  font-size: 24px;
  font-weight: 600;
  color: #37474F;
  margin: 0 0 12px 0;
}

.loading-subtitle {
  font-size: 16px;
  color: #78909C;
  margin: 0;
  line-height: 1.5;
}

/* Responsive Design */
@media (max-width: 480px) {
  .auth-callback-container {
    padding: 16px;
  }
  
  .loading-section {
    padding: 40px 24px;
  }
  
  .loading-title {
    font-size: 20px;
  }
  
  .loading-subtitle {
    font-size: 14px;
  }
}
</style>

