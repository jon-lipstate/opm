<script>
	import axios from 'axios';
	import Slider from '$components/slider.svelte';
	import Select from '$components/select.svelte';
	//
	let attributes = {
		bindings: false,
		freestanding: false,
		stable: false,
		//
		osSpecific: false,
		windows: false,
		darwin: false,
		linux: false
	};

	const submissionKinds = ['Demo / Prototype', 'Library'];

	let useRepos = true; // gist vs repo sources
	let lockRepos = false;
	let selectedKind;
	let selectedLib;
	//
	let libraries = [];
	let showForks = false;
	$: visibleLibraries = libraries.filter((x) => showForks || !x.fork);

	async function handleSelectKind() {
		if (selectedKind > 0) {
			useRepos = true;
			lockRepos = true;
		} else {
			lockRepos = false;
		}
		try {
			const response = await axios.get(`/api/github/${useRepos ? 'getPublicRepos' : 'getPublicGists'}`);
			libraries = [...response.data];
			console.warn(libraries);
		} catch (e) {
			console.error('api call error', e);
		}
	}
	let pkg = null;
	// TODO: PACKAGE PREVIEW
	// TODO: REQUIRED VERIFICATIONS
</script>

<form>
	<div class="my-1">
		<Slider
			bind:checked={useRepos}
			whenTrue={'Source: Repositories'}
			whenFalse={'Source: Gists'}
			disabled={lockRepos}
		/>
	</div>
	<div class="my-1">
		<Select
			label={'Submission Kind:'}
			bind:value={selectedKind}
			prefix={'Select Kind'}
			options={submissionKinds}
			on:change={handleSelectKind}
		/>
	</div>
	<div class="my-1">
		<Select
			disabled={libraries.length == 0}
			label={useRepos ? 'Repo' : 'Gist'}
			bind:value={selectedLib}
			prefix={useRepos ? 'Select Repo' : 'Select Gist'}
			options={visibleLibraries}
			renderAs={useRepos ? (x) => x.name : (x) => x.description}
		/>
	</div>
	<div class=" my-1">
		{#if !!selectedKind}
			<Slider bind:checked={attributes.stable} disabled={!selectedLib} />
			<span>Stable Library</span>
		{/if}

		<Slider bind:checked={attributes.bindings} disabled={!selectedLib} />
		<span>Bindings</span>

		<Slider bind:checked={attributes.freestanding} disabled={!selectedLib} />
		<span>Freestanding</span>

		<Slider bind:checked={attributes.osSpecific} disabled={!selectedLib} />
		<span>OS-specific</span>
		<div>
			{#if attributes.osSpecific}
				<Slider bind:checked={attributes.windows} />
				<span>Windows</span>

				<Slider bind:checked={attributes.linux} />
				<span>Linux</span>

				<Slider bind:checked={attributes.darwin} />
				<span>Darwin</span>
			{/if}
		</div>
	</div>

	<button type="submit" disabled={!selectedLib}>Submit</button>
</form>

<style>
	.my-1 {
		margin-top: 0.25rem;
		margin-bottom: 0.25rem;
	}
	/* .block {
		display: block;
	} */
</style>
