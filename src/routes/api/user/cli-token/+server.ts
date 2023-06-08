import { error, json } from '@sveltejs/kit';
import sql from '$lib/database';
import { getAuth } from '../../auth.js';

export async function GET(event) {
	//@ts-ignore
	const body = await event.request.text();
	const { login, authHeader, session } = await getAuth(event);

	try {
		//@ts-ignore
		const userId = await sql`
  			SELECT public.verify_token(${}, ${user.gh_login}, ${user.gh_access_token}, ${user.gh_avatar}, ${user.gh_created_at})
  		`;
		const userInfo = rows[0].upsert_user;
		if (userInfo.banned) {
			throw error(403, `Banned, Reason: ${userInfo.ban_reason}, Expire: ${userInfo.ban_timeout}`);
		}
		return json(userInfo);
	} catch (err: any) {
		if (err.status == 403) {
			console.error('Banned User:', err.body.message);
			throw error(403, err.body.message);
		}
		console.error('SQL Upsert Error\n', err);
		throw error(500, `SQL Upsert Error:, ${err}`);
	}
}
