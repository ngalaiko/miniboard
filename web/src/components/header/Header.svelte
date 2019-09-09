<script>
    export let api
    export let router

    function onLogout() {
        api.logout()
        router.route('/')
    }
</script>

<span class='header'>
    <span class='menu-left'>
        <button>add</button>
        <button class='offset-left'>search</button>
    </span>
    <span class='menu-right'>
    <button on:click|preventDefault={onLogout}>logout</button>
    </span>
</span>

<style>
    .header {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        margin-top: 10px;
        margin-bottom: 10px;
    }

    .offset-left {
        margin-left: 5px;
    }

    .menu-left {
        display: flex;
        flex-direction: row;
        justify-content: flex-start;
    }

    .menu-rignt {
        display: flex;
        flex-direction: row;
        justify-content: flex-end;
    }

    button {
        background: inherit;
        -webkit-appearance: none;
        -moz-appearance: none;
        font-size: 1.1em;
        cursor: pointer;
        padding: 3px 10px;
        border: 1px solid;
        border-radius: unset;
    }

    button:hover, button:focus {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
