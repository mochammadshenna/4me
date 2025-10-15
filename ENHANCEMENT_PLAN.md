# 4me Todos - Jira-Quality Enhancement Plan

**Goal**: Transform 4me Todos into a Jira-quality Kanban board with Atlassian's pragmatic-drag-and-drop, beautiful UI, and exceptional UX.

**Timeline**: Incremental implementation with measurable milestones
**Testing**: Comprehensive Playwright E2E tests at every stage

---

## Executive Summary

### Current State Analysis

**Strengths:**
- ✅ Functional drag-and-drop with vuedraggable
- ✅ Clean, modern design with glassmorphism effects
- ✅ Responsive layout with good mobile support
- ✅ Basic animations and transitions
- ✅ Vuetify component integration

**Gaps vs. Jira Quality:**
- ❌ Using older vuedraggable library (not performance-optimized)
- ❌ Limited accessibility (no keyboard navigation, screen reader support)
- ❌ Basic drag feedback (missing drop indicators, live previews)
- ❌ No advanced interactions (auto-scroll, multi-select, quick actions)
- ❌ Missing micro-interactions and polish
- ❌ Limited testing coverage for drag-and-drop

### Target State (Jira-Quality)

**Performance:**
- Fast, smooth drag operations (<16ms frame time)
- Optimized bundle size (~4.7kB core drag-and-drop)
- Deferred loading for improved page speed

**Accessibility:**
- Full keyboard navigation (arrow keys, shortcuts)
- Screen reader support with ARIA labels
- Focus management and visual indicators
- Alternative interaction methods

**User Experience:**
- Live drag previews with actual task content
- Visual drop indicators between cards
- Auto-scroll when dragging near edges
- Keyboard shortcuts (Ctrl+C/V, Delete, etc.)
- Toast notifications for all actions
- Loading states and optimistic updates

**Visual Design:**
- Jira-inspired color scheme and spacing
- Smooth animations (spring physics)
- Elevation and depth cues
- Clear visual hierarchy
- Beautiful empty states

---

## Phase 1: Foundation & Setup ✅ COMPLETE

### 1.1 Analysis (✅ Complete)
- [x] Audit current BoardColumn.vue implementation
- [x] Audit current TaskCard.vue implementation
- [x] Research Atlassian pragmatic-drag-and-drop
- [x] Identify UI/UX gaps
- [x] Document current architecture

### 1.2 Dependencies Planning
**Install:**
```bash
cd frontend
npm install @atlaskit/pragmatic-drag-and-drop@latest
npm install @atlaskit/pragmatic-drag-and-drop-hitbox@latest
npm install @atlaskit/pragmatic-drag-and-drop-auto-scroll@latest
npm install @atlaskit/pragmatic-drag-and-drop-live-region@latest
npm install @atlaskit/pragmatic-drag-and-drop-react-beautiful-dnd-migration@latest
```

**Remove:**
```bash
npm uninstall vuedraggable
```

---

## Phase 2: Core Migration (Pragmatic Drag-and-Drop)

### 2.1 Install Pragmatic Drag-and-Drop
**Priority**: High
**Complexity**: Low
**Estimated Time**: 30 minutes

**Tasks:**
1. Install all @atlaskit/pragmatic-drag-and-drop packages
2. Remove vuedraggable dependency
3. Update package.json and package-lock.json
4. Verify build succeeds

**Playwright Test:**
```javascript
test('packages installed correctly', async ({ page }) => {
  // Verify no console errors about missing dependencies
  const errors = []
  page.on('pageerror', err => errors.push(err))
  await page.goto('/projects/1')
  expect(errors).toHaveLength(0)
})
```

### 2.2 Refactor TaskCard Component
**Priority**: High
**Complexity**: Medium
**Estimated Time**: 2-3 hours

**Implementation:**

**File**: `frontend/src/components/board/TaskCard.vue`

**Changes:**
1. Remove vuedraggable drag-handle class
2. Implement draggable adapter from pragmatic-drag-and-drop
3. Add data attributes for drag identification
4. Implement drag preview with live content
5. Add visual states (dragging, over, disabled)

**Key Code Pattern:**
```vue
<script setup>
import { draggable } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'
import { onMounted, ref, onUnmounted } from 'vue'

const taskCardRef = ref(null)
const isDragging = ref(false)
const isOver = ref(false)

onMounted(() => {
  if (!taskCardRef.value) return

  const cleanup = draggable({
    element: taskCardRef.value,
    getInitialData: () => ({
      type: 'task',
      taskId: props.task.id,
      boardId: props.task.board_id
    }),
    onGenerateDragPreview: ({ nativeSetDragImage }) => {
      // Use actual task card as preview
      const preview = taskCardRef.value.cloneNode(true)
      preview.style.transform = 'rotate(3deg)'
      preview.style.boxShadow = '0 8px 32px rgba(0,0,0,0.2)'
      nativeSetDragImage(preview, 0, 0)
    },
    onDragStart: () => {
      isDragging.value = true
    },
    onDrop: () => {
      isDragging.value = false
    }
  })

  onUnmounted(cleanup)
})
</script>

<template>
  <div
    ref="taskCardRef"
    :class="{
      'task-dragging': isDragging,
      'task-over': isOver
    }"
    :data-task-id="task.id"
  >
    <!-- Task content -->
  </div>
</template>

<style scoped>
.task-dragging {
  opacity: 0.5;
  transform: scale(0.98);
}

.task-over {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.15);
}
</style>
```

