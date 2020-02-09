export class Article {
    constructor(body) {
        if (body === undefined) {
            this.body = {}
        } else {
            this.body = body
        }
    }

    setUrl(url) {
        this.body.url = url
    }

    setTitle(title) {
        this.body.title = title
    }

    setIsRead(isRead) {
        this.body.isRead = isRead
    }

    setIsFavorite(isFavorite) {
        this.body.isFavorite = isFavorite
    }

    setCreateTime(seconds) {
        this.body.createTime = seconds
    }

    setName(name) {
        this.body.name = name
    }

    getName() {
        return this.body.name
    }

    getUrl() {
        return this.body.url
    }

    getTitle() {
        return this.body.title
    }

    getIconUrl() {
        return this.body.iconUrl
    }

    getCreateTime() {
        return this.body.createTime
    }

    getContent() {
        return this.body.content
    }

    getIsRead() {
        return this.body.isRead
    }

    getIsFavorite() {
        return this.body.isFavorite
    }

    getSiteName() {
        return this.body.siteName
    }
}

export class Articles {
    constructor(body) {
        this.body = body
        this.articles = []
        body.articles.forEach(a => this.articles.push(new Article(a)))
    }

    getArticlesList() {
        return this.articles
    }

    getNextPageToken() {
        return this.nextPageToken
    }
}

export class ArticlesClient {
    constructor(apiClient) {
        this.apiClient = apiClient
    }

    async get(name) {
        return new Article(await this.apiClient.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`))
    }

    async update(article) {
        return new Article(await this.apiClient.patch(`/api/v1/${article.getName()}`, article.body))
    }

    async delete(name) {
        return new Article(await this.apiClient.delete(`/api/v1/${name}`))
    }

    async list(username, pageSize, from, params) {
        let query = `page_size=${pageSize}`
        if (from !== '') {
            query += `&page_token=${from}`
        }
        if (params) Object.keys(params).forEach(key => {
            query += `&${key}=${params[key]}`
        })
        return new Articles(await this.apiClient.get(`/api/v1/${username}/articles?${query}`))
    }
}
