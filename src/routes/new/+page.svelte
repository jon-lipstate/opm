<script lang="ts">
	import NewGist from './newGist.svelte';
	import NewPackage from './newPackage.svelte';
	import Slider from '$components/slider.svelte';
	import { page } from '$app/stores';
	//@ts-ignore
	//
	let useRepos = true; // gist vs repo sources
	export let repos = $page.data.repos;
	console.warn();
	// TODO: REQUIRED VERIFICATIONS
	// https://github.com/rust-lang/cargo/blob/58b22f00ef76667ec4f93f4f7936831eaf63f5b7/src/cargo/util/restricted_names.rs
</script>

<div class="my-1">
	<Slider bind:checked={useRepos} whenTrue={'Source: Repositories'} whenFalse={'Source: Gists'} />
</div>
<p style="color:red">Submittals not implemented at present.</p>
<pre>{JSON.stringify(repos, null, 2)}</pre>
{#if useRepos}
	<NewPackage />
{:else}
	<NewGist />
{/if}
