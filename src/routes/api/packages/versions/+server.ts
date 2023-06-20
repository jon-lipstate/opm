import { getUserId, getAuth } from '$api/auth';
import sql from '$lib/database';
import { error, json } from '@sveltejs/kit';

export async function DELETE(event) {
	const body = await event.request.json();
	const id = body.id;
	//@ts-ignore
	const { login, session } = await getAuth(event);
	try {
		const userId = await getUserId(login, session.accessToken);
		const res = await sql`DELETE from versions where published_by=${userId} AND id=${id}`;
		return json({ status: 200 });
	} catch (err: any) {
		console.error('SQL New Token Error\n', err);
		if (err.status < 500) {
			throw error(err.status, err.body.message);
		}
		throw error(500, `SQL New Token Error:, ${err}`);
	}
}
