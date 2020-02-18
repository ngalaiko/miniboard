import { ApiClient } from './api'

export class Token {
    token: string

    constructor(body: any) {
		this.token = body.token
    }

    getToken(): string {
		return this.token
    }
}

export class TokensClient {
    private apiClient: ApiClient

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient
    }

    async exchangeCode(code: string): Promise<Token> {
        return new Token(await this.apiClient.post(`/api/v1/tokens`, {
            code: code,
        }))
    }
}
