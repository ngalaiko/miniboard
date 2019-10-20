<script context="module">
    import { writable } from 'svelte/store'
    import { storage } from './storage'

    let allStorage = storage()
    let unreadStorage = storage()
    let starredStorage = storage()

    export const add = async (url, addFunc) => {
        let mock = {
            'url': url,
            'title': url,
            'is_read': false,
            'is_favorite': false,
            'create_time': Date.now(),
            'name': Math.random()
        }

        allStorage.add(mock)
        unreadStorage.add(mock)

        let article = await addFunc(url)

        allStorage.delete(mock.name)
        unreadStorage.delete(mock.name)

        allStorage.add(article)
        unreadStorage.add(article)
    }

    let paneStore = writable('unread')
    export const show = (pane) => paneStore.set(pane)
</script>

<script>
    import { onDestroy } from 'svelte'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'
    import Header from '../header/Header.svelte'

    export let api
    export let articles
    export let router
    export let pane

    onDestroy(paneStore.subscribe(updated => pane = updated))

    const onDeleted = async (name) => {
        allStorage.delete(name)
        unreadStorage.delete(name)
        starredStorage.delete(name)
    }
    const onUpdated = async (updated) => {
        allStorage.update(updated)
        unreadStorage.update(updated)

        if (!updated.is_read) {
            unreadStorage.add(updated)
        } else {
            unreadStorage.delete(updated.name)
        }

        if (updated.is_favorite) {
            starredStorage.add(updated)
        } else {
            starredStorage.delete(updated.name)
        }
    }
</script>

<div>
    <Header
        api={api}
        router={router}
        on:added={(e) => add(e.detail, articles.add)}
        on:search={(e) => console.log('search')}
        on:selected={(e) => show(e.detail)}
    />
    {#if pane === 'starred'}
        <Pagination
            itemsStore={starredStorage.store}
            let:item={article}
            on:loadmore={(e) => starredStorage.loadMoreArticles(articles, e.detail, {'isStarred': true}) }
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
    {#if pane === 'unread'}
        <Pagination
            itemsStore={unreadStorage.store}
            let:item={article}
            on:loadmore={(e) => unreadStorage.loadMoreArticles(articles, e.detail, {'isRead': false}) }
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
