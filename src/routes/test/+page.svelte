<script>
	import axios from 'axios';
	import './ayu-mirage.css'; // or any other style you like
	import { onMount } from 'svelte';

	let readme = '';

	onMount(async () => {
		const url = 'https://raw.githubusercontent.com/odin-lang/Odin/master/examples/demo/demo.odin';
		const res = await axios.post('/api/readme', { url });
		if (res.status != 200) {
			console.error(res.statusText);
		}
		readme = res.data.html;
	});
</script>

<div>
	{@html readme}
</div>

<style>
</style>
