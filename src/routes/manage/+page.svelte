<script>
	import { page } from '$app/stores';
	import Slider from '$components/slider.svelte';
	export let pkgs = $page.data.pkgs ?? [];
	async function deleteVersion(id) {
		console.warn('del id', id);
	}
	async function deletePackage(id) {
		const response = await fetch('/api/packages', {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ id }),
			credentials: 'include'
		});

		if (response.ok) {
			pkgs = await refreshPackages();
		}
	}
	async function refreshPackages() {
		let p = [];
		let response = await fetch(`/api/packages`);
		if (response.ok) {
			p = await response.json();
		}
		return p;
	}
	let enableDeletion = false;
</script>

<main>
	<div>
		<h2>Packages</h2>
		<ul>
			{#each pkgs as pkg, i (pkg.id)}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<details>
					<summary on:click={() => (enableDeletion = false)}>
						<span class="license">{pkg.host_name} / {pkg.owner_name} / {pkg.repo_name}</span>
					</summary>
					<ul class="expanded-list my-4">
						<td>
							<button on:click={() => deletePackage(pkg.id)} disabled={!enableDeletion}>Delete Package</button>
							<Slider bind:checked={enableDeletion} label="Enable Package Deletion" />
							<!-- todo: either place button+slider in component, or find non janky way to fix the paired unlock of delete btn -->
						</td>

						<table>
							<thead>
								<tr>
									<th>Version</th>
									<th>license</th>
									<!-- <th>insecure</th> -->
									<th>Actions</th>
								</tr>
							</thead>
							<tbody>
								{#each pkg.versions as v}
									<tr>
										<td>{v.version}</td>
										<td>{v.license}</td>
										<!-- <td>{v.insecure}</td> -->
										<td title={pkg.versions.length < 2 ? 'Cannot Delete only version of a package' : ''}>
											<button on:click={() => deleteVersion(v.id)} disabled={pkg.versions.length < 2}>Delete</button>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</ul>
				</details>
			{/each}
		</ul>
	</div>
</main>

<style>
	.my-4 {
		margin-top: 1rem;
		margin-bottom: 1rem;
	}
	h2 {
		margin-bottom: 0.25rem;
	}

	table {
		max-width: 90vw;
		width: 100%;
		text-align: center;
		background-color: var(--c-odin-blue-darken-7);
	}
</style>
