<script context='module'>
    import { IndexedDB } from '../indexeddb/IndexedDB.svelte'

    export const Articles = async (api) => {
        let $ = {}

        let db = await IndexedDB()

        $.add = async (url) => {
            let article = await api.post(`/api/v1/${api.subject()}/articles`, {
                url: url
            })

            try {
                db.add(article)
            } finally {
                return article
            }
        }

        $.get = async (name) => {
            try {
                let article = await db.get(name)
                if (!article.content) throw Error('no content')
                return article
            } catch (e) {
                console.error(e)

                let article = await api.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`)
                try {
                    db.update(article)
                } finally {
                    return article
                }
            }
        }

        $.delete = async (name) => {
            try {
                db.delete(name)
            } finally {
                return await api.delete(`/api/v1/${name}`)
            }
        }


        $.update = async (article, mask) => {
            try {
                db.update(article)
            } finally {
                return await api.patch(`/api/v1/${article.name}?update_mask=${mask}`, article)
            }
        }

        let from = ''
        const listApi = async (pageSize) => {
            // if there are no more articles, return en empty list.
            if (from === undefined) {
                return []
            }

            let resp = await api.get(`/api/v1/${api.subject()}/articles?page_size=${pageSize}&page_token=${from}`)
            from = resp.next_page_token // undefined when no more items

            resp.articles.forEach(article => {
                try {
                    db.add(article)
                } finally {
                    return true
                }
            })

            return resp.articles
        }

        let fromLocal = ''
        const listLocal = async (pageSize) => {
            // if there are no more articles, return en empty list.
            if (fromLocal === undefined) {
                return []
            }

            let articles = []
            let add = fromLocal === ''
            await db.forEach('articles', article => {
                if (article.name == fromLocal) add = true
                if (add) articles = [article].concat(articles)
                return articles.length !== pageSize
            })

            // if we got requested number of articles, need to continue next time
            if (articles.length == pageSize) fromLocal = articles[articles.length-1].name
            // if we got less than requested, there are no more articles.
            if (articles.length < pageSize) fromLocal = undefined

            return articles
        }

        $.next = async (pageSize) => {
            try {
                return await listApi(pageSize)
            } catch (e) {
                return await listLocal(pageSize)
            }
        }

        const searchLocal = async (query, limit) => {
            let articles = []
            query = query.toLowerCase()
            await db.forEach('articles', article => {
                if (article.title.toLowerCase().includes(query)) {
                    articles = [article].concat(articles)
                } else if (article.url.toLowerCase().includes(query)) {
                    articles = [article].concat(articles)
                }
                return articles.length !== limit
            })
            return articles
        }

        const searchApi = async (query, limit) => {
            let resp = await api.get(`/api/v1/${api.subject()}/articles:search?query=${query}&page_size=${limit}`)
            return resp.articles
        }

        $.search = async (query, limit) => {
            try {
                return await searchLocal(query, limit)
            } catch (e) {
                return await searchApi(query, limit)
            }
        }

        return $
    }
</script>
