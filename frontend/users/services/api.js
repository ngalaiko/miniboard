const apiUrl = window.location.hostname !== 'localhost'
    ? 'https://api.miniboard.app'
    : 'http://localhost:80';

class Api {
    async get(url) {
        return await _fetch(url)
    }

    async post(url, request) {
        return await _fetch(url, {
            method: 'POST',
            body: JSON.stringify(request),
        })
    }
}

const _fetch = async(url, params) => {
    if (params == undefined) params = {}
    params.credentials = 'include'

    const response = await fetch(apiUrl + url, params)
    const body = await response.json()
    if (response.status !== 200) {
        throw new body.message
    }

    return body
}

export default new Api()
