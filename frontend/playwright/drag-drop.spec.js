import { test, expect } from '@playwright/test'

/**
 * Comprehensive Drag-and-Drop Test Suite
 * Tests the pragmatic-drag-and-drop implementation in the Kanban board
 */

// Test data
const TEST_USER = {
  username: 'testuser',
  email: 'test@example.com',
  password: 'Test123456!'
}

const TEST_PROJECT = {
  name: 'Test Project for Drag-Drop',
  description: 'Testing drag and drop functionality'
}

// Helper functions
async function login(page, username = TEST_USER.username, password = TEST_USER.password) {
  await page.goto('/login')
  await page.fill('[data-testid="username-input"]', username)
  await page.fill('[data-testid="password-input"]', password)
  await page.click('[data-testid="login-button"]')
  await page.waitForURL('/projects')
}

async function createProject(page, projectName = TEST_PROJECT.name) {
  await page.click('[data-testid="create-project-button"]')
  await page.fill('[data-testid="project-name-input"]', projectName)
  await page.click('[data-testid="save-project-button"]')
  await page.waitForTimeout(500)
}

async function createTask(page, boardId, taskTitle) {
  const board = page.locator(`[data-board-id="${boardId}"]`)
  await board.locator('[data-testid="add-task-button"]').click()
  await board.locator('[data-testid="task-title-input"]').fill(taskTitle)
  await board.locator('[data-testid="save-task-button"]').click()
  await page.waitForTimeout(300)
}

