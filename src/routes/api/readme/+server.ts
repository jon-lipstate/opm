import { Remarkable } from 'remarkable';
import axios from 'axios';
import hljs from 'highlight.js';
import createDOMPurify from 'dompurify';
import { JSDOM } from 'jsdom';
import odin from './odin-hl.js';
import { error, json } from '@sveltejs/kit';
// TODO: figure out how to just say parse odin explicitly:
// hljs.listLanguages().forEach((x) => hljs.unregisterLanguage(x));
hljs.registerLanguage('odin', odin);

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
		console.warn('GOT URL');

		if (res.status != 200) {
			throw error(res.status, res.statusText);
		}
		data = '```odin\n' + res.data + '\n```';
		console.info('TODO: REMOVE ODIN TAG WRAPPER');
	} else {
		throw error(400, "readme requires 'data' or 'url' to process.");
	}
	let rendered = md.render(data);
	// const res = hljs.highlight(rendered, { language: 'odin' }, true);
	const dom = new JSDOM(rendered);
	console.info('HLJS');
	try {
		dom.window.document.querySelectorAll('pre').forEach((x) => {
			hljs.highlightElement(x);
			//
		});
	} catch (e) {
		console.error('HLJS-ERR', e);
		throw error(400, e);
	}
	console.warn('AFTER1');
	let html = dom.window.document.documentElement.outerHTML;
	console.info('PURIFY');
	const DOMPurify = createDOMPurify(dom.window);
	html = DOMPurify.sanitize(html);
	console.info('API COMPLETE');
	// let html = res;
	return json({ html });
}
