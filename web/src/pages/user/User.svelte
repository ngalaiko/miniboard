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

    import { Article as article } from '../../clients/articles.js'

    export let user

    export let articlesClient
    export let sourcesClient

    const onAdded = async (url) => {
        const mock = new article()
        mock.setUrl(url)
        mock.setTitle(url)
        mock.setIsRead(false)
        mock.setIsFavorite(false)
        mock.setCreateTime(new Date().toISOString())
        mock.setName(Math.random())

        allStorage.add(mock)
        unreadStorage.add(mock)

        const source = await sourcesClient.createSource(user, url)

        const type = source.getName().replace(user, '')
        switch (true) {
        case type.startsWith('/article'):
            const article = await articlesClient.get(source.getName())
            allStorage.add(article)
            unreadStorage.add(article)
            break;
        }

        allStorage.delete(mock.getName())
        unreadStorage.delete(mock.getName())
    }

    const onDeleted = async (name) => {
        await articlesClient.delete(name)

        allStorage.delete(name)
        unreadStorage.delete(name)
        starredStorage.delete(name)
    }

    const onUpdated = async (updated) => {
        await articlesClient.update(updated)

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
                on:loadmore={(e) => starredStorage.loadMoreArticles(articlesClient, user, e.detail, {'is_favorite': true}) }
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
                on:loadmore={(e) => unreadStorage.loadMoreArticles(articlesClient, user, e.detail, {'is_read': false}) }
            >
                <Article
                    on:deleted={(e) => onDeleted(e.detail)}
                    on:updated={(e) => onUpdated(e.detail)}
                    {article}
                />
            </Articles>
        </Route>
        <Route path="all">
            <Articles
                itemsStore={allStorage.store}
                let:item={article}
                on:loadmore={(e) => allStorage.loadMoreArticles(articlesClient, user, e.detail) }
            >
                <Article
                    on:deleted={(e) => onDeleted(e.detail)}
                    on:updated={(e) => onUpdated(e.detail)}
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