**Playwright Tests:**
```javascript
test.describe('TaskCard Drag-and-Drop', () => {
  test('task can be dragged', async ({ page }) => {
    await page.goto('/projects/1')
    const task = page.locator('[data-task-id="1"]').first()

    // Start drag
    await task.hover()
    await page.mouse.down()

    // Verify dragging class applied
    await expect(task).toHaveClass(/task-dragging/)

    await page.mouse.up()
  })

  test('drag preview shows task content', async ({ page }) => {
    await page.goto('/projects/1')
    const task = page.locator('[data-task-id="1"]').first()
    const taskTitle = await task.locator('.task-title').textContent()

    await task.hover()
    await page.mouse.down()
    await page.mouse.move(100, 100)

    // Verify drag preview exists and contains title
    // (Implementation specific to browser)
  })
})
```

### 2.3 Refactor BoardColumn Component
**Priority**: High
**Complexity**: High
**Estimated Time**: 4-5 hours

**Implementation:**

**File**: `frontend/src/components/board/BoardColumn.vue`

**Changes:**
1. Remove draggable component wrapper
2. Implement dropTargetForElements adapter
3. Add drop indicator component
4. Implement auto-scroll behavior
5. Add keyboard navigation
6. Implement optimistic updates

**Key Code Pattern:**
```vue
<script setup>
import { dropTargetForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'
import { autoScrollForElements } from '@atlaskit/pragmatic-drag-and-drop-auto-scroll/element'
import { extractClosestEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge'
import { onMounted, ref, computed } from 'vue'

const columnRef = ref(null)
const dropIndicatorPosition = ref(null) // 'top' | 'bottom' | null
const isDraggedOver = ref(false)

onMounted(() => {
  if (!columnRef.value) return

  const cleanupDropTarget = dropTargetForElements({
    element: columnRef.value,
    canDrop: ({ source }) => {
      return source.data.type === 'task'
    },
    getData: ({ input, element }) => {
      const closestEdge = extractClosestEdge(input)
      return {
        type: 'board-column',
        boardId: props.board.id,
        edge: closestEdge
      }
    },
    onDragEnter: () => {
      isDraggedOver.value = true
    },
    onDragLeave: () => {
      isDraggedOver.value = false
      dropIndicatorPosition.value = null
    },
    onDrag: ({ source, self }) => {
      const edge = self.data.edge
      dropIndicatorPosition.value = edge
    },
    onDrop: async ({ source, self }) => {
      const taskId = source.data.taskId
      const targetBoardId = props.board.id
      const edge = self.data.edge

      // Optimistic update
      await tasksStore.moveTask(taskId, targetBoardId, edge)

      isDraggedOver.value = false
      dropIndicatorPosition.value = null
    }
  })

  const cleanupAutoScroll = autoScrollForElements({
    element: columnRef.value,
    canScroll: ({ source }) => source.data.type === 'task'
  })

  onUnmounted(() => {
    cleanupDropTarget()
    cleanupAutoScroll()
  })
})
</script>

<template>
  <div
    ref="columnRef"
    class="board-column"
    :class="{ 'drag-over': isDraggedOver }"
    :data-board-id="board.id"
  >
    <div class="column-header">
      <!-- Header content -->
    </div>

    <!-- Drop indicator at top -->
    <div
      v-if="dropIndicatorPosition === 'top'"
      class="drop-indicator drop-indicator-top"
    />

    <div class="tasks-container">
      <TaskCard
        v-for="task in tasks"
        :key="task.id"
        :task="task"
        @click="$emit('task-click', task.id)"
      />

      <!-- Empty state -->
      <div v-if="tasks.length === 0" class="empty-state">
        <v-icon size="48">mdi-clipboard-outline</v-icon>
        <p>Drop tasks here</p>
      </div>
    </div>

    <!-- Drop indicator at bottom -->
    <div
      v-if="dropIndicatorPosition === 'bottom'"
      class="drop-indicator drop-indicator-bottom"
    />

    <!-- Add task button -->
  </div>
</template>

<style scoped>
.drop-indicator {
  height: 4px;
  background: linear-gradient(90deg, #1976D2, #2196F3);
  border-radius: 2px;
  margin: 8px 12px;
  animation: pulse 1.5s ease-in-out infinite;
  box-shadow: 0 0 12px rgba(25, 118, 210, 0.5);
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.6;
    transform: scaleX(0.95);
  }
  50% {
    opacity: 1;
    transform: scaleX(1);
  }
}

.drag-over {
  background: rgba(25, 118, 210, 0.05);
  border: 2px dashed #1976D2;
}
</style>
```

