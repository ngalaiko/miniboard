export class User {
    constructor(body) {
        this.body = body
    }

    getName() {
        return this.body.name
    }
}

export class UsersClient {
    constructor(apiClient) {
        this.apiClient = apiClient
    }

    async me() {
        return new User(await this.apiClient.get('/api/v1/users/me'))
    }
}
