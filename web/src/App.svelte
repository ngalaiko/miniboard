<script>
    import Login from './pages/login/Login.svelte';
    import { Api } from './components/api/api';

    let api = new Api();
</script>

<div class="app">
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
