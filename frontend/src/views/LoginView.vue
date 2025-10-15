<template>
  <v-layout class="login-layout">
    <!-- StudyBuddy Style Login Page -->
    <div class="login-container">
      <!-- Background Elements -->
      <div class="login-bg"></div>
      
      <!-- Login Card -->
      <v-card class="login-card" elevation="0">
        <!-- Header -->
        <div class="login-header">
          <div class="logo-section">
            <v-icon size="40" color="white" class="mb-2">mdi-checkbox-marked-circle</v-icon>
            <h1 class="login-title">4me Todos</h1>
            <p class="login-subtitle">Your Task Hub</p>
          </div>
        </div>
        
        <!-- Form Section -->
        <div class="login-form-section">
          <h2 class="form-title">Welcome Back</h2>
          
          <v-form @submit.prevent="handleLogin" class="login-form">
            <div class="input-group">
              <v-text-field
                v-model="username"
                label="Username"
                prepend-inner-icon="mdi-account"
                variant="outlined"
                :rules="[v => !!v || 'Username is required']"
                class="login-input"
                hide-details="auto"
              />
            </div>
            
            <div class="input-group">
              <v-text-field
                v-model="password"
                label="Password"
                type="password"
                prepend-inner-icon="mdi-lock"
                variant="outlined"
                :rules="[v => !!v || 'Password is required']"
                class="login-input"
                hide-details="auto"
              />
            </div>
            
            <v-alert v-if="error" type="error" class="error-alert" density="compact">
              {{ error }}
            </v-alert>
            
            <v-btn
              type="submit"
              size="large"
              block
              :loading="loading"
              class="login-btn"
            >
              Login
            </v-btn>
            
            <div class="divider-section">
              <v-divider />
              <span class="divider-text">or</span>
              <v-divider />
            </div>
            
            <v-btn
              variant="outlined"
              size="large"
              block
              @click="handleGoogleLogin"
              class="google-btn"
            >
              <v-icon start>mdi-google</v-icon>
              Continue with Google
            </v-btn>
            
            <div class="signup-section">
              <span class="signup-text">Don't have an account?</span>
              <router-link to="/register" class="signup-link">
                Sign up
              </router-link>
            </div>
          </v-form>
        </div>
      </v-card>
    </div>
  </v-layout>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  
  const result = await authStore.login({
    username: username.value,
    password: password.value,
  })
  
  loading.value = false
  
  if (result.success) {
    router.push('/')
  } else {
    error.value = result.error
  }
}

function handleGoogleLogin() {
  authStore.googleLogin()
}
</script>

<style scoped>
/* StudyBuddy Login Page Styling */
.login-layout {
  background: #FFF8E1;
  min-height: 100vh;
}

.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
  position: relative;
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    linear-gradient(135deg, #E3F2FD 0%, #BBDEFB 25%, #FFF8E1 50%, #FFECB3 75%, #E8F5E8 100%),
    radial-gradient(circle at 20% 80%, rgba(227, 242, 253, 0.4) 0%, transparent 60%),
    radial-gradient(circle at 80% 20%, rgba(255, 243, 224, 0.4) 0%, transparent 60%),
    radial-gradient(circle at 40% 40%, rgba(240, 248, 255, 0.3) 0%, transparent 60%);
  z-index: -1;
}

.login-card {
  background: white !important;
  border-radius: 24px !important;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2);
  width: 100%;
  max-width: 400px;
  overflow: hidden;
}

.login-header {
  background: linear-gradient(135deg, #1976D2 0%, #1565C0 50%, #0D47A1 100%);
  padding: 40px 32px 32px;
  text-align: center;
  position: relative;
  overflow: hidden;
}

.login-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 30% 20%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 70% 80%, rgba(255, 255, 255, 0.05) 0%, transparent 50%);
  pointer-events: none;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  z-index: 1;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.login-subtitle {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

.login-form-section {
  padding: 32px;
}

.form-title {
  font-size: 24px;
  font-weight: 600;
  color: #37474F;
  margin: 0 0 24px 0;
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  margin-bottom: 8px;
}

.login-input {
  background: #F8F9FA;
  border-radius: 12px !important;
}

.login-input :deep(.v-field) {
  border-radius: 12px !important;
  background: #F8F9FA !important;
}

.login-input :deep(.v-field--focused) {
  background: white !important;
  box-shadow: 0 0 0 2px #1976D2 !important;
}

.login-btn {
  background: linear-gradient(135deg, #1976D2, #1565C0) !important;
  color: white !important;
  border-radius: 12px !important;
  text-transform: none !important;
  font-weight: 600 !important;
  font-size: 16px !important;
  height: 48px !important;
  margin-top: 8px !important;
  box-shadow: 0 4px 16px rgba(25, 118, 210, 0.3) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  position: relative !important;
  overflow: hidden !important;
}

.login-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.login-btn:hover::before {
  left: 100%;
}

.login-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 24px rgba(25, 118, 210, 0.4) !important;
}

.divider-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 16px 0;
}

.divider-text {
  color: #78909C;
  font-size: 14px;
  font-weight: 500;
}

.google-btn {
  border: 2px solid #E0E0E0 !important;
  border-radius: 12px !important;
  text-transform: none !important;
  font-weight: 500 !important;
  font-size: 16px !important;
  height: 48px !important;
  color: #37474F !important;
  background: white !important;
  transition: all 0.3s ease !important;
}

.google-btn:hover {
  border-color: #1976D2 !important;
  background: #F8F9FA !important;
  transform: translateY(-1px);
}

.error-alert {
  border-radius: 8px !important;
  margin-bottom: 8px !important;
}

.signup-section {
  text-align: center;
  margin-top: 24px;
}

.signup-text {
  color: #78909C;
  font-size: 14px;
}

.signup-link {
  color: #1976D2 !important;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  margin-left: 4px;
  transition: color 0.3s ease;
}

.signup-link:hover {
  color: #1565C0 !important;
}

/* Custom Form Styling */
:deep(.v-text-field--outlined .v-field) {
  border-radius: 12px !important;
}

:deep(.v-field--focused .v-field__outline) {
  --v-field-border-opacity: 1;
  --v-field-border-color: #1976D2;
}

/* Responsive Design */
@media (max-width: 480px) {
  .login-container {
    padding: 16px;
  }
  
  .login-form-section {
    padding: 24px;
  }
  
  .login-header {
    padding: 32px 24px 24px;
  }
  
  .login-title {
    font-size: 24px;
  }
  
  .form-title {
    font-size: 20px;
  }
}
</style>

