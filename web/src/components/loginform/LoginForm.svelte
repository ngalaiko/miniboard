<script>
    import { LoginService } from './login-service'
    import { createEventDispatcher } from 'svelte'

    export let api

    let username = ""
    let password = ""

    let loginService = new LoginService(api)

    let error = ""

    const dispatch = createEventDispatcher()

    function handleClick() {
        error = ""
        if (username == "" || password === "") {
            return
        }
        loginService.login(username, password)
            .then(auth => api.authenticate(auth) )
            .catch(e => { error = e })
            .then(() => dispatch('login', username))
    }
</script>

<form class="form">
    <input type="text" bind:value={username} placeholder="name" required="" />
    <input type="password" bind:value={password} placeholder="password" required="" />
    {#if error != ""}
    <div class="alert">{error}</div>
    {/if}
    <button on:click|preventDefault={handleClick} />
</form>

<style>
    form {
        margin: 25% auto 0;
        max-width: 250px
    }

    button {
        padding-left: 0; padding-right: 0;
        border-left-width: 0; border-right-width: 0;
        white-space: nowrap;
        overflow: hidden;
    }

    input {
        border: 1px solid;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
        margin-bottom: 10px;
    }

    input:focus{
        outline-width: 0
    }

    .alert{
        color: #b94a48;
        background-color: #f2dede;
        border-color: #eed3d7;
        padding: 8px;
        border: 1px solid #fbeed5;
        overflow: auto;
    }
</style>
