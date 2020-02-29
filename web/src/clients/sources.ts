import { ApiClient } from './api'

export class Source {
    name: string
    url: string

    constructor(body: any) {
		this.name = body.name
		this.url = body.url
    }
}

export class SourcesClient {
    private apiClient: ApiClient

    constructor(apiClient: ApiClient) {
		this.apiClient = apiClient
    }

    async create(username: string, url: string): Promise<Source> {
        return new Source(await this.apiClient.post(`/api/v1/${username}/sources`, {
            url: url,
        }))
    }
}
