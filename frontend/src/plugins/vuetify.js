import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import 'vuetify/styles'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#3B82F6',
          secondary: '#64748b',
          accent: '#8b5cf6',
          error: '#ef4444',
          info: '#06b6d4',
          success: '#10b981',
          warning: '#f59e0b',
        },
      },
      dark: {
        colors: {
          primary: '#3B82F6',
          secondary: '#64748b',
          accent: '#8b5cf6',
          error: '#ef4444',
          info: '#06b6d4',
          success: '#10b981',
          warning: '#f59e0b',
        },
      },
    },
  },
})

export default vuetify

