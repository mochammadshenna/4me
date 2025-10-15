<template>
  <AppLayout>
    <div class="dashboard-container">
      <!-- Welcome Section -->
      <div class="welcome-section">
        <h1 class="welcome-title">Welcome back, {{ authStore.user?.username?.toUpperCase() }}! 👋</h1>
        
        <!-- Search Bar -->
        <div class="search-container">
          <v-text-field
            v-model="searchQuery"
            placeholder="Search projects, tasks, comments..."
            variant="outlined"
            class="search-bar"
            prepend-inner-icon="mdi-magnify"
            hide-details
          />
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="stats-grid">
        <!-- Active Projects Card -->
        <div class="stat-card streak-card">
          <div class="stat-icon">
            <v-icon size="24">mdi-chart-line</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ projects.length }}</div>
            <div class="stat-label">
              <v-icon size="16" class="mr-1">mdi-folder</v-icon>
              Active Projects
            </div>
          </div>
        </div>

        <!-- Upcoming Tasks Card -->
        <div class="stat-card upcoming-card">
          <v-icon size="24" class="stat-icon">mdi-calendar</v-icon>
          <div class="stat-content">
            <div class="stat-title">Upcoming</div>
            <div class="stat-subtitle">No upcoming deadlines</div>
          </div>
          <v-icon size="16" class="notification-icon">mdi-bell-outline</v-icon>
        </div>

        <!-- Recent Activity Card -->
        <div class="stat-card recent-card">
          <v-icon size="24" class="stat-icon">mdi-upload</v-icon>
          <div class="stat-content">
            <div class="stat-title">Recent Activity</div>
            <div class="stat-subtitle">No recent activity</div>
          </div>
        </div>
      </div>

      <!-- Projects Section -->
      <div class="projects-section">
        <h2 class="section-title">My Projects</h2>
        
        <div v-if="loading" class="loading-grid">
          <div v-for="i in 3" :key="i" class="loading-card"></div>
        </div>
        
        <div v-else-if="projects.length === 0" class="empty-state">
          <v-icon size="80" color="#B0BEC5">mdi-folder-open-outline</v-icon>
          <h3 class="empty-title">No projects yet</h3>
          <p class="empty-subtitle">Create your first project to get started</p>
        </div>
        
        <div v-else class="projects-grid">
          <div
            v-for="project in projects"
            :key="project.id"
            class="project-card"
            data-testid="project-card"
            @click="goToProject(project.id)"
          >
            <div class="project-header">
              <div class="project-icon" :style="{ backgroundColor: project.color }">
                <v-icon color="white" size="20">mdi-folder</v-icon>
              </div>
              <div class="project-actions">
                <v-btn
                  icon
                  size="small"
                  variant="text"
                  @click.stop="editProject(project)"
                >
                  <v-icon size="16">mdi-pencil</v-icon>
                </v-btn>
                <v-btn
                  icon
                  size="small"
                  variant="text"
                  color="error"
                  @click.stop="confirmDelete(project)"
                >
                  <v-icon size="16">mdi-delete</v-icon>
                </v-btn>
              </div>
            </div>
            <h3 class="project-name">{{ project.name }}</h3>
            <p class="project-description">
              {{ project.description || 'No description provided' }}
            </p>
            <div class="project-footer">
              <div class="project-date">
                <v-icon size="14" class="mr-1">mdi-calendar</v-icon>
                {{ formatDate(project.created_at) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Action Button -->
    <v-btn
      class="fab-button"
      size="large"
      icon
      data-testid="create-project-button"
      @click="showCreateDialog = true"
    >
      <v-icon size="28">mdi-plus</v-icon>
    </v-btn>

    <!-- Create/Edit Project Dialog -->
    <v-dialog v-model="showCreateDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          {{ editingProject ? 'Edit Project' : 'Create New Project' }}
        </v-card-title>
        <v-card-text>
          <v-form ref="projectForm">
            <v-text-field
              v-model="projectForm.name"
              label="Project Name"
              :rules="[v => !!v || 'Name is required']"
              variant="outlined"
              class="mb-4"
              data-testid="project-name-input"
            />
            <v-textarea
              v-model="projectForm.description"
              label="Description (optional)"
              variant="outlined"
              rows="3"
              class="mb-4"
              data-testid="project-description-input"
            />
            <div class="mb-4">
              <label class="text-sm text-gray-600 mb-2 block">Project Color</label>
              <div class="flex gap-2">
                <div
                  v-for="color in colors"
                  :key="color"
                  class="w-10 h-10 rounded-full cursor-pointer border-2 transition-all"
                  :class="projectForm.color === color ? 'border-gray-800 scale-110' : 'border-transparent'"
                  :style="{ backgroundColor: color }"
                  @click="projectForm.color = color"
                />
              </div>
            </div>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="showCreateDialog = false" data-testid="cancel-project-button">Cancel</v-btn>
          <v-btn color="primary" @click="handleSaveProject" :loading="saving" data-testid="save-project-button">
            {{ editingProject ? 'Save' : 'Create' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">Delete Project?</v-card-title>
        <v-card-text>
          Are you sure you want to delete "{{ projectToDelete?.name }}"? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="showDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="handleDeleteProject" :loading="deleting">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </AppLayout>
</template>

<script setup>
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { useProjectsStore } from '@/stores/projects'
import { format } from 'date-fns'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const projectsStore = useProjectsStore()

const projects = computed(() => projectsStore.projects)
const loading = computed(() => projectsStore.loading)

const searchQuery = ref('')
const showCreateDialog = ref(false)
const showDeleteDialog = ref(false)
const editingProject = ref(null)
const projectToDelete = ref(null)
const saving = ref(false)
const deleting = ref(false)

const projectForm = ref({
  name: '',
  description: '',
  color: '#3B82F6',
})

const colors = [
  '#3B82F6', '#10B981', '#F59E0B', '#EF4444',
  '#8B5CF6', '#EC4899', '#06B6D4', '#64748B',
]

const userInitials = computed(() => {
  const name = authStore.user?.username || ''
  return name.substring(0, 2).toUpperCase()
})

onMounted(() => {
  projectsStore.fetchProjects()
})

function formatDate(date) {
  return format(new Date(date), 'MMM d, yyyy')
}

function goToProject(id) {
  router.push(`/projects/${id}`)
}

function editProject(project) {
  editingProject.value = project
  projectForm.value = {
    name: project.name,
    description: project.description || '',
    color: project.color,
  }
  showCreateDialog.value = true
}

function confirmDelete(project) {
  projectToDelete.value = project
  showDeleteDialog.value = true
}

async function handleSaveProject() {
  // Validate form
  if (!projectForm.value.name || projectForm.value.name.trim() === '') {
    console.error('Project name is required')
    return
  }

  saving.value = true
  
  try {
    let result
    if (editingProject.value) {
      result = await projectsStore.updateProject(editingProject.value.id, projectForm.value)
    } else {
      result = await projectsStore.createProject(projectForm.value)
    }
    
    if (result.success) {
      showCreateDialog.value = false
      editingProject.value = null
      projectForm.value = {
        name: '',
        description: '',
        color: '#3B82F6',
      }
      console.log('Project saved successfully')
    } else {
      console.error('Failed to save project:', result.error)
      alert('Failed to save project: ' + result.error)
    }
  } catch (error) {
    console.error('Error saving project:', error)
    alert('Error saving project: ' + error.message)
  } finally {
    saving.value = false
  }
}

async function handleDeleteProject() {
  deleting.value = true
  await projectsStore.deleteProject(projectToDelete.value.id)
  deleting.value = false
  showDeleteDialog.value = false
  projectToDelete.value = null
}

</script>

<style scoped>
/* Dashboard Container */
.dashboard-container {
  width: 100%;
  max-width: 100%;
}

/* Dashboard Specific Styles */
.welcome-section {
  margin-bottom: 32px;
  padding: 0;
  width: 100%;
}

.welcome-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1d29;
  margin: 0 0 24px 0;
  letter-spacing: -0.02em;
}

.search-container {
  max-width: 600px;
  width: 100%;
}

.search-bar {
  background: white;
  border-radius: 12px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

/* Responsive breakpoints for stats grid */
@media (min-width: 768px) and (max-width: 1199px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  }
}

@media (min-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
}

@media (min-width: 1600px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }
}


.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.streak-card {
  background: linear-gradient(135deg, #6C5CE7, #5F4DD3);
  color: white;
}

.upcoming-card, .recent-card {
  position: relative;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  display: flex;
  align-items: center;
}

.stat-title {
  font-size: 16px;
  font-weight: 600;
  color: #37474F;
  margin-bottom: 4px;
}

.stat-subtitle {
  font-size: 14px;
  color: #78909C;
}

.notification-icon {
  position: absolute;
  top: 16px;
  right: 16px;
  color: #78909C;
}

/* Projects Section */
.projects-section {
  margin-bottom: 40px;
  width: 100%;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.section-title {
  font-size: 24px;
  font-weight: 700;
  color: #37474F;
  margin: 0 0 24px 0;
}

.loading-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.loading-card {
  background: white;
  border-radius: 16px;
  height: 200px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #78909C;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  margin: 16px 0 8px 0;
  color: #546E7A;
}

.empty-subtitle {
  font-size: 16px;
  margin: 0;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

/* Responsive breakpoints for projects grid */
@media (min-width: 768px) and (max-width: 1199px) {
  .projects-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}

@media (min-width: 1200px) {
  .projects-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (min-width: 1600px) {
  .projects-grid {
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  }
}

@media (min-width: 2000px) {
  .projects-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  }
}


.project-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #e5e7eb;
}

.project-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  border-color: #6C5CE7;
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.project-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.project-actions {
  display: flex;
  gap: 4px;
}

.project-name {
  font-size: 18px;
  font-weight: 600;
  color: #37474F;
  margin: 0 0 8px 0;
}

.project-description {
  font-size: 14px;
  color: #78909C;
  line-height: 1.5;
  margin: 0 0 16px 0;
}

.project-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.project-date {
  font-size: 12px;
  color: #78909C;
  display: flex;
  align-items: center;
}

/* Floating Action Button */
.fab-button {
  position: fixed !important;
  bottom: 24px !important;
  right: 24px !important;
  background: #4CAF50 !important;
  color: white !important;
  box-shadow: 0 8px 24px rgba(76, 175, 80, 0.4) !important;
  z-index: 1000;
}

.fab-button:hover {
  background: #45A049 !important;
  transform: scale(1.05);
}

/* Dialog Styling */
.v-dialog .v-card {
  border-radius: 16px !important;
  background: white !important;
}

.v-dialog .v-card-title {
  color: #37474F !important;
  font-weight: 600 !important;
  font-size: 20px !important;
}

/* Custom Form Styling */
:deep(.v-text-field--outlined .v-field) {
  border-radius: 12px !important;
}

:deep(.v-textarea--outlined .v-field) {
  border-radius: 12px !important;
}

/* Ensure dashboard content adapts to available space */
.dashboard-container {
  width: 100%;
  max-width: 100%;
  margin: 0 auto;
  overflow-x: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  flex: 1;
  box-sizing: border-box;
  padding: 0 24px;
}

@media (min-width: 768px) {
  .dashboard-container {
    max-width: 1800px;
    padding: 0 32px;
  }
}

/* Responsive padding for different screen sizes */
@media (min-width: 1400px) {
  .dashboard-container {
    padding: 0 48px;
  }
}

@media (min-width: 1800px) {
  .dashboard-container {
    padding: 0 64px;
    max-width: 2000px;
  }
}

/* Prevent content from being cropped */
.stats-grid,
.projects-grid {
  width: 100%;
  max-width: 100%;
  overflow: visible;
}

/* Ensure all dashboard sections use full width */
.dashboard-container > * {
  width: 100%;
  max-width: 100%;
}

/* Additional responsive behavior for collapsed sidebar - handled in grid sections above */

/* Mobile Responsive Design */
@media (max-width: 768px) {
  .dashboard-container {
    width: 100vw !important;
    max-width: 100vw !important;
    padding: 0 16px !important;
    margin: 0 !important;
    box-sizing: border-box !important;
  }
  
  .welcome-section {
    width: 100% !important;
    max-width: 100% !important;
    padding: 16px 0 !important;
    margin-bottom: 24px !important;
  }
  
  .welcome-title {
    font-size: 24px !important;
    margin-bottom: 16px !important;
  }
  
  .search-container {
    width: 100% !important;
    max-width: 100% !important;
  }
  
  .search-bar {
    width: 100% !important;
  }
  
  .stats-grid {
    grid-template-columns: 1fr !important;
    gap: 16px !important;
    margin-bottom: 32px !important;
    width: 100% !important;
    max-width: 100% !important;
  }
  
  .stat-card {
    padding: 20px !important;
    width: 100% !important;
    box-sizing: border-box !important;
  }
  
  .projects-section {
    width: 100% !important;
    max-width: 100% !important;
    margin-bottom: 32px !important;
  }
  
  .section-title {
    font-size: 20px !important;
    margin-bottom: 20px !important;
  }
  
  .projects-grid {
    grid-template-columns: 1fr !important;
    gap: 16px !important;
    width: 100% !important;
    max-width: 100% !important;
  }
  
  .project-card {
    padding: 20px !important;
    width: 100% !important;
    box-sizing: border-box !important;
  }
  
  .empty-state {
    padding: 40px 16px !important;
    width: 100% !important;
    box-sizing: border-box !important;
  }
}

/* Enhanced responsive design for larger screens */
@media (min-width: 768px) and (max-width: 1199px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  }
  
  .projects-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}

/* Optimize for ultra-wide screens when sidebar is collapsed */
@media (min-width: 2000px) {
  .welcome-section {
    max-width: 1800px;
    margin: 0 auto 32px auto;
  }
  
  .search-container {
    max-width: 800px;
  }
  
  .projects-section {
    max-width: 1800px;
    margin: 0 auto 40px auto;
  }
}
</style>

