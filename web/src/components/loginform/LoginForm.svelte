<script>
    import { LoginService } from './login-service';

    export let api;

    let username = "";
    let password = "";

    let loginService = new LoginService(api);

    let error = "";

    function handleClick() {
        error = ""
        if (username == "" || password === "") {
            return
        }
        loginService.login(username, password)
            .then(auth => api.authenticate(auth) )
            .then(() => { location.href = `/users/${username}` } )
            .catch(e => { error = e });
    }
</script>

<form class="form">
    <input type="text" bind:value={username} placeholder="name" required="" />
    <input type="password" bind:value={password} placeholder="password" required="" />
    {#if error != ""}
    <div class="alert">{error}</div>
    {/if}
    <button hidden=true on:click|preventDefault={handleClick} />
</form>

<style>

form {
    margin: 25% auto 0;
    max-width: 250px
}

input {
    border: 1px solid #ccc;
    padding: 3px;
    line-height: 20px;
    width: 250px;
    font-size: 99%;
    margin-bottom: 10px;
    margin-top: 5px;
    -webkit-appearance: none
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
