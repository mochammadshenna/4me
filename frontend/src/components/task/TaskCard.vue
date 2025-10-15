<template>
  <div 
    class="task-card" 
    :class="{ 'high-priority': task.priority === 'high', 'medium-priority': task.priority === 'medium', 'low-priority': task.priority === 'low' }"
    @click="$emit('click')"
    :data-task-id="task.id"
  >
    <!-- Drag Handle -->
    <div class="task-drag-handle">
      <v-icon size="16" color="#BDBDBD">mdi-drag-horizontal</v-icon>
    </div>
    
    <!-- Task Content -->
    <div class="task-content">
      <!-- Task Header -->
      <div class="task-header">
        <div class="task-title">{{ task.title }}</div>
        <div class="task-actions">
          <v-btn 
            icon 
            size="x-small" 
            variant="text" 
            class="task-action-btn"
            @click.stop="$emit('edit')"
          >
            <v-icon size="14">mdi-pencil</v-icon>
          </v-btn>
        </div>
      </div>
      
      <!-- Task Description -->
      <div v-if="task.description" class="task-description">
        {{ truncateDescription(task.description) }}
      </div>
      
      <!-- Task Labels -->
      <div v-if="task.labels && task.labels.length > 0" class="task-labels">
        <div 
          v-for="label in task.labels.slice(0, 3)" 
          :key="label.id"
          class="task-label"
          :style="{ backgroundColor: label.color + '20', color: label.color }"
        >
          {{ label.name }}
        </div>
        <div v-if="task.labels.length > 3" class="more-labels">
          +{{ task.labels.length - 3 }}
        </div>
      </div>
      
      <!-- Task Footer -->
      <div class="task-footer">
        <div class="task-meta">
          <!-- Priority Indicator -->
          <div 
            v-if="task.priority" 
            class="priority-indicator"
            :class="`priority-${task.priority}`"
          >
            <v-icon size="12">
              {{ getPriorityIcon(task.priority) }}
            </v-icon>
          </div>
          
          <!-- Due Date -->
          <div v-if="task.due_date" class="due-date" :class="{ 'overdue': isOverdue(task.due_date) }">
            <v-icon size="12">mdi-calendar</v-icon>
            <span>{{ formatDueDate(task.due_date) }}</span>
          </div>
          
          <!-- Attachments Count -->
          <div v-if="task.attachments && task.attachments.length > 0" class="attachments-count">
            <v-icon size="12">mdi-paperclip</v-icon>
            <span>{{ task.attachments.length }}</span>
          </div>
          
          <!-- Comments Count -->
          <div v-if="task.comments && task.comments.length > 0" class="comments-count">
            <v-icon size="12">mdi-comment</v-icon>
            <span>{{ task.comments.length }}</span>
          </div>
        </div>
        
        <!-- Assignee Avatar -->
        <div class="task-assignee">
          <v-avatar size="24" class="assignee-avatar">
            <v-icon size="16">mdi-account</v-icon>
          </v-avatar>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { format, parseISO } from 'date-fns'

defineProps({
  task: {
    type: Object,
    required: true
  }
})

defineEmits(['click', 'edit'])

function truncateDescription(description) {
  if (description.length > 100) {
    return description.substring(0, 100) + '...'
  }
  return description
}

function getPriorityIcon(priority) {
  const icons = {
    high: 'mdi-arrow-up-bold',
    medium: 'mdi-minus',
    low: 'mdi-arrow-down-bold'
  }
  return icons[priority] || 'mdi-minus'
}

function formatDueDate(dateString) {
  const date = parseISO(dateString)
  const now = new Date()
  const diffInDays = Math.ceil((date - now) / (1000 * 60 * 60 * 24))
  
  if (diffInDays === 0) return 'Today'
  if (diffInDays === 1) return 'Tomorrow'
  if (diffInDays === -1) return 'Yesterday'
  if (diffInDays > 0) return `${diffInDays}d`
  if (diffInDays < 0) return `${Math.abs(diffInDays)}d ago`
  
  return format(date, 'MMM d')
}

function isOverdue(dateString) {
  const date = parseISO(dateString)
  const now = new Date()
  return date < now
}
</script>

<style scoped>
/* Beautiful Task Card Styling */
.task-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(227, 242, 253, 0.3);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.task-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
  border-color: rgba(25, 118, 210, 0.3);
}

.task-card.high-priority {
  border-left: 4px solid #f44336;
}

.task-card.medium-priority {
  border-left: 4px solid #ff9800;
}

.task-card.low-priority {
  border-left: 4px solid #4caf50;
}

.task-drag-handle {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
  cursor: grab;
  padding: 4px;
  border-radius: 4px;
}

.task-card:hover .task-drag-handle {
  opacity: 1;
}

.task-drag-handle:active {
  cursor: grabbing;
}

.task-content {
  padding: 16px;
  padding-right: 32px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.task-title {
  font-size: 14px;
  font-weight: 600;
  color: #37474F;
  line-height: 1.4;
  flex: 1;
  margin-right: 8px;
}

.task-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.task-card:hover .task-actions {
  opacity: 1;
}

.task-action-btn {
  color: #78909C !important;
  transition: all 0.3s ease !important;
}

.task-action-btn:hover {
  color: #1976D2 !important;
  background: rgba(25, 118, 210, 0.08) !important;
}

.task-description {
  font-size: 13px;
  color: #78909C;
  line-height: 1.4;
  margin-bottom: 12px;
}

.task-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.task-label {
  font-size: 11px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 12px;
  background: rgba(25, 118, 210, 0.1);
  color: #1976D2;
}

.more-labels {
  font-size: 11px;
  color: #78909C;
  padding: 4px 8px;
  border-radius: 12px;
  background: rgba(189, 189, 189, 0.1);
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.priority-indicator {
  display: flex;
  align-items: center;
  padding: 4px 6px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 500;
}

.priority-high {
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
}

.priority-medium {
  background: rgba(255, 152, 0, 0.1);
  color: #ff9800;
}

.priority-low {
  background: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.due-date, .attachments-count, .comments-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #78909C;
  padding: 4px 6px;
  border-radius: 8px;
  background: rgba(227, 242, 253, 0.5);
}

.due-date.overdue {
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
}

.task-assignee {
  display: flex;
  align-items: center;
}

.assignee-avatar {
  background: linear-gradient(135deg, #1976D2, #1565C0) !important;
  color: white !important;
}

/* Responsive Design */
@media (max-width: 768px) {
  .task-content {
    padding: 12px;
    padding-right: 28px;
  }
  
  .task-title {
    font-size: 13px;
  }
  
  .task-description {
    font-size: 12px;
  }
  
  .task-meta {
    gap: 8px;
  }
}

/* Animation */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.task-card {
  animation: slideIn 0.3s ease-out;
}
</style>
