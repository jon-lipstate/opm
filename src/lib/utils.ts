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
export function extractHostOwnerAndRepo(urlString) {
	try {
		const url = new URL(urlString);

		// Remove potential credentials
		url.username = '';
		url.password = '';

		const host_name = url.host;

		// The path should look like /user/repo, so after splitting
		// the owner should be all parts up to the last part, and repo the last part
		const pathParts = url.pathname.split('/').filter(Boolean); // filter(Boolean) removes empty strings

		if (pathParts.length < 2) {
			throw new Error('Invalid URL');
		}

		const repo_name = pathParts.pop(); // The last part of the path is the repo
		const owner_name = pathParts.join('/'); // The rest of the path is the owner

		return { host_name, owner_name, repo_name };
	} catch (error: any) {
		const e = `Failed to parse URL: ${error.message}`;
		console.error(e);
		throw Error(e);
	}
}

export function generateRandomString(length: number): string {
	const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
	let result = '';
	for (let i = 0; i < length; i++) {
		const randomIndex = Math.floor(Math.random() * characters.length);
		result += characters[randomIndex];
	}
	return result;
}
