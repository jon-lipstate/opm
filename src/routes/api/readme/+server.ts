import { Remarkable } from 'remarkable';
import axios from 'axios';
import createDOMPurify from 'dompurify';
import { JSDOM } from 'jsdom';
import { error, json } from '@sveltejs/kit';

var md = new Remarkable('commonmark');
md.renderer.rules.code = function (tokens, idx) {
	return '<span class="inline-code">' + tokens[idx].content + '</span>';
};
export async function POST(event) {
	const body = JSON.parse(await event.request.text());
	let data;
	if (body.data) {
		data = body.data;
	} else if (body.url) {
		const res = await axios.get(body.url);

		if (res.status != 200) {
			throw error(res.status, res.statusText);
		}
		data = '```odin\n' + res.data + '\n```';
		console.info('TODO: REMOVE ODIN TAG WRAPPER');
	} else {
		throw error(400, "readme requires 'data' or 'url' to process.");
	}
	let rendered = md.render(data);
	const dom = new JSDOM(rendered);
	let html = dom.window.document.documentElement.outerHTML;
	const DOMPurify = createDOMPurify(dom.window);
	html = DOMPurify.sanitize(html);
	return json({ html });
}
