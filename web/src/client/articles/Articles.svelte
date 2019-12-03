<script context='module'>
    import proto from './articles_service_grpc_web_pb.js'
    import wrappers from 'google-protobuf/google/protobuf/wrappers_pb.js'
    import updateMask from 'google-protobuf/google/protobuf/field_mask_pb.js'

    export const Articles = (apiUrl) => {
        let $ = {}

        const client = new proto.ArticlesServicePromiseClient(apiUrl)

        $.proto = proto

        $.add = async (user, url) => {
            const request = new proto.CreateArticleRequest()
                .setParent(user)
                .setArticle(new proto.Article()
                    .setUrl(url))
            return await client.createArticle(request)
        }

        $.get = async (name) => {
            const request = new proto.GetArticleRequest()
                .setName(name)
            return await client.getArticle(request)
        }

        $.delete = async (name) => {
            const request = new proto.DeleteArticleRequest()
                .setName(name)
            return await client.deleteArticle(request)
        }

        $.update = async (article) => {
            const request = new proto.UpdateArticleRequest()
                .setArticle(article)
                .setUpdateMask(new updateMask.FieldMask().addPaths('is_read').addPaths('is_favorite'))
            return await client.updateArticle(request)
        }

        $.next = async (user, pageSize, from, params) => {
            const request = new proto.ListArticlesRequest()
                .setParent(user)
                .setPageSize(pageSize)
            if (from) request.setPageToken(from)
            if (params) {
                if (params.isRead !== undefined) request.setIsRead(new wrappers.BoolValue().setValue(params.isRead))
                if (params.isStarred !== undefined) request.setIsFavorite(new wrappers.BoolValue().setValue(params.isStarred))
            }
            return await client.listArticles(request)
        }

        return $
    }
</script>
