<script lang="ts">
	import Tags from "$components/tags.svelte"
	import { goto } from '$app/navigation';

	export let pkg: PackageResult;
	let packageName = pkg.name.toLowerCase().replace(/\s/g, '-');

	type PackageResult = {
		version: string;
		name: string;
		kind: string;
		updated: string;
		downloads: number;
		tags: string[];
	};
</script>

<div class="container">
	<h2>
		<a href={`/curated/${packageName}/${pkg.version}`}>
			{pkg.name} <small>v{pkg.version}</small>
		</a>
	</h2>
	<!-- <p><strong>Kind:</strong> {pkg.kind}</p> -->
	<!-- <p><strong>Updated:</strong> {pkg.updated}</p> -->
	<!-- <p><strong>Downloads:</strong> {pkg.downloads.toLocaleString()}</p> -->
	<!-- <p>copy to clipboard</p> -->
	<Tags tags={pkg.tags}/>
</div>

<style>
	.container {
		border: 1px solid #ddd;
		padding: 10px;
		margin: 10px 0;
	}
	
</style>
