<script type="ts">
  import Select from 'svelte-select'

  // @ts-ignore
  import { Category, Categories } from './Category.ts'

	import { createEventDispatcher } from 'svelte'

	const dispatch = createEventDispatcher()

  export let selectedValue: Category
</script>

<div class="select">
  <Select 
    items={Object.values(Categories)} 
    bind:selectedValue={selectedValue} 
    isClearable={false} 
    isSearchable={false}
    listAutoWidth={false}
    on:select={(e) => dispatch('select', e.detail.value)}
  />
</div>

<style>
  .select {
    --border: 0;
    --background: #F5F5F5;
    --height: 100%;
    --indicatorWidth: 0;
  }
</style>
