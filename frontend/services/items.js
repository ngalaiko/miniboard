import Api from '/services/api.js'

class Items {
    async get(id) {
        return await Api.get(`/v1/items/${id}/`)
    }

    async list(params) {
        if (params === undefined) params = {}

        const pageSizeQuery = params.pageSize !== undefined
            ? `&page_size=${params.pageSize}`
            : ''

        const createdLtQuery = params.createdLt !== undefined
            ? `&created_lt=${encodeURIComponent(params.createdLt)}`
            : ''

        let url = ''
        switch (true) {
        case !!params.subscriptionId && !!params.tagId:
            throw new Error('not implemented')
        case !!params.subscriptionId:
            url = `/v1/subscriptions/${params.subscriptionId}/items/?` + pageSizeQuery + createdLtQuery
            break
        case !!params.tagId:
            url = `/v1/tags/${params.tagId}/items/?` + pageSizeQuery + createdLtQuery
            break
        default:
            url = `/v1/items/?` + pageSizeQuery + createdLtQuery
            break
        }

        const body = await Api.get(url)
        return body.items
    }
}

export default new Items()
