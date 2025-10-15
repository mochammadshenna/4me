<template>
  <div class="board-column">
    <!-- Column Header -->
    <div class="column-header">
      <div class="column-title-section">
        <div class="column-color-indicator" :style="{ backgroundColor: getColumnColor(board.name) }"></div>
        <h3 class="column-title">{{ board.name }}</h3>
        <span class="task-count">{{ tasks.length }}</span>
      </div>
      
      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn icon size="small" variant="text" class="column-menu-btn" v-bind="props">
            <v-icon size="18">mdi-dots-horizontal</v-icon>
          </v-btn>
        </template>
        <v-list class="column-menu" density="compact">
          <v-list-item @click="showEditDialog = true">
            <template v-slot:prepend>
              <v-icon size="16">mdi-pencil</v-icon>
            </template>
            <v-list-item-title>Rename</v-list-item-title>
          </v-list-item>
          <v-list-item @click="confirmDelete">
            <template v-slot:prepend>
              <v-icon size="16" color="error">mdi-delete</v-icon>
            </template>
            <v-list-item-title class="text-error">Delete</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </div>
    
    <!-- Tasks Container -->
    <div
      ref="tasksContainerRef"
      class="tasks-container"
      :class="{ 'drag-over': isDragOver, 'drop-target-active': isDropTarget }"
      :data-board-id="board.id"
    >
      <div class="tasks-list">
        <!-- Drop Indicator at top -->
        <div
          v-if="dropIndicatorIndex === 0"
          class="drop-indicator"
          :style="{ height: '4px', marginBottom: '8px' }"
        ></div>

        <template v-for="(task, index) in tasks" :key="task.id">
          <div class="task-wrapper">
            <TaskCard
              :task="task"
              @click="$emit('task-click', task.id)"
            />
          </div>

          <!-- Drop Indicator between tasks -->
          <div
            v-if="dropIndicatorIndex === index + 1"
            class="drop-indicator"
            :style="{ height: '4px', margin: '8px 0' }"
          ></div>
        </template>

        <!-- Empty State -->
        <div v-if="tasks.length === 0" class="empty-state">
          <v-icon size="48" color="#BDBDBD">mdi-clipboard-outline</v-icon>
          <p class="empty-text">No tasks yet</p>
          <p class="empty-subtext">Drag a task here or create a new one</p>
        </div>
      </div>
    </div>
      
      <!-- Add Task Button -->
      <div class="pa-3 border-t">
        <v-btn
          v-if="!showAddTask"
          block
          variant="text"
          @click="showAddTask = true"
        >
          <v-icon start>mdi-plus</v-icon>
          Add Task
        </v-btn>
        
        <div v-else>
          <v-textarea
            v-model="newTaskTitle"
            label="Task title"
            variant="outlined"
            density="compact"
            rows="2"
            autofocus
            @keyup.ctrl.enter="handleAddTask"
            @keyup.esc="cancelAddTask"
          />
          <div class="flex gap-2 mt-2">
            <v-btn color="primary" size="small" @click="handleAddTask" :loading="addingTask">
              Add
            </v-btn>
            <v-btn size="small" @click="cancelAddTask">
              Cancel
            </v-btn>
          </div>
        </div>
      </div>
    
    <!-- Edit Board Dialog -->
    <v-dialog v-model="showEditDialog" max-width="400">
      <v-card>
        <v-card-title>Rename Board</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="editName"
            label="Board Name"
            variant="outlined"
            autofocus
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="showEditDialog = false">Cancel</v-btn>
          <v-btn color="primary" @click="handleUpdateBoard">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title>Delete Board?</v-card-title>
        <v-card-text>
          Are you sure you want to delete "{{ board.name }}"? All tasks in this board will be deleted.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="showDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="handleDeleteBoard">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import apiClient from '@/api/client'
