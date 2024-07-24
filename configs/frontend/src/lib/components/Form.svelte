<script lang="ts">
  import type { FormInput } from '$lib/interfaces';

  export let entry: FormInput;

  let x = entry.default;
  const min = entry.min;
  const max = entry.max;

  const valid = [
    "bg-gray-50", "dark:bg-gray-700", 
    "text-gray-900", "dark:text-white",
    "border-gray-300", "dark:border-gray-600"
  ];
  const invalid = [
    "bg-red-500", "dark:bg-red-600", 
    "text-white", "dark:text-white",
    "border-red-500", "dark:border-red-600"
  ];
  
  const validate = (e) => x = e.target.value; 
  const inRange = (x: number, min: number, max: number) => ((x - min) * (x - max) <= 0);
</script>

<div class="grid grid-cols-2 col-span-2 items-center px-4">
  <label for="{entry.id}" class="text-sm text-gray-900 dark:text-white">{entry.label}</label>
  <input type="number" id="{entry.id}" 
    min="{entry.min}" max="{entry.max}" step="{entry.step}" 
    value="{entry.default}" on:keyup={event => validate(event)}
    class="px-2 py-1 text-sm rounded-lg border focus:ring-0 focus:ring-offset-0 focus:shadow-none
      {inRange(x, min, max) ? valid.join(' ') : invalid.join(' ')}" />
</div>