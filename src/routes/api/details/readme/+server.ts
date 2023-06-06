import { json, error } from '@sveltejs/kit';
import sql from '$lib/database';

export async function POST(event) {
	const body = JSON.parse(await event.request.text());
	const { versionId } = body;
	try {
		let readmeResult = await sql`
  			SELECT readme FROM versions WHERE versions.id = ${versionId}
  		`;
		return json({ readme: readmeResult[0].readme });
	} catch (err) {
		console.error('SQL Readme Error\n', err);
		throw error(500, `SQL Readme Error:, ${err}`);
	}
}
