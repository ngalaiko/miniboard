<script lang="ts">
  import { SearchIcon } from 'svelte-feather-icons'

  export let value: string = ''
</script>

<div class="search">
  <button class="search-button">
    <SearchIcon size="15" />
  </button>
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
</div>

<style>
  .search {
    display: flex;
    padding: 0.2em;
  }

  .search-button {
    background: inherit;
    padding: 0;
    border: 0;
    cursor: pointer;
    font: inherit;
  }

  .search-button:focus {
    outline: none;
  }

  .search-input {
    font: inherit;
    border: 0;
    background: inherit;
    width: 100%;
  }

  .search-input:focus {
    outline: none;
  }


</style>
