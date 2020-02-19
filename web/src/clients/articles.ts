import { ApiClient } from './api'

export class Article {
    url: string
    title: string
    isRead: boolean
    isFavorite: boolean
    name: string
    createTime: number
    iconUrl: string
    content: string
    siteName: string

    constructor(body: any) {
        this.url = body.url 
        this.title = body.title 
        this.isRead = body.isRead 
        this.isFavorite = body.isFavorite 
        this.createTime = body.createTime 
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

    async update(article: Article): Promise<Article> {
        return new Article(await this.apiClient.patch(`/api/v1/${article.name}`, article))
    }

    async delete(name: string): Promise<Article> {
        return new Article(await this.apiClient.delete(`/api/v1/${name}`))
    }

    async list(username: string, pageSize: number, from: string, params?: any): Promise<Articles> {
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