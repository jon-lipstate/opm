<script>
	import axios from 'axios';
	import '$lib/ayu-mirage.css'; // or any other style you like
	import { afterUpdate, onMount } from 'svelte';
	import odin from '$lib/odin-hl.js';
	import hljs from 'highlight.js/lib/core';
	hljs.registerLanguage('odin', odin);

	let readme = '';

	onMount(async () => {
		const url = 'https://raw.githubusercontent.com/odin-lang/Odin/master/examples/demo/demo.odin';
		const res = await axios.post('/api/readme', { url });
		if (res.status != 200) {
			console.error(res.statusText);
		}
		readme = res.data.html;
	});
	afterUpdate(() => {
		const code = document.querySelectorAll('code');
		// debbounce updates with no code blocks
		if (code.length > 0) {
			hljs.highlightAll();
		}
	});
</script>

<div>
	{@html readme}
</div>

<style>
</style>
