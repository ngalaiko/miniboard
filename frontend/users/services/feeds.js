import Api from '/users/services/api.js'

class Feeds {
    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const url = '/v1/feeds?' + pageSizeQuery + createdLtQuery

        const body = await Api.get(url)

        return body.feeds
    }

    async create(params) {
        if (params === undefined) params = {}

        const request = {
            url: params.url,
        }

        if (params.tagIds !== undefined) {
            request.tag_ids = params.tagIds
        }

        return await Api.post('/v1/feeds', request)
    }
}

export default new Feeds()
