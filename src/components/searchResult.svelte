<script lang="ts">
	import Tags from '$components/tags.svelte';
	import { timeAgo } from '$lib/utils';

	export let pkg: App.SearchResult;
	let packageName = pkg.name.toLowerCase().replace(/\s/g, '-');
</script>

<div class="container">
	<h2>
		<a href={`/packages/${pkg.owner}/${packageName}`}>
			{pkg.name}
		</a>
		<span>( {pkg.owner} )</span>
	</h2>
	<span style="display:block">{pkg.description}</span>
	<!-- <pre>{JSON.stringify(pkg, null, 2)}</pre> -->
	<span style="display:block; margin-top:0.5rem"><strong>Updated:</strong>{timeAgo(new Date(pkg.last_updated))}</span>
	<span style="display:block"><strong>Downloads:</strong> {pkg.downloads ?? 0}</span>
	<span style="display:block">v{pkg.version}</span>
	<!-- <p>copy to clipboard</p> -->
	<Tags tags={pkg.keywords} />
</div>

<style>
	h2 {
		padding: 0;
		margin: 0;
	}
	.container {
		border: 1px solid #ddd;
		padding: 10px;
		margin: 10px 0;
	}
</style>
