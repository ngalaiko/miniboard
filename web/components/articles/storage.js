import { writable } from 'svelte/store'

export const storage = () => {
    let $ = {}

    const articlesListStore = writable([])

    $.add = async (article) => {
        articlesListStore.update(list => {
            for (let i in list) {
                if (list[i].name == article.name) return list
            }
            return [article].concat(list).sort((a, b) => {
                return new Date(a.create_time) < new Date(b.create_time) ? 1 : -1
            })
        })
    }
    $.delete = async (name) => {
        articlesListStore.update(list => list.filter(a => a.name != name))
    }
    $.update = async (updated) => {
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

    let from = ''
    $.loadMoreArticles = async (articles, pageSize, isRead) => {
        if (from === undefined) return []

        let resp = await articles.next(pageSize, from, isRead)

        if (resp.articles.length == 0) return 

        from = resp.next_page_token
        articlesListStore.update(list => list.concat(resp.articles))
    }

    $.store = articlesListStore

    return $
}
