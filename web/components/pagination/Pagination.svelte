<script>
    import { createEventDispatcher } from 'svelte'
    import { onDestroy } from 'svelte'
    import { onMount } from 'svelte'

    const dispatch = createEventDispatcher()

    // writable store with list items
    // https://svelte.dev/docs#writable
    export let itemsStore

    let items = []
    let pageSize = 8
    export let pageStart = 0

    const unsubscribeItems = itemsStore.subscribe(value => {
        Array.isArray(items) ? items = value : []
    })

    function loadMore() {
        dispatch('loadmore', pageSize * 2)
    }

    onDestroy(() => {
        unsubscribeItems()
    })

    onMount(() => {
        loadMore()
    })

    function previousPage() {
        pageStart -= pageSize
        dispatch('pagestart', pageStart)
    }

    function nextPage() {
        loadMore()
        pageStart += pageSize
        dispatch('pagestart', pageStart)
    }
</script>

<div class='page'>
    <div class='pagination'>
        {#if pageStart > 0}
            <button class="button-pagination button-previous" on:click|preventDefault={previousPage} >previous</button>
        {/if}
        <div />
        {#if items.length > pageStart + pageSize }
            <button class="button-pagination button-next"  on:click|preventDefault={nextPage} >next</button>
        {/if}
    </div>
    <div class='list'>
        {#each items.slice(pageStart, pageStart + pageSize) as item }
			<slot
				item={item} 
			/>
        {/each}
    </div>
</div>

<style>
    .page {
        margin-bottom: 80px;
    }

    .list {
        display: flex;
        flex-direction: column;
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
