import { error, fail, json } from '@sveltejs/kit';
import sql from '$lib/database';

export async function POST(event) {
	//@ts-ignore
	const body = await event.request.text();

	const query = JSON.parse(body).query;
	const offset = JSON.parse(body).offset ?? 0;
	const limit = JSON.parse(body).offset ?? 25;

	try {
		const res = await sql`
  			SELECT * FROM search_and_get_details(${query}, ${limit}, ${offset})
  		`;
		return json(res);
	} catch (err) {
		console.error('SQL Search Error', err);
		//@ts-ignore
		return error(500, { statusText: `SQL Search Error:, ${err}` });
	}
}
