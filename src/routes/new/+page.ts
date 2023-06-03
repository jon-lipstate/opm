export async function load(pageLoadEvent) {
	// let session = await pageLoadEvent.parent();
	let fetch = pageLoadEvent.fetch;
	const repoRes = await fetch(`/api/github/getPublicRepos`);
	let repos = await repoRes.json();

	return repos;
}
