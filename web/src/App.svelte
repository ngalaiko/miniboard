<script>
    import Articles, { add, search, showSearch } from './components/articles/Articles.svelte'
    import NotFound from './components/notfound/NotFound.svelte'
    import Client from './client/client'
    import Router from './components/router/navaid'
    import LoginForm from './components/loginform/LoginForm.svelte'
    import Header from './components/header/Header.svelte'
    import Reader from './components/reader/Reader.svelte'

    let client = new Client()
    let api = client.api

    let component
    let props

    let router = new Router()
    router
        .on('/', () => {
            component = LoginForm
            props = {
                api: api,
                authorizations: client.authorizations,
                users: client.users,
                router: router,
            }
        })
        .on('/users/:username', () => {
            component = Articles
            props = {
                api: api,
                articles: client.articles,
                labels: client.labels,
            }
        })
        .on('/users/:username/articles/:articleid', (params) => {
            component = Reader
            props = {
                articles: client.articles,
                name: `users/${params.username}/articles/${params.articleid}`,
            }
        })
        .on('*', () => {
            component = NotFound
        })
        .listen()

    if (!api.authorized()) {
        router.route('/')
    }
</script>

<div class="app">
    {#if api.authorized()}
        <Header
            api={api}
            router={router}
            on:added={(e) => add(e.detail, client.articles.add)}
            on:showsearch={(e) => showSearch(e.detail)}
            on:search={(e) => search(e.detail, client.articles.search)}
        />
    {/if}
    <svelte:component this={component} {...props} />
</div>

<style>
    .app {
        padding-left: 5px;
        padding-right: 5px;
        font-family: helvetica neue, Helvetica, Arial, sans-serif;
        text-rendering: optimizeLegibility;
        max-width: 750px;
        margin: auto;
    }
</style>
