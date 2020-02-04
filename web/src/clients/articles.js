import proto from './proto/articles_service_grpc_web_pb.js'
import wrappers from 'google-protobuf/google/protobuf/wrappers_pb.js'
import updateMask from 'google-protobuf/google/protobuf/field_mask_pb.js'
import timestamp from 'google-protobuf/google/protobuf/timestamp_pb.js'

export class Article {
    constructor(protoArticle) {
        if (protoArticle !== undefined) {
            this.proto = protoArticle
        } else {
            this.proto = new proto.Article()
        }
    }

    setUrl(url) {
        this.proto.setUrl(url)
    }

    setTitle(title) {
        this.proto.setTitle(title)
    }

    setIsRead(isRead) {
        this.proto.setIsRead(isRead)
    }

    setIsFavorite(isFavorite) {
        this.proto.setIsFavorite(isFavorite)
    }

    setCreateTime(seconds) {
        const ts = new timestamp.Timestamp()
        ts.setSeconds(seconds)
        this.proto.setCreateTime(ts)
    }

    setName(name) {
        this.proto.setName(name)
    }

    getName() {
        return this.proto.getName()
    }

    getUrl() {
        return this.proto.getUrl()
    }

    getTitle() {
        return this.proto.getTitle()
    }

    getIconUrl() {
        return this.proto.getIconUrl()
    }

    getCreateTime() {
        return this.proto.getCreateTime()
    }

    getContent() {
        return this.proto.getContent()
    }

    getIsRead() {
        return this.proto.getIsRead()
    }

    getIsFavorite() {
        return this.proto.getIsFavorite()
    }

    getSiteName() {
        return this.proto.getSiteName()
    }
}

export class Articles {
    constructor(proto) {
        this.proto = proto

        this.articles = []
        this.proto.getArticlesList().forEach(a => {
            this.articles.push(new Article(a))
        })
    }

    getArticlesList() {
        return this.articles
    }

    getNextPageToken() {
        return this.proto.getNextPageToken()
    }
}

export class ArticlesClient {
    constructor(hostname) {
        this.client = new proto.ArticlesServicePromiseClient(hostname)
    }

    async get(name) {
        const request = new proto.GetArticleRequest()
        request.setName(name)
        request.setView(2)
        return new Article(await this.client.getArticle(request))
    }

    async update(article) {
        const request = new proto.UpdateArticleRequest()
        request.setArticle(article.proto)
        request.setUpdateMask(new updateMask.FieldMask().addPaths('is_read').addPaths('is_favorite'))
        return new Article(await this.client.updateArticle(request))
    }

    async delete(name) {
        const request = new proto.DeleteArticleRequest()
        request.setName(name)
        return new Article(await this.client.deleteArticle(request))
    }

    async list(pageSize, from, params) {
        const request = new proto.ListArticlesRequest()
        request.setPageSize(pageSize)
        if (from) request.setPageToken(from)
        if (params) {
            if (params.isRead !== undefined) request.setIsRead(new wrappers.BoolValue().setValue(params.isRead))
            if (params.isStarred !== undefined) request.setIsFavorite(new wrappers.BoolValue().setValue(params.isStarred))
        }
        return new Articles(await this.client.listArticles(request))
    }
}
