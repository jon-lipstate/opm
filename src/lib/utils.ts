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
export function intoPackageDetails(modpkg: App.ModPkg, githubDescriptor: any, readmeUrl: string): App.PackageDetails {
	let details: Partial<App.PackageDetails> = {
		requirements: { compiler: '' },
		dependsOn: {},
		versions: [], //name/date pair?
		tags: [],
		links: {},
		funding: {}
	};
	details.name = modpkg.name;
	details.version = modpkg.version;
	details.kind = modpkg.kind;
	details.lastUpdated = githubDescriptor.updated_at;
	details.description = modpkg.description;
	details.license = modpkg.license;
	details.readme = readmeUrl;
	details.size_kb = githubDescriptor.size;
	details.owners = [githubDescriptor.owner.login];
	details.archived = githubDescriptor.archived;
	details.links = { url: githubDescriptor.html_url };
	details.tags = modpkg.keywords ?? [];
	details.links.hompage = githubDescriptor.hompage;
	details.requirements!.compiler = modpkg.compiler;
	details.dependsOn = modpkg.dependencies;
	//@ts-ignore
	details.funding = modpkg.funding;

	return details as App.PackageDetails;
}
