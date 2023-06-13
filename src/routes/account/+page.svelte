<script>
	import { page } from '$app/stores';
	import { timeAgo } from '$lib/utils';
	import NewToken from './newToken.svelte';
	export let tokens = $page.data.tokens ?? [];
	let newTokenValue = '';
	// Flag to control visibility of NewToken component
	let isAddingToken = false;

	async function handleTokenGenerated(token) {
		isAddingToken = false;
		newTokenValue = token.detail;
		tokens = await refreshTokens();
	}
	// Function to handle button click
	function handleAddTokenClick() {
		isAddingToken = true;
	}
	async function refreshTokens() {
		let response = await fetch(`/api/user/tokens`);
		let t = [];
		if (response.ok) {
			t = await response.json();
		}
		return t;
	}
	async function deleteToken(id, index) {
		const response = await fetch('/api/user/tokens', {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ id }),
			credentials: 'include'
		});

		if (response.ok) {
			tokens.slice(index, 1);
			tokens = await refreshTokens();
		}
	}
</script>

{#if !isAddingToken}
	<div>
		{#if newTokenValue.length > 0}
			<h2>New Token Value:</h2>
			<h3 class="color-secondary">{newTokenValue}</h3>
			<span class="warning">Warning: Copy immediately, this value will never be shown again.</span>
		{/if}
		<h2>Existing Tokens</h2>
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Used</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each tokens as token, index (token.name)}
					<tr>
						<td>{token.name}</td>
						<td>{!token.last_touched ? 'Never' : timeAgo(token.last_touched)}</td>
						<td>{timeAgo(token.created_at)}</td>
						<td><button on:click={() => deleteToken(token.id, index)}>Delete</button></td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<!-- Show the "Add Token" button only if isAddingToken is false -->
	<button class="addToken" on:click={handleAddTokenClick}>Add Token</button>
{/if}

{#if isAddingToken}
	<div>
		<h2>Add New Token</h2>
		<NewToken on:tokenGenerated={handleTokenGenerated} on:cancel={() => (isAddingToken = false)} />
	</div>
{/if}

<style>
	h2 {
		margin-bottom: 0.25rem;
	}
	.warning {
		color: rgba(200, 0, 0, 0.9);
	}
	h3 {
		padding: 0.5rem;
		margin: 0;
		background-color: gray;
		text-align: center;
	}
	table {
		max-width: 90vw;
		width: 100%;
		text-align: center;
		background-color: var(--c-odin-blue-darken-7);
	}
	button {
		width: 5rem;
	}
	.addToken {
		margin-top: 0.5rem;
	}
</style>
