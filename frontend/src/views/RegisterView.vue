<template>
  <v-layout class="register-layout">
    <!-- StudyBuddy Style Register Page -->
    <div class="register-container">
      <!-- Background Elements -->
      <div class="register-bg"></div>
      
      <!-- Register Card -->
      <v-card class="register-card" elevation="0">
        <!-- Header -->
        <div class="register-header">
          <div class="logo-section">
            <v-icon size="40" color="#1976D2" class="mb-2">mdi-checkbox-marked-circle</v-icon>
            <h1 class="register-title">4me Todos</h1>
            <p class="register-subtitle">Your Task Hub</p>
          </div>
        </div>
        
        <!-- Form Section -->
        <div class="register-form-section">
          <h2 class="form-title">Create Account</h2>
          
          <v-form @submit.prevent="handleRegister" class="register-form">
            <div class="input-group">
              <v-text-field
                v-model="username"
                label="Username"
                prepend-inner-icon="mdi-account"
                variant="outlined"
                :rules="[v => !!v || 'Username is required', v => v.length >= 3 || 'Username must be at least 3 characters']"
                class="register-input"
                hide-details="auto"
              />
            </div>
            
            <div class="input-group">
              <v-text-field
                v-model="email"
                label="Email"
                type="email"
                prepend-inner-icon="mdi-email"
                variant="outlined"
                :rules="[v => !!v || 'Email is required', v => /.+@.+\..+/.test(v) || 'Email must be valid']"
                class="register-input"
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
                :rules="[v => !!v || 'Password is required', v => v.length >= 6 || 'Password must be at least 6 characters']"
                class="register-input"
                hide-details="auto"
              />
            </div>
            
            <div class="input-group">
              <v-text-field
                v-model="confirmPassword"
                label="Confirm Password"
                type="password"
                prepend-inner-icon="mdi-lock-check"
                variant="outlined"
                :rules="[v => !!v || 'Please confirm password', v => v === password || 'Passwords must match']"
                class="register-input"
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
              class="register-btn"
            >
              Sign Up
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
              Sign up with Google
            </v-btn>
            
            <div class="signin-section">
              <span class="signin-text">Already have an account?</span>
              <router-link to="/login" class="signin-link">
                Login
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
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const loading = ref(false)

async function handleRegister() {
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  
  error.value = ''
  loading.value = true
  
  const result = await authStore.register({
    username: username.value,
    email: email.value,
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
/* StudyBuddy Register Page Styling */
.register-layout {
  background: #FFF8E1;
  min-height: 100vh;
}

.register-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
  position: relative;
}

.register-bg {
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

.register-card {
  background: white !important;
  border-radius: 24px !important;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2);
  width: 100%;
  max-width: 450px;
  overflow: hidden;
}

.register-header {
  background: linear-gradient(135deg, #E3F2FD 0%, #BBDEFB 100%);
  padding: 40px 32px 32px;
  text-align: center;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.register-title {
  font-size: 28px;
  font-weight: 700;
  color: #1976D2;
  margin: 0 0 8px 0;
}

.register-subtitle {
  font-size: 14px;
  color: #546E7A;
  margin: 0;
}

.register-form-section {
  padding: 32px;
}

.form-title {
  font-size: 24px;
  font-weight: 600;
  color: #37474F;
  margin: 0 0 24px 0;
  text-align: center;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  margin-bottom: 8px;
}

.register-input {
  background: #F8F9FA;
  border-radius: 12px !important;
}

.register-input :deep(.v-field) {
  border-radius: 12px !important;
  background: #F8F9FA !important;
}

.register-input :deep(.v-field--focused) {
  background: white !important;
  box-shadow: 0 0 0 2px #1976D2 !important;
}

.register-btn {
  background: linear-gradient(135deg, #1976D2, #1565C0) !important;
  color: white !important;
  border-radius: 12px !important;
  text-transform: none !important;
  font-weight: 600 !important;
  font-size: 16px !important;
  height: 48px !important;
  margin-top: 8px !important;
  box-shadow: 0 4px 16px rgba(25, 118, 210, 0.3) !important;
  transition: all 0.3s ease !important;
}

.register-btn:hover {
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

.signin-section {
  text-align: center;
  margin-top: 24px;
}

.signin-text {
  color: #78909C;
  font-size: 14px;
}

.signin-link {
  color: #1976D2 !important;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  margin-left: 4px;
  transition: color 0.3s ease;
}

.signin-link:hover {
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
  .register-container {
    padding: 16px;
  }
  
  .register-form-section {
    padding: 24px;
  }
  
  .register-header {
    padding: 32px 24px 24px;
  }
  
  .register-title {
    font-size: 24px;
  }
  
  .form-title {
    font-size: 20px;
  }
}
</style>

