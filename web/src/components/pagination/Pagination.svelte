<script>
    import PullToRefresh from '../pulltorefresh/PullToRefresh.js'
    import { createEventDispatcher } from 'svelte'
    import {
        onMount,
        onDestroy,
    } from 'svelte'

    const dispatch = createEventDispatcher()

    // writable store with list items
    // https://svelte.dev/docs#writable
    export let itemsStore

    let items = []
    let pageSize = 8

    const unsubscribeItems = itemsStore.subscribe(value => {
        Array.isArray(items) ? items = value : []
    })

    const loadMore = () => dispatch('loadmore', pageSize * 2)

    const ptr = PullToRefresh()
    ptr.mainElement = '#list'
    ptr.shouldPullToRefresh = () => {
        return document.querySelector('#list').scrollTop === 0
    }

    onDestroy(() => {
        ptr.destroy()
        unsubscribeItems()
        window.onscroll = () => {}
    })

    onMount(() => {
        loadMore()
        ptr.init()
    })

    const isBottom = () => {
        return (window.innerHeight + window.scrollY) >= document.body.offsetHeight
    }

    window.onscroll = () => {
        if (!isBottom()) {
            return
        }
        loadMore()
    }
</script>

<div id='page'>
    <div id='list'>
        {#each items as item, i (item.getName()) }
            <slot
                item={item}
            />
        {/each}
    </div>
</div>

<style>
    #page {
        display: flex;
        flex-direction: column;
        height: 95%;
    }

    #list {
        overflow-y: scroll;
    }

    .pagination {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
    }

    .button-pagination {
        background: inherit;
        -webkit-appearance: none;
        -moz-appearance: none;
        font-size: 1.1em;
        cursor: pointer;
        border: 0px;
    }

    .button-next::after {
        content: " »";
    }

    .button-next {
        align-self: flex-end;
    }

    .button-previous {
        align-self: flex-start;
    }

    .button-previous::before {
        content: "« ";
    }

    button:hover, button:focus, input:hover, input:focus {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