import TaskCard from '@/components/task/TaskCard.vue'
import { useTasksStore } from '@/stores/tasks'
import { onMounted, onUnmounted, ref } from 'vue'
import { dropTargetForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'
import { autoScrollForElements } from '@atlaskit/pragmatic-drag-and-drop-auto-scroll/element'
import { attachClosestEdge, extractClosestEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge'

const props = defineProps({
  board: {
    type: Object,
    required: true,
  },
  projectId: {
    type: Number,
    required: true,
  },
})

const emit = defineEmits(['task-click', 'update-board', 'delete-board'])

const tasksStore = useTasksStore()

const tasks = ref([])
const showAddTask = ref(false)
const newTaskTitle = ref('')
const addingTask = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const editName = ref(props.board.name)
const isDragOver = ref(false)
const isDropTarget = ref(false)
const dropIndicatorIndex = ref(null)
const tasksContainerRef = ref(null)
let cleanupDropTarget = null
let cleanupAutoScroll = null

// Column colors for visual distinction
function getColumnColor(columnName) {
  const colors = {
    'To Do': '#4CAF50',
    'In Progress': '#FF9800', 
    'Done': '#2196F3',
    'Backlog': '#9C27B0',
    'Review': '#FF5722'
  }
  return colors[columnName] || '#607D8B'
}

// Calculate drop position based on pointer location
function getDropIndex(source, event) {
  const taskElements = Array.from(
    tasksContainerRef.value?.querySelectorAll('.task-wrapper') || []
  )

  if (taskElements.length === 0) return 0

  const pointerY = event.clientY || event.pageY

  for (let i = 0; i < taskElements.length; i++) {
    const rect = taskElements[i].getBoundingClientRect()
    const midpoint = rect.top + rect.height / 2

    if (pointerY < midpoint) {
      return i
    }
  }

  return taskElements.length
}

// Setup drop target and auto-scroll
onMounted(async () => {
  await loadTasks()

  if (!tasksContainerRef.value) return

  // Setup drop target
  cleanupDropTarget = dropTargetForElements({
    element: tasksContainerRef.value,
    canDrop: ({ source }) => {
      return source.data.type === 'task'
    },
    getData: ({ input }) => {
      return attachClosestEdge(
        { boardId: props.board.id },
        {
          element: tasksContainerRef.value,
          input,
          allowedEdges: ['top', 'bottom']
        }
      )
    },
    onDragEnter: ({ source }) => {
      isDropTarget.value = true
      isDragOver.value = true
    },
    onDrag: ({ source, location }) => {
      const closestEdge = extractClosestEdge(location.current.dropTargets[0]?.data)

      // Calculate drop indicator position
      const sourceTaskId = source.data.taskId
      const sourceIndex = tasks.value.findIndex(t => t.id === sourceTaskId)

      // Show indicator at appropriate position
      if (closestEdge === 'top') {
        dropIndicatorIndex.value = 0
      } else if (closestEdge === 'bottom') {
        dropIndicatorIndex.value = tasks.value.length
      }
    },
    onDragLeave: () => {
      isDropTarget.value = false
      isDragOver.value = false
      dropIndicatorIndex.value = null
    },
    onDrop: async ({ source, location }) => {
      isDropTarget.value = false
      isDragOver.value = false
      dropIndicatorIndex.value = null

      const taskId = source.data.taskId
      const sourceBoardId = source.data.boardId
      const targetBoardId = props.board.id

      if (!taskId) return

      // Find the dropped task
      const task = source.data.task

      // Calculate new position
      const closestEdge = extractClosestEdge(location.current.dropTargets[0]?.data)
      let newIndex = closestEdge === 'top' ? 0 : tasks.value.length

      // If moving within same board, handle reordering
      if (sourceBoardId === targetBoardId) {
        const oldIndex = tasks.value.findIndex(t => t.id === taskId)
        if (oldIndex !== -1) {
          tasks.value.splice(oldIndex, 1)
          tasks.value.splice(newIndex, 0, task)
          await tasksStore.updateTask(taskId, { position: newIndex })
        }
      } else {
        // Moving to different board
        tasks.value.splice(newIndex, 0, task)
        await tasksStore.moveTask(taskId, targetBoardId, newIndex)
      }
    }
  })

  // Setup auto-scroll
  cleanupAutoScroll = autoScrollForElements({
    element: tasksContainerRef.value,
    canScroll: () => true
  })
})

onUnmounted(() => {
  if (cleanupDropTarget) {
    cleanupDropTarget()
  }
  if (cleanupAutoScroll) {
    cleanupAutoScroll()
  }
})

async function loadTasks() {
  try {
    const response = await apiClient.get(`/boards/${props.board.id}/tasks`)
    tasks.value = response.data || []
    tasksStore.setTasks(props.board.id, tasks.value)
  } catch (error) {
    console.error('Failed to load tasks:', error)
    tasks.value = []
  }
}

async function handleAddTask() {
  if (!newTaskTitle.value.trim()) return
  
  addingTask.value = true
  const result = await tasksStore.createTask(props.board.id, {
    title: newTaskTitle.value,
    priority: 'medium',
  })
  
  if (result.success) {
    tasks.value.push(result.data)
    newTaskTitle.value = ''
    showAddTask.value = false
  }
  addingTask.value = false
}

function cancelAddTask() {
  newTaskTitle.value = ''
  showAddTask.value = false
}

function handleUpdateBoard() {
  emit('update-board', props.board.id, { name: editName.value })
  showEditDialog.value = false
}

function confirmDelete() {
  showDeleteDialog.value = true
}

function handleDeleteBoard() {
  emit('delete-board', props.board.id)
  showDeleteDialog.value = false
}
</script>

<style scoped>
/* Beautiful Board Column Styling */
.board-column {
  width: 320px;
  min-width: 320px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(227, 242, 253, 0.3);
  backdrop-filter: blur(10px);
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 200px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.board-column:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.column-header {
  padding: 20px 16px 16px;
  border-bottom: 1px solid rgba(227, 242, 253, 0.3);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.column-title-section {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.column-color-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.column-title {
  font-size: 16px;
  font-weight: 600;
  color: #37474F;
  margin: 0;
}

.task-count {
  background: rgba(227, 242, 253, 0.8);
  color: #1976D2;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  min-width: 20px;
  text-align: center;
}

.column-menu-btn {
  color: #78909C !important;
  transition: all 0.3s ease !important;
}

.column-menu-btn:hover {
  color: #1976D2 !important;
  background: rgba(25, 118, 210, 0.08) !important;
}

.column-menu {
  border-radius: 12px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12) !important;
  border: 1px solid rgba(0, 0, 0, 0.08) !important;
}

.tasks-container {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  transition: all 0.3s ease;
}

.tasks-container.drag-over {
  background: rgba(25, 118, 210, 0.05);
}

.tasks-container.drop-target-active {
  border: 2px dashed #1976D2;
  border-radius: 12px;
  background: rgba(25, 118, 210, 0.08);
}

/* Drop Indicator */
.drop-indicator {
  background: linear-gradient(90deg, #1976D2 0%, #2196F3 100%);
  border-radius: 2px;
  box-shadow: 0 0 8px rgba(25, 118, 210, 0.5);
  animation: pulse 1s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.8;
  }
  50% {
    opacity: 1;
  }
}

.tasks-list {
  min-height: 100px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.task-wrapper {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.task-wrapper:hover {
  transform: translateY(-2px);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  color: #9E9E9E;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  margin: 12px 0 4px 0;
  color: #757575;
}

.empty-subtext {
  font-size: 14px;
  margin: 0;
  color: #BDBDBD;
}

/* Drag and Drop States */
:deep(.task-ghost) {
  opacity: 0.5;
  transform: rotate(5deg);
  background: rgba(25, 118, 210, 0.1) !important;
  border: 2px dashed #1976D2 !important;
}

:deep(.task-chosen) {
  transform: scale(1.02);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2) !important;
  z-index: 1000;
}

:deep(.task-drag) {
  opacity: 0.8;
  transform: rotate(2deg);
}

/* Scrollbar Styling */
.tasks-container::-webkit-scrollbar {
  width: 6px;
}

.tasks-container::-webkit-scrollbar-track {
  background: transparent;
}

.tasks-container::-webkit-scrollbar-thumb {
  background: rgba(187, 222, 251, 0.5);
  border-radius: 3px;
}

.tasks-container::-webkit-scrollbar-thumb:hover {
  background: rgba(187, 222, 251, 0.8);
}

/* Responsive Design */
@media (max-width: 768px) {
  .board-column {
    width: 280px;
    min-width: 280px;
  }
  
  .column-header {
    padding: 16px 12px 12px;
  }
  
  .tasks-container {
    padding: 8px;
  }
}

/* Animation for new tasks */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.task-wrapper {
  animation: slideIn 0.3s ease-out;
}
</style>

