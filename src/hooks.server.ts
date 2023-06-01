import { SvelteKitAuth } from '@auth/sveltekit';
import GitHub from '@auth/core/providers/github';
import { GITHUB_ID, GITHUB_SECRET, AUTH_SECRET } from '$env/static/private';
import axios from 'axios';
import type { DBUser } from './routes/api/user/+server';
//
let eventFetch;

const svaData = {
	providers: [
		//@ts-ignore
		GitHub({
			clientId: GITHUB_ID,
			clientSecret: GITHUB_SECRET,
			authorization: {
				params: { scope: 'read:user' }
			}
		})
	],
	secret: AUTH_SECRET,
	trustHost: true,
	callbacks: {
		async jwt({ token, account }) {
			if (account) {
				// TODO: why do i need to do this here?? resolve this api better
				try {
					let userRes = await axios.get(`https://api.github.com/user`, {
						headers: {
							Authorization: `token ${account.access_token}`
						}
					});
					const user = userRes.data;

					const db_user: DBUser = {
						gh_login: user.login,
						gh_access_token: account.access_token!,
						gh_avatar: user.avatar_url,
						gh_id: user.id,
						gh_created_at: user.created_at
					};
					const dbRes = await eventFetch(`/api/user`, {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify(db_user)
					});
					if (dbRes.status != 200) {
						const msg = await dbRes.json();
						token = null; // how to pass to the client ..?
						throw msg;
					} else {
						// Save the access token and refresh token in the JWT on the initial login
						const augmentedToken = {
							...token,
							login: user.login,
							access_token: account.access_token
						};
						delete augmentedToken.name;
						delete augmentedToken.sub; // not sure what this is?
						return augmentedToken;
					}
				} catch (e) {
					console.error('USER FAILED TO AUTHENTICATE IN JWT', e);
				}
			}
			return token;
		},
		async session({ session, token }) {
			//@ts-ignore
			session.user.name = token.login;
			//@ts-ignore
			session.error = token.error;
			//@ts-ignore
			session.accessToken = token.access_token;
			return session;
		}
	}
};

export function handle(event) {
	//@ts-ignore
	eventFetch = event.event.fetch;
	return sva(event);
}
//@ts-ignore
const sva = SvelteKitAuth(svaData);
