import { error, json } from '@sveltejs/kit';
import axios from 'axios';
import { getAuth } from '../../auth.js';

export async function GET(event) {
	const { login, authHeader, session } = await getAuth(event);
	if (!login) {
		throw error(403, 'Not Logged in');
	}
	try {
		const repoRes = await axios.get(`https://api.github.com/user/repos?per_page=500`, authHeader);
		let repos = repoRes.data;

		// TODO: move this to client side and __dont__ await it?
		// let repos = repoRes.data.filter((x) => x.language == 'Odin');
		// const otherRepos = repoRes.data.filter((x) => x.language != 'Odin');

		// for (const repo of otherRepos) {
		// 	const url = `https://api.github.com/repos/${login}/${repo.name}/languages`;
		// 	const res = await axios.get(url, authHeader);
		// 	if ('Odin' in res.data) {
		// 		repos.push(repo);
		// 	}
		// }
		let cleanedRepos = repos.map((r) => {
			let license = { ...r.license };
			let value = {
				name: r.name,
				owner: r.owner.login,
				html_url: r.html_url,
				ssh_url: r.ssh_url,
				description: r.description,
				fork: r.fork,
				created_at: r.created_at,
				updated_at: r.updated_at,
				homepage: r.homepage,
				language: r.language,
				size_kb: r.size,
				archived: r.archived,
				license_id: license.spdx_id,
				license_name: license.name,
				visibility: r.visibility,
				default_branch: r.default_branch
			};
			return value;
		});

		return json({ repos: cleanedRepos });
	} catch (e: any) {
		if (e.status == 403) {
			throw e;
		}
		console.error(`>>> getPublicGists: "${session}"`);
		throw error(503, `api: getPublicRepos, err: ${e}`);
	}
}
const delOwnerKeys = [
	'url',
	'html_url',
	'followers_url',
	'following_url',
	'gists_url',
	'starred_url',
	'subscriptions_url',
	'organizations_url',
	'repos_url',
	'events_url',
	'received_events_url'
];
const delRepoKeys = [
	'forks_url',
	'keys_url',
	'collaborators_url',
	'teams_url',
	'hooks_url',
	'issue_events_url',
	'events_url',
	'assignees_url',
	'branches_url',
	'tags_url',
	'blobs_url',
	'git_tags_url',
	'git_refs_url',
	'trees_url',
	'statuses_url',
	'languages_url',
	'stargazers_url',
	'contributors_url',
	'subscribers_url',
	'subscription_url',
	'commits_url',
	'git_commits_url',
	'comments_url',
	'issue_comment_url',
	'contents_url',
	'compare_url',
	'merges_url',
	'archive_url',
	'downloads_url',
	'issues_url',
	'pulls_url',
	'milestones_url',
	'notifications_url',
	'labels_url',
	'releases_url',
	'deployments_url',
	'git_url',
	'clone_url',
	'svn_url',
	'mirror_url'
];
