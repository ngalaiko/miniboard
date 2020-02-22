<script lang="ts">
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'

  // @ts-ignore
  import { ArticlesClient, ListParams } from '../../clients/articles.ts'

  export let username: string
  export let articlesClient: ArticlesClient

  enum Selected {
    All = 0,
    Unread = 1,
    Favorite = 2,
  }

  let selected: Selected = Selected.Unread

  let selectedArticleName: string|null = null
</script>

<div class="user">
  <div class="menu column">
    <Menu
      on:unread={() => { selected = Selected.Unread }}
      on:favorite={() => { selected = Selected.Favorite }}
      on:all={() => { selected =  Selected.All }}
    />
  </div>
  <div class="list column">
    {#if selected == Selected.All}
      <List 
        username={username}
        articlesClient={articlesClient}
        listParams={new ListParams()}
        on:selected={(e) => selectedArticleName = e.detail}
      />
    {:else if selected == Selected.Unread}
      <List 
        username={username}
        articlesClient={articlesClient}
        listParams={new ListParams(undefined, false)}
        on:selected={(e) => selectedArticleName = e.detail}
      />
    {:else if selected == Selected.Favorite}
      <List 
        username={username}
        articlesClient={articlesClient}
        listParams={new ListParams(true, undefined)}
        on:selected={(e) => selectedArticleName = e.detail}
      />
    {/if}
  </div>
  <div class="reader column">
    <Reader
      articleName={selectedArticleName}
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
    max-width:  64%;
  }

  .column {
    border-left: 1px solid;
  }
</style>
