export async function load({ params, fetch, url }) {
	let pkgs = [];

	let response = await fetch(`/api/packages`);
	if (response.ok) {
		pkgs = await response.json();
	}
	return { pkgs };
}
