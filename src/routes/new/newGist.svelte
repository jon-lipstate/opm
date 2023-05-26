<script lang="ts">
	import axios from 'axios';
	import Slider from '$components/slider.svelte';
	import Select from '$components/select.svelte';
	import Gist from '$components/gist.svelte';
	import { onMount } from 'svelte';
	//@ts-ignore
	//
	let selectedIndex;
	//
	let gists = [];

	onMount(async () => {
		try {
			const gistRes = await axios.get(`/api/github/getPublicGists`);
			//@ts-ignore
			gists = gistRes.data;
		} catch (e) {
			console.error('api call error', e);
		}
	});

	async function handleSelectLib() {
		selGist = gists[selectedIndex];
	}
	let selGist: any;

	let gistAttrs = {
		readme: '',
		isBindings: false,
		isFreestanding: false,
		isOsSpecific: false,
		isWindows: false,
		isLinux: false,
		isDarwin: false
	};
	// TODO: REQUIRED VERIFICATIONS
</script>

<form action="newGist">
	<div class="my-1">
		<Select
			label={'Gist'}
			bind:value={selectedIndex}
			prefix={'Select Gist'}
			options={gists}
			renderAs={(x) => x.description}
			on:change={handleSelectLib}
		/>
	</div>

	<div class=" my-1">
		<h5>Gist Meta Data</h5>
		<div>
			<textarea bind:value={gistAttrs.readme} cols="60" rows="4" placeholder="Readme for the Gist" />
		</div>
		<Slider bind:checked={gistAttrs.isBindings} disabled={selectedIndex === ''} />
		<span>Bindings</span>

		<Slider bind:checked={gistAttrs.isFreestanding} disabled={selectedIndex === ''} />
		<span>Freestanding</span>

		<Slider bind:checked={gistAttrs.isOsSpecific} disabled={selectedIndex === ''} />
		<span>OS-specific</span>
		<div>
			{#if gistAttrs.isOsSpecific}
				<Slider bind:checked={gistAttrs.isWindows} />
				<span>Windows</span>

				<Slider bind:checked={gistAttrs.isLinux} />
				<span>Linux</span>

				<Slider bind:checked={gistAttrs.isDarwin} />
				<span>Darwin</span>
			{/if}
		</div>
	</div>
	<button type="submit" disabled={selectedIndex === ''}>Submit</button>
</form>

<h3>Submission Preview</h3>
{#if selGist}
	<Gist gist={selGist} attrs={gistAttrs} />
	<pre>{JSON.stringify(selGist, null, 2)}</pre>
{/if}

<style>
</style>
