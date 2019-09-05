<svelte:window on:pushstate={x => pathname = x.target.location.pathname} on:popstate={x => pathname = x.target.location.pathname}/>

<script>
    import Articles from './components/articles/Articles.svelte';
    import NotFound from './components/notfound/NotFound.svelte';
    import Client from './components/client/client';
    import Router from "./components/router/navaid";
    import LoginForm from './components/loginform/LoginForm.svelte'
    import Header from './components/header/Header.svelte'

    let client = new Client()
    let api = client.api
    let user = null

    let component
    let props

    if (api.authorized()) {
        setUser(api.subject())
    }

    let router = new Router()
    router.on('/', () => {
        component = Articles
        props = {
            api: api,
        }
    })
    .on('*', () => {
        component = NotFound
    })
    .listen()

    function setUser(username) {
        user = api.get(`/api/v1/${username}`)
    }

    function set(obj, key, value) {
        obj[key] = value
        return obj
    }
</script>

<div class="app">
    {#if user == null }
        <LoginForm api={api} on:login={event => setUser(`users/${event.detail}`)} />
    {:else}
        {#await user}
            logging in...
        {:then user}
            <Header api={api} on:logout={() => { user = null } }/>
            <svelte:component this={component} {...set(props, 'user', user)} />
        {/await}
    {/if}
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
