import LocalStorage from './localstorage/localstorage'

export default class Api {
    constructor() {
        this.localStorage = new LocalStorage()
    }

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

    send(url, method, body) {
        let options = {
            method: method,
            headers: {
                'Content-Type':  'application/json',
                'Authorization': this.authorization()
            },
        }
        if (body !== undefined ) {
            options.body = JSON.stringify(body)
        }
        return fetch(url, options)
            .then(response => response.json())
        // todo: handle errors
    }

    authenticate(auth) {
        this.localStorage.set('authentication.access_token', auth.access_token)
        this.localStorage.set('authentication.token_type', auth.token_type)
    }

    authorization() {
        let access_token = this.localStorage.get('authentication.access_token')
        let token_type = this.localStorage.get('authentication.token_type')
        if (access_token !== null && token_type !== null) {
            return `${token_type} ${access_token} `
        }
        return ''
    }

    authorized() {
        return this.authorization() != ''
    }

    logout() {
        this.localStorage.remove('authentication.access_token')
        this.localStorage.remove('authentication.token_type')
    }

    subject() {
        let token = this.localStorage.get('authentication.access_token')
        if (token == null) {
            return ''
        }
        let base64Url = token.split('.')[1]
        let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        var jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
        }).join(''));
        return JSON.parse(jsonPayload).sub
    }
}
