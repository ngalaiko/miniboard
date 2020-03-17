<script lang="ts">
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
  import { Router, Route, navigate } from 'svelte-routing'
  // @ts-ignore
  import { ArticlesClient } from '../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient } from '../../clients/sources.ts'

  export let articlesClient: ArticlesClient
  export let sourcesClient: SourcesClient
</script>

<div class="articles">
  <Router>
    <Route path="*articleid" let:params>
      <div class="list column { params.articleid === '' ? 'full-screen' : 'hidden' }">
        <List
          selectedArticleName="users/{params.userid}/articles/{params.articleid}"
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          on:select={(e) => navigate(`/${e.detail}${location.search}`, { replace: true })}
        />
      </div>
      <div class="reader column { params.articleid === '' ? 'hidden' : 'full-screen' }">
        <Reader
          articleName="users/{params.userid}/articles/{params.articleid}"
          articlesClient={articlesClient}
          on:close={() => navigate(`/users/${params.userid}/articles${location.search}`, { replace: true })}
        />
      </div>
    </Route>
  </Router>
</div>

<style>
  .articles {
    display: flex;
    max-width: 100%;
    min-width: 100%;
  }

  .list {
    flex-basis: 25%;
    max-width:  25%;
    min-width: 9em;
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
