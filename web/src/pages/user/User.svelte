<script lang="ts">
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
  import { Router, Route } from 'svelte-routing'
  // @ts-ignore
  import { ArticlesClient, ListParams } from '../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient } from '../../clients/sources.ts'

  export let articlesClient: ArticlesClient
  export let sourcesClient: SourcesClient

  let selectedArticleName: string = ''
  $: selectedArticleName = location.hash.slice(1)
</script>

<div class="user">
  <Router>
    <div class="menu column">
      <Menu />
    </div>
    <div class="list column">
      <Route path="all" let:params>
        <List
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          listParams={new ListParams()}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </Route>
      <Route path="unread" let:params>
        <List
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          listParams={new ListParams().withRead(false)}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </Route>
      <Route path="favorite" let:params>
        <List
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          listParams={new ListParams().withFavorite(true)}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </Route>
    </div>
    <div class="reader column">
      <Reader
        articleName={selectedArticleName}
        articlesClient={articlesClient}
      />
    </div>
  </Router>
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
    max-width:  65%;
    min-width:  65%;
  }

  .column {
    border-left: 1px solid;
  }
</style>
