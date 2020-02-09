export class Source {
    constructor(body) {
		this.body = body
    }

    getName() {
		return this.body.name
    }

    getUrl() {
		return this.body.url
    }
}

export class SourcesClient {
    constructor(apiClient) {
		this.apiClient = apiClient
    }

    async createSource(username, url) {
        return new Source(await this.apiClient.post(`/api/v1/${username}/sources`, {
            url: url,
        }))
    }
}
