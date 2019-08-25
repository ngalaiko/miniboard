<svelte:window on:pushstate={x => pathname = x.target.location.pathname} on:popstate={x => pathname = x.target.location.pathname}/>

<script>
    import Login from './pages/login/Login.svelte';
    import NotFound from './pages/notfound/NotFound.svelte';
    import { Api } from './components/api/api';
    import navid from "./components/navaid/navaid";

    let api = new Api();
    let router = navid();

    let pathname = location.pathname;
</script>

<div class="app">
   // todo: router component
  <Login api={api} />
</div>

<style>

.app {
    padding-left: 5px;
    padding-right: 5px;
    font-family: helvetica neue, Helvetica, Arial, sans-serif;
    text-rendering: optimizeLegibility
}

</style>
