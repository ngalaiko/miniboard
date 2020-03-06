import { ApiClient } from './api'

export class ListParams {
    isFavorite?: boolean
    isRead?: boolean
    title?: string

    withFavorite(isFavorite: boolean): ListParams {
        this.isFavorite = isFavorite
        return this
    }

    withRead(isRead: boolean): ListParams {
        this.isRead = isRead
        return this
    }

    withTitle(title?: string): ListParams {
        switch (true) {
            case !title:
                this.title = undefined
                break
            case title && title.length < 3:
                break
            default:
                this.title = title
                break
        }
        return this
    }
}

export class Article {
    url: string
    title: string
    isRead: boolean
    isFavorite: boolean
    name: string
    updateTime: number
    iconUrl: string
    content: string
    siteName: string

    constructor(body: any) {
        this.url = body.url
        this.title = body.title
        this.isRead = body.isRead
        this.isFavorite = body.isFavorite
        this.updateTime = body.updateTime
        this.name = body.name
        this.iconUrl = body.iconUrl
        this.content = body.content
        this.siteName = body.siteName
    }
}

export class Articles {
    articles: Article[] = []
    nextPageToken: string = ""

    constructor(body: any) {
        body.articles.forEach((a: any) => this.articles.push(new Article(a)))
        this.nextPageToken = body.nextPageToken
    }
}

export class ArticlesClient {
    private apiClient: ApiClient

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient
    }

    async get(name: string): Promise<Article> {
        return new Article(await this.apiClient.get(`/api/v1/${name}?view=ARTICLE_VIEW_FULL`))
    }

    async update(articleName: string, diff: any): Promise<Article> {
        return new Article(await this.apiClient.patch(`/api/v1/${articleName}`, diff))
    }

    async delete(name: string): Promise<Article> {
        return new Article(await this.apiClient.delete(`/api/v1/${name}`))
    }

    async list(username: string, pageSize: number, from: string, params: ListParams): Promise<Articles> {
        let query = `page_size=${pageSize}`
        if (from !== '') {
            query += `&page_token=${from}`
        }
        if (params.isFavorite !== undefined) {
            query += `&isFavorite=${params.isFavorite}`
        }
        if (params.isRead !== undefined) {
            query += `&isRead=${params.isRead}`
        }
        if (params.title !== undefined) {
            query += `&title=${params.title}`
        }
        return new Articles(await this.apiClient.get(`/api/v1/${username}/articles?${query}`))
    }
}
