<script>
	export let checked = true;
	export let whenTrue = '';
	export let whenFalse = '';
	export let label = '';
	export let disabled = false;
</script>

<label class="switch">
	<input type="checkbox" bind:checked on:click={() => (checked = !checked)} {disabled} />
	<span class="slider" />
</label>
{#if !!label}
	<span>{label}</span>
{/if}
<span>{checked ? whenTrue : whenFalse}</span>

<style>
	.switch {
		position: relative;
		display: inline-block;
		width: 30px;
		height: 15px;
	}

	.switch input {
		opacity: 0;
		width: 0;
		height: 0;
	}

	.slider {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: var(--color-theme-4);
		-webkit-transition: 0.4s;
		transition: 0.4s;
	}

	.slider:before {
		position: absolute;
		content: '';
		height: 12px;
		width: 12px;
		left: 2px;
		bottom: 1.5px;
		background-color: var(--color-text);
		transition: 0.4s;
	}
	input:disabled + .slider {
		background-color: #777 !important;
	}
	input:disabled + .slider:before {
		background-color: #aaa !important;
	}
	input:checked + .slider {
		background-color: var(--color-theme-3);
		/* #2196f3; */
	}

	input:focus + .slider {
		box-shadow: 0 0 1px #2196f3;
	}

	input:checked + .slider:before {
		transform: translateX(14px);
	}
</style>
