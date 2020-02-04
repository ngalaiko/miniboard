<script>
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
    let component

    const unsubscribeItems = itemsStore.subscribe(value => {
        Array.isArray(items) ? items = value : []
    })

    const loadMore = () => dispatch('loadmore', pageSize * 2)

	const onScroll = e => {
        const offset = e.target.scrollHeight - e.target.clientHeight - e.target.scrollTop
        if (offset > 100) {
			return
        }
        loadMore()
    }

    onMount(() => {
        loadMore()

        component.addEventListener("scroll", onScroll)
        component.addEventListener("resize", onScroll)
    })

    onDestroy(() => {
        unsubscribeItems()

        component.addEventListener("scroll", null)
        component.addEventListener("resize", null)
    })
</script>

<div id='page'>
    <div id='list' bind:this={component} >
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
