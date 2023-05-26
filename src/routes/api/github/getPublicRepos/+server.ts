import { json } from '@sveltejs/kit';
import axios from 'axios';
import { getAuth } from '../../auth.js';

export async function GET(event) {
	const { login, authHeader, session } = await getAuth(event);
	try {
		const repoRes = await axios.get(`https://api.github.com/user/repos`, authHeader);

		const repos = repoRes.data.filter((x) => x.language == 'Odin');
		const otherRepos = repoRes.data.filter((x) => x.language != 'Odin');

		for (const repo of otherRepos) {
			const url = `https://api.github.com/repos/${login}/${repo.name}/languages`;
			const res = await axios.get(url, authHeader);
			if ('Odin' in res.data) {
				repos.push(repo);
			}
		}
		return json(repos);
	} catch (e) {
		console.error(`>>> getPublicGists: "${session}"`);
		//@ts-ignore
		return fail(503, `api: getPublicRepos, err: ${e}`);
	}
}
