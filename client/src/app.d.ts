// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface Platform {}

	type PackageResult = {
		version: string;
		name: string;
		kind: string;
		updated: string;
		downloads: number;
		tags: string[];
	};
	type PackageDetails = {
		name: string;
		version: string;
		description: string;
		tags: string[];
		readme: string; // markdown, html formatted
		versions: string[];
		funding:string[]; // github patreon etc?
		dependsOn: string[];
		usedBy: string[];
		requirements: { minCompilierVersion: string };
		links: Record<string,string>;
		lastUpdated: string;
		license: string;
		kind: "unstable"|"community"|"curated"|"demo",
		size: string; // kb
		owners: { name: string; username: string }[];
		stats: { allTimeDownloads: number };
	};
	}
}

export {};
