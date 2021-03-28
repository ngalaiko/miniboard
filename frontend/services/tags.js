import Api from '/services/api.js'

let tagsByTitle = undefined

class Tags {
    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const url = '/v1/tags/?' + pageSizeQuery + createdLtQuery

        const body = await Api.get(url)

        return body.tags
    }

    async create(params) {
        if (params === undefined) params = {}

        const request = {
            title: params.title,
        }

        return await Api.post('/v1/tags/', request)
    }
}

export default new Tags()
