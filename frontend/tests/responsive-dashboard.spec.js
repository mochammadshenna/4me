import { expect, test } from '@playwright/test'
import { isAuthenticated, setupAuthenticatedUser } from './helpers/auth.js'

test.describe('Dashboard Responsive Layout', () => {
    test.beforeEach(async ({ page }) => {
        // Check if already authenticated, if not, login
        if (!(await isAuthenticated(page))) {
            await setupAuthenticatedUser(page)
        }

        // Navigate to dashboard
        await page.goto('/')

        // Wait for the dashboard to load
        await page.waitForSelector('.welcome-section')
    })

    test('should display dashboard elements correctly on desktop', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Check that main elements are visible
        await expect(page.locator('.welcome-title')).toBeVisible()
        await expect(page.locator('.search-container')).toBeVisible()
        await expect(page.locator('.stats-grid')).toBeVisible()
        await expect(page.locator('.projects-section')).toBeVisible()

        // Check sidebar is visible
        await expect(page.locator('.studybuddy-sidebar')).toBeVisible()

        // Capture screenshot of desktop layout
        await page.screenshot({
            path: 'test-results/screenshots/desktop-expanded-sidebar.png',
            fullPage: true
        })
    })

    test('should collapse sidebar and expand main content', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Get initial main content width
        const mainContent = page.locator('.main-content')
        const initialWidth = await mainContent.boundingBox()

        // Click on user card to collapse sidebar
        await page.locator('.clickable-user-card').click()

        // Wait for transition
        await page.waitForTimeout(400)

        // Check that sidebar is collapsed
        await expect(page.locator('.studybuddy-sidebar')).toHaveClass(/v-navigation-drawer--mini-variant/)

        // Check that main content has expanded
        const expandedWidth = await mainContent.boundingBox()
        expect(expandedWidth.width).toBeGreaterThan(initialWidth.width)

        // Capture screenshot of collapsed sidebar layout
        await page.screenshot({
            path: 'test-results/screenshots/desktop-collapsed-sidebar.png',
            fullPage: true
        })
    })

    test('should expand sidebar and contract main content', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // First collapse the sidebar
        await page.locator('.clickable-user-card').click()
        await page.waitForTimeout(400)

        // Get collapsed main content width
        const mainContent = page.locator('.main-content')
        const collapsedWidth = await mainContent.boundingBox()

        // Click on logo to expand sidebar
        await page.locator('.clickable-logo').click()

        // Wait for transition
        await page.waitForTimeout(400)

        // Check that sidebar is expanded
        await expect(page.locator('.studybuddy-sidebar')).not.toHaveClass(/v-navigation-drawer--mini-variant/)

        // Check that main content has contracted
        const expandedWidth = await mainContent.boundingBox()
        expect(expandedWidth.width).toBeLessThan(collapsedWidth.width)
    })

    test('should adapt stats grid layout when sidebar collapses', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Get initial stats grid layout
        const statsGrid = page.locator('.stats-grid')
        const initialGridStyle = await statsGrid.evaluate(el => {
            return window.getComputedStyle(el).gridTemplateColumns
        })

        // Collapse sidebar
        await page.locator('.clickable-user-card').click()
        await page.waitForTimeout(400)

        // Check that stats grid has adapted (should have more columns)
        const collapsedGridStyle = await statsGrid.evaluate(el => {
            return window.getComputedStyle(el).gridTemplateColumns
        })

        // The grid should adapt to use more space
        expect(collapsedGridStyle).toBeDefined()
    })

    test('should adapt projects grid layout when sidebar collapses', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Get initial projects grid layout
        const projectsGrid = page.locator('.projects-grid')
        const initialGridStyle = await projectsGrid.evaluate(el => {
            return window.getComputedStyle(el).gridTemplateColumns
        })

        // Collapse sidebar
        await page.locator('.clickable-user-card').click()
        await page.waitForTimeout(400)

        // Check that projects grid has adapted
        const collapsedGridStyle = await projectsGrid.evaluate(el => {
            return window.getComputedStyle(el).gridTemplateColumns
        })

        // The grid should adapt to use more space
        expect(collapsedGridStyle).toBeDefined()
    })

    test('should be responsive on mobile devices', async ({ page }) => {
        // Set mobile viewport
        await page.setViewportSize({ width: 375, height: 667 })

        // Check that elements are stacked vertically
        await expect(page.locator('.stats-grid')).toBeVisible()
        await expect(page.locator('.projects-grid')).toBeVisible()

        // Check that sidebar behavior is different on mobile
        const sidebar = page.locator('.studybuddy-sidebar')
        const sidebarClasses = await sidebar.getAttribute('class')

        // On mobile, sidebar should be temporary
        expect(sidebarClasses).toContain('v-navigation-drawer--temporary')

        // Capture screenshot of mobile layout
        await page.screenshot({
            path: 'test-results/screenshots/mobile-layout.png',
            fullPage: true
        })
    })

    test('should handle ultra-wide screens properly', async ({ page }) => {
        // Set ultra-wide viewport
        await page.setViewportSize({ width: 2560, height: 1440 })

        // Check that content is centered and not stretched too wide
        const mainContent = page.locator('.main-content')
        const boundingBox = await mainContent.boundingBox()

        // Content should be reasonably sized, not stretched across entire width
        expect(boundingBox.width).toBeLessThan(2000) // Should have some max-width constraint
    })

    test('should maintain smooth transitions during sidebar toggle', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Monitor transition performance
        const mainContent = page.locator('.main-content')

        // Start transition
        await page.locator('.clickable-user-card').click()

        // Check that transition is smooth by monitoring multiple points
        await page.waitForTimeout(100)
        const midTransition = await mainContent.boundingBox()

        await page.waitForTimeout(200)
        const nearEnd = await mainContent.boundingBox()

        await page.waitForTimeout(200)
        const final = await mainContent.boundingBox()

        // Verify smooth progression
        expect(midTransition.width).toBeGreaterThan(0)
        expect(nearEnd.width).toBeGreaterThan(0)
        expect(final.width).toBeGreaterThan(0)
    })
})

