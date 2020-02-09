<script>
    import User  from './pages/user/User.svelte'
    import Codes from './pages/codes/Codes.svelte'
    import Login from './pages/login/Login.svelte'
    import NotFound from './pages/notfound/NotFound.svelte'
    import Reader from './pages/reader/Reader.svelte'
    import { Router, Route, navigate } from 'svelte-routing'

    import { ApiClient } from './clients/api.js'
    import { ArticlesClient } from './clients/articles.js'
    import { CodesClient } from './clients/codes.js'
    import { UsersClient } from './clients/users.js'
    import { TokensClient } from './clients/tokens.js'
    import { SourcesClient } from './clients/sources.js'

    const apiClient = new ApiClient()
    const usersClient = new UsersClient(apiClient)
    const articlesClient = new ArticlesClient(apiClient)
    const codesClient = new CodesClient(apiClient)
    const tokensClient = new TokensClient(apiClient)
    const sourcesClient = new SourcesClient(apiClient)

    export let url = ""

    if ('serviceWorker' in navigator) {
        window.addEventListener('load', () => {
            navigator.serviceWorker.register(`/sw.js`, {
                scope: '/'
            }).then((registration) => {
                console.log('ServiceWorker registration successful')
            }, (err) => {
                console.log(`ServiceWorker registration failed: ${err}`)
            })
        })
    }

    if (location.pathname == "/") {
        usersClient.me()
            .then(user => navigate(`/${user.getName()}/unread`))
            .catch(e => { /* ignore */ })
    }
</script>

<div class="app">
    <Router url="{url}">
        <Route path="/codes/:code" let:params>
            <Codes tokensClient={tokensClient} code="{params.code}" />
        </Route>
        <Route path="/users/:userid/articles/:articleid" let:params>
            <Reader
                name="users/{params.userid}/articles/{params.articleid}"
                articlesClient={articlesClient}
            />
        </Route>
        <Route path="/users/:userid/*" let:params>
            <User
                user="users/{params.userid}"
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
        padding-left: 5px;
        padding-right: 5px;
        font-family: -apple-system, BlinkMacSystemFont, helvetica neue, Helvetica, Arial, sans-serif;
        max-width: 800px;
        margin: auto;
    }

    :global(body) {
        margin: 0px;
    }
</style>
