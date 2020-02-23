<script lang="ts">
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
  import { Router, Route, navigate } from 'svelte-routing'

  // @ts-ignore
  import { ArticlesClient, ListParams } from '../../clients/articles.ts'

  export let username: string
  export let articlesClient: ArticlesClient

  let selectedArticleName: string|null = null
</script>

<div class="user">
  <div class="menu column">
    <Menu
      on:unread={() => { navigate(`/${username}/unread`) }}
      on:favorite={() => { navigate(`/${username}/favorite`) }}
      on:all={() => { navigate(`/${username}/all`) }}
    />
  </div>
  <div class="list column">
    <Router>
      <Route path="all">
        <List
          username={username}
          articlesClient={articlesClient}
          listParams={new ListParams()}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </Route>
      <Route path="favorite">
        <List
          username={username}
          articlesClient={articlesClient}
          listParams={new ListParams(true, undefined)}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </Route>
      <Route path="*"> <!-- unread -->
        <List
          username={username}
          articlesClient={articlesClient}
          listParams={new ListParams(undefined, false)}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </Route>
    </Router>
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
    max-width:  65%;
    min-width:  65%;
  }

  .column {
    border-left: 1px solid;
  }
</style>
