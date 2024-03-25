<script lang="ts">
  import Form from '$lib/components/Form.svelte';
  import Theme from '$lib/components/Theme.svelte';
  import Services from '$lib/components/Services.svelte';

  export const toggleSidebar = () => {
		document.getElementById('sidebar')?.classList.toggle('-translate-x-full');
		document.getElementById('open-sidebar')?.classList.toggle('hidden');
		document.getElementById('close-sidebar')?.classList.toggle('hidden');
	};

  const serviceOptions = ["VLL", "VPLS", "VPRN"];
  const link = [
    { "name": "Top Link", "src": "PE1", "dest": "P3" },
    { "name": "Bottom Link", "src": "PE1", "dest": "P4" }
  ]
  const formInputs = [
    { "name": "size", "label": "Size (bytes)", "min": 0, "max": 8000, "step": 1, "default": 2000 },
    { "name": "interval", "label": "Interval (secs)", "min": 0.01, "max": 1, "step": 0.01, "default": 0.01 }
  ]
</script>

<!-- NAVBAR -->
<nav class="fixed top-0 z-30 px-3 py-4 w-screen select-none text-sm font-nunito bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
	<div class="flex justify-between">
		<!-- navbar left item -->
		<div class="flex items-center space-x-2">
      <button type="button" class="flex dark:text-gray-200" on:click={toggleSidebar}>
        <svg id="open-sidebar" class="w-5 h-5 hidden" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 25 25" xmlns="http://www.w3.org/2000/svg">
          <path d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path>
        </svg>
        <svg id="close-sidebar" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 25 25" xmlns="http://www.w3.org/2000/svg">
          <path d="M6 18L18 6M6 6l12 12"></path>
        </svg>
      </button>
      <!--<div class="flex"><img src="/images/containerlab.svg" alt="Logo" width="35" /></div>-->
			<div class="flex px-2"><img src="/images/navbar-logo.png" alt="Logo" width="25" /></div>
      <!--<div class="flex px-2"><img src="/images/{darkMode ? 'nwhite' : 'nblue'}.png" alt="Logo" width="70" /></div>-->
		</div>
		<!-- navbar centre item -->
		<div class="md:flex text-center">
			<p class="dark:text-gray-200">ANYSec Lab</p>
      <p class="dark:text-gray-200">&nbsp;Automation Panel</p>
		</div>
		<!-- navbar right item -->
    <div class="flex items-center mr-3">
		  <Theme/>
    </div>
	</div>
</nav>

<!-- SIDEBAR -->
<div class="relative md:flex h-screen overflow-hidden">
  <aside id="sidebar" class="text-sm font-nunito pb-4 overflow-y-auto scroll-light dark:scroll-dark z-20 w-[300px] transition-transform -translate-x-0 h-screen sticky bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
    <div class="px-4 space-y-4 pt-[87px] md:pt-[72px] ">
      <div class="rounded-lg border border-gray-200 dark:border-gray-600">
        <p class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600">ICMP State</p>
        <div class="py-3 space-y-3">
          {#each formInputs as entry}
            <Form entry={entry} />
          {/each}
          <hr class="h-px my-8 bg-gray-200 border-0 dark:bg-gray-700">
          <Services services={serviceOptions}/>
        </div>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-600 md:min-w-40">
        <p class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600">Link State</p>
        <div class="py-3 space-y-3">
          {#each link as entry}
            <div class="flex items-center justify-between px-4">
              <div>
                <p class="text-gray-900 dark:text-white">{entry.name}</p>
                <p class="text-gray-900 dark:text-white text-xs">{entry.src} &mdash; {entry.dest}</p>
              </div>
              <div class="flex">
                <label class="inline-flex items-center cursor-pointer">
                  <input type="checkbox" value="" class="sr-only peer">
                  <div class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                </label>
              </div>
            </div>
          {/each}
        </div>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-600 md:min-w-40">
        <p class="px-4 py-2 text-center rounded-t-lg font-bold text-gray-900 dark:text-white bg-gray-200 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600">ANYSec State</p>
        <div class="py-3 space-y-3">
          <Services services={serviceOptions}/>
        </div>
      </div>
    </div>
  </aside>
</div>
