import Api from '/users/services/api.js'

class Users {
    async create(params) {
        if (params === undefined) params = {}

        const request = {
            username: params.username,
            password: params.password
        }

        return await Api.post('/v1/users', request)
    }
}

export default new Users()