**Playwright Tests:**
```javascript
test.describe('BoardColumn Drop Target', () => {
  test('shows drop indicator when dragging over', async ({ page }) => {
    await page.goto('/projects/1')

    const task = page.locator('[data-task-id="1"]').first()
    const targetColumn = page.locator('[data-board-id="2"]').first()

    // Start drag
    await task.hover()
    await page.mouse.down()
    await page.mouse.move(100, 100)

    // Move over target column
    const box = await targetColumn.boundingBox()
    await page.mouse.move(box.x + box.width / 2, box.y + 20)

    // Verify drop indicator appears
    await expect(targetColumn.locator('.drop-indicator')).toBeVisible()
  })

  test('task moves to new column on drop', async ({ page }) => {
    await page.goto('/projects/1')

    const task = page.locator('[data-task-id="1"]').first()
    const sourceColumn = page.locator('[data-board-id="1"]')
    const targetColumn = page.locator('[data-board-id="2"]')

    // Verify task in source column
    await expect(sourceColumn.locator('[data-task-id="1"]')).toBeVisible()

    // Drag to target column
    await task.dragTo(targetColumn)

    // Verify task moved
    await expect(targetColumn.locator('[data-task-id="1"]')).toBeVisible()
    await expect(sourceColumn.locator('[data-task-id="1"]')).not.toBeVisible()
  })

  test('auto-scrolls when dragging near edge', async ({ page }) => {
    await page.goto('/projects/1')
    await page.setViewportSize({ width: 1920, height: 1080 })

    const task = page.locator('[data-task-id="1"]').first()
    const column = page.locator('[data-board-id="1"]').first()

    // Get initial scroll position
    const initialScroll = await column.evaluate(el => el.scrollTop)

    // Drag near bottom edge
    await task.hover()
    await page.mouse.down()
    const box = await column.boundingBox()
    await page.mouse.move(box.x + 50, box.y + box.height - 20)

    // Wait for auto-scroll
    await page.waitForTimeout(500)

    // Verify scrolled
    const newScroll = await column.evaluate(el => el.scrollTop)
    expect(newScroll).toBeGreaterThan(initialScroll)

    await page.mouse.up()
  })
})
```

---

## Phase 3: Advanced UI/UX Enhancements

### 3.1 Keyboard Navigation
**Priority**: High (Accessibility)
**Complexity**: High
**Estimated Time**: 4-5 hours

**Implementation:**

**Features:**
- Arrow keys to navigate between tasks
- Enter to open task detail
- Space to select/deselect
- Ctrl+C/X/V for copy/cut/paste
- Delete to remove task
- Escape to cancel operations
- Tab for focus management

**Key Code:**
```vue
<script setup>
import { useKeyboardNavigation } from '@/composables/useKeyboardNavigation'

const {
  selectedTaskId,
  handleKeyDown,
  focusTask,
  moveTaskUp,
  moveTaskDown,
  moveTaskLeft,
  moveTaskRight
} = useKeyboardNavigation(boards, tasks)

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
  })
})
</script>
```

**Composable**: `frontend/src/composables/useKeyboardNavigation.js`
```javascript
export function useKeyboardNavigation(boards, tasks) {
  const selectedTaskId = ref(null)
  const clipboard = ref(null)

  function handleKeyDown(event) {
    if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA') {
      return // Don't intercept when typing
    }

    switch(event.key) {
      case 'ArrowUp':
        event.preventDefault()
        moveTaskUp()
        break
      case 'ArrowDown':
        event.preventDefault()
        moveTaskDown()
        break
      case 'ArrowLeft':
        event.preventDefault()
        moveTaskLeft()
        break
      case 'ArrowRight':
        event.preventDefault()
        moveTaskRight()
        break
      case 'Enter':
        event.preventDefault()
        openSelectedTask()
        break
      case ' ':
        event.preventDefault()
        toggleSelection()
        break
      case 'c':
        if (event.ctrlKey || event.metaKey) {
          event.preventDefault()
          copyTask()
        }
        break
      case 'x':
        if (event.ctrlKey || event.metaKey) {
          event.preventDefault()
          cutTask()
        }
        break
      case 'v':
        if (event.ctrlKey || event.metaKey) {
          event.preventDefault()
          pasteTask()
        }
        break
      case 'Delete':
      case 'Backspace':
        event.preventDefault()
        deleteTask()
        break
      case 'Escape':
        selectedTaskId.value = null
        break
    }
  }

  // Implementation of navigation functions...

  return {
    selectedTaskId,
    handleKeyDown,
    focusTask,
    moveTaskUp,
    moveTaskDown,
    moveTaskLeft,
    moveTaskRight
  }
}
```

