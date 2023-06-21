<script>
	import { generateRandomString } from '$lib/utils';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	//
	export let options = [];
	export let renderAs = (a) => a;
	export let value = ''; // int - leave as '' for the prefix
	export let label = '';
	export let prefix = '';
	export let disabled = false;
	//
	function handleChange(event) {
		dispatch('change', event.target.value);
	}
	function handleClick(event) {
		dispatch('click', event);
	}
	const name = `${label}_${generateRandomString(8)}`;
</script>

<label for={name}>{label}</label>
<select bind:value {name} on:change={handleChange} on:click={handleClick} {disabled}>
	{#if !!prefix}
		<option value="" disabled>-- {prefix} -- </option>
	{/if}
	{#each options as option, i}
		<option value={i}>{renderAs(option)}</option>
	{/each}
</select>

<style lang="scss">
	select {
		background-color: var(--color-bg-1);
		color: var(--color-text);
	}
</style>
