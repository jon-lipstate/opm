// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	type auth = { login: string; authHeader: any; session: any };
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface Platform {}

		type SearchResult = {
			package_id: number;
			name: string;
			description: string;
			version: string;
			last_updated: string;
			downloads: number;
			all_downloads: number;
			stars: number;
			keywords: string[];
		};
		type PackageDetails = {
			name: string;
			version: string;
			description: string;
			archived: boolean;
			tags: string[];
			readme: string; // markdown, html formatted
			versions: string[];
			funding: Record<string, string>; // github patreon etc?
			dependsOn: NamedVersion;
			// usedBy: string[];
			requirements: { compiler: string };
			links: Record<string, string>;
			lastUpdated: string;
			license: string;
			kind: string; //"unstable"|"community"|"curated"|"demo",
			size_kb: string; // kb
			owners: { name: string; username: string }[];
			// stats: { allTimeDownloads: number };
		};
		type ModPkg = {
			name: string;
			version: string;
			description: string;
			authors: string[];
			repository: URL;
			license: string;
			keywords: string[];
			funding: Record<string, string>?;
			kind: 'Demo' | 'Library';
			os: string['Linux' | 'Windows' | 'Darwin' | 'Essence'];
			compiler: string;
			dependencies: NamedVersion;
		};
		type NamedVersion = Record<string, string>; // {name:version}
	}
}

export {};
