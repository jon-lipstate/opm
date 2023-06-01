export async function load(pageLoadEvent) {
	console.warn('PAGE-LOAD');
	// let session = await pageLoadEvent.parent();
	let fetch = pageLoadEvent.fetch;
	const repoRes = await fetch(`/api/github/getPublicRepos`);
	const repos = await repoRes.json();

	return repos;
}
