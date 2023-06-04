import sql from '$lib/database';
import { extractUserAndProject, generateSlug, isValidSemver } from '$lib/utils.js';
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
	const slug = generateSlug(pkg.name);
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

	if (pkg.name.length < 1) {
		errors.push('Name must be at least one character.');
	}

	if (slug.length < 1) {
		errors.push('Name must be encodable as a slug with >1 char.');
	}

	if (description.length < 10) {
		errors.push('Description must have at least 10 chars.');
	}
	// TODO: cli path needs to fetch token from db to make this call
	const { login, authHeader } = await getAuth(event);
	const { user, repo } = extractUserAndProject(pkg.url)!; // todo: should fail on null?
	if (user != login) {
		errors.push('User and login do not match, note that only github is supported at present.');
	}
	const contentsRes = await axios.get(`https://api.github.com/repos/${login}/${repo}/contents`, authHeader);
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
			let [name, slug] = x.split('/');
			let version = pkg.dependencies[x];
			return { name, slug, version };
		});
		//@ts-ignore
		const pkgIdsRes = await sql`select * from get_package_ids(${deps})`;
		for (let i = 0; i < deps.length; i++) {
			const id = pkgIdsRes[i].package_id;
			if (!id) {
				errors.push(`Invalid package ${deps[i].name}/${deps[i].slug}`);
			}
			//@ts-ignore
			deps[i].id = id;
		}
		//@ts-ignore
		const versionIdsRes = await sql`select * from get_version_ids(${deps})`;
		for (let i = 0; i < deps.length; i++) {
			const vid = versionIdsRes[i].version_id;
			const dep = deps[i];
			if (!vid) {
				errors.push(`Invalid Version Id ${dep.name}/${dep.slug}`);
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
		${pkg.name}::TEXT,
		${slug}::TEXT,
		${description}::TEXT,
		${readme}::TEXT,
		${pkg.url}::TEXT,
		${pkg.version.replace('v', '')}::TEXT,
		${pkg.license}::TEXT,
		${size_kb}::INTEGER,
		${pkg.compiler ?? 'not specified'}::TEXT,
		ARRAY[${sql.array(pkg.keywords)}],
		ARRAY[${sql.unsafe(depVersionIds.join(','))}]::INTEGER[]
		)`;
	} catch (e: any) {
		console.warn('sql error:', e);
	}

	return json({ msg: 'Not Implemented' }, { status: 201 });
	// TODO: NAME MUST BE SLUGGABLE IMPLICITLY

	// await make_package_table();

	try {
		console.warn('INSERTING PKG');
		await insert_new_package({
			name: 'Test',
			updated_at: new Date(),
			created_at: new Date(),
			downloads: 42,
			description: 'at a time',
			readme: 'stuff and things',
			repository: 'https://google.com',
			size_kb: 42
		});

		return json(null, { status: 201 });
	} catch (error) {
		console.error('Error inserting package details', error);
		return json(null, { status: 500 });
	}
}

const insert_new_package = async function ({
	name,
	updated_at,
	created_at,
	downloads,
	description,
	readme,
	repository,
	size_kb
}) {
	console.warn(name, updated_at, created_at, downloads, description, readme, repository, size_kb);
	try {
		await sql`INSERT INTO packages (name, updated_at, created_at, downloads, description, readme, repository, size_kb)
				VALUES (${name}, ${updated_at}, ${created_at}, ${downloads}, ${description}, ${readme}, ${repository}, ${size_kb})`;

		console.log('Package inserted successfully');
	} catch (error) {
		console.error('Error inserting package: ', error);
	}
};
