<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	export let tabs: string[];
	export let selectedTab: number = 0;
	const dispatch = createEventDispatcher();

	function selectTab(tabIndex: number) {
		selectedTab = tabIndex;
		dispatch('select', tabIndex);
	}
</script>

<ul>
	{#each tabs as tab, i (i)}
		<!-- svelte-ignore a11y-click-events-have-key-events -->
		<li class:selected={i === selectedTab} on:click={() => selectTab(i)}>{tab}</li>
	{/each}
</ul>

<style>
	ul {
		display: flex;
		justify-content: space-around;
		width: 100%;
		background-color: rgba(0, 0, 0, 0.1);
		margin: 0;
		padding-left: 0;
		margin-bottom: 1rem;
		flex-direction: column;
	}

	li {
		display: list-item;
		list-style-type: none;
		font-size: 1.5rem;
		cursor: pointer;
		text-align: center;
	}
	.selected {
		color: var(--color-theme-1);
	}
	@media (min-width: 550px) {
		ul {
			flex-direction: row;
		}
	}
</style>
