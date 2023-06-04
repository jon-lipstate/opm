import { error, json } from '@sveltejs/kit';
import sql from '$lib/database';

export async function POST(event) {
	//@ts-ignore
	const body = await event.request.text();

	const query = JSON.parse(body).query;
	const offset = JSON.parse(body).offset ?? 0;
	const limit = JSON.parse(body).offset ?? 25;

	try {
		const res = await sql`
  			SELECT * FROM search_and_get_results(${query}, ${limit}, ${offset})
  		`;
		return json(res);
	} catch (err) {
		console.error('SEARCH');
		console.error('SQL Search Error\n', err);
		throw error(500, `SQL Search Error:, ${err}`);
	}
}
