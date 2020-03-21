<script lang="ts">
  import { SearchIcon } from 'svelte-feather-icons'
  import { XIcon } from 'svelte-feather-icons'

  export let value: string = ''
</script>

<div class="search">
  <div class="icon">
    <SearchIcon size="15" />
  </div>
  <input
    class="search-input"
    placeholder="Search articles"
    bind:value={value}
    on:change
    on:input
    on:cut
    on:copy
    on:paste
  />
  <button style="{value ? '' : 'display: none;'}" class="search-close" on:click={() => value = ''}>
    <XIcon size="15"/>
  </button>
</div>

<style>
  .search {
    display: flex;
    padding: 0 0.2em;
    align-items: center;
    border: 1px solid black;
    border-radius: 0.3em;
    margin-top: 0.2em;
  }

  .icon {
    display: flex;
    align-items: center;
  }

  .search-input {
    font: inherit;
    border: 0;
    margin: 0;
    background: inherit;
    padding-left: 0.2em;
    width: 100%;
  }

  .search-close {
    display: flex;
    align-items: center;
    font: inherit;
    padding: 0;
    margin: 0;
    border: 0;
    background: inherit;
    cursor: pointer;
  }

  .search-close:active {
    color: inherit;
  }

  .search-input:focus {
    outline: none;
  }
</style>
