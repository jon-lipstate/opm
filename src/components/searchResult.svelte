<script lang="ts">
	import Tags from '$components/tags.svelte';
	import { goto } from '$app/navigation';
	import { timeAgo } from '$lib/utils';

	export let pkg: App.SearchResult;
	let packageName = pkg.name.toLowerCase().replace(/\s/g, '-');
</script>

<div class="container">
	<h2>
		<a href={`/packages/${packageName}/${pkg.version}`}>
			{pkg.name} <small>v{pkg.version}</small>
		</a>
	</h2>
	<span><strong>Updated:</strong>{timeAgo(new Date(pkg.last_updated))}</span>
	<p><strong>Downloads:</strong> {pkg.downloads ?? 0} (all versions: {pkg.all_downloads ?? 0})</p>
	<p><strong>Stars:</strong> {pkg.stars ?? 0}</p>
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