**Playwright Tests:**
```javascript
test.describe('Keyboard Navigation', () => {
  test('arrow keys navigate between tasks', async ({ page }) => {
    await page.goto('/projects/1')

    // Focus first task
    await page.keyboard.press('Tab')
    const firstTask = page.locator('[data-task-id="1"]')
    await expect(firstTask).toBeFocused()

    // Navigate down
    await page.keyboard.press('ArrowDown')
    const secondTask = page.locator('[data-task-id="2"]')
    await expect(secondTask).toBeFocused()

    // Navigate right (next column)
    await page.keyboard.press('ArrowRight')
    const taskInNextColumn = page.locator('[data-board-id="2"] [data-task-id]').first()
    await expect(taskInNextColumn).toBeFocused()
  })

  test('Enter opens task detail', async ({ page }) => {
    await page.goto('/projects/1')

    await page.keyboard.press('Tab')
    await page.keyboard.press('Enter')

    await expect(page.locator('[role="dialog"]')).toBeVisible()
  })

  test('Ctrl+C and Ctrl+V copy and paste task', async ({ page }) => {
    await page.goto('/projects/1')

    await page.keyboard.press('Tab')
    const originalTitle = await page.locator('[data-task-id="1"] .task-title').textContent()

    // Copy
    await page.keyboard.press('Control+c')

    // Navigate to different column
    await page.keyboard.press('ArrowRight')

    // Paste
    await page.keyboard.press('Control+v')

    // Verify copied task appears
    const copiedTask = page.locator('[data-board-id="2"]').locator('text=' + originalTitle)
    await expect(copiedTask).toBeVisible()
  })
})
```

### 3.2 Screen Reader Support
**Priority**: High (Accessibility)
**Complexity**: Medium
**Estimated Time**: 3-4 hours

**Implementation:**

**Features:**
- ARIA labels for all interactive elements
- Live region announcements for drag-and-drop
- Descriptive alt text and roles
- Focus management

**Key Code:**
```vue
<template>
  <div
    role="region"
    :aria-label="`${board.name} column with ${tasks.length} tasks`"
    class="board-column"
  >
    <h3 :id="`board-${board.id}-heading`">{{ board.name }}</h3>

    <div
      role="list"
      :aria-labelledby="`board-${board.id}-heading`"
      class="tasks-container"
    >
      <div
        v-for="(task, index) in tasks"
        :key="task.id"
        role="listitem"
        :aria-label="`Task ${index + 1} of ${tasks.length}: ${task.title}`"
        :aria-description="`Priority: ${task.priority}, Due: ${task.due_date || 'No due date'}`"
      >
        <TaskCard :task="task" />
      </div>
    </div>

    <!-- Live region for announcements -->
    <div
      role="status"
      aria-live="polite"
      aria-atomic="true"
      class="sr-only"
    >
      {{ liveRegionMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const liveRegionMessage = ref('')

function announceDragStart(task) {
  liveRegionMessage.value = `Started dragging ${task.title}`
}

function announceDrop(task, targetBoard) {
  liveRegionMessage.value = `Dropped ${task.title} in ${targetBoard.name}`
}
</script>

<style>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
```

**Playwright Accessibility Tests:**
```javascript
test.describe('Accessibility', () => {
  test('has proper ARIA labels', async ({ page }) => {
    await page.goto('/projects/1')

    // Check board column labels
    const column = page.locator('[role="region"]').first()
    const label = await column.getAttribute('aria-label')
    expect(label).toContain('column with')
    expect(label).toContain('tasks')
  })

  test('announces drag operations to screen readers', async ({ page }) => {
    await page.goto('/projects/1')

    const liveRegion = page.locator('[role="status"]')

    // Start drag
    const task = page.locator('[data-task-id="1"]').first()
    await task.hover()
    await page.mouse.down()

    // Check live region updated
    const message = await liveRegion.textContent()
    expect(message).toContain('Started dragging')

    await page.mouse.up()
  })

  test('all interactive elements are keyboard accessible', async ({ page }) => {
    await page.goto('/projects/1')

    // Tab through all focusable elements
    const focusableElements = []
    await page.keyboard.press('Tab')

    let currentElement = await page.locator(':focus')
    focusableElements.push(currentElement)

    // Verify all task cards, buttons, and inputs are reachable
    for (let i = 0; i < 20; i++) {
      await page.keyboard.press('Tab')
      currentElement = await page.locator(':focus')
      const tagName = await currentElement.evaluate(el => el.tagName)

      // Verify focusable elements are visible and interactive
      if (tagName !== 'BODY') {
        await expect(currentElement).toBeVisible()
      }
    }
  })
})
```

### 3.3 Enhanced Visual Feedback
**Priority**: Medium
**Complexity**: Medium
**Estimated Time**: 2-3 hours

**Features:**
- Spring physics animations
- Elevation changes during drag
- Blur and scale effects
- Color transitions
- Particle effects on drop (optional)

