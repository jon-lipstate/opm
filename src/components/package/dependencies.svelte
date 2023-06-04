<script lang="ts">
	import { goto } from '$app/navigation';
	import { timeAgo } from '$lib/utils';
	import type { FlatDependencies, LicenseSummary } from 'src/routes/api/details/dependencies/+server';
	//
	function navTo(owner, slug) {
		// pageParamsStore.set({ owner, slug });
		goto(`/packages/${owner}/${slug}`);
	}
	//
	export let licenses: LicenseSummary[] = [];
	export let flat: FlatDependencies[] = [];
</script>

<section>
	<!-- <pre>{JSON.stringify(licenses, null, 2)}</pre>
	<pre>{JSON.stringify(flat, null, 2)}</pre> -->
	<!-- License Summary -->
	<div>
		<h2>Licenses</h2>
		<ul>
			{#each licenses as lic, i (lic)}
				<li>
					<span class="license">{lic.license} </span>
					<span class="primary">:: </span>
					<span> {lic.packages.join(', ')}</span>
				</li>
			{/each}
		</ul>
	</div>
	<!-- Display Table -->
	<div>
		<h2>Dependency Summary</h2>
		<!-- <pre>{JSON.stringify(flat, null, 2)}</pre> -->

		<table>
			<thead>
				<tr>
					<th>Package</th>
					<th>Version</th>
					<th>License</th>
					<th>Updated</th>
					<th>Security</th>
				</tr>
			</thead>
			<tbody>
				{#each flat as dep, i (i)}
					<tr>
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<td on:click={() => navTo(dep.owner, dep.slug)} class="link">{dep.package_name}</td>
						<td>{dep.version}</td>
						<td>{dep.license}</td>
						{#if dep.state == 'archived'}
							<td class="archived">ARCHIVED</td>
						{:else}
							<td>{timeAgo(dep.last_updated)}</td>
						{/if}
						{#if dep.insecure}
							<td class="insecure-warning">INSECURE</td>
						{:else}
							<td />
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>

<style>
	.link {
		color: var(--color-theme-1);
	}
	.link:hover {
		text-decoration: underline;
		cursor: pointer;
	}
	.primary {
		color: var(--color-theme-1);
	}
	.license {
		color: var(--c-odin-blue-lighten-4);
	}
	.archived {
		color: orange;
		font-weight: 900;
	}
	.insecure-warning {
		color: red;
		font-weight: 900;
	}
	h2 {
		/* color: var(--color-theme-1); */
		margin-bottom: 0.2rem;
	}
	table {
		width: 100%;
		text-align: left;
	}
	ul {
		margin: 0;
		padding: 0;
	}
	li {
		margin: 0;
		list-style-type: none;
	}
</style>
