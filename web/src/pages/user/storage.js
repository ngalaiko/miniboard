import { writable } from 'svelte/store'

export const storage = () => {
    let $ = {}

    const articlesListStore = writable([])

    $.add = async (article) => {
        articlesListStore.update(list => {
            if (list === undefined) list = []

            for (let i in list) {
                if (list[i].getName() == article.getName()) return list
            }

            return [article].concat(list).sort((a, b) => {
                const d1 = a.getCreateTime()
                const d2 = b.getCreateTime()
                return d1 < d2 ? 1 : -1
            })
        })
    }
    $.delete = async (name) => {
        articlesListStore.update(list => list.filter(a => a.getName() != name))
    }

    $.update = async (updated) => {
        articlesListStore.update(list => {
            for (let i in list) {
                if (list[i].getName() == updated.getName()) {
                    list[i] = updated
                    break
                }
            }
            return list
        })
    }

    let from = ''
    $.loadMoreArticles = async (articlesClient, username, pageSize, params) => {
        if (from === undefined) return []

        let resp = await articlesClient.list(username, pageSize, from, params)

        from = resp.getNextPageToken()

        if (from == '') from = undefined

        if (resp.getArticlesList().length == 0) return 

        resp.getArticlesList().forEach(article => $.add(article))
    }

    $.store = articlesListStore

    return $
}
