import { ApiClient } from './api'

export class Source {
    name: string

    constructor(body: any) {
		this.name = body.name
    }
}

export class SourcesClient {
    private apiClient: ApiClient

    constructor(apiClient: ApiClient) {
		this.apiClient = apiClient
    }

    async addLink(username: string, url: string): Promise<Source> {
        return new Source(await this.apiClient.post(`/api/v1/${username}/sources`, {
            url: url,
        }))
    }

    async addFile(username: string, file: File) {
        return new Source(await this.apiClient.post(`/api/v1/${username}/sources`, {
            raw: file,
        }))
    }
}
