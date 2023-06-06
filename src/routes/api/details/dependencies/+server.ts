import { json, error } from '@sveltejs/kit';
import sql from '$lib/database';

export type LicenseSummary = {
	license: string;
	packages: string[];
};
export type FlatDependencies = {
	host_name: string;
	owner_name: string;
	repo_name: string;
	version: string;
	license: string;
	last_updated: string;
	state: string;
	insecure: boolean;
};

export async function POST(event) {
	const body = JSON.parse(await event.request.text());
	const { versionId } = body;
	let licenseSummary: LicenseSummary[];
	let flat: FlatDependencies[];
	try {
		licenseSummary = await sql`
  			SELECT * FROM get_all_dependency_licenses(${versionId})
  		`;
		flat = await sql`
        SELECT * FROM get_dependencies_flat(${versionId})
        `;
		return json({ flat, licenses: licenseSummary });
	} catch (err) {
		console.error('API/DEPENDENCIES/DETAILS');
		console.error('SQL Search Error\n', err);
		//@ts-ignore
		if (err.status == 404) throw err;
		else throw error(500, `SQL Search Error:, ${err}`);
	}
}
