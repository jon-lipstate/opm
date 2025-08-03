<template>
	<q-card class="q-mt-lg">
		<q-card-section>
			<div class="row items-center q-mb-md">
				<q-icon name="description" size="24px" class="q-mr-sm" />
				<span class="text-h6">README</span>
			</div>
			<q-separator class="q-mb-md" />

			<div v-if="loading" class="text-center q-py-xl">
				<q-spinner size="40px" color="primary" />
				<p class="text-caption q-mt-sm">Loading README...</p>
			</div>

			<div v-else-if="error" class="text-center q-py-xl text-grey-6">
				<q-icon name="error_outline" size="48px" />
				<p class="q-mt-sm">Failed to load README</p>
			</div>

			<div v-else v-html="renderedHtml" class="markdown-body"></div>
		</q-card-section>
	</q-card>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import odinLang from 'src/lib/odin-hl.js'

// Import common languages
import javascript from 'highlight.js/lib/languages/javascript'
import bash from 'highlight.js/lib/languages/bash'
import json from 'highlight.js/lib/languages/json'
import xml from 'highlight.js/lib/languages/xml'
import css from 'highlight.js/lib/languages/css'
import plaintext from 'highlight.js/lib/languages/plaintext'

// Register languages
hljs.registerLanguage('odin', odinLang)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('js', javascript)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('sh', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('css', css)
hljs.registerLanguage('plaintext', plaintext)
hljs.registerLanguage('text', plaintext)

// Props
const props = defineProps({
	content: {
		type: String,
		default: '',
	},
	loading: {
		type: Boolean,
		default: false,
	},
	error: {
		type: Boolean,
		default: false,
	},
})

// State
const renderedHtml = ref('')

// Configure marked
marked.setOptions({
	breaks: true, // Convert \n to <br>
	gfm: true, // GitHub Flavored Markdown
	sanitize: false, // We trust our content
	highlight: function (code, lang) {
		// Try to highlight with the specified language
		if (lang && hljs.getLanguage(lang)) {
			try {
				const result = hljs.highlight(code, { language: lang })
				return result.value
			} catch (err) {
				console.error('Highlight error:', err)
			}
		}

		// Try to auto-detect language
		try {
			const result = hljs.highlightAuto(code, ['odin', 'javascript', 'bash', 'json'])
			// console.log('Auto-detected language:', result.language)
			return result.value
		} catch (err) {
			console.error('Highlight auto error:', err)
		}

		// Return unhighlighted code
		return code
	},
})

// Parse markdown when content changes
const parseMarkdown = async (markdown) => {
	if (!markdown) {
		renderedHtml.value = ''
		return
	}

	try {
		// Parse markdown to HTML
		const html = marked(markdown)
		// console.log('Parsed HTML:', html)
		renderedHtml.value = html

		// Apply syntax highlighting after render
		await nextTick()

		// Highlight any code blocks that weren't processed by marked
		const codeBlocks = document.querySelectorAll('.markdown-body pre code')
		// console.log('Found code blocks:', codeBlocks.length)

		codeBlocks.forEach((block) => {
			// Check if already highlighted
			if (!block.classList.contains('hljs')) {
				// Get language from class
				const match = block.className.match(/language-(\w+)/)
				// if (match) {
				// 	console.log('Found language class:', match[1])
				// }
				hljs.highlightElement(block)
				// console.log('Highlighted block. Classes:', block.className)
			}
		})
	} catch (err) {
		console.error('Error parsing markdown:', err)
		renderedHtml.value = '<p>Error rendering README</p>'
	}
}

// Watch for content changes
watch(
	() => props.content,
	(newContent) => {
		parseMarkdown(newContent)
	},
	{ immediate: true },
)
</script>

<style lang="scss">
@import 'src/css/ayu-mirage.css';

.markdown-body {
	line-height: 1.6;

	h1,
	h2,
	h3 {
		margin-top: 24px;
		margin-bottom: 16px;
		font-weight: 600;
	}

	h1 {
		font-size: 2em;
	}
	h2 {
		font-size: 1.5em;
	}
	h3 {
		font-size: 1.25em;
	}

	p {
		margin-bottom: 16px;
	}

	pre {
		background-color: #242936;
		padding: 16px;
		border-radius: 6px;
		overflow-x: auto;
		margin-bottom: 16px;
	}

	pre code {
		background-color: transparent;
		padding: 0;
		border-radius: 0;
	}

	body.body--dark & {
		pre {
			background-color: var(--dark-code, #242936);
		}
	}

	a {
		color: var(--odin-blue);
		text-decoration: none;

		&:hover {
			text-decoration: underline;
		}
	}

	ul,
	ol {
		margin-bottom: 16px;
		padding-left: 2em;
	}

	blockquote {
		margin: 0 0 16px 0;
		padding: 0 1em;
		color: #6a737d;
		border-left: 0.25em solid #dfe2e5;
	}

	table {
		margin-bottom: 16px;
		border-spacing: 0;
		border-collapse: collapse;

		th,
		td {
			padding: 6px 13px;
			border: 1px solid #dfe2e5;
		}

		th {
			font-weight: 600;
			background-color: #f6f8fa;
		}

		tr {
			background-color: #fff;
			border-top: 1px solid #c6cbd1;

			&:nth-child(2n) {
				background-color: #f6f8fa;
			}
		}
	}

	body.body--dark & table {
		th {
			background-color: #161b22;
		}

		tr {
			background-color: #0d1117;

			&:nth-child(2n) {
				background-color: #161b22;
			}
		}

		th,
		td {
			border-color: #30363d;
		}
	}
}
</style>
