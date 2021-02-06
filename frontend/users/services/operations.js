import Api from '/users/services/api.js'

class Operations {
    async get(id) {
        return await Api.get(`/v1/operations/${id}`)
    }

    async wait(id) {
        const operation = await this.get(id)

        switch (true) {
        case !operation.done:
            await new Promise(r => setTimeout(r, 1000))
            return await this.wait(id)
        case operation.result.error !== undefined:
            throw operation.result.error.message
        case operation.result.response !== undefined:
            return operation.result.response
        default:
            throw `invalid state for operation ${id}`
        }
    }
}

export default new Operations()
