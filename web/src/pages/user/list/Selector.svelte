<script type="ts">
  // @ts-ignore
  import { Category, Categories } from './Category.ts'

	import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let selectedValue: Category = Categories['all']
</script>

<div class="select">
  <div 
    class="option border border-left { selectedValue.value === Categories['all'].value ? 'selected' : '' }" 
    on:click={() => dispatch('select', Categories['all'].value)}
  >
    All
  </div>
  <div 
    class="option border border-right { selectedValue.value === Categories['unread'].value ? 'selected' : '' }"
    on:click={() => dispatch('select', Categories['unread'].value)}
  >
    Unread
  </div>
</div>

<style>
  .select {
    display: flex;
    justify-content: center;
    padding: 0.2em;
  }

  .option {
    display: flex;
    justify-content: center;
    width: 3em;
    padding: 0.25em 1em;
    cursor: pointer;
  }

  .border {
    border-top: 1px solid black;
    border-bottom: 1px solid black;
  }

  .border-left {
    border-left: 1px solid black;
  }

  .border-right {
    border-right: 1px solid black;
  }

  .selected {
    background: black;
    color: white;
  }
</style>
