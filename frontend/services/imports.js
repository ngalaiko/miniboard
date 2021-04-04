import Api from '/services/api.js'

class Imports {
    async create(raw) {
        return await Api.fetch('/v1/imports/', {
            method: 'POST',
            body: raw,
            headers: new Headers({
                "Content-Type": "application/xml",
            }),
        })
    }
}

export default new Imports()
