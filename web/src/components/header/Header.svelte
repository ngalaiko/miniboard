<script>
    import { slide } from 'svelte/transition'
    import { quintOut } from 'svelte/easing'
    import { createEventDispatcher } from 'svelte'

    import Add from '../../icons/Add.svelte'
    import Logout from '../../icons/Logout.svelte'

    import Navigation from './navigation/Navigation.svelte'

    const dispatch = createEventDispatcher()

    export let router

    const onLogout = async () => {
        await fetch("/logout")
        router.route('/')
    }

    let showAdd = false
    let url = ''
    let query = ''

    let typingTimerID
</script>

<div class='header'>
    <span class='menu'>
        <span class='menu-left'>
            <button on:click|preventDefault={() => {
                showAdd = !showAdd
            }}>
                <Add />
            </button>
        </span>
        <span class='menu-middle'>
            <Navigation on:selected />
        </span>
        <span class='menu-right'>
            <button on:click|preventDefault={onLogout}>
                <Logout />
            </button>
        </span>
    </span>
    {#if showAdd}
        <form transition:slide='{{ duration: 300, easing: quintOut }}' class='add-form'>
            <input
                class='add-input'
                type="text"
                bind:value={url}
                placeholder="https://..."
                required=""
            />
            <button
                class='button offset-left'
                on:click|preventDefault={() => {
                    dispatch('added', url)
                    url = ''
                }} />
        </form>
    {/if}
</div>

<style>
    .header {
        display: flex;
        flex-direction: column;
    }

    .add-form {
        margin: 0;
        display: flex;
        flex-direction: row;
    }

    .menu {
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

    .menu-middle {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    .menu-rignt {
        display: flex;
        flex-direction: row;
        justify-content: flex-end;
    }

    .add-input {
        border: 0px;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
    }

    span {
        align-content: center;
        display: flex;
    }

    button {
        margin: 0px;
        padding: 0px;
        border: 0px;
        background: inherit;
        -webkit-appearance: none;
        -moz-appearance: none;
        cursor: pointer;
    }

    input:focus {
        outline-width: 0;
    }

    button:hover, .button:focus {
        outline-width: 0;
    }

    button:hover, button:focus {
        outline-width: 0;
    }
</style>
