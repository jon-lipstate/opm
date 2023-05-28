import { error, json } from '@sveltejs/kit';
import sql from '$lib/database';

export async function POST(event) {
	//@ts-ignore
	const body = await event.request.text();
	const offset = JSON.parse(body).offset ?? 0;
	const limit = JSON.parse(body).limit ?? 100;
	try {
		const count = await sql`SELECT COUNT(*) FROM packages;`;
		const results = await sql`
  			SELECT * FROM browse_packages(${limit}, ${offset})
  		`;
		return json({ ...count[0], values: results });
	} catch (err) {
		console.error('SQL Search Error', err);
		//@ts-ignore
		return error(500, { statusText: `SQL Search Error:, ${err}` });
	}
}
