<script context='module'>
    export const Articles = async (api) => {
        let $ = {}

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


        $.update = async (article, mask) => {
            return await api.patch(`/api/v1/${article.name}?update_mask=${mask}`, article)
        }

        $.next = async (pageSize, from, isRead) => {
            let url = `/api/v1/${api.subject()}/articles?page_size=${pageSize}`
            if (from) url += `&page_token=${from}`
            if (isRead !== undefined) url += `&is_read=${isRead}`
            let resp = await api.get(url)
            return resp

        }

        $.search = async (query, limit) => {
            let resp = await api.get(`/api/v1/${api.subject()}/articles:search?query=${query}&page_size=${limit}`)
            return resp.articles
        }

        return $
    }
</script>
