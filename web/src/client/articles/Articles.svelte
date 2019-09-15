<script context='module'>
    import { IndexedDB } from '../indexeddb/IndexedDB.svelte'

    export const Articles = async (api) => {
        let $ = {}

        let from = ''
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
                if (article.content == undefined) {
                    let article = await api.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`)
                    db.add(article)
                    return article
                }
                return article
            } catch (e) {
                console.error(e)
                return await api.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`)
            }
        }

        $.delete = async (name) => {
            try {
                db.delete(name)
            } finally {
                return await api.delete(`/api/v1/${name}`)
            }
        }

        $.next = async (pageSize) => {
            // if there are no more articles, return en empty list.
            if (from === undefined) {
                return []
            }

            let articles = []
            await db.forEach('articles', article => {
                articles = [article].concat(articles)
                return articles.length !== pageSize
            })

            if (articles.length !== 0) return articles

            let resp = await api.get(`/api/v1/${api.subject()}/articles?page_size=${pageSize}&page_token=${from}`)
            from = resp.next_page_token // undefined when no more items

            resp.articles.forEach(article => {
                try {
                    db.add(article)
                } catch (e) {
                    return false
                }
                return true
            })

            return resp.articles
        }

        $.search = async (query, limit) => {
            let resp = await api.get(`/api/v1/${api.subject()}/articles:search?query=${query}&page_size=${limit}`)
            return resp.articles
        }

        return $
    }
</script>
