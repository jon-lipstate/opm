<script lang="ts">
	import { goto } from '$app/navigation';

	export let name: string;
</script>

<div style="display:inline">
	<button class="nametag" on:click={() => goto(`https://github.com/${name}`)}>
		{name}
	</button>
</div>

<style>
	.nametag {
		border: none;
		background-color: rgba(220, 22, 0, 0.839);
		color: var(--c-odin-orange-lighten-8);
		padding: 3px 7px;
		text-align: center;
		text-decoration: none;
		display: inline-block;
		font-size: 0.7rem;
		margin: 2px 2px;
		cursor: pointer;
		border-radius: 16px;
		transition: background-color 0.3s ease;
	}

	.nametag:hover {
		background-color: var(--c-odin-orange-darken-1);
	}
</style>
