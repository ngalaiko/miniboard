<script lang="ts">
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'

  // @ts-ignore
  import { ArticlesClient } from '../../clients/articles.ts'

  export let username: string
  export let articlesClient: ArticlesClient

  let selected: string|null = null
</script>

<div class="user">
  <div class="menu column">
    <Menu />
  </div>
  <div class="list column">
    <List 
      username={username}
      articlesClient={articlesClient}
      on:selected={(e) => selected = e.detail}
    />
  </div>
  <div class="reader column">
    <Reader
      articleName={selected}
      articlesClient={articlesClient}
    />
  </div>
</div>

<style>
  .user {
    display: flex;
    max-width: 100%;
    min-width: 100%;
  }

  .menu {
    flex-basis: 15%;
    max-width:  15%;
    min-width:  15%;
  }

  .list {
    flex-basis: 20%;
    max-width:  20%;
    min-width:  20%;
  }

  .reader {
  }

  .column {
    border-left: 1px solid;
  }
</style>
