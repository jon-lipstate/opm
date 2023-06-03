export function generateRandomName(length: number): string {
	const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
	let result = '';
	for (let i = 0; i < length; i++) {
		result += chars.charAt(Math.floor(Math.random() * chars.length));
	}
	return result;
}
export function timeAgo(dateString): string {
	const now = new Date();
	const past = new Date(dateString);
	//@ts-ignore
	const diffMs = now - past; // difference in milliseconds
	if (diffMs < 0) {
		return `In the future`;
	}
	const seconds = diffMs / 1000;
	const minutes = seconds / 60;
	const hours = minutes / 60;
	const days = hours / 24;
	const months = days / 30.44; // Average number of days in a month
	const years = days / 365.25; // Number of days in a year, accounting for leap years

	// Return a string describing the difference
	if (seconds < 60) {
		return `${Math.round(seconds)} seconds ago`;
	} else if (minutes < 60) {
		return `${Math.round(minutes)} minutes ago`;
	} else if (hours < 24) {
		return `${Math.round(hours)} hours ago`;
	} else if (days < 30.44) {
		return `${Math.round(days)} days ago`;
	} else if (months < 12) {
		return `${Math.round(months)} months ago`;
	} else {
		return `${Math.round(years)} years ago`;
	}
}
/**
 * is >90 days old
 */
export function isStale(date): boolean {
	const ninetyDaysInMilliseconds = 90 * 24 * 60 * 60 * 1000;
	const difference = Date.now() - new Date(date).getTime();
	if (difference < 0) return false; // dates in future
	return difference > ninetyDaysInMilliseconds;
}

export function generateSlug(str: string): string {
	return str
		.toLowerCase() // convert to lower case
		.replace(/\s+/g, '-') // replace spaces with hyphens
		.replace(/[^\w\-]+/g, '') // non-word [a-z0-9_], non-hyphen characters
		.replace(/\-\-+/g, '-') // replace multiple hyphens with a single hyphen
		.replace(/^-+/, '') //  leading hyphens
		.replace(/-+$/, ''); //  trailing hyphens
}

export function isValidSemver(version) {
	// ban the v, we want to store without
	const semverRegex = new RegExp(
		/^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$/
	);
	return semverRegex.test(version);
}
/**
 *
 * @param urlString ASSUMES: https: hostname / user / repository
 * @returns user, repo
 */
export function extractUserAndProject(urlString) {
	try {
		const url = new URL(urlString);
		const pathParts = url.pathname.split('/');

		// The path should look like /user/repo, so after splitting
		// the user should be in index 1 and repo in index 2
		if (pathParts.length < 3) {
			throw new Error('Invalid URL');
		}

		const user = pathParts[1];
		const repo = pathParts[2];

		return { user, repo };
	} catch (error: any) {
		console.error(`Failed to parse URL: ${error.message}`);
		return null;
	}
}
