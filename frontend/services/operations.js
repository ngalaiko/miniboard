import Api from '/services/api.js'

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
            throw new Error(operation.result.error.message)
        case operation.result.response !== undefined:
            return operation.result.response
        default:
            throw new Error(`invalid state for operation ${id}`)
        }
    }
}

export default new Operations()
