<svelte:window on:pushstate={x => pathname = x.target.location.pathname} on:popstate={x => pathname = x.target.location.pathname}/>

<script>
    import Login from './pages/login/Login.svelte';
    import { Api } from './components/api/api';
    import { Router } from "./components/router/router";

    let api = new Api()

    let router = new Router()

    router.register("/", Login, {
        api: api,
    })

    let pathname = location.pathname;
</script>

<div class="app">
  {#each router.route(pathname) as { component, props } }
    <svelte:component this={ component } { ...props }/>
  {/each}
</div>

<style>

.app {
    padding-left: 5px;
    padding-right: 5px;
    font-family: helvetica neue, Helvetica, Arial, sans-serif;
    text-rendering: optimizeLegibility
}

</style>
