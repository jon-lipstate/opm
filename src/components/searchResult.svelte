<script lang="ts">
	import Tags from "$components/tags.svelte"
	import { goto } from '$app/navigation';

	export let pkg: App.PackageResult;
	let packageName = pkg.name.toLowerCase().replace(/\s/g, '-');

</script>

<div class="container">
	<h2>
		<a href={`/packages/${packageName}/${pkg.version}`}>
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
