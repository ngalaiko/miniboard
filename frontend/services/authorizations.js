import Api from '/services/api.js'

class Authorizations {
    async create(params) {
        if (params === undefined) params = {}

        const request = {
            username: params.username,
            password: params.password
        }

        return await Api.post('/v1/authorizations', request)
    }
}

export default new Authorizations()
