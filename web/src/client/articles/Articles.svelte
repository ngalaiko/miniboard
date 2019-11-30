<script context='module'>
    import proto from './articles_service_grpc_web_pb.js'

        const client = new proto.ArticlesServiceClient('http://localhost:8080')
        const request = new proto.CreateArticleRequest()
        const article = new proto.Article()
        request.setArticle(article)

        const call = client.createArticle(request, {
            'grpc-encoding': 'gzip',
        }, (err, response) => {
            console.log(response);
            console.log(err);
        })

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

        $.next = async (pageSize, from, params) => {
            let url = `/api/v1/${api.subject()}/articles?page_size=${pageSize}`
            if (from) url += `&page_token=${from}`
            if (params) {
                if (params.isRead !== undefined) url += `&is_read=${params.isRead}`
                if (params.isStarred !== undefined) url += `&is_favorite=${params.isStarred}`
            }
            let resp = await api.get(url)
            return resp

        }

        return $
    }
</script>
