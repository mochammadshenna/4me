import { expect } from '@playwright/test'

/**
 * Authentication helper for Playwright tests
 * Handles login with demo user credentials and state persistence
 */

const DEMO_USER = {
  username: 'admin',
  password: 'password'
}

/**
 * Logs in with demo user credentials
 * @param {import('@playwright/test').Page} page - Playwright page object
 */
export async function loginWithDemoUser(page) {
  // Navigate to login page
  await page.goto('/login')

  // Wait for login form to load
  await page.waitForSelector('input[type="text"], input[type="email"]')

  // Fill in credentials
  await page.fill('input[type="text"], input[type="email"]', DEMO_USER.username)
  await page.fill('input[type="password"]', DEMO_USER.password)

  // Submit login form
  await page.click('button[type="submit"], .v-btn:has-text("Login"), .v-btn:has-text("Sign In")')

  // Wait for redirect to dashboard
  await page.waitForURL('**/dashboard', { timeout: 10000 })
  await page.waitForSelector('.welcome-section', { timeout: 5000 })

  // Verify we're logged in
  const welcomeText = await page.textContent('.welcome-title')
  expect(welcomeText).toContain('ADMIN')
}

/**
 * Saves authentication state to file for reuse across tests
 * @param {import('@playwright/test').Page} page - Playwright page object
 * @param {string} filePath - Path to save auth state
 */
export async function saveAuthState(page, filePath = './tests/auth-state.json') {
  await page.context().storageState({ path: filePath })
}

/**
 * Loads authentication state from file
 * @param {import('@playwright/test').BrowserContext} context - Browser context
 * @param {string} filePath - Path to auth state file
 */
export async function loadAuthState(context, filePath = './tests/auth-state.json') {
  try {
    await context.addCookies(JSON.parse(require('fs').readFileSync(filePath, 'utf8')))
  } catch (error) {
    console.log('Auth state file not found, will need to login')
  }
}

/**
 * Sets up authenticated state for tests
 * @param {import('@playwright/test').Page} page - Playwright page object
 */
export async function setupAuthenticatedUser(page) {
  try {
    // Try to login with demo user
    await loginWithDemoUser(page)

    // Save auth state for future tests
    await saveAuthState(page)

    console.log('Successfully authenticated with demo user')
  } catch (error) {
    console.error('Failed to authenticate:', error)
    throw error
  }
}

/**
 * Checks if user is already authenticated
 * @param {import('@playwright/test').Page} page - Playwright page object
 * @returns {Promise<boolean>}
 */
export async function isAuthenticated(page) {
  try {
    await page.goto('/dashboard')
    await page.waitForSelector('.welcome-section', { timeout: 3000 })
    const welcomeText = await page.textContent('.welcome-title')
    return welcomeText && welcomeText.includes('ADMIN')
  } catch {
    return false
  }
}
