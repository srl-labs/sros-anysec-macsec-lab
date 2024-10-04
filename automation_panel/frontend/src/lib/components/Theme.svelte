<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	let darkMode = false;

	onMount(() => {
		if (browser) {
			darkMode =
				localStorage.theme === 'dark' ||
				(!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);

			themeIconSwitch();
		}
	});

	const themeIconSwitch = () => {
		if (darkMode) {
			document.documentElement.classList.add('dark');
			document.getElementById('toggle-dark-icon')?.classList.add('hidden');
			document.getElementById('toggle-light-icon')?.classList.remove('hidden');
		} else {
			document.documentElement.classList.remove('dark');
			document.getElementById('toggle-dark-icon')?.classList.remove('hidden');
			document.getElementById('toggle-light-icon')?.classList.add('hidden');
		}
	};

	const toggleDarkMode = () => {
		darkMode = !darkMode;
		localStorage.setItem('theme', darkMode ? 'dark' : 'light');
		themeIconSwitch();
	};
</script>

<button
	class="flex w-8 h-5 ml-2 rounded-full items-center transition duration-30 {darkMode
		? 'bg-gray-500'
		: 'bg-gray-300'}"
	on:click={toggleDarkMode}
>
	<div
		id="switch-toggle"
		class="w-5 h-5 p-1 relative rounded-full transition duration-500 transform {darkMode
			? 'bg-gray-600 translate-x-3'
			: 'bg-white border border-gray-300 -translate-x-0'}"
	>
		<svg
			id="toggle-dark-icon"
			class="text-black"
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
			/>
		</svg>
		<svg
			id="toggle-light-icon"
			class="text-white hidden"
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
			/>
		</svg>
	</div>
</button>
