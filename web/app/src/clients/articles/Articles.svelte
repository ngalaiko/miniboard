<script context='module'>
    import proto from './articles_service_grpc_web_pb.js'
    import wrappers from 'google-protobuf/google/protobuf/wrappers_pb.js'
    import updateMask from 'google-protobuf/google/protobuf/field_mask_pb.js'

    export const Articles = (apiUrl) => {
        let $ = {}

        const client = new proto.ArticlesServicePromiseClient(apiUrl)

        $.proto = proto

        $.get = async (name) => {
            const request = new proto.GetArticleRequest()
            request.setName(name)
            request.setView(2)
            return await client.getArticle(request)
        }

        $.delete = async (name) => {
            const request = new proto.DeleteArticleRequest()
            request.setName(name)
            return await client.deleteArticle(request)
        }

        $.update = async (article) => {
            const request = new proto.UpdateArticleRequest()
            request.setArticle(article)
            request.setUpdateMask(new updateMask.FieldMask().addPaths('is_read').addPaths('is_favorite'))
            return await client.updateArticle(request)
        }

        $.next = async (pageSize, from, params) => {
            const request = new proto.ListArticlesRequest()
            request.setPageSize(pageSize)
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
