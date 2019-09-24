<script context="module">
    import { writable } from 'svelte/store'

    const searchListStore = writable([])
    export const search = async (query, searchFunc) => {
        query === ''
            ? searchListStore.set([])
            : searchListStore.set(await searchFunc(query, 100))
    }
</script>

<script>
    import { get } from 'svelte/store'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'

    export let api
    export let articles
    export let router

    async function onDeleted(name) {
        searchListStore.update(list => list.filter(article => article.name != name))
    }
</script>

<div>
    <Pagination
        itemsStore={searchListStore}
        let:item={article}
    >
        <Article
            on:deleted={(e) => onDeleted(e.detail)}
            articles={articles}
            router={router}
            {...article}
        />
    </Pagination>
</div>
