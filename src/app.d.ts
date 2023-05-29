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
			owner: string;
			description: string;
			version: string;
			last_updated: string;
			downloads: number; // all time downloads
			stars: number;
			keywords: string[];
		};

		//select * from get_all_dependency_licenses(4);
		//"license"	"packages"
		//"BSD 3-Clause"	"{""async runtime""}"
		//"MIT"	"{""http server""}"

		//select * from get_dependencies_flat(4)
		//"packagename"	"version"	"license"	"lastupdated"	"archived"	"insecure"
		//"http server"	"1.0.0"	"MIT"	"2023-05-28 22:59:16.899138"	false	true
		//"async runtime"	"1.2.3"	"BSD 3-Clause"	"2023-05-28 22:59:16.899138"	false	false
		// select * from get_package_details(1);
		// "name"	"description"	"archived"	"keywords"	"stars"	"repository"	"readme"	"owner"	"authors"
		// "http server"	"a cool http/1.1 server"	false	"{fancy,pants}"	0	"https://repository1"	"readme1"	"jon"	"{}"

		type PackageDetails = {
			id: number;
			name: string;
			description: string;
			archived: boolean;
			keywords: string[];
			stars: number;
			repository: string; // http repo url
			readme: string; // markdown, html formatted
			owner: string;
			authors: string[];
			versions: VersionDetails[]; // appended by seperate query
			usedBy: string[]?; // appended by seperate query
		};
		//select * from get_version_details(4);
		//"version"	"isinsecure"	"createdat"	"size_kb"	"dependencycount"	"compiler"	"license"	"insecuredependency"
		//"99.99.99"	false	"2023-05-28 23:02:46.186889"	99999	2	"dev-2023-05"	"GPL"	true
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
