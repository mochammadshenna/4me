<template>
  <v-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    max-width="900"
    scrollable
  >
    <v-card v-if="task" class="task-detail-dialog">
      <v-card-title class="d-flex align-center pa-4 border-b">
        <v-text-field
          v-model="task.title"
          variant="plain"
          density="compact"
          class="text-h5 font-weight-bold"
          @blur="handleUpdateTask"
        />
        <v-btn icon @click="close">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      
      <v-card-text class="pa-6">
        <v-row>
          <v-col cols="12" md="8">
            <!-- Description -->
            <div class="mb-6">
              <h3 class="text-lg font-semibold mb-2 flex items-center gap-2">
                <v-icon>mdi-text</v-icon>
                Description
              </h3>
              <v-textarea
                v-model="task.description"
                variant="outlined"
                rows="4"
                placeholder="Add a description..."
                @blur="handleUpdateTask"
              />
            </div>
            
            <!-- Comments Section -->
            <div class="mb-6">
              <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
                <v-icon>mdi-comment-outline</v-icon>
                Comments ({{ comments.length }})
              </h3>
              
              <!-- Add Comment -->
              <div class="mb-4">
                <v-textarea
                  v-model="newComment"
                  variant="outlined"
                  rows="2"
                  placeholder="Write a comment..."
                  density="compact"
                />
                <v-btn
                  color="primary"
                  size="small"
                  class="mt-2"
                  @click="handleAddComment"
                  :disabled="!newComment.trim()"
                  :loading="addingComment"
                >
                  Add Comment
                </v-btn>
              </div>
              
              <!-- Comments List -->
              <div class="space-y-3">
                <v-card
                  v-for="comment in comments"
                  :key="comment.id"
                  variant="outlined"
                  class="pa-3"
                >
                  <div class="flex items-start justify-between mb-2">
                    <div class="flex items-center gap-2">
                      <v-avatar size="32" color="primary">
                        <v-img v-if="comment.user?.avatar_url" :src="comment.user.avatar_url" />
                        <span v-else class="text-white text-sm">
                          {{ comment.user?.username?.substring(0, 2).toUpperCase() }}
                        </span>
                      </v-avatar>
                      <div>
                        <div class="font-semibold">{{ comment.user?.username }}</div>
                        <div class="text-xs text-gray-500">{{ formatDate(comment.created_at) }}</div>
                      </div>
                    </div>
                    
                    <v-menu v-if="isOwnComment(comment)">
                      <template v-slot:activator="{ props }">
                        <v-btn icon size="x-small" variant="text" v-bind="props">
                          <v-icon>mdi-dots-vertical</v-icon>
                        </v-btn>
                      </template>
                      <v-list density="compact">
                        <v-list-item @click="startEditComment(comment)">
                          <v-list-item-title>Edit</v-list-item-title>
                        </v-list-item>
                        <v-list-item @click="handleDeleteComment(comment.id)">
                          <v-list-item-title class="text-error">Delete</v-list-item-title>
                        </v-list-item>
                      </v-list>
                    </v-menu>
                  </div>
                  
                  <v-textarea
                    v-if="editingComment?.id === comment.id"
                    v-model="editingCommentText"
                    variant="outlined"
                    rows="2"
                    density="compact"
                  />
                  <p v-else class="text-sm">{{ comment.content }}</p>
                  
                  <div v-if="editingComment?.id === comment.id" class="flex gap-2 mt-2">
                    <v-btn size="small" color="primary" @click="handleUpdateComment">Save</v-btn>
                    <v-btn size="small" @click="cancelEditComment">Cancel</v-btn>
                  </div>
                </v-card>
              </div>
            </div>
            
            <!-- Attachments Section -->
            <div class="mb-6">
              <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
                <v-icon>mdi-paperclip</v-icon>
                Attachments ({{ attachments.length }})
              </h3>
              
              <input
                ref="fileInput"
                type="file"
                class="hidden"
                @change="handleFileUpload"
              />
              
              <v-btn
                color="primary"
                variant="outlined"
                size="small"
                class="mb-3"
                @click="$refs.fileInput.click()"
                :loading="uploadingFile"
              >
                <v-icon start>mdi-upload</v-icon>
                Upload File
              </v-btn>
              
              <div class="space-y-2">
                <v-card
                  v-for="attachment in attachments"
                  :key="attachment.id"
                  variant="outlined"
                  class="pa-3"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3">
                      <v-icon>{{ getFileIcon(attachment.file_type) }}</v-icon>
                      <div>
                        <a :href="attachment.file_url" target="_blank" class="text-primary hover:underline">
                          {{ attachment.filename }}
                        </a>
                        <div class="text-xs text-gray-500">
                          {{ formatFileSize(attachment.size) }} • {{ formatDate(attachment.uploaded_at) }}
                        </div>
                      </div>
                    </div>
                    <v-btn
                      icon
                      size="small"
                      variant="text"
                      color="error"
                      @click="handleDeleteAttachment(attachment.id)"
                    >
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                  </div>
                </v-card>
              </div>
            </div>
            
            <!-- Activity/History -->
            <div>
              <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
                <v-icon>mdi-history</v-icon>
                Activity
              </h3>
              
              <v-timeline density="compact" side="end">
                <v-timeline-item
                  v-for="item in history"
                  :key="item.id"
                  dot-color="primary"
                  size="small"
                >
                  <div class="text-sm">
                    <strong>{{ item.user?.username }}</strong> {{ item.action }}
                    <div class="text-xs text-gray-500">{{ formatDate(item.created_at) }}</div>
                    <pre v-if="item.changes" class="text-xs text-gray-600 mt-1">{{ formatChanges(item.changes) }}</pre>
                  </div>
                </v-timeline-item>
              </v-timeline>
            </div>
          </v-col>
          
          <!-- Sidebar -->
          <v-col cols="12" md="4">
            <div class="space-y-4">
              <!-- Priority -->
              <div>
                <label class="text-sm font-semibold text-gray-600 mb-2 block">Priority</label>
                <v-select
                  v-model="task.priority"
                  :items="['low', 'medium', 'high', 'urgent']"
                  variant="outlined"
                  density="compact"
                  @update:model-value="handleUpdateTask"
                />
              </div>
              
              <!-- Due Date -->
              <div>
                <label class="text-sm font-semibold text-gray-600 mb-2 block">Due Date</label>
                <input
                  type="date"
                  :value="dueDate"
                  @change="handleDueDateChange"
                  class="w-full px-3 py-2 border rounded"
                />
              </div>
              
              <!-- Labels -->
              <div>
                <label class="text-sm font-semibold text-gray-600 mb-2 block">Labels</label>
                <v-select
                  v-model="selectedLabelIds"
                  :items="labelsStore.labels"
                  item-title="name"
                  item-value="id"
                  variant="outlined"
                  density="compact"
                  multiple
                  chips
                  @update:model-value="handleUpdateTask"
                >
                  <template v-slot:chip="{ item, props }">
                    <v-chip v-bind="props" :color="item.raw.color" size="small" />
                  </template>
                </v-select>
              </div>
              
              <!-- Delete Task -->
              <v-btn
                color="error"
                variant="outlined"
                block
                @click="showDeleteDialog = true"
              >
                <v-icon start>mdi-delete</v-icon>
                Delete Task
              </v-btn>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    
    <!-- Delete Confirmation -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title>Delete Task?</v-card-title>
        <v-card-text>
          Are you sure you want to delete this task? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="showDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="handleDeleteTask">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-dialog>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useLabelsStore } from '@/stores/labels'
