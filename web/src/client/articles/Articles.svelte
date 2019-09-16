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
                let article =  await api.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`)
                try {
                    db.add(article)
                } finally {
                    return article
                }
            } catch (e) {
                console.error(e)
                return await db.get(name)
            }
        }

        $.delete = async (name) => {
            try {
                db.delete(name)
            } finally {
                return await api.delete(`/api/v1/${name}`)
            }
        }

        let from = ''
        const fromApi = async (pageSize) => {
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
        const fromLocalStorage = async (pageSize) => {
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
                return await fromApi(pageSize)
            } catch (e) {
                return fromLocalStorage(pageSize)
            }
        }

        $.search = async (query, limit) => {
            // TODO: use client side
            let resp = await api.get(`/api/v1/${api.subject()}/articles:search?query=${query}&page_size=${limit}`)
            return resp.articles
        }

        return $
    }
</script>
