<script lang="ts">
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
  import { Router, Route } from 'svelte-routing'
  // @ts-ignore
  import { ArticlesClient } from '../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient } from '../../clients/sources.ts'

  export let articlesClient: ArticlesClient
  export let sourcesClient: SourcesClient

  export let username: string = ''

  let selectedArticleName: string = ''
  $: selectedArticleName = location.hash.slice(1)
</script>

<div class="user">
  <div class="list column {selectedArticleName === '' ? 'full-screen' : 'hidden'}">
    <List
      username={username}
      articlesClient={articlesClient}
      sourcesClient={sourcesClient}
      on:select={(e) => selectedArticleName = e.detail}
    />
  </div>
  <div class="reader column { selectedArticleName === '' ? 'hidden' : 'full-screen' }">
    <Reader
      articleName={selectedArticleName}
      articlesClient={articlesClient}
      on:close={() => selectedArticleName = ''}
    />
  </div>
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
