export async function load({ params, fetch, url }) {
	let owner = params.owner;
	let slug = params.slug;

	let flat;
	let licenses;
	let details;
	let error = '';

	let response = await fetch(`/api/details`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ owner, slug })
	});
	if (response.ok) {
		details = (await response.json()).pkg;
		// console.warn(details);--c-odin-blue-lighten-3
	} else {
		error = (await response.json()).message;
	}

	const versionTag = url.searchParams.get('version');

	// default: newest index, which has highed id-key
	let versionIndex = details.versions.reduce((max, celm, cidx) => {
		return celm.id > details.versions[max].id ? cidx : max;
	}, 0);
	if (versionTag) {
		const qIndex = details.versions.findIndex((x) => x.version == versionTag);
		if (qIndex == -1) {
			error += `\nInvalid Version on query-string: ${versionTag}`;
			console.error(error);
		} else {
			versionIndex = qIndex;
		}
	}
	const versionId = details.versions[versionIndex].id;
	response = await fetch(`/api/details/dependencies`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ versionId })
	});
	if (response.ok) {
		const data = await response.json();
		flat = data.flat;
		licenses = data.licenses;
	} else {
		error += ' ' + (await response.json()).message;
	}
	return { details, flat, licenses, error, versionIndex };
}
