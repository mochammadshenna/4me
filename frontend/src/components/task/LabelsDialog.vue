<template>
  <v-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    max-width="500"
  >
    <v-card>
      <v-card-title class="d-flex align-center pa-4 border-b">
        <span class="text-h5">Manage Labels</span>
        <v-spacer />
        <v-btn icon @click="close">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      
      <v-card-text class="pa-4">
        <!-- Add New Label -->
        <div class="mb-4">
          <v-text-field
            v-model="newLabel.name"
            label="Label Name"
            variant="outlined"
            density="compact"
            placeholder="Enter label name"
          />
          
          <div class="mb-3">
            <label class="text-sm text-gray-600 mb-2 block">Color</label>
            <div class="flex gap-2">
              <div
                v-for="color in colors"
                :key="color"
                class="w-8 h-8 rounded-full cursor-pointer border-2 transition-all"
                :class="newLabel.color === color ? 'border-gray-800 scale-110' : 'border-transparent'"
                :style="{ backgroundColor: color }"
                @click="newLabel.color = color"
              />
            </div>
          </div>
          
          <v-btn
            color="primary"
            block
            @click="handleAddLabel"
            :disabled="!newLabel.name.trim()"
            :loading="adding"
          >
            <v-icon start>mdi-plus</v-icon>
            Add Label
          </v-btn>
        </div>
        
        <v-divider class="my-4" />
        
        <!-- Existing Labels -->
        <div class="space-y-2">
          <v-card
            v-for="label in labels"
            :key="label.id"
            variant="outlined"
            class="pa-3"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3 flex-1">
                <div
                  class="w-6 h-6 rounded-full flex-shrink-0"
                  :style="{ backgroundColor: label.color }"
                />
                
                <div v-if="editingLabel?.id === label.id" class="flex-1">
                  <v-text-field
                    v-model="editingLabelData.name"
                    density="compact"
                    variant="outlined"
                    hide-details
                  />
                </div>
                <span v-else class="font-medium">{{ label.name }}</span>
              </div>
              
              <div class="flex gap-2">
                <v-btn
                  v-if="editingLabel?.id === label.id"
                  icon
                  size="small"
                  color="primary"
                  @click="handleUpdateLabel"
                >
                  <v-icon>mdi-check</v-icon>
                </v-btn>
                <v-btn
                  v-if="editingLabel?.id === label.id"
                  icon
                  size="small"
                  @click="cancelEdit"
                >
                  <v-icon>mdi-close</v-icon>
                </v-btn>
                
                <v-btn
                  v-if="!editingLabel"
                  icon
                  size="small"
                  @click="startEdit(label)"
                >
                  <v-icon>mdi-pencil</v-icon>
                </v-btn>
                <v-btn
                  v-if="!editingLabel"
                  icon
                  size="small"
                  color="error"
                  @click="handleDeleteLabel(label.id)"
                >
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </div>
            </div>
            
            <div v-if="editingLabel?.id === label.id" class="flex gap-2 mt-3">
              <div
                v-for="color in colors"
                :key="color"
                class="w-6 h-6 rounded-full cursor-pointer border-2 transition-all"
                :class="editingLabelData.color === color ? 'border-gray-800' : 'border-transparent'"
                :style="{ backgroundColor: color }"
                @click="editingLabelData.color = color"
              />
            </div>
          </v-card>
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { useLabelsStore } from '@/stores/labels'
import { computed, ref } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  projectId: Number,
})

const emit = defineEmits(['update:modelValue'])

const labelsStore = useLabelsStore()

const labels = computed(() => labelsStore.labels)

const newLabel = ref({
  name: '',
  color: '#3B82F6',
})

const editingLabel = ref(null)
const editingLabelData = ref({
  name: '',
  color: '',
})

const adding = ref(false)

const colors = [
  '#3B82F6', '#10B981', '#F59E0B', '#EF4444',
  '#8B5CF6', '#EC4899', '#06B6D4', '#64748B',
  '#6366F1', '#14B8A6', '#F97316', '#84CC16',
]

async function handleAddLabel() {
  adding.value = true
  const result = await labelsStore.createLabel(props.projectId, newLabel.value)
  
  if (result.success) {
    newLabel.value = {
      name: '',
      color: '#3B82F6',
    }
  }
  adding.value = false
}

function startEdit(label) {
  editingLabel.value = label
  editingLabelData.value = {
    name: label.name,
    color: label.color,
  }
}

function cancelEdit() {
  editingLabel.value = null
  editingLabelData.value = {
    name: '',
    color: '',
  }
}

async function handleUpdateLabel() {
  await labelsStore.updateLabel(editingLabel.value.id, editingLabelData.value)
  cancelEdit()
}

async function handleDeleteLabel(labelId) {
  await labelsStore.deleteLabel(labelId)
}

function close() {
  emit('update:modelValue', false)
}
</script>

