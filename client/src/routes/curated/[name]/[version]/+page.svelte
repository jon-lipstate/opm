<script lang="ts">
	import { page } from '$app/stores';
	let query = $page.url.searchParams.get('query') || 'No query provided';

	const results = [
		{
			name: "Yuki's ECS",
			version: '1.2.3',
			description: 'an ecs game engine',
			updated: '1/1/1901',
			downloads: 100000,
			tags: ['ecs', 'engine'],
			kind: 'Curated Library'
		}
	];

	type PackageDetails = {
		name: string;
		version: string;
		description: string;
		tags: string[];
		readme: string; // markdown, html formatted
		// versions: string[];
		// funding:string[]; // github patreon etc?
		// dependencies: string[];
		// dependents: string[];
		requirements: { minCompilierVersion: string };
		docsUrl: string;
		repositoryUrl: string;
		lastUpdated: string;
		license: string;
		size: string; // kb
		// qualityRating:number; // 1-100 ?
		owners: { name: string; username: string }[];
		// stats: { allTimeDownloads: number };
	};
</script>

<svelte:head>
	<title>Search Results for '{query}'</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<section>
	<h1>hi</h1>
</section>

<style>
</style>
