<script>
    import Articles  from './components/articles/Articles.svelte'
    import NotFound from './components/notfound/NotFound.svelte'
    import { Client } from './client/Client.svelte'
    import { Router } from './components/router/Router.svelte'
    import LoginForm from './components/loginform/LoginForm.svelte'
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
    })
</script>

<div class="app">
    {#await clientPromise}
    {:then}
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
        font-family: -apple-system, BlinkMacSystemFont, helvetica neue, Helvetica, Arial, sans-serif;
        text-rendering: optimizeLegibility;
        max-width: 800px;
        margin: auto;
    }
</style>