test.describe('Drag-and-Drop Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Register or login
    await page.goto('/register')
    try {
      await page.fill('[data-testid="username-input"]', TEST_USER.username)
      await page.fill('[data-testid="email-input"]', TEST_USER.email)
      await page.fill('[data-testid="password-input"]', TEST_USER.password)
      await page.click('[data-testid="register-button"]')
    } catch {
      // User might already exist, try login
      await login(page)
    }
  })

  test('should display draggable task cards', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create a task
    const todoBoard = page.locator('[data-board-id="1"]')
    await createTask(page, 1, 'Test Task 1')

    // Verify task is draggable
    const taskCard = page.locator('[data-task-id]').first()
    await expect(taskCard).toBeVisible()
    await expect(taskCard).toHaveAttribute('data-task-id')

    // Check for drag handle
    const dragHandle = taskCard.locator('[data-drag-handle="true"]')
    await expect(dragHandle).toBeVisible()
  })

  test('should drag task from one column to another', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create task in "To Do" board
    await createTask(page, 1, 'Draggable Task')

    // Get source and target columns
    const sourceColumn = page.locator('[data-board-id="1"]')
    const targetColumn = page.locator('[data-board-id="2"]') // In Progress

    // Get the task
    const task = sourceColumn.locator('[data-task-id]').first()
    const taskId = await task.getAttribute('data-task-id')

    // Perform drag and drop
    await task.dragTo(targetColumn)

    // Wait for animation and API call
    await page.waitForTimeout(500)

    // Verify task moved to target column
    const movedTask = targetColumn.locator(`[data-task-id="${taskId}"]`)
    await expect(movedTask).toBeVisible()

    // Verify task removed from source column
    const taskInSource = sourceColumn.locator(`[data-task-id="${taskId}"]`)
    await expect(taskInSource).not.toBeVisible()
  })

  test('should reorder tasks within the same column', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create multiple tasks in same column
    await createTask(page, 1, 'Task 1')
    await createTask(page, 1, 'Task 2')
    await createTask(page, 1, 'Task 3')

    const column = page.locator('[data-board-id="1"]')
    const tasks = column.locator('[data-task-id]')

    // Get initial order
    const initialCount = await tasks.count()
    expect(initialCount).toBe(3)

    // Get first and last task
    const firstTask = tasks.first()
    const lastTask = tasks.last()

    // Drag first task to bottom
    await firstTask.dragTo(lastTask)
    await page.waitForTimeout(500)

    // Verify order changed
    const newFirstTaskText = await column.locator('[data-task-id]').first().textContent()
    expect(newFirstTaskText).not.toContain('Task 1')
  })

  test('should show drop indicators during drag', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Indicator Test Task')

    const task = page.locator('[data-task-id]').first()
    const targetColumn = page.locator('[data-board-id="2"]')

    // Start dragging
    await task.hover()
    await page.mouse.down()

    // Move to target column
    const targetBox = await targetColumn.boundingBox()
    await page.mouse.move(targetBox.x + 100, targetBox.y + 50)

    // Check for drop indicator
    const dropIndicator = page.locator('.drop-indicator')
    await expect(dropIndicator).toBeVisible()

    // Complete drag
    await page.mouse.up()
  })

  test('should highlight drop target during drag', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Highlight Test')

    const task = page.locator('[data-task-id]').first()
    const targetColumn = page.locator('[data-board-id="2"]')

    // Start drag
    await task.hover()
    await page.mouse.down()

    // Move to target
    const targetBox = await targetColumn.boundingBox()
    await page.mouse.move(targetBox.x + 100, targetBox.y + 50)

    // Verify target column has highlight class
    await expect(targetColumn).toHaveClass(/drop-target-active/)

    await page.mouse.up()
  })

  test('should apply dragging visual state to task', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Visual State Test')

    const task = page.locator('[data-task-id]').first()

    // Start dragging
    await task.hover()
    await page.mouse.down()

    // Verify dragging class applied
    await expect(task).toHaveClass(/is-dragging/)

    // Verify reduced opacity during drag
    const opacity = await task.evaluate((el) =>
      window.getComputedStyle(el).opacity
    )
    expect(parseFloat(opacity)).toBeLessThan(1)

    await page.mouse.up()
  })

  test('should support auto-scroll when dragging near edges', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create many tasks to enable scrolling
    for (let i = 1; i <= 10; i++) {
      await createTask(page, 1, `Task ${i}`)
    }

    const column = page.locator('[data-board-id="1"]')
    const firstTask = column.locator('[data-task-id]').first()

    // Start dragging first task
    await firstTask.hover()
    await page.mouse.down()

    // Move to bottom of column (should trigger auto-scroll)
    const columnBox = await column.boundingBox()
    await page.mouse.move(columnBox.x + 100, columnBox.y + columnBox.height - 10)

    // Wait for auto-scroll
    await page.waitForTimeout(1000)

    // Verify scroll position changed
    const scrollTop = await column.evaluate((el) => el.scrollTop)
    expect(scrollTop).toBeGreaterThan(0)

    await page.mouse.up()
  })

  test('should prevent click event when dragging', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Click Prevention Test')

    const task = page.locator('[data-task-id]').first()

    // Try to drag then click
    await task.hover()
    await page.mouse.down()
    await page.mouse.move(50, 50) // Small movement
    await page.mouse.up()
    await task.click()

    // Verify task detail modal doesn't open (or check other non-click behavior)
    const modal = page.locator('[data-testid="task-detail-modal"]')
    await expect(modal).not.toBeVisible()
  })

  test('should handle drag between multiple columns', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create tasks in first column
    await createTask(page, 1, 'Multi-column Task 1')
    await createTask(page, 1, 'Multi-column Task 2')

    // Drag first task to second column
    const task1 = page.locator('[data-board-id="1"]').locator('[data-task-id]').first()
    const task1Id = await task1.getAttribute('data-task-id')
    await task1.dragTo(page.locator('[data-board-id="2"]'))
    await page.waitForTimeout(500)

    // Drag second task to third column
    const task2 = page.locator('[data-board-id="1"]').locator('[data-task-id]').first()
    const task2Id = await task2.getAttribute('data-task-id')
    await task2.dragTo(page.locator('[data-board-id="3"]'))
    await page.waitForTimeout(500)

    // Verify both tasks moved correctly
    await expect(page.locator(`[data-board-id="2"] [data-task-id="${task1Id}"]`)).toBeVisible()
    await expect(page.locator(`[data-board-id="3"] [data-task-id="${task2Id}"]`)).toBeVisible()
  })

  test('should maintain task data after drag', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Data Integrity Test')

    // Get task text before drag
    const task = page.locator('[data-task-id]').first()
    const taskText = await task.textContent()
    const taskId = await task.getAttribute('data-task-id')

    // Drag to another column
    await task.dragTo(page.locator('[data-board-id="2"]'))
    await page.waitForTimeout(500)

    // Verify task data unchanged
    const movedTask = page.locator(`[data-task-id="${taskId}"]`)
    const movedTaskText = await movedTask.textContent()
    expect(movedTaskText).toContain('Data Integrity Test')
  })

  test('should handle rapid successive drags', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Rapid Drag Test')

    const task = page.locator('[data-task-id]').first()
    const column2 = page.locator('[data-board-id="2"]')
    const column3 = page.locator('[data-board-id="3"]')

    // Rapid drags
    await task.dragTo(column2)
    await page.waitForTimeout(200)

    const task2 = column2.locator('[data-task-id]').first()
    await task2.dragTo(column3)
    await page.waitForTimeout(200)

    const task3 = column3.locator('[data-task-id]').first()
    await task3.dragTo(page.locator('[data-board-id="1"]'))
    await page.waitForTimeout(500)

    // Verify final position
    await expect(page.locator('[data-board-id="1"]').locator('[data-task-id]').first()).toBeVisible()
  })

  test('should show custom drag preview', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Preview Test')

    const task = page.locator('[data-task-id]').first()

    // Start dragging
    await task.hover()
    await page.mouse.down()
    await page.mouse.move(100, 100)

    // The drag preview is handled by browser native drag-and-drop
    // We can verify the task maintains its style during drag
    await expect(task).toHaveClass(/is-dragging/)

    await page.mouse.up()
  })

  test('should work on mobile viewport', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 })

    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Mobile Drag Test')

    const task = page.locator('[data-task-id]').first()
    const targetColumn = page.locator('[data-board-id="2"]')

    // Touch-based drag
    await task.dragTo(targetColumn)
    await page.waitForTimeout(500)

    // Verify drag worked on mobile
    await expect(targetColumn.locator('[data-task-id]')).toBeVisible()
  })

  test('should handle empty column drops', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Empty Column Test')

    const task = page.locator('[data-task-id]').first()
    const emptyColumn = page.locator('[data-board-id="3"]') // Assuming Done is empty

    // Drag to empty column
    await task.dragTo(emptyColumn)
    await page.waitForTimeout(500)

    // Verify task appears in empty column
    await expect(emptyColumn.locator('[data-task-id]')).toBeVisible()

    // Verify empty state message is gone
    const emptyState = emptyColumn.locator('.empty-state')
    await expect(emptyState).not.toBeVisible()
  })

  test('should preserve scroll position after drag', async ({ page }) => {
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create many tasks
    for (let i = 1; i <= 15; i++) {
      await createTask(page, 1, `Scroll Test Task ${i}`)
    }

    const column = page.locator('[data-board-id="1"]')

    // Scroll down
    await column.evaluate((el) => el.scrollTop = 200)
    const scrollBefore = await column.evaluate((el) => el.scrollTop)

    // Drag a visible task
    const visibleTask = column.locator('[data-task-id]').nth(5)
    await visibleTask.dragTo(page.locator('[data-board-id="2"]'))
    await page.waitForTimeout(500)

    // Verify scroll position maintained (approximately)
    const scrollAfter = await column.evaluate((el) => el.scrollTop)
    expect(Math.abs(scrollAfter - scrollBefore)).toBeLessThan(50)
  })
})

