<script lang="ts">
	import { afterUpdate, onMount } from 'svelte';
	import '$lib/ayu-mirage.css'; // or any other style you like
	import odin from '$lib/odin-hl.js';
	import hljs from 'highlight.js/lib/core';
	hljs.registerLanguage('odin', odin);

	export let readme: string;
	/* !! TODO: SANITIZE THE README !!*/
	afterUpdate(() => {
		const code = document.querySelectorAll('code');
		// debbounce updates with no code blocks
		if (code.length > 0) {
			hljs.highlightAll();
		}
	});
</script>

<section>
	{@html readme}
</section>

<style>
	section {
		max-width: 90vw;
		overflow: auto;
	}
</style>