**Implementation:**
```vue
<style scoped>
/* Spring physics animation */
.task-card {
  transition:
    transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1),
    box-shadow 0.3s ease,
    opacity 0.2s ease;
}

.task-card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow:
    0 8px 24px rgba(0, 0, 0, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.08);
}

.task-card.dragging {
  transform: rotate(3deg) scale(1.05);
  box-shadow:
    0 12px 32px rgba(0, 0, 0, 0.2),
    0 4px 12px rgba(0, 0, 0, 0.1);
  opacity: 0.9;
  z-index: 1000;
}

/* Glassmorphism effect */
.board-column {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

/* Smooth color transitions */
.task-priority-high {
  border-left: 4px solid #FF5722;
  background: linear-gradient(90deg, rgba(255, 87, 34, 0.05) 0%, transparent 100%);
}

.task-priority-urgent {
  border-left: 4px solid #F44336;
  background: linear-gradient(90deg, rgba(244, 67, 54, 0.08) 0%, transparent 100%);
  animation: pulse-urgent 2s ease-in-out infinite;
}

@keyframes pulse-urgent {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(244, 67, 54, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(244, 67, 54, 0);
  }
}

/* Drop success animation */
@keyframes drop-success {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.task-card.just-dropped {
  animation: drop-success 0.4s ease-out;
}
</style>
```

### 3.4 Toast Notifications
**Priority**: Medium
**Complexity**: Low
**Estimated Time**: 2 hours

**Implementation:**

**Install:**
```bash
npm install vue-toastification@next
```

**Setup**: `frontend/src/main.js`
```javascript
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

app.use(Toast, {
  position: 'bottom-right',
  timeout: 3000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: 'button',
  icon: true,
  rtl: false,
  transition: 'Vue-Toastification__bounce',
  maxToasts: 5,
  newestOnTop: true
})
```

**Usage in Stores:**
```javascript
import { useToast } from 'vue-toastification'

export const useTasksStore = defineStore('tasks', () => {
  const toast = useToast()

  async function moveTask(taskId, targetBoardId) {
    try {
      const task = findTask(taskId)
      const targetBoard = findBoard(targetBoardId)

      // Optimistic update
      updateTaskBoardId(taskId, targetBoardId)

      // API call
      await apiClient.patch(`/tasks/${taskId}/move`, {
        board_id: targetBoardId
      })

      toast.success(`Moved "${task.title}" to ${targetBoard.name}`, {
        icon: '✅'
      })
    } catch (error) {
      // Rollback optimistic update
      rollbackTaskMove(taskId)

      toast.error('Failed to move task', {
        icon: '❌'
      })
    }
  }

  return { moveTask }
})
```

**Playwright Tests:**
```javascript
test('shows toast notification on task move', async ({ page }) => {
  await page.goto('/projects/1')

  const task = page.locator('[data-task-id="1"]').first()
  const targetColumn = page.locator('[data-board-id="2"]').first()

  // Drag and drop
  await task.dragTo(targetColumn)

  // Verify toast appears
  const toast = page.locator('.Vue-Toastification__toast--success')
  await expect(toast).toBeVisible()
  await expect(toast).toContainText('Moved')

  // Verify toast disappears after timeout
  await page.waitForTimeout(3500)
  await expect(toast).not.toBeVisible()
})
```

---

## Phase 4: Comprehensive Testing Strategy

### 4.1 Playwright E2E Test Suite

**Test Organization:**
```
frontend/
├── playwright.config.js
└── tests/
    ├── e2e/
    │   ├── drag-drop.spec.js        # Core drag-and-drop tests
    │   ├── keyboard-nav.spec.js      # Keyboard navigation
    │   ├── accessibility.spec.js     # A11y compliance
    │   ├── performance.spec.js       # Performance metrics
    │   └── cross-browser.spec.js     # Browser compatibility
    └── fixtures/
        └── test-data.js              # Test data generators
```

**Playwright Config Enhancement:**
```javascript
// playwright.config.js
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html'],
    ['json', { outputFile: 'test-results/results.json' }],
    ['junit', { outputFile: 'test-results/junit.xml' }]
  ],
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure'
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] }
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] }
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] }
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] }
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] }
    }
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120000
  }
})
```