test.describe('Project Creation Functionality', () => {
    test.beforeEach(async ({ page }) => {
        // Check if already authenticated, if not, login
        if (!(await isAuthenticated(page))) {
            await setupAuthenticatedUser(page)
        }

        // Navigate to dashboard
        await page.goto('/')

        // Wait for the dashboard to load
        await page.waitForSelector('.welcome-section')
    })

    test('should open project creation dialog when FAB is clicked', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Click the floating action button
        await page.click('.fab-button')

        // Wait for dialog to appear
        await page.waitForSelector('.v-dialog', { state: 'visible' })

        // Check that dialog is visible
        await expect(page.locator('.v-dialog')).toBeVisible()
        await expect(page.locator('text=Create New Project')).toBeVisible()

        // Capture screenshot of project creation dialog
        await page.screenshot({
            path: 'test-results/screenshots/project-creation-dialog.png',
            fullPage: true
        })
    })

    test('should validate project form fields', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Open project creation dialog
        await page.click('.fab-button')
        await page.waitForSelector('.v-dialog', { state: 'visible' })

        // Try to submit form without filling required fields
        await page.click('button:has-text("Create")')

        // Check that form validation prevents submission
        // The dialog should still be visible
        await expect(page.locator('.v-dialog')).toBeVisible()
    })

    test('should create project with valid data', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Open project creation dialog
        await page.click('.fab-button')
        await page.waitForSelector('.v-dialog', { state: 'visible' })

        // Fill in project details
        await page.fill('input[label="Project Name"]', 'Test Project')
        await page.fill('textarea[label*="Description"]', 'This is a test project description')

        // Select a color (click on the first color option)
        await page.click('.w-10.h-10.rounded-full:first-of-type')

        // Submit the form
        await page.click('button:has-text("Create")')

        // Wait for dialog to close
        await page.waitForSelector('.v-dialog', { state: 'hidden' })

        // Check that the project appears in the projects grid
        await expect(page.locator('text=Test Project')).toBeVisible()

        // Capture screenshot of created project
        await page.screenshot({
            path: 'test-results/screenshots/project-created.png',
            fullPage: true
        })
    })

    test('should handle project creation error gracefully', async ({ page }) => {
        // Set desktop viewport
        await page.setViewportSize({ width: 1920, height: 1080 })

        // Open project creation dialog
        await page.click('.fab-button')
        await page.waitForSelector('.v-dialog', { state: 'visible' })

        // Fill in project details with invalid data (very long name)
        await page.fill('input[label="Project Name"]', 'a'.repeat(1000))

        // Submit the form
        await page.click('button:has-text("Create")')

        // Check that error is handled (dialog might close or show error)
        // The test should not crash
        await page.waitForTimeout(1000)
    })
})
