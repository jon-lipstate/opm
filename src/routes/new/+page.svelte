<script lang="ts">
	import NewGist from './newGist.svelte';
	import NewPackage from './newPackage.svelte';
	import Slider from '$components/slider.svelte';
	//@ts-ignore
	//
	let useRepos = true; // gist vs repo sources

	// TODO: REQUIRED VERIFICATIONS
</script>

<div class="my-1">
	<Slider bind:checked={useRepos} whenTrue={'Source: Repositories'} whenFalse={'Source: Gists'} />
</div>
{#if useRepos}
	<NewPackage />
{:else}
	<NewGist />
{/if}
