<script lang="ts">
	// TODO: make the gist and repo seperate pages/components?
	import axios from 'axios';
	import Slider from '$components/slider.svelte';
	import Select from '$components/select.svelte';
	import Package from '$components/package/_package.svelte';
	import { page } from '$app/stores';
	import json5 from 'json5';
	//
	export let repos = [];
	//
	let session = $page.data.session;
	let selectedIndex;
	let showForks = false;
	let missingPackageFile = false;
	let pkgFile: App.ModPkg | null = null;
	//@ts-ignore
	let details: App.PackageDetails = null;

	//@ts-ignore
	$: visibleRepos = repos.filter((x) => showForks || !x.fork);

	async function handleSelectLib() {
		missingPackageFile = false;
		selLib = visibleRepos[selectedIndex];
		//@ts-ignore
		const name = visibleRepos[selectedIndex].name;
		try {
			// const contentsResponse = await axios.post(`/api/github/getPkgFile`, { name });
			// pkgFile = json5.parse(contentsResponse.data);
			//@ts-ignore
			const readmeResponse = await axios.get(`https://api.github.com/repos/${session.user.name}/${name}/readme`);
			const download_url = readmeResponse.data.download_url;
			// details = intoPackageDetails(pkgFile!, selLib, readmeResponse.data.download_url);
			console.warn('readmeResponse', readmeResponse.data);
		} catch (e) {
			console.error(e);
			//@ts-ignore
			if (e.response.status == 404) {
				missingPackageFile = true;
			} else {
				console.error('api call error', e);
			}
		}
	}
	let selLib: any;
	async function submitPackageDetails(e) {
		e.preventDefault();
		try {
			const response = await axios.post('/api/packages/new', details);
			console.log(response.data);
		} catch (error) {
			console.error('Error submitting package details', error, details);
		}
	}
	async function subTest() {
		const data = {
			name: 'regex',
			version: '0.1.0',
			authors: ['jon-lipstate'],
			description: 'regex is a nfa regex engine',
			url: 'https://github.com/jon-lipstate/odin-regex',
			readme: 'readme.md',
			license: 'BSD-3',
			keywords: ['regex', 'nfa'],
			kind: 'Library',
			os: ['Linux', 'Windows', 'Darwin', 'Essence'],
			compiler: 'dev-2023-05',
			dependencies: {
				'jon/http-server': '1.0.0',
				'odie/async-runtime': '1.2.3'
			}
		};
		const response = await axios.post('/api/packages/new', data);
		console.log(response.data);
	}
</script>

<button on:click={subTest}>TEST</button>

<form on:submit={submitPackageDetails}>
	<div class="my-1">
		<Select
			label={'Repo'}
			bind:value={selectedIndex}
			prefix={'Select Repo'}
			options={visibleRepos}
			renderAs={(x) => x.name}
			on:change={handleSelectLib}
		/>
		<Slider bind:checked={showForks} />
		<span>{showForks ? 'Forks Shown' : 'Forks Hidden'}</span>
	</div>

	<button type="submit" disabled={selectedIndex === '' || missingPackageFile}>Submit</button>
</form>
{#if missingPackageFile}
	<div style="color:red">Missing mod.pkg file</div>
{/if}
<h3>Submission Preview</h3>
{#if details && !missingPackageFile}
	<Package {details} versionIndex={0} flat={null} licenses={null} />
	<!-- <pre>{JSON.stringify(details, null, 2)}</pre> -->
{/if}
