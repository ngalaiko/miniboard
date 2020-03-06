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

<div class="user">
  <Router>
    <Route path="*articleName" let:params>
      <div class="list column { params.articleName === '' ? 'full-screen' : 'hidden' }">
        <List
          selectedArticleName="users/{params.userid}/{params.articleName}"
          username="users/{params.userid}"
          articlesClient={articlesClient}
          sourcesClient={sourcesClient}
          on:select={(e) => navigate(`/${e.detail}${location.search}`, { replace: true })}
        />
      </div>
      <div class="reader column { params.articleName === '' ? 'hidden' : 'full-screen' }">
        <Reader
          articleName="users/{params.userid}/{params.articleName}"
          articlesClient={articlesClient}
          on:close={() => navigate(`/users/${params.userid}${location.search}`, { replace: true })}
        />
      </div>
    </Route>
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
