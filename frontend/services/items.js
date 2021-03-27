import Api from '/services/api.js'

class Items {
    async get(id) {
        return await Api.get(`/v1/items/${id}`)
    }

    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        const subscriptionIdQuery = params.subscriptionIdEq !== undefined
            ? `&subscription_id_eq=${encodeURIComponent(params.subscriptionIdEq)}`
            : ''

        const tagIdQuery = params.tagIdEq !== undefined
            ? `&tag_id_eq=${encodeURIComponent(params.tagIdEq)}`
            : ''

        const url = '/v1/items?' + pageSizeQuery + createdLtQuery + subscriptionIdQuery + tagIdQuery

        const body = await Api.get(url)

        return body.items
    }
}

export default new Items()
