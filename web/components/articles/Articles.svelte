<script context="module">
    import { writable } from 'svelte/store'
    import { storage } from './storage'

    let allStorage = storage()
    let unreadStorage = storage()

    export const add = async (url, addFunc) => {
        let mock = {
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'name': Math.random()
        }

        allStorage.add(mock)
        unreadStorage.add(mock)

        let article = await addFunc(url)

        allStorage.delete(mock.name)
        unreadStorage.delete(mock.name)

        allStorage.add(article)
        unreadStorage.delete(article)
    }

    let paneStore = writable('unread')
    export const show = (pane) => paneStore.set(pane)
</script>

<script>
    import { onDestroy } from 'svelte'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'

    export let api
    export let articles
    export let router
    export let pane

    onDestroy(paneStore.subscribe(updated => pane = updated))

    const onDeleted = async (name) => {
        allStorage.delete(name)
        unreadStorage.delete(name)
    }
    const onUpdated = async (updated) => {
        allStorage.update(updated)
        unreadStorage.update(updated)
    }
</script>

<div>
    {#if pane === 'unread'}
        <Pagination
            itemsStore={unreadStorage.store}
            let:item={article}
            on:loadmore={(e) => unreadStorage.loadMoreArticles(articles, e.detail, false) }
        >
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                on:updated={(e) => onUpdated(e.detail)}
                router={router}
                articles={articles}
                {...article}
            />
        </Pagination>
    {/if}
    {#if pane === 'all'}
        <Pagination
            itemsStore={allStorage.store}
            let:item={article}
            on:loadmore={(e) => allStorage.loadMoreArticles(articles, e.detail) }
        >
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                on:updated={(e) => onUpdated(e.detail)}
                router={router}
                articles={articles}
                {...article}
            />
        </Pagination>
    {/if}
</div>
