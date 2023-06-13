<script>
	import { createEventDispatcher } from 'svelte';
	import Slider from '$components/slider.svelte';

	const dispatch = createEventDispatcher();
	let tokenName = '';
	let publishScope = false;
	let fullScope = true;
	let updateScope = false;
	let authorsScope = false;
	$: publishScope,
		() => {
			if (publishScope) {
				updateScope = true;
			}
		};
	$: scopes = ['publish', 'update', 'invite'];
	$: readyToGenerate = tokenName.length > 0;
	// TODO: update slider to dispatch event so i can set update and disable properly
	async function generateToken() {
		const response = await fetch('/api/user/tokens', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				tokenName,
				scopes
			}),
			credentials: 'include'
		});

		if (response.ok) {
			const token = await response.json();
			dispatch('tokenGenerated', token.token);
		} else {
			console.error('Failed to generate token');
		}
	}
</script>

<div>
	<label>
		Token Name:
		<input type="text" bind:value={tokenName} />
	</label>
	<h4>Scopes:</h4>
	<div title="Only full rights are granted at present.">
		<Slider bind:checked={fullScope} label={'Full Scopes'} disabled={true} />
	</div>
	{#if !fullScope}
		<div>
			<Slider bind:checked={publishScope} label={'Publish'} />
		</div>
		<div>
			<Slider checked={updateScope} label={'Update'} disabled={publishScope} />
		</div>
		<div title="Feature Not Implemented">
			<Slider bind:checked={authorsScope} label={'Invite Collaborators'} disabled={true} />
		</div>
	{/if}
	<div class="buttons">
		<span title={readyToGenerate ? '' : 'Assign a Token Name'}>
			<button on:click={generateToken} disabled={!readyToGenerate}>Generate Token</button>
		</span>
		<button on:click={() => dispatch('cancel')}>Cancel</button>
	</div>
</div>

<style>
	h4 {
		margin: 0.2rem 0;
	}
	.buttons {
		margin: 0.2rem 0;
	}
</style>