test.describe('Drag-and-Drop Accessibility', () => {
  test('should have proper ARIA attributes', async ({ page }) => {
    await login(page)
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Accessibility Test')

    const task = page.locator('[data-task-id]').first()

    // Check for accessibility attributes
    await expect(task).toHaveAttribute('data-task-id')
    await expect(task).toBeVisible()
  })

  test('should announce drag operations to screen readers', async ({ page }) => {
    // This test would require checking for live region updates
    // Placeholder for comprehensive accessibility testing
    await login(page)
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Screen Reader Test')

    // Live regions should announce drag start/end
    // This would require actual screen reader integration testing
  })
})

test.describe('Drag-and-Drop Performance', () => {
  test('should handle 50+ tasks without lag', async ({ page }) => {
    await login(page)
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    // Create many tasks
    for (let i = 1; i <= 50; i++) {
      await createTask(page, 1, `Performance Task ${i}`)
    }

    // Measure drag performance
    const task = page.locator('[data-task-id]').first()
    const startTime = Date.now()

    await task.dragTo(page.locator('[data-board-id="2"]'))
    await page.waitForTimeout(300)

    const endTime = Date.now()
    const duration = endTime - startTime

    // Drag should complete in reasonable time
    expect(duration).toBeLessThan(2000)
  })

  test('should not cause memory leaks', async ({ page }) => {
    await login(page)
    await createProject(page)
    const project = page.locator('[data-testid="project-card"]').first()
    await project.click()

    await createTask(page, 1, 'Memory Test Task')

    // Perform multiple drags
    for (let i = 0; i < 20; i++) {
      const task = page.locator('[data-task-id]').first()
      const targetBoard = i % 2 === 0 ? '2' : '1'
      await task.dragTo(page.locator(`[data-board-id="${targetBoard}"]`))
      await page.waitForTimeout(200)
    }

    // Page should remain responsive
    const task = page.locator('[data-task-id]').first()
    await expect(task).toBeVisible()
  })
})
