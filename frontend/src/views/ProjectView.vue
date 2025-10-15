<template>
  <AppLayout>
    <!-- Project Header -->
    <div class="project-header">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center">
          <v-btn icon @click="goBack" class="back-btn">
            <v-icon>mdi-arrow-left</v-icon>
          </v-btn>
          <div class="ml-4">
            <h1 class="project-title">
              <span v-if="project">{{ project.name }}</span>
              <v-skeleton-loader v-else type="text" width="200" />
            </h1>
            <p v-if="project" class="project-description">{{ project.description }}</p>
          </div>
        </div>
        
        <div class="project-actions">
          <v-btn icon @click="showLabelsDialog = true" class="action-btn">
            <v-icon>mdi-label</v-icon>
          </v-btn>
        </div>
      </div>
    </div>
    
    <!-- Kanban Board -->
    <div class="kanban-board">
      <div class="kanban-container">
        <div class="kanban-columns">
          <!-- Board Columns -->
          <BoardColumn
            v-for="board in boards"
            :key="board.id"
            :board="board"
            :project-id="projectId"
            @task-click="openTaskDetail"
            @update-board="handleUpdateBoard"
            @delete-board="handleDeleteBoard"
          />
          
          <!-- Add Board Button -->
          <div class="flex-shrink-0 w-80">
            <v-card
              v-if="!showAddBoard"
              class="pa-4 cursor-pointer hover:bg-gray-100"
              elevation="0"
              @click="showAddBoard = true"
            >
              <div class="flex items-center text-gray-600">
                <v-icon class="mr-2">mdi-plus</v-icon>
                <span>Add Board</span>
              </div>
            </v-card>
            
            <v-card v-else class="pa-4" elevation="2">
              <v-text-field
                v-model="newBoardName"
                label="Board Name"
                variant="outlined"
                density="compact"
                autofocus
                @keyup.enter="handleAddBoard"
                @keyup.esc="cancelAddBoard"
              />
              <div class="flex gap-2 mt-2">
                <v-btn color="primary" size="small" @click="handleAddBoard">
                  Add
                </v-btn>
                <v-btn size="small" @click="cancelAddBoard">
                  Cancel
                </v-btn>
              </div>
            </v-card>
          </div>
        </div>
      </div>

      <!-- Task Detail Dialog -->
    <TaskDetailDialog
      v-model="showTaskDialog"
      :task-id="selectedTaskId"
      :project-id="projectId"
      @task-updated="handleTaskUpdated"
      @task-deleted="handleTaskDeleted"
    />
    
      <!-- Labels Management Dialog -->
      <LabelsDialog
        v-model="showLabelsDialog"
        :project-id="projectId"
      />
    </div>
  </AppLayout>
</template>

<script setup>
import BoardColumn from '@/components/board/BoardColumn.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import LabelsDialog from '@/components/task/LabelsDialog.vue'
import TaskDetailDialog from '@/components/task/TaskDetailDialog.vue'
import { useAuthStore } from '@/stores/auth'
import { useBoardsStore } from '@/stores/boards'
import { useLabelsStore } from '@/stores/labels'
import { useProjectsStore } from '@/stores/projects'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const projectsStore = useProjectsStore()
const boardsStore = useBoardsStore()
const labelsStore = useLabelsStore()

const projectId = computed(() => parseInt(route.params.id))
const project = computed(() => projectsStore.currentProject)
const boards = computed(() => boardsStore.boards)

const showAddBoard = ref(false)
const newBoardName = ref('')
const showTaskDialog = ref(false)
const selectedTaskId = ref(null)
const showLabelsDialog = ref(false)

const userInitials = computed(() => {
  const name = authStore.user?.username || ''
  return name.substring(0, 2).toUpperCase()
})

onMounted(async () => {
  await projectsStore.fetchProject(projectId.value)
  await boardsStore.fetchBoards(projectId.value)
  await labelsStore.fetchLabels(projectId.value)
})

function goBack() {
  router.push('/')
}

async function handleAddBoard() {
  if (!newBoardName.value.trim()) return
  
  const position = boards.value.length
  await boardsStore.createBoard(projectId.value, {
    name: newBoardName.value,
    position,
  })
  
  newBoardName.value = ''
  showAddBoard.value = false
}

function cancelAddBoard() {
  newBoardName.value = ''
  showAddBoard.value = false
}

async function handleUpdateBoard(id, data) {
  await boardsStore.updateBoard(id, data)
}

async function handleDeleteBoard(id) {
  await boardsStore.deleteBoard(id)
}

function openTaskDetail(taskId) {
  selectedTaskId.value = taskId
  showTaskDialog.value = true
}

function handleTaskUpdated() {
  // Refresh boards to get updated tasks
  boardsStore.fetchBoards(projectId.value)
}

function handleTaskDeleted() {
  showTaskDialog.value = false
  boardsStore.fetchBoards(projectId.value)
}

</script>

<style scoped>
/* Project View Specific Styles */
.project-header {
  margin-bottom: 32px;
}

.back-btn {
  background: white !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
  border-radius: 12px !important;
}

.project-title {
  font-size: 28px;
  font-weight: 700;
  color: #37474F;
  margin: 0 0 8px 0;
}

.project-description {
  font-size: 16px;
  color: #78909C;
  margin: 0;
}

.project-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  background: white !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
  border-radius: 12px !important;
}

.kanban-board {
  background: linear-gradient(135deg, #F8F9FA 0%, #E3F2FD 100%);
  border-radius: 20px;
  padding: 24px;
  min-height: calc(100vh - 200px);
  position: relative;
  overflow: hidden;
}

.kanban-board::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 20%, rgba(227, 242, 253, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(255, 243, 224, 0.2) 0%, transparent 50%);
  pointer-events: none;
}

.kanban-container {
  position: relative;
  z-index: 1;
  height: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.kanban-columns {
  display: flex;
  gap: 20px;
  padding: 12px;
  min-height: calc(100vh - 300px);
  align-items: flex-start;
}

/* Responsive Design */
@media (max-width: 768px) {
  .project-header {
    margin-bottom: 24px;
  }
  
  .project-title {
    font-size: 24px;
  }
  
  .kanban-board {
    padding: 16px;
  }
}
</style>

