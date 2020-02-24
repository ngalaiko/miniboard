<script lang="ts">
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
  import { Router, Route } from 'svelte-routing'

  // @ts-ignore
  import { ArticlesClient, ListParams } from '../../clients/articles.ts'

  export let username: string
  export let articlesClient: ArticlesClient

  let selectedArticleName: string|null = null
</script>

<div class="user">
  <Router>
    <Route path="all">
      <div class="menu column">
        <Menu
          current='all'
        />
      </div>
      <div class="list column">
        <List
          username={username}
          articlesClient={articlesClient}
          listParams={new ListParams()}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </div>
    </Route>
    <Route path="favorite">
      <div class="menu column">
        <Menu
          current='favorite'
        />
      </div>
      <div class="list column">
        <List
          username={username}
          articlesClient={articlesClient}
          listParams={new ListParams().withFavorite(true)}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </div>
    </Route>
    <Route path="*"> <!-- unread -->
      <div class="menu column">
        <Menu
          current='unread'
        />
      </div>
      <div class="list column">
        <List
          username={username}
          articlesClient={articlesClient}
          listParams={new ListParams().withRead(false)}
          on:selected={(e) => selectedArticleName = e.detail}
        />
      </div>
    </Route>
  </Router>
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
