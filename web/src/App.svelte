<svelte:window on:pushstate={x => pathname = x.target.location.pathname} on:popstate={x => pathname = x.target.location.pathname}/>

<script>
    import Login from './pages/login/Login.svelte';
    import User from './pages/user/User.svelte';
    import { Api } from './components/api/api';
    import { Router } from "./components/router/router";
    import NotFound from './pages/notfound/NotFound.svelte'
    import LoginForm from './components/loginform/LoginForm.svelte'
    import Header from './components/header/Header.svelte'

    let api = new Api()
    let user = null

    if (api.authorized()) {
        setUser(api.subject())
    }

    let router = new Router()
    router.register('/', User, {
        api: api,
    })
    router.register('*', NotFound)
    router.listen()

    let pathname = location.pathname

    function setUser(username) {
      user = api.get(`/api/v1/${username}`)
    }

    function set(obj, key, value) {
      obj[key] = value
      return obj
    }
</script>

<div class="app">
  {#each router.current() as { component, props } }
    {#if user == null }
      <LoginForm api={api} on:loggedin={event => setUser(`users/${event.detail}`)} />
    {:else}
      {#await user}
        logging in...
      {:then user}
        <Header api={api} />
        <svelte:component this={component} {...set(props, 'user', user)} />
      {/await}
    {/if}
  {/each}
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
