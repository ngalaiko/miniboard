import Api from '/users/services/api.js'

class Tags {
    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const url = '/v1/tags?' + pageSizeQuery + createdLtQuery

        const body = await Api.get(url)

        return body.tags
    }
}

export default new Tags()
