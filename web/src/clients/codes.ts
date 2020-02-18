import { ApiClient } from './api'

export class Code {
    body: any

    constructor(body: any) {
        this.body = body
    }
}

export class CodesClient {
    private apiClient: ApiClient

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient
    }

    async sendCode(email: string): Promise<Code> {
        return new Code(await this.apiClient.post(`/api/v1/codes`, {
            email: email,
        }))
    }
}
