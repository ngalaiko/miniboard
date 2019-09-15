<script>
    import Articles, { add } from './components/articles/Articles.svelte'
    import Search, { search } from './components/search/Search.svelte'
    import NotFound from './components/notfound/NotFound.svelte'
    import { Client } from './client/Client.svelte'
    import Router from './components/router/navaid'
    import LoginForm from './components/loginform/LoginForm.svelte'
    import Header from './components/header/Header.svelte'
    import Reader from './components/reader/Reader.svelte'

    let apiClient
    let router
    let component
    let props

    let clientPromise = Client().then(client => {
        apiClient = client
        router = new Router()

        router
            .on('/', () => {
                component = LoginForm
                props = {
                    api: client.api,
                    authorizations: client.authorizations,
                    users: client.users,
                    router: router,
                }
            })
            .on('/users/:username', () => {
                component = Articles
                props = {
                    api: client.api,
                    articles: client.articles,
                }
            })
            .on('/users/:username/search', () => {
                component = Search
                props = {
                    api: client.api,
                    articles: client.articles,
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

        if (!client.api.authorized()) {
            router.route('/')
        } 
    })
</script>

<div class="app">
    {#await clientPromise}
    {:then}
        {#if apiClient.api.authorized()}
            <Header
                api={apiClient.api}
                router={router}
                on:added={(e) => add(e.detail, apiClient.articles.add)}
                on:search={(e) => search(e.detail, apiClient.articles.search)}
            />
        {/if}
        <svelte:component
            this={component}
            {...props}
        />
    {:catch e}
        {e}
    {/await}
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
