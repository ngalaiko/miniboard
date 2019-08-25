<script>
    import { LoginService } from './login-service';

    export let api;

    let username = "";
    let password = "";

    let loginService = new LoginService(api);

    function handleClick() {
        loginService.login(username, password);
    }
</script>

<form class="form">
    <input type="text" bind:value={username} placeholder="name" required="" />
    <input type="password" bind:value={password} placeholder="password" required="" />
    <button hidden=true on:click|preventDefault={handleClick}>Login</button>
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

</style>
