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
			slug: string;
			owner: string;
			description: string;
			version: string;
			last_updated: string;
			downloads: number; // all time downloads
			bookmarks: number;
			license: string;
			dependency_count: string;
			keywords: string[];
		};

		type PackageDetails = {
			id: number;
			name: string;
			slug: string;
			description: string;
			archived: boolean;
			keywords: string[];
			bookmarks: number;
			url: string; // http repo url
			readme: string; // markdown, html formatted
			owner: string;
			authors: string[]?;
			versions: VersionDetails[]; // appended by seperate query
			usedBy: string[]?; // appended by seperate query
		};
		type VersionDetails = {
			id: number;
			version: string;
			insecure: boolean;
			createdat: string;
			size_kb: number;
			dependency_count: number;
			compiler: string; // eg DEV-05-23
			license: string;
			has_insecure_dependency: boolean;
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
