<script lang="ts">
	import NewGist from './newGist.svelte';
	import NewPackage from './newPackage.svelte';
	import Slider from '$components/slider.svelte';
	import { page } from '$app/stores';
	//@ts-ignore
	//
	let useRepos = true; // gist vs repo sources
	export let repos = $page.data.repos;
	// TODO: REQUIRED VERIFICATIONS
	// https://github.com/rust-lang/cargo/blob/58b22f00ef76667ec4f93f4f7936831eaf63f5b7/src/cargo/util/restricted_names.rs
	/*
	  {
    "name": "pico",
    "owner": "jon-lipstate",
    "html_url": "https://github.com/jon-lipstate/pico",
    "ssh_url": "git@github.com:jon-lipstate/pico.git",
    "description": "a tiny text editor",
    "fork": false,
    "created_at": "2023-05-06T19:42:56Z",
    "updated_at": "2023-06-01T15:24:38Z",
    "homepage": "",
    "language": "Odin",
    "size_kb": 22,
    "archived": false,
    "license_id": "BSD-3-Clause",
    "license_name": "BSD 3-Clause \"New\" or \"Revised\" License",
    "visibility": "public",
    "default_branch": "master"
  },
	*/
</script>

<div class="my-1">
	<Slider bind:checked={useRepos} whenTrue={'Source: Repositories'} whenFalse={'Source: Gists'} />
</div>
<p style="color:red">Submittals not implemented at present.</p>
<!-- <pre>{JSON.stringify(repos, null, 2)}</pre> -->
{#if useRepos}
	<NewPackage {repos} />
{:else}
	<NewGist />
{/if}