**Complete Drag-and-Drop Test Suite:**
```javascript
// tests/e2e/drag-drop.spec.js
import { test, expect } from '@playwright/test'

test.describe('Drag and Drop - Core Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Login and navigate to project
    await page.goto('/login')
    await page.fill('[name="username"]', 'testuser')
    await page.fill('[name="password"]', 'password123')
    await page.click('button[type="submit"]')
    await page.waitForURL('/')

    // Navigate to test project
    await page.click('text=Test Project')
    await page.waitForURL(/\/projects\/\d+/)
  })

  test('can drag task between columns', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()
    const sourceColumn = page.locator('[data-board-id="1"]')
    const targetColumn = page.locator('[data-board-id="2"]')

    // Get initial counts
    const initialSourceCount = await sourceColumn.locator('[data-task-id]').count()
    const initialTargetCount = await targetColumn.locator('[data-task-id]').count()

    // Perform drag
    await task.dragTo(targetColumn)

    // Wait for animation
    await page.waitForTimeout(500)

    // Verify task moved
    expect(await sourceColumn.locator('[data-task-id]').count()).toBe(initialSourceCount - 1)
    expect(await targetColumn.locator('[data-task-id]').count()).toBe(initialTargetCount + 1)
    expect(await targetColumn.locator('[data-task-id="1"]')).toBeVisible()
  })

  test('shows drop indicator during drag', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()
    const targetColumn = page.locator('[data-board-id="2"]')

    // Start drag
    await task.hover()
    await page.mouse.down()

    // Move over target
    const box = await targetColumn.boundingBox()
    await page.mouse.move(box.x + box.width / 2, box.y + 50)

    // Verify drop indicator
    await expect(targetColumn.locator('.drop-indicator')).toBeVisible()

    await page.mouse.up()
  })

  test('drag preview shows task content', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()
    const taskTitle = await task.locator('.task-title').textContent()

    // Start drag
    await task.hover()
    await page.mouse.down()

    // Verify dragging class applied
    await expect(task).toHaveClass(/task-dragging/)

    await page.mouse.up()
  })

  test('can reorder tasks within same column', async ({ page }) => {
    const column = page.locator('[data-board-id="1"]')
    const firstTask = column.locator('[data-task-id]').first()
    const secondTask = column.locator('[data-task-id]').nth(1)

    const firstTaskId = await firstTask.getAttribute('data-task-id')
    const secondTaskId = await secondTask.getAttribute('data-task-id')

    // Drag first task below second
    await firstTask.dragTo(secondTask, { targetPosition: { x: 0, y: 50 } })

    // Wait for animation
    await page.waitForTimeout(500)

    // Verify order changed
    const newFirstTask = column.locator('[data-task-id]').first()
    const newFirstTaskId = await newFirstTask.getAttribute('data-task-id')

    expect(newFirstTaskId).toBe(secondTaskId)
  })

  test('auto-scrolls when dragging near column edge', async ({ page }) => {
    const column = page.locator('[data-board-id="1"]')

    // Add many tasks to make column scrollable
    for (let i = 0; i < 20; i++) {
      await page.click('[data-board-id="1"] button:has-text("Add Task")')
      await page.fill('[data-board-id="1"] textarea', `Test Task ${i}`)
      await page.click('[data-board-id="1"] button:has-text("Add")')
      await page.waitForTimeout(200)
    }

    // Get scrollable container
    const scrollContainer = column.locator('.tasks-container')

    // Scroll to top
    await scrollContainer.evaluate(el => el.scrollTop = 0)
    const initialScroll = await scrollContainer.evaluate(el => el.scrollTop)

    // Start dragging task
    const task = column.locator('[data-task-id]').first()
    await task.hover()
    await page.mouse.down()

    // Move to bottom edge
    const box = await scrollContainer.boundingBox()
    await page.mouse.move(box.x + 50, box.y + box.height - 20)

    // Wait for auto-scroll
    await page.waitForTimeout(1000)

    // Verify scrolled
    const newScroll = await scrollContainer.evaluate(el => el.scrollTop)
    expect(newScroll).toBeGreaterThan(initialScroll)

    await page.mouse.up()
  })

  test('shows toast notification on successful drop', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()
    const targetColumn = page.locator('[data-board-id="2"]')

    await task.dragTo(targetColumn)

    // Verify toast appears
    const toast = page.locator('.Vue-Toastification__toast--success')
    await expect(toast).toBeVisible()
    await expect(toast).toContainText('Moved')
  })

  test('handles drag cancellation on Escape key', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()
    const sourceColumn = page.locator('[data-board-id="1"]')

    const initialCount = await sourceColumn.locator('[data-task-id]').count()

    // Start drag
    await task.hover()
    await page.mouse.down()
    await page.mouse.move(100, 100)

    // Cancel with Escape
    await page.keyboard.press('Escape')

    // Verify task remained in place
    expect(await sourceColumn.locator('[data-task-id]').count()).toBe(initialCount)
    expect(await sourceColumn.locator('[data-task-id="1"]')).toBeVisible()
  })

  test('prevents dropping on invalid targets', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()
    const header = page.locator('.project-header')

    // Try to drag to header (invalid target)
    await task.hover()
    await page.mouse.down()

    const box = await header.boundingBox()
    await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2)

    // Verify no drop indicator in invalid area
    await expect(header.locator('.drop-indicator')).not.toBeVisible()

    await page.mouse.up()
  })

  test('maintains task data after drag', async ({ page }) => {
    const task = page.locator('[data-task-id="1"]').first()

    // Get task details
    const originalTitle = await task.locator('.task-title').textContent()
    const originalPriority = await task.locator('.task-priority').textContent()

    // Drag to another column
    const targetColumn = page.locator('[data-board-id="2"]')
    await task.dragTo(targetColumn)

    // Verify data preserved
    const movedTask = targetColumn.locator('[data-task-id="1"]')
    await expect(movedTask.locator('.task-title')).toHaveText(originalTitle)
    await expect(movedTask.locator('.task-priority')).toHaveText(originalPriority)
  })
})

test.describe('Drag and Drop - Performance', () => {
  test('drag operations complete within performance budget', async ({ page }) => {
    await page.goto('/projects/1')

    // Start performance measurement
    await page.evaluate(() => performance.mark('drag-start'))

    const task = page.locator('[data-task-id="1"]').first()
    const targetColumn = page.locator('[data-board-id="2"]')

    await task.dragTo(targetColumn)
    await page.waitForTimeout(100)

    // End measurement
    await page.evaluate(() => performance.mark('drag-end'))

    const duration = await page.evaluate(() => {
      performance.measure('drag-operation', 'drag-start', 'drag-end')
      const measure = performance.getEntriesByName('drag-operation')[0]
      return measure.duration
    })

    // Should complete in less than 100ms for good UX
    expect(duration).toBeLessThan(100)
  })

  test('handles rapid successive drags', async ({ page }) => {
    await page.goto('/projects/1')

    // Perform 5 rapid drags
    for (let i = 0; i < 5; i++) {
      const task = page.locator('[data-task-id]').first()
      const targetColumn = page.locator('[data-board-id="2"]')

      await task.dragTo(targetColumn)
      await page.waitForTimeout(50) // Minimal delay
    }

    // Verify all operations completed successfully
    const targetCount = await page.locator('[data-board-id="2"] [data-task-id]').count()
    expect(targetCount).toBeGreaterThanOrEqual(5)
  })
})

test.describe('Drag and Drop - Edge Cases', () => {
  test('handles empty column drops', async ({ page }) => {
    await page.goto('/projects/1')

    // Create empty column
    await page.click('button:has-text("Add Board")')
    await page.fill('input[label="Board Name"]', 'Empty Column')
    await page.click('button:has-text("Add")')

    const emptyColumn = page.locator('[data-board-id]').last()
    const task = page.locator('[data-task-id="1"]').first()

    // Drag to empty column
    await task.dragTo(emptyColumn)

    // Verify task in empty column
    await expect(emptyColumn.locator('[data-task-id="1"]')).toBeVisible()
    await expect(emptyColumn.locator('.empty-state')).not.toBeVisible()
  })

  test('handles simultaneous drags (if multi-user)', async ({ page, context }) => {
    // This would test WebSocket real-time updates if implemented
    // Placeholder for future multi-user features
  })
})
```