import { useTasksStore } from '@/stores/tasks'
import { format, parseISO } from 'date-fns'
import { computed, ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  taskId: Number,
  projectId: Number,
})

const emit = defineEmits(['update:modelValue', 'task-updated', 'task-deleted'])

const tasksStore = useTasksStore()
const labelsStore = useLabelsStore()
const authStore = useAuthStore()

const task = ref(null)
const comments = ref([])
const attachments = ref([])
const history = ref([])
const newComment = ref('')
const addingComment = ref(false)
const editingComment = ref(null)
const editingCommentText = ref('')
const uploadingFile = ref(false)
const showDeleteDialog = ref(false)
const fileInput = ref(null)

const selectedLabelIds = computed({
  get: () => task.value?.labels?.map(l => l.id) || [],
  set: (value) => {
    if (task.value) {
      task.value.label_ids = value
    }
  },
})

const dueDate = computed(() => {
  if (!task.value?.due_date) return ''
  try {
    return format(parseISO(task.value.due_date), 'yyyy-MM-dd')
  } catch {
    return ''
  }
})

watch(() => props.taskId, async (newId) => {
  if (newId && props.modelValue) {
    await loadTaskData()
  }
})

watch(() => props.modelValue, async (isOpen) => {
  if (isOpen && props.taskId) {
    await loadTaskData()
  }
})

