<script>
    import { slide } from 'svelte/transition'
    import { quintOut } from 'svelte/easing'
    import { createEventDispatcher } from 'svelte'

    const dispatch = createEventDispatcher()

    export let api
    export let router

    function onLogout() {
        api.logout()
        router.route('/')
    }

    let showAdd = false
    let showSearch = false
    let url = ''

    let typingTimerID
</script>

<div class='header'>
    <span class='menu'>
        <span class='menu-left'>
        <button on:click|preventDefault={() => {
            showAdd = !showAdd
            showSearch = false
        }}>add</button>
        <button on:click|preventDefault={() => {
            showSearch = !showSearch
            showAdd = false
        }} class='offset-left'>search</button>
        </span>
        <span class='menu-right'>
            <button on:click|preventDefault={onLogout}>logout</button>
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
                class='offset-left'
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
                bind:value={url}
                placeholder="search..."
                required=""
                on:input={() => {
                    clearTimeout(typingTimerID)
                    typingTimerID = setTimeout(() => {
                        console.log('search', this.value)
                    }, 300)
                }}/>
        </form>
    {/if}
</div>

<style>
    .header {
        display: flex;
        flex-direction: column;
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
        border: 1px solid;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
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

    input:focus {
        outline-width: 0;
    }

    button:hover, button:focus {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
