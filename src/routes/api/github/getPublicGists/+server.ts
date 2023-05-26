import { fail, json, type ActionFailure } from '@sveltejs/kit';
import axios from 'axios';
import { getAuth } from '../../auth.js';

export async function GET(event) {
	const { login, authHeader, session } = await getAuth(event);

	try {
		const gistRes = await axios.get(`https://api.github.com/users/${login}/gists`, authHeader);
		// Currently, the GitHub API doesn't support filtering gists by language directly
		// so we have to fetch all gists, then filter them manually.
		// This could be a lot slower than getting repos if the user has a lot of gists.
		const gists = gistRes.data;
		gistRes.data.filter((gist) => {
			for (const file of Object.values(gist.files)) {
				//@ts-ignore
				if (file.language === 'Odin') {
					return true;
				}
			}
			return false;
		});

		return json(gists);
	} catch (e) {
		console.error(`>>> getPublicGists: "${session}"`);
		//@ts-ignore
		return fail(503, `api: getPublicGists, err: ${e}`);
	}
}