async function loadTaskData() {
  task.value = await tasksStore.fetchTask(props.taskId)
  comments.value = await tasksStore.fetchComments(props.taskId)
  attachments.value = await tasksStore.fetchAttachments(props.taskId)
  history.value = await tasksStore.fetchTaskHistory(props.taskId)
}

async function handleUpdateTask() {
  if (!task.value) return
  
  await tasksStore.updateTask(task.value.id, {
    title: task.value.title,
    description: task.value.description,
    priority: task.value.priority,
    due_date: task.value.due_date,
    label_ids: selectedLabelIds.value,
  })
  
  emit('task-updated')
}

function handleDueDateChange(event) {
  if (task.value) {
    task.value.due_date = event.target.value ? new Date(event.target.value).toISOString() : null
    handleUpdateTask()
  }
}

async function handleAddComment() {
  addingComment.value = true
  const result = await tasksStore.addComment(props.taskId, newComment.value)
  
  if (result.success) {
    comments.value.push(result.data)
    newComment.value = ''
  }
  addingComment.value = false
}

function isOwnComment(comment) {
  return comment.user_id === authStore.user?.id
}

function startEditComment(comment) {
  editingComment.value = comment
  editingCommentText.value = comment.content
}

function cancelEditComment() {
  editingComment.value = null
  editingCommentText.value = ''
}

async function handleUpdateComment() {
  const result = await tasksStore.updateComment(editingComment.value.id, editingCommentText.value)
  
  if (result.success) {
    const index = comments.value.findIndex(c => c.id === editingComment.value.id)
    if (index !== -1) {
      comments.value[index].content = editingCommentText.value
    }
    cancelEditComment()
  }
}

async function handleDeleteComment(commentId) {
  await tasksStore.deleteComment(commentId)
  comments.value = comments.value.filter(c => c.id !== commentId)
}

async function handleFileUpload(event) {
  const file = event.target.files[0]
  if (!file) return
  
  uploadingFile.value = true
  const result = await tasksStore.uploadAttachment(props.taskId, file)
  
  if (result.success) {
    attachments.value.push(result.data)
  }
  uploadingFile.value = false
  
  // Clear input
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

async function handleDeleteAttachment(attachmentId) {
  await tasksStore.deleteAttachment(attachmentId)
  attachments.value = attachments.value.filter(a => a.id !== attachmentId)
}

async function handleDeleteTask() {
  await tasksStore.deleteTask(props.taskId)
  emit('task-deleted')
  close()
}

function close() {
  emit('update:modelValue', false)
}

function formatDate(date) {
  if (!date) return ''
  try {
    return format(parseISO(date), 'MMM d, yyyy h:mm a')
  } catch {
    return ''
  }
}

function formatFileSize(bytes) {
  if (!bytes) return '0 B'
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
}

function getFileIcon(fileType) {
  if (!fileType) return 'mdi-file'
  if (fileType.includes('image')) return 'mdi-file-image'
  if (fileType.includes('pdf')) return 'mdi-file-pdf'
  if (fileType.includes('word') || fileType.includes('document')) return 'mdi-file-word'
  if (fileType.includes('excel') || fileType.includes('spreadsheet')) return 'mdi-file-excel'
  return 'mdi-file'
}

function formatChanges(changes) {
  if (!changes) return ''
  return JSON.stringify(changes, null, 2)
}
</script>

<style scoped>
.hidden {
  display: none;
}
</style>

