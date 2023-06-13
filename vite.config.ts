import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	server: {
		host: true
	},
	plugins: [sveltekit()],
	resolve: {
		alias: {
			$components: '/src/components',
			$stores: '/src/stores',
			$api: '/src/routes/api'
		}
	}
});
