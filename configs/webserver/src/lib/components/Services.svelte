<script lang="ts">
  export let key: string;
  export let services: string[];
</script>

{#each services as svc}
  <div class="flex items-center justify-between px-4">
    <div class="flex">
      <p class="text-gray-900 dark:text-white">{svc}</p>
    </div>
    <div class="flex">
      <label class="inline-flex items-center cursor-pointer">
        <input type="checkbox" class="sr-only peer" id="{key}-{svc.toLowerCase()}">
        <div class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
      </label>
    </div>
  </div>
{/each}