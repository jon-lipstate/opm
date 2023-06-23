import { getAuth } from '$api/auth';
import sql from '$lib/database';
import { extractHostOwnerAndRepo, isValidSemver } from '$lib/utils.js';
import { error, json } from '@sveltejs/kit';
import axios from 'axios';

async function isValidUrl(url) {
	try {
		const response = await axios.head(url);
		return response.status === 200;
	} catch (error) {
		return false;
	}
}
export async function GET(event) {
	// TODO: allow for query string to select other users, eg browse by owner
	//@ts-ignore
	const { session } = await getAuth(event);
	let pkgs;
	try {
		const userId = session.user.id;
		pkgs = await sql`SELECT id,host_name,owner_name,repo_name,description FROM packages where owner_id=${userId}`;
		for (let i = 0; i < pkgs.length; i++) {
			const versions = await sql`
						SELECT 
							v.id,
							v.version,
							v.license,
							v.commit_hash,
							v.insecure,
							EXISTS (
								SELECT 1 FROM package_dependencies pd WHERE pd.depends_on_id = v.id
							) AS has_deps
						FROM 
							versions v
						WHERE 
							v.package_id=${pkgs[i].id}
					`;
			pkgs[i].versions = versions ?? [];
		}
		return json(pkgs);
	} catch (err: any) {
		console.error('SQL Get Packages Error\n', err);
		if (err.status < 500) {
			throw error(err.status, err.body.message);
		}
		throw error(500, `SQL Get Packages Error:, ${err}`);
	}
}

export async function POST(event) {
	let pkg;
	let userId = -1;
	let host_name: String;
	let owner_name: String;
	let repo_name: String;
	try {
		console.warn(event.request);
		pkg = await event.request.json();
	} catch (e) {
		throw error(400, `Failed to parse JSON: ${e}.`);
	}
	try {
		const token = pkg.token;
		const idRes = await sql`select * from public.verify_token(${token})`;
		userId = idRes[0].verify_token;
	} catch (e) {
		throw error(401, 'Invalid or Revoked CLI-Token');
	}
	try {
		({ host_name, owner_name, repo_name } = extractHostOwnerAndRepo(pkg.userData.url));
	} catch (e) {
		throw error(400, 'Unable to parse url, expected host/owner/repo (eg `https://the_host.com/the_owner/the_repo`)');
	}
	try {
		const urlResult = await isValidUrl(pkg.userData.url);
		if (!urlResult) {
			throw Error('Invalid URL');
		}
	} catch (e) {
		throw error(400, 'Invalid Package URL');
	}
	try {
		const { size_kb, compiler, commit_hash, readme_contents } = pkg;
		const { url, readme, description, version, license, keywords, dependencies } = pkg.userData;
		//
		let errors: string[] = [];
		//
		if (!isValidSemver(version)) {
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
		if (!license) {
			errors.push('Packages without licenses are prohibited.');
		}
		if (!readme_contents || !readme) {
			errors.push('Expected a readme file.');
		}
		if (!commit_hash) {
			errors.push('Commit Hash Missing.');
		}
		if (!compiler) {
			errors.push('Compiler Info Missing.');
		}
		const res = await event.fetch(`/api/preprocess-readme`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ readme_contents })
		});
		if (res.status != 200) {
			console.error(res.statusText);
			errors.push(`Readme Parse Error: ${res.statusText}`);
		}
		const preparedReadme = (await res.json()).html;
		//
		let depVersionIds: Number[] = [];
		// this verifies deps are good:
		if (dependencies) {
			const deps = dependencies.map((x) => {
				try {
					const { host_name, owner_name, repo_name } = extractHostOwnerAndRepo(x.name);
					let version = x.version;
					return { host_name, owner_name, repo_name, version };
				} catch (ex) {
					throw error(400, `Unable to parse dependancy: ${x.name}`);
				}
			});
			//@ts-ignore
			const pkgIdsRes = await sql`SELECT * FROM get_package_ids(${deps})`;
			for (let i = 0; i < deps.length; i++) {
				const id = pkgIdsRes[i]?.package_id;
				if (!id) {
					errors.push(`Invalid package ${JSON.stringify(deps[i])}`);
				}
				//@ts-ignore
				deps[i].id = id;
			}
			//@ts-ignore
			const versionIdsRes = await sql`SELECT * FROM get_version_ids(${deps})`;
			for (let i = 0; i < deps.length; i++) {
				const vid = versionIdsRes[i]?.version_id;
				const dep = deps[i];
				if (!vid && vid !== 0) {
					errors.push(`Invalid Version Id ${JSON.stringify(dep)}`);
				}
			}
			depVersionIds = versionIdsRes.map((x) => x.version_id);
		}

		if (errors.length != 0) {
			console.warn(errors);
			return json({ message: errors.join('\n') }, { status: 400 });
		}
		//@ts-ignore
		const createRes = await sql`call upsert_full_package(
			${userId}::INTEGER,
			${host_name}::TEXT,
			${owner_name}::TEXT,
			${repo_name}::TEXT,
			${description}::TEXT,
			${preparedReadme}::TEXT,
			${url}::TEXT,
			${version.replace('v', '')}::TEXT,
			${license}::TEXT,
			${size_kb}::INTEGER,
			${compiler ?? 'unknown'}::TEXT,
			${commit_hash}::TEXT,
			ARRAY[${sql.array(keywords)}],
			${depVersionIds}::INTEGER[]
			)`;
		return json({ message: 'Successful Upsert' }, { status: 201 });
		//
	} catch (err: any) {
		if (err.status == 400) {
			throw err;
		}
		const errStr = String(err).replace('PostgresError: ', '');
		if (errStr.includes('Version')) {
			throw error(400, { message: errStr }); // cannot replace a version
		} else {
			throw error(503, { message: errStr });
		}
	}
}

export async function DELETE(event) {
	const body = await event.request.json();
	const id = body.id;
	//@ts-ignore
	const { session } = await getAuth(event);
	try {
		const userId = session.user.id;

		await sql`DELETE from packages where owner_id=${userId} AND id=${id}`;
		return json({ status: 200 });
	} catch (err: any) {
		console.error('SQL New Token Error\n', err);
		if (err.status < 500) {
			throw error(err.status, err.body.message);
		}
		throw error(500, `SQL New Token Error:, ${err}`);
	}
}
