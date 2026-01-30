// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightThemeGalaxy from 'starlight-theme-galaxy'

// https://astro.build/config
export default defineConfig({
	markdown: {
		smartypants: false,
	},
	integrations: [
		starlight({
			plugins: [starlightThemeGalaxy()],
			title: 'Boiler',
			description: 'CLI tool for managing code snippets and project stacks',
			logo: {
				src: './src/assets/logo.svg',
			},
			favicon: '/favicon.svg',
			customCss: ['./src/styles/custom.css'],
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/rishiyaduwanshi/boiler' }],
			components: {
				SocialIcons: './src/components/CustomHeader.astro',
			},
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'guides/introduction' },
						{ label: 'Installation', slug: 'guides/installation' },
						{ label: 'Quick Start', slug: 'guides/quickstart' },
						{ label: 'Use Cases', slug: 'guides/usecases' },
					],
				},
				{
					label: 'Commands',
					autogenerate: { directory: 'commands' },
				},
				{
					label: 'CLI Reference',
					collapsed: true,
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
