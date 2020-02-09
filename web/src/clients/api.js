import { navigate } from "svelte-routing"

export class ApiClient {
    get(url) {
        return this.send(url, 'GET')
    }

    post(url, data) {
        return this.send(url, 'POST', data)		
    }

    patch(url, data) {
        return this.send(url, 'PATCH', data)
    }

    delete(url) {
        return this.send(url, 'DELETE')
    }

    async send(url, method, body) {
        const options = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
        }
        
        if (body !== undefined) {
            options.body = JSON.stringify(body)
        }

        const resp = await fetch(url, options)
        
        if (resp.status / 100 === 2) {
            return resp.json()
        }

        throw `${method} ${url} - status code: ${resp.status}`		
    }
}
