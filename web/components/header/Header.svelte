<script>
    import { slide } from 'svelte/transition'
    import { quintOut } from 'svelte/easing'
    import { createEventDispatcher } from 'svelte'

    import Add from '../../icons/Add.svelte'
    import Logout from '../../icons/Logout.svelte'
    import Search from '../../icons/Search.svelte'

    const dispatch = createEventDispatcher()

    export let api
    export let router

    function onLogout() {
        api.logout()
        router.route('/')
    }

    let showAdd = false
    export let showSearch = false
    let url = ''
    let query = ''

    let typingTimerID
</script>

<div class='header'>
    <span class='menu'>
        <span class='menu-left'>
        <button on:click|preventDefault={() => {
            showAdd = !showAdd
            showSearch = false
            router.route(`/${api.subject()}`)
            }}>
            <Add />
        </button>
        <button on:click|preventDefault={() => {
            showSearch = !showSearch
            showAdd = false
            showSearch
                ? router.route(`/${api.subject()}/search`)
                : router.route(`/${api.subject()}`)
            }} class='offset-left'>
            <Search />
        </button>
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
                }}>+</button>
        </form>
    {/if}
    {#if showSearch}
        <form transition:slide='{{ duration: 300, easing: quintOut }}' class='search-form'>
            <input
                class='search-input'
                type="text"
                bind:value={query}
                placeholder="search..."
                required=""
                on:input={() => {
                    clearTimeout(typingTimerID)
                    typingTimerID = setTimeout(() => {
                        dispatch('search', query)
                    }, 300)
                }}/>
        </form>
    {/if}
</div>

<style>
    .header {
        display: flex;
        flex-direction: column;
        margin-bottom: 5px;
    }

    .add-form, .search-form {
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

    .menu-rignt {
        display: flex;
        flex-direction: row;
        justify-content: flex-end;
    }

    .add-input, .search-input {
        border: 0px;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
    }

    button {
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

    button:hover {
        text-decoration: underline;
    }
</style>
