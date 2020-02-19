import { ApiClient } from './api'

export class User {
    name: string

    constructor(body: any) {
        this.name = body.name
    }
}

export class UsersClient {
    private apiClient: ApiClient

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient
    }

    async me(): Promise<User> {
        return new User(await this.apiClient.get('/api/v1/users/me'))
    }
}
