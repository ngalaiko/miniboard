<script context="module">
    import { writable } from 'svelte/store'
    import { storage } from './storage'

    let allStorage = storage()
    let unreadStorage = storage()
    let starredStorage = storage()

    let paneStore = writable('unread')
    export const show = (pane) => paneStore.set(pane)
</script>

<script>
    import { onDestroy } from 'svelte'

    import Article from '../article/Article.svelte'
    import Pagination from '../pagination/Pagination.svelte'
    import Header from '../header/Header.svelte'

    import timestamp from 'google-protobuf/google/protobuf/timestamp_pb.js'

    export let user

    export let articles
    export let router
    export let pane

    const onAdded = async (user, url) => {
        let mock = new articles.proto.Article()
            .setUrl(url)
            .setTitle(url)
            .setIsRead(false)
            .setIsFavorite(false)
            .setCreateTime(new timestamp.Timestamp().setSeconds(new Date().getTime()))
            .setName(Math.random())

        allStorage.add(mock)
        unreadStorage.add(mock)

        let article = await articles.add(user, url)

        allStorage.delete(mock.getName())
        unreadStorage.delete(mock.getName())

        allStorage.add(article)
        unreadStorage.add(article)
    }

    onDestroy(paneStore.subscribe(updated => pane = updated))

    const onDeleted = async (name) => {
        await articles.delete(name)

        allStorage.delete(name)
        unreadStorage.delete(name)
        starredStorage.delete(name)
    }
    const onUpdated = async (updated) => {
        await articles.update(updated)

        allStorage.update(updated)
        unreadStorage.update(updated)

        if (!updated.getIsRead()) {
            unreadStorage.add(updated)
        } else {
            unreadStorage.delete(updated.getName())
        }

        if (updated.getIsFavorite()) {
            starredStorage.add(updated)
        } else {
            starredStorage.delete(updated.getName())
        }
    }
</script>

<div class='articles'>
    <Header
        router={router}
        on:added={(e) => onAdded(user, e.detail)}
        on:search={(e) => console.log('search')}
        on:selected={(e) => show(e.detail)}
    />
    {#if pane === 'starred'}
        <Pagination
            itemsStore={starredStorage.store}
            let:item={article}
            on:loadmore={(e) => starredStorage.loadMoreArticles(user, articles, e.detail, {'isStarred': true}) }
        >
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                on:updated={(e) => onUpdated(e.detail)}
                router={router}
                {article}
            />
        </Pagination>
    {/if}
    {#if pane === 'unread'}
        <Pagination
            itemsStore={unreadStorage.store}
            let:item={article}
            on:loadmore={(e) => unreadStorage.loadMoreArticles(user, articles, e.detail, {'isRead': false}) }
        >
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                on:updated={(e) => onUpdated(e.detail)}
                router={router}
                articles={articles}
                {article}
            />
        </Pagination>
    {/if}
    {#if pane === 'all'}
        <Pagination
            itemsStore={allStorage.store}
            let:item={article}
            on:loadmore={(e) => allStorage.loadMoreArticles(user, articles, e.detail) }
        >
            <Article
                on:deleted={(e) => onDeleted(e.detail)}
                on:updated={(e) => onUpdated(e.detail)}
                router={router}
                articles={articles}
                {article}
            />
        </Pagination>
    {/if}
</div>

<style>
    .articles {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
    }
</style>
