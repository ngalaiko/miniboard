<script context="module">
    import { writable } from 'svelte/store'

    const articlesListStore = writable([])

    export const add = async (article) => {
        articlesListStore.update(list => [article].concat(list))
    }
    export const remove = async (name) => {
        articlesListStore.update(list => list.filter(a => a.name != name))
    }
    export const update = async (updated) => {
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

    const fromStore = writable('')
</script>

<script>
    import { onDestroy } from 'svelte'
    import { createEventDispatcher } from 'svelte'

    import Article from '../../article/Article.svelte'
    import Pagination from '../../pagination/Pagination.svelte'

    const dispatch = createEventDispatcher()

    export let articles
    export let router

    let from 
    const unsubscribe = fromStore.subscribe(updated => from = updated)
    onDestroy(unsubscribe)

    const loadMoreArticles = async (pageSize) => {
        if (from === undefined) return []

        let resp = await articles.next(pageSize, from, false)

        if (resp.articles.length == 0) return 

        fromStore.set(resp.next_page_token)
        articlesListStore.update(list => list.concat(resp.articles))
    }
</script>

<div>
    unread
    <Pagination
        itemsStore={articlesListStore}
        let:item={article}
        on:loadmore={(e) => loadMoreArticles(e.detail) }
    >
        <Article
            on:deleted
            on:updated
            router={router}
            articles={articles}
            {...article}
        />
    </Pagination>
</div>
