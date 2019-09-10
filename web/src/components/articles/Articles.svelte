<script context="module">
    import { writable } from 'svelte/store'
    const articlesListStore = writable([])
    const foundListStore = writable([])
</script>

<script>
    import { onDestroy } from 'svelte';
    import Article from '../article/Article.svelte'
    import Header from '../header/Header.svelte'
    import Pagination, { addItem, deleteItem } from '../pagination/Pagination.svelte'

    export let api
    export let articles
    export let labels

    let articlesList
    let foundList

    const unsubscribeArticlesList = articlesListStore.subscribe(value => {
        articlesList = value
    })
    const unsubscribeFoundList = foundListStore.subscribe(value => {
        foundList = value
    })

    onDestroy(() => {
        unsubscribeArticlesList()
        unsubscribeFoundList()
    })

    async function loadMore(pageSize) {
        return await articles.next(pageSize * 2)
    }

    async function onSearch(e) {
        let query = e.detail
        showSearch = true

        query == ''
            ? foundListStore.set([])
            : foundListStore.set(await articles.search(query, pageSize))
    }

    async function onAdd(e) {
        let url = e.detail
        let rnd = Math.random()

        addItem({
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'random': rnd
        })

        let article = await articles.add(url)

        deleteItem(article => article.random != rnd)
        addItem(article)
    }

    async function onDeleted(name) {
        await articles.delete(name)
        deleteItem(article => article.name != name)
    }

    let showSearch = false
</script>

<div>
    <Header
        on:added={onAdd}
        on:searchstop={() => showSearch = false}
        on:searchstart={onSearch}
    />
    <Pagination
        items={articlesList}
        loadItems={loadMore}
        let:item={article}
    >
        <Article
            on:deleted={(e) => onDeleted(e.detail)}
            articles={articles}
            labels={labels}
            {...article}
        />
    </Pagination>
</div>

<style>
    .add-input {
        border: 1px solid;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
    }

    form {
        display: flex;
        flex-direction: row;
        justify-content: flex-start;
        margin: 0px;
        margin-bottom: 5px;
    }

    .button-add {
        width: 0;
        height: 0;
        padding-left: 0; padding-right: 0;
        border-left-width: 0; border-right-width: 0;
        white-space: nowrap;
        overflow: hidden;
    }

    button:hover, button:focus, input:hover, input:focus {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
