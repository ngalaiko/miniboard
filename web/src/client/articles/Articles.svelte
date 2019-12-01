<script context='module'>
    import proto from './articles_service_grpc_web_pb.js'
    import wrappers from 'google-protobuf/google/protobuf/wrappers_pb.js'

    export const Articles = async (api) => {
        let $ = {}

        const client = new proto.ArticlesServiceClient('http://localhost:8080')

        $.add = async (url) => {
            const request = new proto.CreateArticleRequest()
                .setArticle(new proto.Article()
                    .setUrl(url))

            let error
            let response
            await client.createArticle(request, {}, (err, resp) => {
                response = resp
                error = err
            })

            if (error) throw error

            return response
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
            const request = new proto.ListArticlesRequest()
            if (from) request.setPageToken(from)
            if (params) {
                if (params.isRead !== undefined) request.setIsRead(new wrappers.BoolValue().setValue(params.isRead))
                if (params.isStarred !== undefined) request.setIsStarred(new wrappers.BoolValue().setValue(params.isStarred))
            }

            let error
            let response
            await client.listArticles(request, {}, (err, resp) => {
                response = resp
                error = err
            })

            if (error) throw error

            return response
        }

        return $
    }
</script>
