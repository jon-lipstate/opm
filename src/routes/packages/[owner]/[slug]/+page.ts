// export async function load({ fetch, url }) {
export async function load({ params, fetch }) {
	let owner = params.owner;
	let slug = params.slug;

	let details;
	let error;

	const response = await fetch(`/api/details`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ owner, slug })
	});
	if (response.ok) {
		details = (await response.json()).pkg;
	} else {
		error = (await response.json()).message;
	}
	return { details, error };
}
