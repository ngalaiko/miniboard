<script lang="ts">
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

  .list {
    flex-basis: 25%;
    max-width:  25%;
    min-width:  25%;
  }

  .reader {
    max-width:  75%;
    min-width:  75%;
  }

  .column {
    border-left: 1px solid;
  }

  @media screen and (max-width: 414px) {
    .list {
      max-width:  100%;
      min-width:  100%;
    }

    .reader {
      display: none;
    }
  }
</style>
