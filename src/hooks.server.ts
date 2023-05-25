import { SvelteKitAuth } from '@auth/sveltekit';
import GitHub from '@auth/core/providers/github';
import { GITHUB_ID, GITHUB_SECRET, AUTH_SECRET } from '$env/static/private';
import axios from 'axios';

export const handle = SvelteKitAuth({
	providers: [
		//@ts-ignore
		GitHub({
			clientId: GITHUB_ID,
			clientSecret: GITHUB_SECRET,
			authorization: {
				params: { scope: 'read:org read:user user:email' }
			}
		})
	],
	secret: AUTH_SECRET,
	trustHost: true,
	callbacks: {
		async jwt({ token, account }) {
			if (account) {
				// TODO: why do i need to do this here?? resolve this api better
				let user = await axios.get(`https://api.github.com/user`, {
					headers: {
						Authorization: `token ${account.access_token}`
					}
				});
				// Save the access token and refresh token in the JWT on the initial login
				const augmentedToken = {
					...token,
					login: user.data.login,
					access_token: account.access_token
				};
				delete augmentedToken.name;
				delete augmentedToken.sub; // not sure what this is?
				return augmentedToken;
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
});
