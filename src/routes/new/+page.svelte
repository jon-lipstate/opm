<script lang="ts">
	import axios from 'axios';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import AuthRequired from '$components/authRequired.svelte';

	export let data; // self+orgs
	let showForks = false;
	$: repos = data.repos;
	$: repoNames = repos.filter((x) => showForks || !x.fork).map((x) => x.name);
	let selectedRepo: string;
	let entitySelectDisabled = false;
	// onMount(async () => {
	// 	entities = data.entities as string[];
	// 	if (data.entities.length < 2) {
	// 		entitySelectDisabled = true;
	// 		if (data.entities.length == 1) {
	// 			selectedOwner = data.entities[0];
	// 			handleSelection({ target: { value: selectedOwner } });
	// 		}
	// 	}
	// });
	// console.log(data.repos.filter((x) => x.language != 'Odin'));

	async function handleSelection(event) {
		selectedRepo = event.target.value;
		// const response = await axios.post(
		// 	'?/getPublicPackages',
		// 	{ owner: selectedOwner },
		// 	{
		// 		headers: {
		// 			'Content-Type': 'application/x-www-form-urlencoded'
		// 		}
		// 	}
		// );
		// console.log(response.data);
	}
	function handleShowForks() {
		showForks = !showForks;
		console.log(showForks);
	}
</script>

{#if $page.data.session}
	<main>
		<h1>New Package</h1>
		<select bind:value={selectedRepo} on:change={handleSelection} disabled={entitySelectDisabled}>
			<option disabled selected value="">Select a Repo</option>
			{#each repoNames as repo (repo)}
				<option value={repo}>{repo}</option>
			{/each}
		</select>
		<label>
			<input type="checkbox" bind:checked={showForks} on:click={handleShowForks} />
			Show forks
		</label>

		<pre>{JSON.stringify(repos.find((x) => x.name == selectedRepo) ?? '', null, 2)}</pre>

		<!-- 
		{#if githubRepo == null}
			<form on:submit|preventDefault={onSubmit}>
				<input type="text" name="fetchRepo" id="fetchRepo" placeholder="user/project (eg odin-lang/odin)" />
				<button type="submit">Fetch</button>
			</form>
		{:else}
			<button on:click={() => (githubRepo = null)}>Reset</button>
			<form id="repoForm" on:submit|preventDefault={onSubmit}>
				<input type="text" name="repo" id="repo" placeholder="user/project" />
				<select bind:value={branches[selectedBranch]} on:change={handleSelectionChange}>
					{#each branches as branch, index (index)}
						<option>{branch}</option>
					{/each}
				</select>
				<button type="submit">Submit</button>
			</form>
		{/if}

		<h3>New Package Generation Workflow:</h3>
		<ol>
			<li><strong>User Authentication</strong>: Require github oauth</li>
			<li>
				<strong>Display Pulldown of User/Organizations</strong>: The user selects from a dropdown of their username and
				organizations they are a part of.
			</li>
			<li>
				<strong>Display Pulldown of Projects</strong>: Depending on the selected user/organization, display a dropdown
				of available projects/repositories.
			</li>
			<li>
				<strong>Fetch Project Data</strong>: Once the user selects a project and clicks on the "Fetch" button, retrieve
				the data for the selected project.
			</li>
			<li><strong>Hydrate Package</strong>: Pre-fill from github response</li>
			<li>
				<strong>Verification</strong>: Verify important contents for package creation:
				<ul>
					<li><strong>License Existence</strong>: Check that the project has a license.</li>
					<li><strong>Readme Existence</strong>: Check that the project has a readme.</li>
					<li><strong>Submitter Rights</strong>: Validate the user is the owner / authorized collaborator.</li>
					<li><strong>Tags</strong>: Ensure tags exist and are valid.</li>
					<li><strong>Pkg File Declaration</strong>: e.g. package.json is declared and properly configured.</li>
				</ul>
			</li>
			<li>
				<strong>Form Editing</strong>: Edit Values that are mutable
				<ul>
					<li>Select Package Kind (Lib,Demo,??)</li>
					<li>License Alias? (TODO: license.key on github gives likely correct type?? eg bsd-3-clause, prefer it?)</li>
					<li>Add Topic Tags, require 1, restrict new user to use existing?</li>
				</ul>
			</li>
			<li><strong>Validation and Error Handling</strong>: Client + Server Side Verification</li>
			<li>
				<strong>Submission</strong>: Save to db, perhaps allow for pre-published state, or save for later with a 30day
				expiry?
			</li>
			<li><strong>Confirmation</strong>: Redirect to new package details page</li>
			<li><strong>Notifications</strong>: Opt in to Notifications (? future todo)</li>
		</ol> -->
	</main>
{:else}
	<AuthRequired />
{/if}

<style>
</style>
