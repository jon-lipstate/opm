<script lang="ts">
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
	<div>
		{#each pkg.tags as tag}
			<button class="tag" on:click={() => goto(`/search?tag=${encodeURIComponent(tag)}`)}>
				{tag}
			</button>
		{/each}
	</div>
</div>

<style>
	.container {
		border: 1px solid #ddd;
		padding: 10px;
		margin: 10px 0;
	}
	.tag {
		border: none;
		background-color: var(--color-theme-2);
		color: #ddd;
		padding: 5px 10px;
		text-align: center;
		text-decoration: none;
		display: inline-block;
		font-size: 14px;
		margin: 2px 2px;
		cursor: pointer;
		border-radius: 16px;
		transition: background-color 0.3s ease;
	}

	.tag:hover {
		background-color: #bbb;
	}
</style>
