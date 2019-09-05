import LocalStorage from './localstorage/localstorage'

export default function api() {
    let $ = {}

    let localStorage = new LocalStorage()

    $.get = (url) => {
        return send(url, 'GET')
    }

    $.post = (url, data) => {
        return send(url, 'POST', data)
    }

    $.patch = (url, data) => {
        return send(url, 'PATCH', data)
    }

    $.delete = (url) => {
        return send(url, 'DELETE')
    }

    $.authenticate = (auth) => {
        localStorage.set('authentication.access_token', auth.access_token)
        localStorage.set('authentication.token_type', auth.token_type)
    }

    $.authorized = () => {
        return authorization() != ''
    }

    $.logout = () => {
        localStorage.remove('authentication.access_token')
        localStorage.remove('authentication.token_type')
    }

    $.subject = () => {
        let token = localStorage.get('authentication.access_token')
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

    let send = async (url, method, body) => {
        let options = {
            method: method,
            headers: {
                'Content-Type':  'application/json',
                'Authorization': authorization()
            },
        }
        if (body !== undefined ) {
            options.body = JSON.stringify(body)
        }
        let resp = await fetch(url, options)

        // todo: handle errors
        return resp.json()
    }

    let authorization = () => {
        let access_token = localStorage.get('authentication.access_token')
        let token_type = localStorage.get('authentication.token_type')
        if (access_token !== null && token_type !== null) {
            return `${token_type} ${access_token} `
        }
        return ''
    }

    return $
}
