<template>
  <v-card
    ref="taskCardRef"
    class="task-card cursor-pointer hover:shadow-md transition-all"
    :class="{ 'is-dragging': isDragging, 'drag-over': isDraggedOver }"
    elevation="1"
    :data-task-id="task.id"
    :data-board-id="task.board_id"
    @click="handleClick"
  >
    <div class="drag-handle cursor-move pa-3" :data-drag-handle="true">
      <div class="flex items-start justify-between mb-2">
        <h4 class="text-base font-medium flex-1">{{ task.title }}</h4>
        <v-icon size="small" class="text-gray-400 ml-2">
          mdi-drag-vertical
        </v-icon>
      </div>
      
      <p v-if="task.description" class="text-sm text-gray-600 mb-2 line-clamp-2">
        {{ task.description }}
      </p>
      
      <!-- Labels -->
      <div v-if="task.labels && task.labels.length > 0" class="flex flex-wrap gap-1 mb-2">
        <v-chip
          v-for="label in task.labels.slice(0, 3)"
          :key="label.id"
          :color="label.color"
          size="x-small"
          label
        >
          {{ label.name }}
        </v-chip>
        <v-chip v-if="task.labels.length > 3" size="x-small" variant="text">
          +{{ task.labels.length - 3 }}
        </v-chip>
      </div>
      
      <!-- Task Metadata -->
      <div class="flex items-center gap-3 text-xs text-gray-500">
        <!-- Priority Badge -->
        <v-chip
          :color="priorityColor"
          size="x-small"
          variant="flat"
        >
          {{ priorityText }}
        </v-chip>
        
        <!-- Due Date -->
        <div v-if="task.due_date" class="flex items-center gap-1">
          <v-icon size="small" :color="isDueDatePast ? 'error' : 'default'">
            mdi-calendar
          </v-icon>
          <span :class="{ 'text-error': isDueDatePast }">
            {{ formattedDueDate }}
          </span>
        </div>
        
        <!-- Comments Count -->
        <div v-if="task.comments_count" class="flex items-center gap-1">
          <v-icon size="small">mdi-comment-outline</v-icon>
          <span>{{ task.comments_count }}</span>
        </div>
        
        <!-- Attachments Count -->
        <div v-if="task.attachments_count" class="flex items-center gap-1">
          <v-icon size="small">mdi-paperclip</v-icon>
          <span>{{ task.attachments_count }}</span>
        </div>
      </div>
    </div>
  </v-card>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { format, isPast, parseISO } from 'date-fns'
import { draggable } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'
import { setCustomNativeDragPreview } from '@atlaskit/pragmatic-drag-and-drop/element/set-custom-native-drag-preview'
import { pointerOutsideOfPreview } from '@atlaskit/pragmatic-drag-and-drop/element/pointer-outside-of-preview'

const props = defineProps({
  task: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['click'])

// Drag state
const taskCardRef = ref(null)
const isDragging = ref(false)
const isDraggedOver = ref(false)
let cleanupDraggable = null

// Handle click - prevent firing when dragging
const handleClick = (event) => {
  if (!isDragging.value) {
    emit('click', event)
  }
}

const priorityColor = computed(() => {
  const colors = {
    low: 'blue-grey',
    medium: 'blue',
    high: 'orange',
    urgent: 'red',
  }
  return colors[props.task.priority] || 'blue-grey'
})

const priorityText = computed(() => {
  return props.task.priority?.toUpperCase() || 'MEDIUM'
})

const formattedDueDate = computed(() => {
  if (!props.task.due_date) return ''
  try {
    return format(parseISO(props.task.due_date), 'MMM d')
  } catch {
    return ''
  }
})

const isDueDatePast = computed(() => {
  if (!props.task.due_date) return false
  try {
    return isPast(parseISO(props.task.due_date))
  } catch {
    return false
  }
})

// Setup draggable behavior
onMounted(() => {
  if (!taskCardRef.value) return

  cleanupDraggable = draggable({
    element: taskCardRef.value.$el || taskCardRef.value,
    getInitialData: () => ({
      type: 'task',
      taskId: props.task.id,
      boardId: props.task.board_id,
      task: props.task
    }),
    onGenerateDragPreview: ({ nativeSetDragImage }) => {
      setCustomNativeDragPreview({
        nativeSetDragImage,
        getOffset: pointerOutsideOfPreview({
          x: '16px',
          y: '8px'
        }),
        render: ({ container }) => {
          const preview = taskCardRef.value.$el.cloneNode(true)
          preview.style.width = `${taskCardRef.value.$el.offsetWidth}px`
          preview.style.transform = 'rotate(3deg)'
          preview.style.opacity = '0.9'
          preview.style.boxShadow = '0 12px 48px rgba(0, 0, 0, 0.25)'
          container.appendChild(preview)
        }
      })
    },
    onDragStart: () => {
      isDragging.value = true
    },
    onDrop: () => {
      isDragging.value = false
    }
  })
})

onUnmounted(() => {
  if (cleanupDraggable) {
    cleanupDraggable()
  }
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.task-card {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.task-card.is-dragging {
  opacity: 0.5;
  transform: scale(0.98);
}

.task-card.drag-over {
  border: 2px dashed rgba(var(--v-theme-primary), 0.5);
  background-color: rgba(var(--v-theme-primary), 0.05);
}

.drag-handle {
  touch-action: none;
}
</style>

