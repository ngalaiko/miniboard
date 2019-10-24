<script>
    export let authorizations
    export let users
    export let codes
    export let router

    document.title = "Miniboard"

    let email = ''
    let message = ''
    let error = ''

    const handleClick = async () => {
        if (email == '' ) {
            return
        }
        let resp = codes.sendCode(email)
        message = 'check your email'
        error = ''
    }
</script>

<form class='form'>
    <input name='email' type='email' bind:value={email} placeholder='email' required tabindex='1' />
    <button on:click|preventDefault={handleClick} />
    {#if error != ''}
        <div class='alert'>{error}</div>
    {/if}
    {#if message != ''}
        <div class='info'>{message}</div>
    {/if}
</form>

<style>
    .form {
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
    }

    input:focus{
        outline: none;
        outline-width: 0;
    }

    .info {
        text-align: center;
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
