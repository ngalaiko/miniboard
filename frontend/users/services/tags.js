import Api from '/users/services/api.js'

let tagsByTitle = undefined

class Tags {
    async listAll() {
        const listAllRec = async(pageSize, createdLt) => {
            const params = {}

            if (pageSize === undefined) pageSize = 100
            if (pageSize !== undefined) params.pageSize = pageSize
            if (createdLt !== undefined) params.createdLt = createdLt

            const tags = await this.list(params)

            // possible race condition, I don't care
            if (tagsByTitle === undefined) tagsByTitle = {}
            tags.forEach(tag => tagsByTitle[tag.title] = tag)

            if (tags.length < pageSize) {
                return tags
            }

            params.createdLt = tags[tags.length - 1].created

            return tags.concat(await TagsService.list(params))
        }
        return await listAllRec()
    }

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

    async getOrCreate(params) {
        if (params === undefined) params = {}

        if (tagsByTitle === undefined) await this.list()

        const found = tagsByTitle[params.title]
        if (found !== undefined) return found

        return this.create(params)
    }

    async create(params) {
        if (params === undefined) params = {}

        const request = {
            title: params.title,
        }

        return await Api.post('/v1/tags', request)
    }
}

export default new Tags()
