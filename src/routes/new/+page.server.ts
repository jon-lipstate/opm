import { json } from '@sveltejs/kit';
import axios from 'axios';

let authHeader;
export async function load(event) {
	const data = await event.parent();
	//@ts-ignore
	const { user, accessToken } = data.session;
	authHeader = {
		headers: {
			Authorization: `token ${accessToken}`
		}
	};
	const repoRes = await axios.get(`https://api.github.com/user/repos`, authHeader);

	const repos = repoRes.data.filter((x) => x.language == 'Odin');
	const otherRepos = repoRes.data.filter((x) => x.language != 'Odin');
	for (const repo of otherRepos) {
		const url = `https://api.github.com/repos/${user.name}/${repo.name}/languages`;
		const res = await axios.get(url, authHeader);
		if ('Odin' in res.data) {
			repos.push(repo);
		}
	}

	return {
		repos
	};
}

export const actions = {
	getPublicPackages: async (event) => {
		const form = await event.request.formData();
		const owner = form.get('owner');
		const packages = await axios.get(`https://api.github.com/orgs/${owner}/repos`, authHeader);
		console.warn(packages);
	}
};
