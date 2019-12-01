<script context='module'>
    import proto from './articles_service_grpc_web_pb.js'
    import wrappers from 'google-protobuf/google/protobuf/wrappers_pb.js'

    export const Articles = () => {
        let $ = {}

        const client = new proto.ArticlesServicePromiseClient('http://localhost:8080')

        $.add = async (url) => {
            const request = new proto.CreateArticleRequest()
                .setArticle(new proto.Article()
                    .setUrl(url))
            const response = await client.createArticle(request)
            return response.toObject()
        }

        $.get = async (name) => {
            const request = new proto.GetArticleRequest()
                .setName(name)
            const response = await client.getArticle(request)
            return response.toObject()
        }

        $.delete = async (name) => {
            const request = new proto.DeleteArticleRequest()
                .setName(name)
            const response = await client.deleteArticle(request)
            return response.toObject()
        }

        $.update = async (article, mask) => {
            return {}
            //return await api.patch(`/api/v1/${article.name}?update_mask=${mask}`, article)
        }

        $.next = async (pageSize, from, params) => {
            const request = new proto.ListArticlesRequest()
            if (from) request.setPageToken(from)
            if (params) {
                if (params.isRead !== undefined) request.setIsRead(new wrappers.BoolValue().setValue(params.isRead))
                if (params.isStarred !== undefined) request.setIsStarred(new wrappers.BoolValue().setValue(params.isStarred))
            }
            const response = await client.listArticles(request)
            return response.toObject()
        }

        return $
    }
</script>
