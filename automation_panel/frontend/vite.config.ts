import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit()],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	},

	server: {
		proxy: {
			// svelte uses the relative /api path to reach the backend
			// can be substituted with the actual proxy like nginx
			'/api': {
				target: 'http://localhost:8080',
				rewrite: (path) => path.replace(/^\/api/, ''),
			}
		}
	}
});
