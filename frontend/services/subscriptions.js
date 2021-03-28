import Api from '/services/api.js'

class Subscriptions {
    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const url = '/v1/subscriptions/?' + pageSizeQuery + createdLtQuery

        const body = await Api.get(url)

        return body.subscriptions
    }

    async create(params) {
        if (params === undefined) params = {}

        const request = {
            url: params.url,
        }

        if (params.tagIds !== undefined) {
            request.tag_ids = params.tagIds
        }

        return await Api.post('/v1/subscriptions/', request)
    }
}

export default new Subscriptions()
