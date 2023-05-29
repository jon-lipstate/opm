export async function load({ fetch, url }) {
	const query = url.searchParams.get('query');
	let results = [];
	if (query) {
		const response = await fetch(`/api/search`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ query })
		});
		if (response.ok) {
			results = await response.json();
		}
	}
	return { results };
}
