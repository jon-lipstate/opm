import { error, json } from '@sveltejs/kit';
import sql from '$lib/database';
import { getAuth, getUserId } from '$api/auth.js';
import { generateRandomString } from '$lib/utils.js';

export async function GET(event) {
	//@ts-ignore
	const { login, session } = await getAuth(event);
	let tokens;
	try {
		const userId = await getUserId(login, session.accessToken);
		tokens = await sql`SELECT id, name, created_at, last_touched, revoked FROM api_tokens where user_id=${userId}`;
		return json(tokens);
	} catch (err: any) {
		console.error('SQL New Token Error\n', err);
		if (err.status < 500) {
			throw error(err.status, err.body.message);
		}
		throw error(500, `SQL New Token Error:, ${err}`);
	}
}

export async function POST(event) {
	const body = await event.request.json();
	const { login, session } = await getAuth(event);
	const tokenName = body.tokenName;
	const tokenValue = generateRandomString(32);
	const tokenScopes = body.scopes; // todo: configure scopes
	if (!tokenName) {
		throw error(400, 'No Token Name Supplied');
	}
	try {
		const userId = await getUserId(login, session.accessToken);
		await sql`SELECT * FROM insert_token(${userId},${tokenName},${tokenValue})`;
		return json({ token: tokenValue });
	} catch (err: any) {
		console.error('SQL New Token Error\n', err);
		if (err.status < 500) {
			throw error(err.status, err.body.message);
		}
		throw error(500, `SQL New Token Error:, ${err}`);
	}
}

export async function DELETE(event) {
	const body = await event.request.json();
	const id = body.id;
	//@ts-ignore
	const { login, session } = await getAuth(event);
	try {
		const userId = await getUserId(login, session.accessToken);
		await sql`DELETE from api_tokens where user_id=${userId} AND id=${id}`;
		return json({ status: 200 });
	} catch (err: any) {
		console.error('SQL New Token Error\n', err);
		if (err.status < 500) {
			throw error(err.status, err.body.message);
		}
		throw error(500, `SQL New Token Error:, ${err}`);
	}
}
