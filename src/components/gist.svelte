<script lang="ts">
	import Tags from '$components/tags.svelte';
	import { timeAgo } from '$lib/utils';
	export let gist: any;
	export let attrs: gistAttrs;
	$: files = Object.keys(gist.files).filter((x) => !!x);
	let tags: string[] = [];

	$: attrs, setTags();

	function setTags() {
		tags = [];
		if (attrs.isBindings) tags.push('Bindings');
		if (attrs.isFreestanding) tags.push('No-Std');
		if (attrs.isOsSpecific) {
			if (attrs.isDarwin) tags.push('Darwin');
			if (attrs.isLinux) tags.push('Linux');
			if (attrs.isWindows) tags.push('Windows');
		}
	}
	type gistAttrs = {
		readme: string;
		isBindings: boolean;
		isFreestanding: boolean;
		isOsSpecific: boolean;
		isWindows: boolean;
		isLinux: boolean;
		isDarwin: boolean;
	};
</script>

<div>
	<a href={gist.html_url}>{gist.description}</a>
	<div>Updated: {timeAgo(gist.updated_at)}</div>
	<!-- <pre>{JSON.stringify(attrs, null, 2)}</pre> -->
	<Tags {tags} />
	<div>
		<h5>Readme:</h5>
		<div class="readme">{attrs.readme}</div>
	</div>
	<div>
		Files:
		<ul>
			{#each files as filename}
				<li class="file">
					<a href={gist.files[filename].raw_url}>{filename}</a> ({gist.files[filename].size} bytes)
				</li>
			{/each}
		</ul>
	</div>
</div>

<style>
	div,
	h5,
	ul {
		margin: 0.1rem;
	}
	ul,
	.readme {
		padding-left: 1rem;
	}
</style>
