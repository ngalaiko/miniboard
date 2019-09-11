<script context="module">
    import { writable } from 'svelte/store'

    const articlesListStore = writable([])
</script>

<script>
    import { onDestroy } from 'svelte';
    import Article from '../article/Article.svelte'
    import Header from '../header/Header.svelte'
    import Pagination  from '../pagination/Pagination.svelte'

    export let api
    export let articles
    export let labels

    onDestroy(() => {
        unsubscribeArticlesList()
    })

    async function loadMoreArticles(pageSize) {
        let nextPage = await articles.next(pageSize)
        if (nextPage.length == 0) {
            return 
        }
        articlesListStore.update(list => list.concat(nextPage))
    }

    async function onAdd(e) {
        let url = e.detail
        let rnd = Math.random()

        articlesListStore.update(list => [{
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'random': rnd
        }].concat(list))

        let article = await articles.add(url)

        articlesListStore.update(list => list.filter(article => article.random != rnd))
        articlesListStore.update(list => [article].concat(list))
    }

    async function onDeleted(name) {
        await articles.delete(name)
        articlesListStore.update(list => list.filter(article => article.name != name))
    }
</script>

<div>
    <Header
        on:added={onAdd}
    />
    // todo: connect search to pagination
    <Pagination
        itemsStore={articlesListStore}
        let:item={article}
        on:loadmore={(e) => loadMoreArticles(e.detail) }
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