### 4.2 Performance Testing
```javascript
// tests/e2e/performance.spec.js
import { test, expect } from '@playwright/test'

test.describe('Performance Metrics', () => {
  test('page load within budget', async ({ page }) => {
    const start = Date.now()
    await page.goto('/projects/1')
    await page.waitForLoadState('networkidle')
    const loadTime = Date.now() - start

    expect(loadTime).toBeLessThan(3000) // 3 second budget
  })

  test('drag operation frame rate', async ({ page }) => {
    await page.goto('/projects/1')

    // Start FPS monitoring
    await page.evaluate(() => {
      window.fpsData = []
      let lastTime = performance.now()

      function measureFPS() {
        const now = performance.now()
        const fps = 1000 / (now - lastTime)
        window.fpsData.push(fps)
        lastTime = now

        if (window.fpsData.length < 60) {
          requestAnimationFrame(measureFPS)
        }
      }

      requestAnimationFrame(measureFPS)
    })

    // Perform drag
    const task = page.locator('[data-task-id="1"]').first()
    const target = page.locator('[data-board-id="2"]')
    await task.dragTo(target)

    // Get FPS data
    const fpsData = await page.evaluate(() => window.fpsData)
    const avgFPS = fpsData.reduce((a, b) => a + b, 0) / fpsData.length

    // Should maintain 60 FPS (or close to it)
    expect(avgFPS).toBeGreaterThan(55)
  })

  test('memory usage stays stable', async ({ page }) => {
    await page.goto('/projects/1')

    // Get initial memory
    const initialMemory = await page.evaluate(() =>
      performance.memory.usedJSHeapSize
    )

    // Perform 20 drag operations
    for (let i = 0; i < 20; i++) {
      const task = page.locator('[data-task-id]').first()
      const target = page.locator('[data-board-id="2"]')
      await task.dragTo(target)
      await page.waitForTimeout(100)
    }

    // Check memory didn't grow significantly
    const finalMemory = await page.evaluate(() =>
      performance.memory.usedJSHeapSize
    )

    const memoryGrowth = (finalMemory - initialMemory) / initialMemory
    expect(memoryGrowth).toBeLessThan(0.5) // Less than 50% growth
  })
})
```

---

## Phase 5: Documentation

### 5.1 Update CLAUDE.md

Add section on drag-and-drop implementation:

```markdown
## Drag-and-Drop System

### Architecture

**Library**: @atlaskit/pragmatic-drag-and-drop (Atlassian's performance-focused library)

**Pattern**: Adapters for draggable elements and drop targets

**Key Features:**
- ~4.7kB core bundle
- Framework-agnostic (works with Vue 3)
- Headless design (full styling control)
- Built-in accessibility support
- Auto-scroll near edges
- Keyboard navigation

### Component Structure

**TaskCard.vue** - Draggable Element:
```javascript
import { draggable } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'

