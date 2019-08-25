<script>
    import { LoginService } from './login-service';

    let loginService = new LoginService();

    let username = "";
    let password = "";

    function handleClick() {
        loginService.login(username, password);
    }
</script>

<form class="form">
    <label for="form-username">Username</label>
    <input type="text" bind:value={username} required="" />
    <label for="form-password">Password</label>
    <input type="password" bind:value={password} required="" />
    <div class="buttons">
        <button on:click|preventDefault={handleClick}>Login</button>
    </div>
</form>

<style>

form {
    margin: 50px auto 0;
    max-width: 280px
}

label {
    cursor: pointer;
    display: block
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

.buttons {
    margin-top: 10px;
    margin-bottom: 20px
}

button {
    display: inline-block;
    -webkit-appearance: none;
    -moz-appearance: none;
    font-size: 1.1em;
    cursor: pointer;
    padding: 3px 10px;
    border: 1px solid;
    border-radius: unset;
    border-color: #3079ed;
    background: #4d90fe;
    color: #fff
}

button:hover, .button:focus {
    border-color: #2f5bb7;
    background: #357ae8
}

</style>
