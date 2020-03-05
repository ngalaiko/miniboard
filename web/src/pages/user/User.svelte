<script lang="ts">
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
  import { Router, Route, navigate } from 'svelte-routing'
  // @ts-ignore
  import { ArticlesClient } from '../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient } from '../../clients/sources.ts'
  // @ts-ignore
  import { Categories } from './list/Category.ts'

  export let articlesClient: ArticlesClient
  export let sourcesClient: SourcesClient

  let selectedArticleName: string = ''
  $: selectedArticleName = location.hash.slice(1)
  $: console.log(selectedArticleName)
</script>

<div class="user">
  <Router>
    <div class="list column {selectedArticleName === '' ? 'full-screen' : 'hidden'}">
      <Route path="all" let:params>
        <List
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          category={Categories.All}
          on:selected={(e) => selectedArticleName = e.detail}
          on:select_category={(e) => navigate(e.detail)}
        />
      </Route>
      <Route path="unread" let:params>
        <List
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          category={Categories.Unread}
          on:selected={(e) => selectedArticleName = e.detail}
          on:select_category={(e) => navigate(e.detail)}
        />
      </Route>
      <Route path="favorite" let:params>
        <List
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          category={Categories.Favorite}
          on:selected={(e) => selectedArticleName = e.detail}
          on:select_category={(e) => navigate(e.detail)}
        />
      </Route>
    </div>
    <div class="reader column { selectedArticleName === '' ? 'hidden' : 'full-screen' }">
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
    .full-screen {
      max-width:  100%;
      min-width:  100%;
    }

    .hidden {
      display: none;
    }
  }
</style>
