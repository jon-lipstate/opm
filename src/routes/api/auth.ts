import sql from '$lib/database';
import { error } from '@sveltejs/kit';

export async function getAuth(event): Promise<auth> {
	const session = await event.locals.getSession();
	//@ts-ignore
	const login = session?.user.name;
	const authHeader = {
		headers: {
			//@ts-ignore
			Authorization: `token ${session?.accessToken}`
		}
	};
	return {
		login,
		authHeader,
		session
	};
}
// TODO: move this to be part of the user's auth token...
export async function getUserId(userName, authToken): Promise<number> {
	try {
		const userRes = await sql`
			SELECT id FROM users
			WHERE gh_login = ${userName}
			AND gh_access_token = ${authToken}
		`;
		return userRes[0].id;
	} catch (err) {
		console.error('User Validation Error\n', err);
		throw error(401, `User Auth Error:, ${err}`);
	}
}
