<script>
    import Articles  from './pages/articles/Articles.svelte'
    import Codes from './pages/codes/Codes.svelte'
    import Login from './pages/login/Login.svelte'
    import NotFound from './pages/notfound/NotFound.svelte'
    import Reader from './pages/reader/Reader.svelte'

    import { Articles as ArticlesClient } from './clients/articles/Articles.svelte'
    import { Codes as CodesClient } from './clients/codes/Codes.svelte'
    import { Users } from './clients/users/Users.svelte'
    import { Tokens } from './clients/tokens/Tokens.svelte'

    import navaid from 'navaid'

    const apiUrl = '__API_URL__'

    let router = navaid()
    let component
    let props

    const users = Users(apiUrl)
    const articles = ArticlesClient(apiUrl)
    const codes = CodesClient(apiUrl)
    const tokens = Tokens(apiUrl)

    router
        .on('/', () => {
            component = Login
            props = {
                codes: codes,
                users: users,
                router: router,
            }
        })
        .on('/codes/:code', (params) => {
            component = Codes
            props = {
                tokens: tokens,
                code: params.code,
                router: router,
            }
        })
        .on('/users/:id', (params) => {
            component = Articles
            props = {
                user: `users/${params.id}`,
                articles: articles,
                router: router,
            }
        })
        .on('/users/:id/articles/:articleid', (params) => {
            component = Reader
            props = {
                articles: articles,
                name: `users/${params.id}/articles/${params.articleid}`,
            }
        })
        .on('*', () => {
            component = NotFound
        })
        .listen()

    if ('serviceWorker' in navigator) {
        window.addEventListener('load', () => {
            navigator.serviceWorker.register('/service_worker/service_worker.js').then((registration) => {
                console.log('ServiceWorker registration successful')
            }, (err) => {
                console.log(`ServiceWorker registration failed: ${err}`)
            })
        })
    }

    if (location.pathname == "/") {
        users.me()
            .then(user => router.route(`/${user.getName()}`))
            .catch(e => { /* ignore */ })
    }
</script>

<div class="app">
    <svelte:component
        this={component}
        {...props}
    />
</div>

<style>
    .app {
        display: flex;
        height: 100%;
        padding-left: 5px;
        padding-right: 5px;
        font-family: -apple-system, BlinkMacSystemFont, helvetica neue, Helvetica, Arial, sans-serif;
        text-rendering: optimizeLegibility;
        max-width: 800px;
        margin: auto;
    }

    :global(body) {
        margin: 0px;
    }
</style>
