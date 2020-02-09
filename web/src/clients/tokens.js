export class Token {
    constructor(body) {
		this.token = body.token
    }

    getToken() {
		return this.token
    }
}

export class TokensClient {
    constructor(apiClient) {
        this.apiClient = apiClient
    }

    async exchangeCode(code) {
        return new Token(await this.apiClient.post(`/api/v1/tokens`, {
            code: code,
        }))
    }
}
