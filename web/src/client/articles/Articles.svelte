<script context='module'>
    export const Articles = (api) => {
        let $ = {}

        let from = ''

        $.add = async (url) => {
            return await api.post(`/api/v1/${api.subject()}/articles`, {
                url: url
            })
        }

        $.get = async (name) => {
            return await api.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`)
        }

        $.delete = async (name) => {
            return await api.delete(`/api/v1/${name}`)
        }

        $.updateLabels = async (article) => {
            return api.patch(`/api/v1/${article.name}?update_mask=label_ids`, {
                label_ids: article.label_ids,
            })
        }

        $.next = async (pageSize) => {
            // if there are no more articles, return en empty list.
            if (from === undefined) {
                return []
            }
            let resp = await api.get(`/api/v1/${api.subject()}/articles?page_size=${pageSize}&page_token=${from}`)
            from = resp.next_page_token // undefined when no more items
            return resp.articles
        }

        $.search = async (query, limit) => {
            let resp = await api.get(`/api/v1/${api.subject()}/articles:search?query=${query}&page_size=${limit}`)
            return resp.articles
        }

        return $
    }
</script>
