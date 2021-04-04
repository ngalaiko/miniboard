const apiUrl = window.location.hostname !== 'localhost'
    ? 'https://api.miniboard.app'
    : 'http://localhost:80';

class Api {
    async get(url) {
        return await this.fetch(url)
    }

    async post(url, request) {
        return await this.fetch(url, {
            method: 'POST',
            body: JSON.stringify(request),
            headers: new Headers({
                "Content-Type": "application/json",
            }),
        })
    }

    async fetch(url, params) {
        if (params == undefined) params = {}
        params.credentials = 'include'

        const response = await fetch(apiUrl + url, params)
        const body = await response.json()
        if (response.status !== 200) {
            throw new Error(body.message)
        }

        return body
    }
}

export default new Api()
