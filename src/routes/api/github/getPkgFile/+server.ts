import { json } from '@sveltejs/kit';
import axios from 'axios';
import { getAuth } from '../../auth.js';

export async function POST(event) {
	const { login, authHeader } = await getAuth(event);
	const body = JSON.parse(await event.request.text());
	try {
		const contentsRes = await axios.get(`https://api.github.com/repos/${login}/${body.name}/contents`, authHeader);
		const pkgFile = contentsRes.data.find((x) => (x.name as string).includes('.pkg'));
		if (!pkgFile) {
			console.error(`>>> FAILED TO FIND PKG: "${login}/${body.name}"`);
			return json({ error: 'mod.pkg file is missing.' }, { status: 404 });
		}
		const file = await axios.get(pkgFile.download_url);
		try {
			// make sure it parses:
			const str = JSON.stringify(file.data);
			const _js = JSON.parse(str);
		} catch (e) {
			console.error(`>>> getPkgFile - malformed pkg file: "${login}/${body.name}"`);
			return json({ error: `malformed pkg file: ${e}` }, { status: 400 });
		}
		return json(file.data);
	} catch (e) {
		console.error(`>>> getPkgFile: "${login}/${body.name}"`);
		return json({ error: `api: getPkgFile, err: ${e}` }, { status: 503 });
	}
}

// GET /repos/{owner}/{repo}/readme'
// GET /repos/{owner}/{repo}/contents/{path}
//
