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

    const fromStore = writable('')
</script>

<script>
    import { onDestroy } from 'svelte'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'

    export let api
    export let articles
    export let router

    let from
    const unsubscribe = fromStore.subscribe(updated => from = updated)
    onDestroy(unsubscribe)

    const loadMoreArticles = async (pageSize) => {
        let resp = await articles.next(pageSize, from)
        if (resp.articles.length == 0) {
            return 
        }
        fromStore.set(resp.next_page_token)
        articlesListStore.update(list => list.concat(resp.articles))
    }

    const onDeleted = async (name) => {
        articlesListStore.update(list => list.filter(article => article.name != name))
    }

    const onUpdated = async (updated) => {
        articlesListStore.update(list => {
            for (let i in list) {
                if (list[i].name == updated.name) {
                    list[i] = updated
                    break
                }
            }
            return list
        })
    }
</script>

<div>
    <Pagination
        itemsStore={articlesListStore}
        let:item={article}
        on:loadmore={(e) => loadMoreArticles(e.detail) }
    >
        <Article
            on:deleted={(e) => onDeleted(e.detail)}
            on:updated={(e) => onUpdated(e.detail)}
            router={router}
            articles={articles}
            {...article}
        />
    </Pagination>
</div>
