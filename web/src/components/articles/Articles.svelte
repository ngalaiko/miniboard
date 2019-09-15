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

    const pageStartStore = writable(0)
</script>

<script>
    import { get } from 'svelte/store'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'

    export let api
    export let articles

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
    <Pagination
        itemsStore={articlesListStore}
        let:item={article}
        pageStart={get(pageStartStore)}
        on:loadmore={(e) => loadMoreArticles(e.detail) }
        on:pagestart={(e) => pageStartStore.set(e.detail)}
    >
        <Article
            on:deleted={(e) => onDeleted(e.detail)}
            articles={articles}
            on:labeladded
            {...article}
        />
    </Pagination>
</div>
