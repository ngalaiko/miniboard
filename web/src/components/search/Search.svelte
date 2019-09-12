<script context="module">
    import { writable } from 'svelte/store'

    const searchListStore = writable([])
    export const search = async (query, searchFunc) => {
        query === ''
            ? searchListStore.set([])
            : searchListStore.set(await searchFunc(query, 100))
    }

    const pageStartStore = writable(0)
</script>

<script>
    import { get } from 'svelte/store'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'

    export let api
    export let articles
    export let labels

    async function onDeleted(name) {
        await articles.delete(name)
        searchListStore.update(list => list.filter(article => article.name != name))
    }
</script>

<div>
    <Pagination
        itemsStore={searchListStore}
        let:item={article}
        pageStart={get(pageStartStore)}
        on:pagestart={(e) => pageStartStore.set(e.detail)}
    >
        <Article
            on:deleted={(e) => onDeleted(e.detail)}
            articles={articles}
            labels={labels}
            {...article}
        />
    </Pagination>
</div>
