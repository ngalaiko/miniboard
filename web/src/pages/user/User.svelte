<script lang="ts">
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'

  // @ts-ignore
  import { ArticlesClient } from '../../clients/articles.ts'

  export let username: string
  export let articlesClient: ArticlesClient
</script>

<div class="user">
  <div class="menu column">
    <Menu />
  </div>
  <div class="list column">
    <List 
      username={username}
      articlesClient={articlesClient}
    />
  </div>
  <div class="reader column">
    <Reader />
  </div>
</div>

<style>
  .user {
    display: flex;
    width: 100%;
  }

  .list {
    flex-basis: 20%;
  }

  .menu {
    flex-basis: 20%;
  }

  .reader {
    flex-basis: 60%;
  }

  .column {
    border-left: 1px solid;
    border-right: 1px solid;
  }
</style>
