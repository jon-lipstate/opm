<script lang="ts">
	export let details: App.PackageDetails;
	// TODO: Scan all deps to get a Set of Licenses. Add some Rules on GPL and Missing to turn red?
</script>

<section>
	<!-- License Summary -->
	<div>
		<h3>Licenses</h3>
		<span style="color:red">todo: Traverse the graph, summarize licenses as a Set</span>
		<ul>
			<li>
				<span>MIT : Something, Another</span>
			</li>
			<li>
				<span>GPL : SOME_THING_ELSE</span>
			</li>
		</ul>
	</div>
	<!-- Display Table -->
	<div>
		<h3>Dependency Summary</h3>
		<table>
			<thead>
				<tr>
					<th>Package</th>
					<th>Version</th>
					<th>License</th>
					<th>Depends On</th>
					<th>Active</th>
				</tr>
			</thead>
			<tbody>
				{#each Object.keys(details.dependsOn) as dep, i (i)}
					<tr>
						<td>{dep}</td>
						<td>{details.dependsOn[dep]} (@todo)</td>
						<td>??MIT</td>
						<td><a href="/{dep}/{details.dependsOn[dep]}">#?</a></td>
						<td>(x days ago) | archived</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>

<style>
	h3 {
		color: var(--color-theme-1);
		margin-bottom: 0.2rem;
	}
	table {
		width: 100%;
		text-align: left;
	}
</style>
