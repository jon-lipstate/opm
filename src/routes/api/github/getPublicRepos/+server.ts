import { json } from '@sveltejs/kit';
import axios from 'axios';
import { getAuth } from '../../auth.js';

export async function GET(event) {
	try {
		const { login, authHeader } = await getAuth(event);

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
		//@ts-ignore
		return fail(503, `api: getPublicRepos, err: ${e}`);
	}
}
