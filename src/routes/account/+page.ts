export async function load({ params, fetch, url }) {
	let tokens = [];

	let response = await fetch(`/api/user/tokens`);
	if (response.ok) {
		tokens = await response.json();
	}
	return { tokens };
}
