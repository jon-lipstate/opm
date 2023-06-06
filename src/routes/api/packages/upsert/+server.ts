import sql from '$lib/database';
import { extractHostOwnerAndRepo, generateSlug, isValidSemver } from '$lib/utils.js';
import { error, json } from '@sveltejs/kit';
import { getAuth } from '../../auth.js';
import axios from 'axios';

async function getUserId(userName, authToken): Promise<number> {
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

export async function POST(event) {
	const pkg = JSON.parse(await event.request.text());
	let userId = -1;
	const { host_name, owner_name, repo_name } = extractHostOwnerAndRepo(pkg.url);
	console.warn(host_name, owner_name, repo_name);
	const description = pkg.description;
	let readme: string = '';
	let size_kb = 42;
	// console.warn(pkg);

	if (false) {
		// todo: cli token path
	} else {
		let auth = await getAuth(event); // website path
		userId = await getUserId(auth.session.user.name, auth.session.accessToken);
	}
	if (userId == -1) {
		console.error('Invalid User Identity', await event.locals.getSession());
		throw error(401, 'User Identity Unknown.');
	}

	let errors: string[] = [];

	if (!isValidSemver(pkg.version)) {
		errors.push("Invalid 'version' field. Must comply with semver, do not include 'v'");
	}

	if (!host_name || host_name.length < 5) {
		// todo: ping the host? smallest i think is something like `sr.ht`
		errors.push('Invalid Host.');
	}

	if (!owner_name || owner_name.length < 3) {
		// todo: find min gh length?
		errors.push('Owner name invalid.');
	}
	if (!repo_name || repo_name.length < 3) {
		// todo: find min gh length?
		errors.push('Repo name invalid.');
	}

	if (description.length < 10) {
		errors.push('Description must have at least 10 chars.');
	}
	// TODO: cli path needs to fetch token from db to make this call
	const { login, authHeader } = await getAuth(event);
	if (owner_name != login) {
		errors.push('Require Owner to match login, to be retired once CLI completed (this issue bars groups for now too).');
	}
	const contentsRes = await axios.get(`https://api.github.com/repos/${login}/${repo_name}/contents`, authHeader);
	if (!pkg.license) {
		errors.push('Packages without licenses are prohibited.');
	}

	let readmeData = contentsRes.data.find((x) => (x.name as string).includes(pkg.readme));
	if (!pkg.readme && !readmeData) {
		readmeData = contentsRes.data.find((x) => (x.name as string).includes('readme.md'));
	}
	if (!readmeData) {
		errors.push('Cannot locate readme, expected `pkg.readme: the_file.ext` or `readme.md`');
	} else {
		const url = readmeData.download_url;
		const res = await event.fetch(`/api/readme`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ url })
		});
		if (res.status != 200) {
			console.error(res.statusText);
			errors.push(`Readme Parse Error: ${res.statusText}`);
		}
		readme = (await res.json()).html;
	}
	let depVersionIds: Number[] = [];
	// this verifies deps are good:
	if (pkg.dependencies) {
		const deps = Object.keys(pkg.dependencies).map((x) => {
			let [host_name, owner_name, repo_name] = x.split('/');
			let version = pkg.dependencies[x];
			return { host_name, owner_name, repo_name, version };
		});
		//@ts-ignore
		const pkgIdsRes = await sql`SELECT * FROM get_package_ids(${deps})`;
		for (let i = 0; i < deps.length; i++) {
			const id = pkgIdsRes[i].package_id;
			if (!id) {
				errors.push(`Invalid package ${deps}`);
			}
			//@ts-ignore
			deps[i].id = id;
		}
		//@ts-ignore
		const versionIdsRes = await sql`SELECT * FROM get_version_ids(${deps})`;
		for (let i = 0; i < deps.length; i++) {
			const vid = versionIdsRes[i].version_id;
			const dep = deps[i];
			if (!vid && vid !== 0) {
				errors.push(`Invalid Version Id ${dep}`);
			}
		}
		depVersionIds = versionIdsRes.map((x) => x.version_id);
	}

	if (errors.length != 0) {
		console.warn(errors);
		throw error(400, errors.join('\n'));
	}
	try {
		//@ts-ignore
		const createRes = await sql`call upsert_full_package(
		${userId}::INTEGER,
		${host_name}::TEXT,
		${owner_name}::TEXT,
		${repo_name}::TEXT,
		${description}::TEXT,
		${readme}::TEXT,
		${pkg.url}::TEXT,
		${pkg.version.replace('v', '')}::TEXT,
		${pkg.license}::TEXT,
		${size_kb}::INTEGER,
		${pkg.compiler ?? 'unknown'}::TEXT,
		${pkg.commit_hash}::TEXT,
		ARRAY[${sql.array(pkg.keywords)}],
		${depVersionIds}::INTEGER[]
		)`;
		//
	} catch (e: any) {
		console.warn('sql error:', e);
	}

	return json({ msg: 'Not Implemented' }, { status: 201 });
	try {
		return json(null, { status: 201 });
	} catch (error) {
		console.error('Error inserting package details', error);
		return json(null, { status: 500 });
	}
}
