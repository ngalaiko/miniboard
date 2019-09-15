<script context='module'>
    import { LocalStorage } from './localstorage/LocalStorage.svelte'

    export const Api = () => {
        let $ = {}

        let localStorage = LocalStorage()

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
            localStorage.set('authentication.refresh_token', auth.refresh_token)
            localStorage.set('authentication.token_type', auth.token_type)
        }

        $.authorized = () => {
            return authorization() != ''
        }

        $.logout = () => {
            localStorage.remove('authentication.access_token')
            localStorage.remove('authentication.refresh_token')
            localStorage.remove('authentication.token_type')
        }

        $.subject = () => {
            let token = localStorage.get('authentication.access_token')
            if (token == null) {
                return ''
            }
            let base64Url = token.split('.')[1]
            let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
            let payload = JSON.parse(decodeURIComponent(atob(base64).split('').map(c => {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
            }).join('')));

            return payload.sub
        }

        const send = async (url, method, body) => {
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

            if (resp.status != 401) {
                return resp.json()
            }

            // if not expired, return an error
            if (new Date($.subject().exp * 1000) > new Date()) {
                return resp.json()
            }

            let auth = await $.post(`/api/v1/authorizations`, {
                refresh_token: localStorage.get('authentication.refresh_token'),
                grant_type: 'refresh_token',
            })

            $.authenticate(auth)

            return send(url, method, body)
        }

        const authorization = () => {
            let access_token = localStorage.get('authentication.access_token')
            let token_type = localStorage.get('authentication.token_type')
            if (access_token !== null && token_type !== null) {
                return `${token_type} ${access_token} `
            }
            return ''
        }

        return $
    }
</script>
