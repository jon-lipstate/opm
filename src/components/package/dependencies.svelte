<script lang="ts">
	import { goto } from '$app/navigation';
	import { timeAgo } from '$lib/utils';
	import type { FlatDependencies, LicenseSummary } from 'src/routes/api/details/dependencies/+server';
	//
	function navTo(host, owner, repo) {
		goto(`/${host}/${owner}/${repo}`);
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
		<h2>License Groups</h2>
		<ul>
			{#each licenses as lic, i (lic)}
				<details>
					<summary>
						<span class="license">{lic.license}</span>
					</summary>
					<ul class="expanded-list">
						{#each lic.packages as pkg}
							<li>{pkg}</li>
						{/each}
					</ul>
				</details>
			{/each}
		</ul>
	</div>
	<!-- Display Table -->
	<div>
		<h2>Flat Dependency List</h2>
		<!-- <pre>{JSON.stringify(flat, null, 2)}</pre> -->

		<table>
			<thead>
				<tr>
					<th>Package</th>
					<th>Owner</th>
					<th>Host</th>
					<th>Version</th>
					<th>License</th>
					<th>Updated</th>
					<!-- <th>Security</th> -->
				</tr>
			</thead>
			<tbody>
				{#each flat as dep, i (i)}
					<tr>
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<td on:click={() => navTo(dep.host_name, dep.owner_name, dep.repo_name)} class="link">{dep.repo_name}</td>
						<td>{dep.owner_name}</td>
						<td>{dep.host_name}</td>
						<td>{dep.version}</td>
						<td>{dep.license}</td>
						{#if dep.state == 'archived'}
							<td class="archived">ARCHIVED</td>
						{:else}
							<td>{timeAgo(dep.last_updated)}</td>
						{/if}
						<!-- {#if dep.insecure}
							<td class="insecure-warning">INSECURE</td>
						{:else}
							<td>No Reports</td>
						{/if} -->
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>

<style>
	.expanded-list {
		padding-left: 1rem;
	}
	.link {
		color: var(--color-theme-1);
	}
	.link:hover {
		text-decoration: underline;
		cursor: pointer;
	}
	/* .primary {
		color: var(--color-theme-1);
	} */
	.license {
		color: var(--c-odin-blue-lighten-4);
	}
	.archived {
		color: orange;
		font-weight: 900;
	}
	/* .insecure-warning {
		color: red;
		font-weight: 900;
	} */
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

	@media (max-width: 650px) {
		thead {
			display: none;
		}
		table {
			width: 100%;
			table-layout: auto;
			border-spacing: 1rem 0;
		}

		tbody tr {
			/* Display rows as blocks for small screens */
			display: block;
			padding: 1rem;
			margin-bottom: 1rem;
			border: 1px solid #ccc;
			border-radius: 10px;
			box-shadow: 2px 2px 6px 0 rgba(0, 0, 0, 0.2);
		}

		td {
			display: block;
			font-size: 1.2rem;

			padding: 0.1rem;
		}

		td:before {
			content: attr(data-label);
			float: left;
			color: var(--c-odin-blue-lighten-4);
			text-transform: uppercase;
			margin-right: 0.2rem;
		}

		td:nth-of-type(2):before {
			content: 'Owner :: ';
		}
		td:nth-of-type(3):before {
			content: 'Host :: ';
		}
		td:nth-of-type(4):before {
			content: 'Version :: ';
		}
		td:nth-of-type(3):before {
			content: 'License :: ';
		}
		td:nth-of-type(5):before {
			content: 'Updated :: ';
		}
		td:nth-of-type(6):before {
			content: 'Security :: ';
		}
	}
</style>
