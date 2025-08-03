import { defineBoot } from '#q-app/wrappers'
import { Dark } from 'quasar'

export default defineBoot(() => {
	// Load dark mode preference from localStorage
	const savedDarkMode = localStorage.getItem('darkMode')
	
	// If no preference saved, default to true (dark mode)
	const isDarkMode = savedDarkMode === null ? true : savedDarkMode === 'true'
	
	Dark.set(isDarkMode)
})
