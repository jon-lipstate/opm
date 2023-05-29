export async function load({ fetch, url }) {
	const offset = url.searchParams.get('offset') ?? 0;
	let results = [];
	const response = await fetch(`/api/browse`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ offset })
	});
	if (response.ok) {
		results = await response.json();
	}
	return { results };
}
