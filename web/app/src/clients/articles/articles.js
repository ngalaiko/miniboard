import proto from './articles_service_grpc_web_pb.js'
import wrappers from 'google-protobuf/google/protobuf/wrappers_pb.js'
import updateMask from 'google-protobuf/google/protobuf/field_mask_pb.js'

export class ArticlesClient {
    constructor(hostname) {
        this.client = new proto.ArticlesServicePromiseClient(hostname)
        this.proto = proto
    }

    async get(name) {
        const request = new proto.GetArticleRequest()
        request.setName(name)
        request.setView(2)
        return await this.client.getArticle(request)
    }

    async update(article) {
        const request = new proto.UpdateArticleRequest()
        request.setArticle(article)
        request.setUpdateMask(new updateMask.FieldMask().addPaths('is_read').addPaths('is_favorite'))
        return await this.client.updateArticle(request)
    }

    async delete(name) {
        const request = new proto.DeleteArticleRequest()
        request.setName(name)
        return await this.client.deleteArticle(request)
    }

    async list(pageSize, from, params) {
        const request = new proto.ListArticlesRequest()
        request.setPageSize(pageSize)
        if (from) request.setPageToken(from)
        if (params) {
            if (params.isRead !== undefined) request.setIsRead(new wrappers.BoolValue().setValue(params.isRead))
            if (params.isStarred !== undefined) request.setIsFavorite(new wrappers.BoolValue().setValue(params.isStarred))
        }
        return await this.client.listArticles(request)
    }
}
