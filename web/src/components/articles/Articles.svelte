<script context="module">
    let articlesList = []
    let pageStart = 0
</script>

<script>
    import Article from '../article/Article.svelte'
    import { onMount } from 'svelte'

    export let api
    export let articles
    export let labels

    let pageSize = 5

    async function loadMore() {
        let list = await articles.next(pageSize * 2)
        articlesList = articlesList.concat(list)
	    pageSize = getPageSize()
    }
    loadMore()

    let url = ''
    async function onAdd() {
        let rnd = Math.random()
        articlesList = [{
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'random': rnd
        }].concat(articlesList)

        let article = await articles.add(url)

        articlesList = [article].concat(articlesList.filter(article => article.random != rnd ))

        url = ''
    }

    async function onDeleted(name) {
        await articles.delete(name)
        articlesList = articlesList.filter(article => article.name != name )
    }

    function previousPage() {
        pageStart -= pageSize
    }

    function nextPage() {
        loadMore().then(() => { pageStart += pageSize })
    }

	// change pageSize on window resize
	window.onresize = function(event) {
		let newSize = getPageSize()
		if (newSize != pageSize) {
			pageSize = newSize
		}
	}

	function getPageSize() {
		let list = document.getElementsByClassName('articles list')[0]
		let size = Math.floor((window.innerHeight - list.offsetTop) / 100)
		return size > 1 ? size : 1
	}
</script>

<div>
    <form>
        <input type="text" bind:value={url} placeholder="https://..." required="" />
        <button class="button-add" on:click|preventDefault={onAdd} />
    </form>
    <div class='pagination'>
        {#if pageStart != 0}
            <button class="button-pagination button-previous" on:click|preventDefault={previousPage} >previous</button>
        {/if}
        <div />
        {#if articlesList.length > pageStart + pageSize }
            <button class="button-pagination button-next"  on:click|preventDefault={nextPage} >next</button>
        {/if}
    </div>
    <div class='articles list'>
        {#each articlesList.slice(pageStart, pageStart+pageSize) as article, i (article.name) }
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                articles={articles}
                labels={labels}
                {...article}
            />
        {/each}
    </div>
</div>

<style>
    .list {
      display: flex;
      flex-direction: column;
    }

    input {
        border: 1px solid;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
    }

    input:focus{
        outline-width: 0
    }

    form {
        display: flex;
        flex-direction: row;
        margin: 0px;
        margin-bottom: 20px;
    }

    .button-add {
        width: 0;
        height: 0;
        padding-left: 0; padding-right: 0;
        border-left-width: 0; border-right-width: 0;
        white-space: nowrap;
        overflow: hidden;
    }

    .pagination {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
	  margin-bottom: 5px;
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

    .button-pagination:hover, .button-pagination:focus {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
