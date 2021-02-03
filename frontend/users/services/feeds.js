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

        let tagIdsQuery = ''
        if (params.tagIds !== undefined) params.tagIds.forEach(tagId => {
            tagIdsQuery += `&tag_ids=${tagId}`
        })

        const url = '/v1/feeds?' + pageSizeQuery + createdLtQuery + tagIdsQuery

        const body = await Api.get(url)

        return body.feeds
    }
}

export default new Feeds()
