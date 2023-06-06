import { json, error } from '@sveltejs/kit';
import sql from '$lib/database';

export async function POST(event) {
	const body = JSON.parse(await event.request.text());
	const { host, owner, repo } = body;
	try {
		const idResult = await sql`
  			SELECT * FROM get_package_id(${host},${owner}, ${repo})
  		`;
		const id = idResult[0]?.get_package_id;
		const packageResult = await sql`SELECT * FROM get_package_details(${id})`;
		const pkg = packageResult[0];

		if (pkg) {
			const versionResult = await sql`SELECT * FROM get_version_details(${id})`;
			const version = versionResult;
			pkg.versions = version ?? [];
			return json({ pkg });
		} else {
			console.error('PKG', pkg);
			throw error(404, `Package not Found: ${host}/${owner}/${repo}`);
		}
	} catch (err) {
		console.error('API/DETAILS');
		console.error('SQL Search Error\n', err);
		//@ts-ignore
		if (err.status == 404) throw err;
		else throw error(500, `SQL Search Error:, ${err}`);
	}
}
