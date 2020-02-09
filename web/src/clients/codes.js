export class Code {
    constructor(json) {
    }
}

export class CodesClient {
    constructor(apiClient) {
        this.apiClient = apiClient
    }

    async sendCode(email) {
        return new Code(await this.apiClient.post(`/api/v1/codes`, {
            email: email,
        }))
    }
}
