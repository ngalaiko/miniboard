<script lang="ts">
  import Articles  from './pages/articles/Articles.svelte'
  import Codes from './pages/codes/Codes.svelte'
  import Login from './pages/login/Login.svelte'
  import NotFound from './pages/notfound/NotFound.svelte'
  import { Router, Route, navigate } from 'svelte-routing'

  // @ts-ignore
  import { ApiClient } from './clients/api.ts'
  // @ts-ignore
  import { ArticlesClient } from './clients/articles.ts'
  // @ts-ignore
  import { CodesClient } from './clients/codes.ts'
  // @ts-ignore
  import { UsersClient, User } from './clients/users.ts'
  // @ts-ignore
  import { TokensClient } from './clients/tokens.ts'
  // @ts-ignore
  import { SourcesClient } from './clients/sources.ts'

  const apiClient = new ApiClient()
  const usersClient = new UsersClient(apiClient)
  const articlesClient = new ArticlesClient(apiClient)
  const codesClient = new CodesClient(apiClient)
  const tokensClient = new TokensClient(apiClient)
  const sourcesClient = new SourcesClient(apiClient)

  if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
      navigator.serviceWorker.register(`/sw.js`, {
        scope: '/'
      }).then(() => {
        console.log('ServiceWorker registration successful')
      }, (err: Error) => {
        console.log(`ServiceWorker registration failed: ${err}`)
      })
    })
  }

  if (location.pathname == "/") {
    usersClient.me()
      .then((user: User) => navigate(`/${user.name}/articles`))
      .catch(() => { /* ignore */ })
  }
</script>

<svelte:head>
  <title>Miniboard</title>
</svelte:head>

<div class="app">
  <Router>
    <Route path="/codes/:code" let:params>
      <Codes tokensClient={tokensClient} code="{params.code}" />
    </Route>
    <Route path="/users/:userid/articles/*" let:params>
      <Articles
        articlesClient={articlesClient}
        sourcesClient={sourcesClient}
      />
    </Route>
    <Route path="/">
      <Login codesClient={codesClient} />
    </Route>
    <Route path="*" component={NotFound} />
  </Router>
</div>

<style>
  .app {
    display: flex;
    height: 100%;
    font-family: Helvetica neue, Helvetica, Arial, sans-serif;
    background: #F5F5F5;
    font-size: 18px;
  }

  :global(html) {
    height: 100%;
  }

  :global(body) {
    margin: 0;
    height: 100%;
    overflow: hidden;
  }
</style>
