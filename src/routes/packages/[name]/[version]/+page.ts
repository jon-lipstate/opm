// export async function load({ fetch, url }) {
export async function load({ params, fetch }) {
	let name = params.name;
	let version = params.version;
	// empty details:
	let details = {
		name: '',
		version: '',
		description: '',
		tags: [],
		versions: [],
		funding: [],
		dependsOn: [],
		usedBy: [],
		requirements: { minCompilierVersion: '' },
		links: { repo: '', discord: '' },
		lastUpdated: '',
		license: '',
		size: '',
		kind: '',
		owners: [{ name: '', username: '' }],
		stats: { allTimeDownloads: 0 },
		readme: ''
	};

	// Fetch data from the API
	const response = await fetch(`/api/details`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name, version })
	});
	if (response.ok) {
		details = await response.json();
	} else {
		console.error(response.status, response.statusText);
	}
	return { details };
}
/*
{
	url: URL {
	  href: 'http://localhost:5173/packages/http-server/1.0.0',
	  origin: 'http://localhost:5173',
	  protocol: 'http:',
	  username: '',
	  password: '',
	  host: 'localhost:5173',
	  hostname: 'localhost',
	  port: '5173',
	  pathname: '/packages/http-server/1.0.0',
	  search: '',
	  searchParams: URLSearchParams {},
	  hash: ''
	},
	params: { name: 'http-server', version: '1.0.0' },
	data: null,
	route: { id: '/packages/[name]/[version]' },
	fetch: [AsyncFunction (anonymous)],
	setHeaders: [Function: setHeaders],
	depends: [Function: depends],
	parent: [AsyncFunction: parent]
  }
  */
