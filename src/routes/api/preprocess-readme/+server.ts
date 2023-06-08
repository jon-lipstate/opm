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
	const body = await event.request.json();
	let data = body.readme_contents;
	//
	let rendered = md.render(data);
	//
	const dom = new JSDOM(rendered);
	let html = dom.window.document.documentElement.outerHTML;
	const DOMPurify = createDOMPurify(dom.window);
	html = DOMPurify.sanitize(html);
	//
	return json({ html });
}
