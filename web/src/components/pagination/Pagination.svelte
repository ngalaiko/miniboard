<script context="module">
    import { writable } from 'svelte/store'

    const itemsStore = writable([])
    const pageSizeStore = writable(5)
    const pageStartStore = writable(0) 

    export const addItem = (item) => {
        itemsStore.update(items => [item].concat(items))
    }

    export const deleteItem = (filterFunc) => {
        itemsStore.update(items => items.filter(filterFunc))
    }
</script>

<script>
    import { onDestroy } from 'svelte'

    let items
    let pageSize
    let pageStart

    // (pageSize) => Promise([items])
    export let loadItems

    const unsubscribeItems = itemsStore.subscribe(value => {
        items = value
    })
    const unsubscribePageSize = pageSizeStore.subscribe(value => {
        pageSize = value
    })
    const unsubscribePageStart = pageStartStore.subscribe(value => {
        pageStart = value
    })

    async function loadMore() {
        itemsStore.set(items.concat(await loadItems(pageSize)))
    }
    loadMore().then(updatePageSize)

    onDestroy(() => {
        unsubscribeItems()
        unsubscribePageSize()
        unsubscribePageStart()
    })

    function previousPage() {
        pageStartStore.set(pageStart - pageSize)
        updatePageSize()
    }

    function nextPage() {
        loadMore().then(() => pageStartStore.set(pageStart + pageSize))
        updatePageSize()
    }

	function getPageSize() {
		let list = document.getElementsByClassName('list')[0]
		let size = Math.floor((window.innerHeight - list.offsetTop) / 100)
		return size > 1 ? size : 1
	}

    function updatePageSize() {
		let newSize = getPageSize()
		if (newSize != pageSize) {
			pageSizeStore.set(newSize)
		}
    }
	window.onresize = updatePageSize
</script>

<div>
    <div class='pagination'>
        {#if pageStart != 0}
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
