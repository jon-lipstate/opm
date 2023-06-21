import { error, json } from '@sveltejs/kit';
import sql from '$lib/database';
export type DBUser = {
	gh_login: string;
	gh_access_token: string;
	gh_avatar: URL;
	gh_id: number;
	gh_created_at: string;
	id: number;
};
export async function POST(event) {
	//@ts-ignore
	const body = await event.request.text();
	const user = JSON.parse(body) as DBUser;

	try {
		//@ts-ignore
		const rows = await sql`
  			SELECT public.upsert_user(${user.gh_id}, ${user.gh_login}, ${user.gh_access_token}, ${user.gh_avatar}, ${user.gh_created_at})
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
