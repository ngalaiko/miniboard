<script context="module">
    import { writable } from 'svelte/store'

    const articlesListStore = writable([])

    export const add = async (url, addFunc) => {
        let rnd = Math.random()

        articlesListStore.update(list => [{
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'random': rnd
        }].concat(list))

        let article = await addFunc(url)

        articlesListStore.update(list => list.filter(article => article.random != rnd))
        articlesListStore.update(list => [article].concat(list))
    }

    let displaySearchStore = writable(false)
    export const showSearch = (value) => {
        displaySearchStore.set(value)
    }

    const searchListStore = writable([])
    export const search = async (query, searchFunc) => {
        query === ''
            ? searchListStore.set([])
            : searchListStore.set(await searchFunc(query, 100))
    }
</script>

<script>
    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'
    import { onDestroy } from 'svelte'

    export let api
    export let articles
    export let labels

    let displaySearch
    const unsubscribeSearchStore = displaySearchStore.subscribe(value => {
        if (value != displaySearch) {
            displaySearch = value
        }
    })
    onDestroy(() => unsubscribeSearchStore())

    async function loadMoreArticles(pageSize) {
        let nextPage = await articles.next(pageSize)
        if (nextPage.length == 0) {
            return 
        }
        articlesListStore.update(list => list.concat(nextPage))
    }

    async function onDeleted(name) {
        await articles.delete(name)
        articlesListStore.update(list => list.filter(article => article.name != name))
    }
</script>

<div>
    {#if displaySearch}
        <Pagination
            itemsStore={searchListStore}
            let:item={article}
        >
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                articles={articles}
                labels={labels}
                {...article}
            />
        </Pagination>
    {:else}
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
    {/if}
</div>
