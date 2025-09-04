<script lang="ts">
	import type { AllState, LinkState, PageData, ServiceState } from '$lib/interfaces';

	import { onMount } from 'svelte';

	import Form from '$lib/components/Form.svelte';
	import Theme from '$lib/components/Theme.svelte';
	import Services from '$lib/components/Services.svelte';

	export let data: PageData;
	const urlHost = data.urlHost;
	const allState = data.state;
	const fetchUrl = data.fetchUrl;

	export const toggleSidebar = () => {
		document.getElementById('sidebar')?.classList.toggle('-translate-x-full');
		document.getElementById('open-sidebar')?.classList.toggle('hidden');
		document.getElementById('close-sidebar')?.classList.toggle('hidden');
	};

	const packetCaptureLinks = ['PE1-1/1/c1/1', 'PE1-1/1/c2/1'];
	const serviceOptions = ['VLL', 'VPLS', 'VPRN'];
	const link = [
		{ id: 'top', title: 'Top Link', src: 'PE1', dest: 'P3' },
		{ id: 'bottom', title: 'Bottom Link', src: 'PE1', dest: 'P4' }
	];
	const formInputs = [
		{ id: 'size', label: 'Size (bytes)', min: 0, max: 8000, step: 1, default: 2000 },
		{ id: 'interval', label: 'Interval (secs)', min: 0.01, max: 1, step: 0.01, default: 0.01 }
	];

	function edgesharkLink(link: string) {
		let [neName, ifcName] = link.split('-');
		let baseUrl = `packetflix:ws://${urlHost}:5001/capture?`;
		let urlParams = `container={"network-interfaces":["${ifcName}"],"name":"${neName.toLocaleLowerCase()}","type":"docker","prefix":""}&nif=${ifcName}`;
		return baseUrl + urlParams;
	}

	const updateToggle = (state: AllState) => {
		const toggle = (section: string, data: ServiceState | LinkState) => {
			for (const [key, val] of Object.entries(data)) {
				const id = `${section}-${key}`;
				const element = document.getElementById(id) as HTMLInputElement;
				element.disabled = false;
				element.checked = val;
			}
		};
		toggle('icmp', state.icmp);
		toggle('link', state.link);
		toggle('anysec', state.anysec);
	};

	async function fetchState() {
		try {
			const resp = await fetch(fetchUrl);
			const data = await resp.json();
			updateToggle(data);
		} catch (e) {
			window.location.reload();
		}
	}

	function triggerSet(e: any, module: string, service: string) {
		let currentState = e.target.checked;
		let targetState = currentState ? 'enable' : 'disable';
		setState(module, service, currentState, targetState);
	}

	async function setState(
		module: string,
		service: string,
		currentState: boolean,
		targetState: string
	) {
		try {
			const url = `/api/set/${module}/${service}/${targetState}`;
			const resp = await fetch(url, { method: 'POST' });
			const status = await resp.text();
			console.log(`${module}/${service}/${targetState} - ${status}`);
		} catch (e) {
			console.log(e);
			const element = document.getElementById(`${module}-${service}`) as HTMLInputElement;
			element.checked = currentState;
		}
	}

	onMount(async () => {
		fetchState();
		setInterval(() => fetchState(), 5000); // update state every 5 seconds
	});
</script>

<svelte:head>
	<title>SROS ANYSec Lab - Automation Panel</title>
</svelte:head>

<nav class="fixed w-screen top-0 z-30 text-sm font-nunito">
	<div
		class="flex items-center justify-between p-4 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700"
	>
		<div class="flex items-center">
			<button type="button" class="pr-4 dark:text-gray-200 lg:hidden" on:click={toggleSidebar}>
				<svg
					id="open-sidebar"
					class="w-5 h-5"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					viewBox="0 0 25 25"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path>
				</svg>
				<svg
					id="close-sidebar"
					class="w-5 h-5 hidden"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					viewBox="0 0 25 25"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
			<div><img src="/images/navbar-logo.png" alt="Logo" width="25" /></div>
		</div>
		<div class="md:flex text-center">
			<p class="dark:text-gray-200">ANYSec Lab</p>
			<p class="dark:text-gray-200">&nbsp;Automation Panel</p>
		</div>
		<div class="flex pr-4 items-center">
			<Theme />
		</div>
	</div>
</nav>

<div class="flex text-sm font-nunito">
	<aside
		id="sidebar"
		class="h-screen fixed lg:sticky top-0 left-0 w-[300px] pt-[73px] md:pt-[58px] bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transition-transform duration-300 -translate-x-full lg:-translate-x-0"
	>
		<div class="flex flex-col h-full">
			<div class="flex-grow overflow-y-auto scroll-light dark:scroll-dark">
				<div class="p-4 space-y-4">
					<div class="rounded-lg border border-gray-200 dark:border-gray-600">
						<p
							class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600"
						>
							ICMP State
						</p>
						<div class="py-3 space-y-3">
							{#each formInputs as entry}
								<Form {entry} />
							{/each}
							<hr class="h-px my-8 bg-gray-200 border-0 dark:bg-gray-700" />
							<Services key="icmp" services={serviceOptions} state={allState.icmp} />
						</div>
					</div>
					<div class="rounded-lg border border-gray-200 dark:border-gray-600 md:min-w-40">
						<p
							class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600"
						>
							Link State
						</p>
						<div class="py-3 space-y-3">
							{#each link as entry}
								<div class="flex items-center justify-between px-4">
									<div>
										<p class="text-gray-900 dark:text-white">{entry.title}</p>
										<p class="text-gray-900 dark:text-white text-xs">
											{entry.src} &mdash; {entry.dest}
										</p>
									</div>
									<div class="flex">
										<label class="inline-flex items-center cursor-pointer">
											<input
												type="checkbox"
												class="sr-only peer"
												id="link-{entry.id}"
												checked={allState.link[entry.id]}
												on:change={(event) => triggerSet(event, 'link', entry.id)}
											/>
											<div
												class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"
											></div>
										</label>
									</div>
								</div>
							{/each}
						</div>
					</div>
					<div class="rounded-lg border border-gray-200 dark:border-gray-600 md:min-w-40">
						<p
							class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600"
						>
							ANYSec State
						</p>
						<div class="py-3 space-y-3">
							<Services key="anysec" services={serviceOptions} state={allState.anysec} />
						</div>
					</div>
					<div class="rounded-lg border border-gray-200 dark:border-gray-600 md:min-w-40">
						<p
							class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600"
						>
							Packet Capture
						</p>
						<div class="flex flex-col space-y-4 p-4">
							{#each packetCaptureLinks as entry}
								{@const [neName, ifcname] = entry.split('-')}
								<a
									href={edgesharkLink(entry)}
									target="_blank"
									class="px-3 py-1 flex items-center justify-between rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 border dark:border-gray-600"
								>
									<span>{neName}</span>
									<span>{ifcname}</span>
									<span class="w-6"><img src="/images/edgeshark.png" alt="fin" /></span>
								</a>
							{/each}
						</div>
					</div>
				</div>
			</div>
		</div>
	</aside>
	<main class="flex-1 bg-gray-400 pt-[73px] md:pt-[58px]">
		<!-- svelte-ignore a11y-missing-attribute -->
		<iframe
			class="h-[calc(100vh-55px)] w-full"
			src="http://{urlHost}:3000/d/kawVPD-Gk/anysec-telemetry"
		></iframe>
	</main>
</div>
