<script>
    import Articles from './components/articles/Articles.svelte';
    import NotFound from './components/notfound/NotFound.svelte';
    import Client from './components/client/client';
    import Router from "./components/router/navaid";
    import LoginForm from './components/loginform/LoginForm.svelte'
    import Header from './components/header/Header.svelte'

    let client = new Client()
    let api = client.api

    let component
    let props

    let router = new Router()
    router.on('/', () => {
        component = Articles
        props = {
            api: api,
            articles: client.articles,
            labels: client.labels,
        }
    })
    router.on('/login', () => {
        component = LoginForm
        props = {
            api: api,
            authorizations: client.authorizations,
            users: client.users,
            router: router,
        }
    })
    .on('*', () => {
        component = NotFound
    })
    .listen()

    if (!api.authorized()) {
        router.route("/login")
    }
</script>

<div class="app">
    {#if api.authorized() }
        <Header api={api} router={router} />
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
