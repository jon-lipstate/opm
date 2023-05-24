<script lang="ts">
	import { goto } from '$app/navigation';

	export let tags: string[];
</script>

<div>
    {#each tags as tag}
        <button class="tag" on:click={() => goto(`/search?tag=${encodeURIComponent(tag)}`)}>
            {tag}
        </button>
    {/each}
</div>

<style>

	.tag {
		border: none;
		background-color: var(--color-theme-3);
		color: #ddd;
		padding: 5px 10px;
		text-align: center;
		text-decoration: none;
		display: inline-block;
		font-size: 14px;
		margin: 2px 2px;
		cursor: pointer;
		border-radius: 16px;
		transition: background-color 0.3s ease;
	}

	.tag:hover {
		background-color: #bbb;
	}
</style>
