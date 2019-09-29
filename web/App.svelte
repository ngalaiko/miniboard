<script>
    import Articles, { add, show } from './components/articles/Articles.svelte'
    import NotFound from './components/notfound/NotFound.svelte'
    import { Client } from './client/Client.svelte'
    import { Router } from './components/router/Router.svelte'
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
                if (client.api.authorized()) {
                    router.route(`/${client.api.subject()}`)
                    return
                }
                component = LoginForm
                props = {
                    api: client.api,
                    authorizations: client.authorizations,
                    codes: client.codes,
                    users: client.users,
                    router: router,
                }
            })
            .on('/users/:username', () => {
                component = Articles
                props = {
                    api: client.api,
                    articles: client.articles,
                    router: router,
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
                on:search={(e) => console.log('search')}
                on:selected={(e) => show(e.detail)}
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