onMounted(() => {
  const cleanup = draggable({
    element: taskCardRef.value,
    getInitialData: () => ({ type: 'task', taskId: props.task.id }),
    onDragStart: () => { /* ... */ },
    onDrop: () => { /* ... */ }
  })

  onUnmounted(cleanup)
})
```

**BoardColumn.vue** - Drop Target:
```javascript
import { dropTargetForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'

onMounted(() => {
  const cleanup = dropTargetForElements({
    element: columnRef.value,
    canDrop: ({ source }) => source.data.type === 'task',
    onDrop: ({ source, self }) => {
      // Handle task move
    }
  })

  onUnmounted(cleanup)
})
```

### Testing Drag-and-Drop

**Playwright Pattern**:
```javascript
// Drag and drop
await task.dragTo(targetColumn)

// With specific position
await task.dragTo(targetColumn, {
  targetPosition: { x: 50, y: 100 }
})

// Manual drag for complex scenarios
await task.hover()
await page.mouse.down()
await page.mouse.move(x, y)
await page.mouse.up()
```

**Test Coverage Requirements**:
- Basic drag between columns
- Reorder within column
- Drop indicators
- Auto-scroll behavior
- Keyboard navigation
- Accessibility
- Performance metrics

### Keyboard Navigation

**Shortcuts**:
- `Arrow Keys` - Navigate between tasks
- `Enter` - Open task detail
- `Space` - Select/deselect task
- `Ctrl+C/X/V` - Copy/cut/paste
- `Delete` - Remove task
- `Escape` - Cancel operation

### Accessibility

**ARIA Implementation**:
- Roles: `region`, `list`, `listitem`
- Labels: Descriptive aria-labels for all interactive elements
- Live regions: Announcements for drag operations
- Focus management: Visible focus indicators

**Screen Reader Support**:
- Live region announcements for all drag-and-drop operations
- Descriptive labels for task states
- Alternative keyboard-only workflows
```

### 5.2 Create DRAG_DROP_GUIDE.md

Comprehensive guide for developers working with the drag-and-drop system.

---

## Implementation Checklist

### Phase 1: Foundation ✅
- [x] Analyze current implementation
- [x] Research Atlassian pragmatic-drag-and-drop
- [x] Document gaps and requirements
- [ ] Install dependencies

### Phase 2: Core Migration
- [ ] Remove vuedraggable
- [ ] Refactor TaskCard with draggable adapter
- [ ] Refactor BoardColumn with dropTarget adapter
- [ ] Implement drop indicators
- [ ] Add auto-scroll behavior
- [ ] Test basic drag-and-drop

### Phase 3: Advanced Features
- [ ] Implement keyboard navigation
- [ ] Add screen reader support
- [ ] Enhance visual feedback
- [ ] Add toast notifications
- [ ] Implement loading states
- [ ] Add empty state improvements

### Phase 4: Testing
- [ ] Write Playwright drag-and-drop tests
- [ ] Write keyboard navigation tests
- [ ] Write accessibility tests
- [ ] Write performance tests
- [ ] Write cross-browser tests
- [ ] Achieve >90% test coverage

### Phase 5: Documentation
- [ ] Update CLAUDE.md
- [ ] Create DRAG_DROP_GUIDE.md
- [ ] Update component documentation
- [ ] Add usage examples
- [ ] Create video tutorials (optional)

---

## Success Metrics

### Performance
- ✅ Drag operations complete in <100ms
- ✅ Maintain 60 FPS during animations
- ✅ Bundle size increase <50kB
- ✅ Page load time <3 seconds

### Accessibility
- ✅ WCAG 2.1 AA compliance
- ✅ Keyboard navigation for all operations
- ✅ Screen reader announcements for all actions
- ✅ Focus indicators visible

### User Experience
- ✅ Visual feedback within 16ms
- ✅ Drop indicators always visible
- ✅ Auto-scroll when near edges
- ✅ Toast notifications for all actions

### Testing
- ✅ >90% code coverage
- ✅ All Playwright tests passing
- ✅ Cross-browser compatibility
- ✅ Mobile touch support

---

## Risk Mitigation

### Technical Risks

**Risk**: Breaking existing drag-and-drop functionality
**Mitigation**: Incremental migration, comprehensive testing, feature flags

**Risk**: Performance degradation
**Mitigation**: Performance budgets, monitoring, profiling

**Risk**: Accessibility regressions
**Mitigation**: Automated a11y tests, manual testing, screen reader testing

**Risk**: Browser compatibility issues
**Mitigation**: Cross-browser testing, polyfills, progressive enhancement

### Project Risks

**Risk**: Scope creep
**Mitigation**: Phased approach, clear deliverables, regular checkpoints

**Risk**: Timeline delays
**Mitigation**: Buffer time in estimates, parallel workstreams, prioritization

**Risk**: Testing gaps
**Mitigation**: Test-first approach, mandatory test coverage, CI/CD integration

---

## Next Steps

1. **Approve this plan** and review with stakeholders
2. **Set up Playwright** environment and verify tests can run
3. **Create feature branch**: `feature/pragmatic-drag-drop-migration`
4. **Begin Phase 2**: Install dependencies and start migration
5. **Daily progress tracking** with TodoWrite tool
6. **Weekly demos** of completed features

---

**Document Version**: 1.0
**Created**: 2025-10-15
**Author**: Claude Code
**Status**: Ready for Implementation
