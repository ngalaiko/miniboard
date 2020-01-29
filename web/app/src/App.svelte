<script>
    import User  from './pages/user/User.svelte'
    import Codes from './pages/codes/Codes.svelte'
    import Login from './pages/login/Login.svelte'
    import NotFound from './pages/notfound/NotFound.svelte'
    import Reader from './pages/reader/Reader.svelte'
    import { Router, Route, navigate } from 'svelte-routing'

    import { Articles as ArticlesClient } from './clients/articles/Articles.svelte'
    import { Codes as CodesClient } from './clients/codes/Codes.svelte'
    import { Users } from './clients/users/Users.svelte'
    import { Tokens } from './clients/tokens/Tokens.svelte'
    import { SourcesClient } from './clients/sources/sources.js'

    const apiUrl = location.origin

    const users = Users(apiUrl)
    const articles = ArticlesClient(apiUrl)
    const codes = CodesClient(apiUrl)
    const tokens = Tokens(apiUrl)
    const sourcesClient = new SourcesClient(apiUrl)

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
        users.me()
            .then(user => navigate(`/${user.getName()}/unread`))
            .catch(e => { /* ignore */ })
    }
</script>

<div class="app">
    <Router url="{url}">
        <Route path="/codes/:code" let:params>
            <Codes tokens={tokens} code="{params.code}" />
        </Route>
        <Route path="/users/:userid/articles/:articleid" let:params>
            <Reader
                name="users/{params.userid}/articles/{params.articleid}"
                articles={articles}
            />
        </Route>
        <Route path="/users/:userid/*" let:params>
            <User
                user="users/{params.userid}"
                articles={articles}
                sourcesClient={sourcesClient}
            />
        </Route>
        <Route path="/">
            <Login codes={codes} users={users} />
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
