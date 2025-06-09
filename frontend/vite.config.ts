import { defineConfig } from 'vite'

export default defineConfig({
	base: '/telegram/',
	build: {
		emptyOutDir: true,
		outDir: '../backend/telegram/static',
	},
})


