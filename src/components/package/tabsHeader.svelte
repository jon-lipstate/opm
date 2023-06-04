<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	type Tab = {
		name: string;
		disabled?: boolean;
	};
	export let tabs: Tab[];
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
		<li
			class:selected={i === selectedTab}
			class:disabled={tab.disabled}
			on:click={() => {
				if (!tab.disabled) {
					selectTab(i);
				}
			}}
		>
			{tab.name}
		</li>
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
	li.disabled {
		color: var(--color-theme-4);
		cursor: not-allowed;
	}
	.selected {
		color: var(--c-odin-blue-lighten-3);
	}
	@media (min-width: 550px) {
		ul {
			flex-direction: row;
		}
	}
</style>
