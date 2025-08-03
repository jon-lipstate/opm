export function timeAgo(input) {
	const now = new Date()

	// Parse the input into a Date object
	const past = new Date(input)

	// Handle invalid dates
	if (isNaN(past)) {
		return 'invalid date'
	}

	// Check if the input is a date-only string (e.g., "YYYY-MM-DD")
	const isDateOnly = /^\d{4}-\d{2}-\d{2}$/.test(input)

	// If it's a date-only input and matches today's date, return "today"
	if (isDateOnly && now.toDateString() === past.toDateString()) {
		return 'today'
	}

	// Calculate the difference in seconds
	const diff = (now - past) / 1000

	// Define time intervals in seconds
	const intervals = {
		year: 365 * 24 * 60 * 60,
		month: 30 * 24 * 60 * 60,
		week: 7 * 24 * 60 * 60,
		day: 24 * 60 * 60,
		hour: 60 * 60,
		minute: 60,
		second: 1,
	}

	if (diff < 60) {
		return 'just now'
	}

	// Loop through intervals to find the largest applicable time interval
	for (const key in intervals) {
		const value = Math.floor(diff / intervals[key])
		if (value > 0) {
			return `${value} ${key}${value > 1 ? 's' : ''} ago`
		}
	}

	return 'just now' // Fallback for very recent timestamps
}

export function devLog(...args) {
	if (process.env.DEV) console.log(...args)
}

import { useUserStore } from 'src/stores/user-store'
export function expectAuth() {
	if (process.env.VITE_ENV != 'dev' && process.env.VITE_ENV != 'preview') {
		return true // dont run in production
	}
	const userStore = useUserStore()
	// Check both isAuthenticated and user existence
	if (!userStore.isAuthenticated || !userStore.user) {
		const stack = new Error().stack
		// Extract the call-site information (file and line number)
		const callSite =
			stack
				?.split('\n')[2] // Third line in stack trace (function caller)
				?.trim() || 'Unknown location'
		console.error(`Expected User Auth. Call-site: ${callSite}`)
		return false
	}
	return true
}
