<script>
    import Articles  from './components/articles/Articles.svelte'
    import NotFound from './components/notfound/NotFound.svelte'
    import { Articles as ArticlesClient } from './client/articles/Articles.svelte'
    import { Codes as CodesClient } from './client/codes/Codes.svelte'
    import { Users } from './client/users/Users.svelte'
    import { Tokens } from './client/tokens/Tokens.svelte'
    import navaid from 'navaid'
    import LoginForm from './components/loginform/LoginForm.svelte'
    import Reader from './components/reader/Reader.svelte'
    import Codes from './components/codes/Codes.svelte'

    const apiUrl = 'http://localhost:8080'

    let router = navaid()
    let component
    let props

    const users = Users(apiUrl)
    const articles = ArticlesClient(apiUrl)
    const codes = CodesClient(apiUrl)
    const tokens = Tokens(apiUrl)

    router
        .on('/', () => {
            component = LoginForm
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
            navigator.serviceWorker.register('/serviceWorker.js').then((registration) => {
                console.log('ServiceWorker registration successful')
            }, (err) => {
                console.log(`ServiceWorker registration failed: ${err}`)
            })
        })
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
</style>
