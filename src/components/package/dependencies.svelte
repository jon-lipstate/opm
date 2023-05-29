<script lang="ts">
	import { timeAgo } from '$lib/utils';
	import type { FlatDependencies, LicenseSummary } from 'src/routes/api/details/dependencies/+server';
	import { onMount } from 'svelte';
	//
	export let versionId;
	let licenses: LicenseSummary[] = [];
	let flat: FlatDependencies[] = [];
	let error;
	//
	onMount(async () => {
		const response = await fetch(`/api/details/dependencies`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ versionId })
		});
		if (response.ok) {
			const pageData = await response.json();
			flat = pageData.flat;
			licenses = pageData.licenses;
		} else {
			error = (await response.json()).message;
		}
	});
</script>

<section>
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
						<td><a href="/packages/{dep.owner}/{dep.slug}">{dep.package_name}</a></td>
						<td>{dep.version}</td>
						<td>{dep.license}</td>
						{#if dep.archived}
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
	.primary {
		color: var(--color-theme-1);
	}
	.license {
		color: greenyellow;
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
