<script lang="ts">
  import { SearchIcon } from 'svelte-feather-icons'

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
</div>

<style>
  .search {
    display: flex;
    padding: 0.2em;
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

  .search-input:focus {
    outline: none;
  }
</style>
