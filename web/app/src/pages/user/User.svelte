<script context="module">
    import { writable } from 'svelte/store'
    import { storage } from './storage'
    import { Router, Route } from 'svelte-routing'

    let allStorage = storage()
    let unreadStorage = storage()
    let starredStorage = storage()
</script>

<script>
    import { onDestroy } from 'svelte'

    import Article from './article/Article.svelte'
    import Articles from './articles/Articles.svelte'
    import Header from './header/Header.svelte'

    import timestamp from 'google-protobuf/google/protobuf/timestamp_pb.js'

    export let user

    export let articles

    const onAdded = async (url) => {
        const ts = new timestamp.Timestamp()
        ts.setSeconds(new Date() / 1000)

        const mock = new articles.proto.Article()
        mock.setUrl(url)
        mock.setTitle(url)
        mock.setIsRead(false)
        mock.setIsFavorite(false)
        mock.setCreateTime(ts)
        mock.setName(Math.random())

        allStorage.add(mock)
        unreadStorage.add(mock)

        let article = await articles.add(user, url)

        allStorage.delete(mock.getName())
        unreadStorage.delete(mock.getName())

        allStorage.add(article)
        unreadStorage.add(article)
    }

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
        username={user}
        on:added={(e) => onAdded(e.detail)}
        on:search={(e) => console.log('search')}
    />
    <Router>
        <Route path="starred">
            <Articles
                itemsStore={starredStorage.store}
                let:item={article}
                on:loadmore={(e) => starredStorage.loadMoreArticles(user, articles, e.detail, {'isStarred': true}) }
            >
                <Article
                    on:deleted={(e) => onDeleted(e.detail)}
                    on:updated={(e) => onUpdated(e.detail)}
                    {article}
                />
            </Articles>
        </Route>
        <Route path="*">
            <Articles
                itemsStore={unreadStorage.store}
                let:item={article}
                on:loadmore={(e) => unreadStorage.loadMoreArticles(user, articles, e.detail, {'isRead': false}) }
            >
                <Article
                    on:deleted={(e) => onDeleted(e.detail)}
                    on:updated={(e) => onUpdated(e.detail)}
                    articles={articles}
                    {article}
                />
            </Articles>
        </Route>
        <Route path="all">
            <Articles
                itemsStore={allStorage.store}
                let:item={article}
                on:loadmore={(e) => allStorage.loadMoreArticles(user, articles, e.detail) }
            >
                <Article
                    on:deleted={(e) => onDeleted(e.detail)}
                    on:updated={(e) => onUpdated(e.detail)}
                    articles={articles}
                    {article}
                />
            </Articles>
        </Route>
    </Router>
</div>

<style>
    .articles {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
    }
</style>
