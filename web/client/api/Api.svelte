<script context='module'>
    export const Api = () => {
        let $ = {}

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

        $.logout = () => {
            send('/logout', 'GET')
        }

        let subject = ''
        $.subject = () => {
            if (subject !== '') {
                return subject
            }
            let parts = location.pathname.split('/')
            return `${parts[1]}/${parts[2]}`
        }

        const send = async (url, method, body) => {
            let u = new URL(location.href)
            let authCode = u.searchParams.get('authorization_code')
            if (authCode != null) {
                url += `&authorization_code=${authCode}`
                window.history.replaceState({}, document.Title, `${u.origin}${u.pathname}`);
            }

            let options = {
                method: method,
                headers: {
                    'Content-Type':  'application/json',
                },
            }
            if (body !== undefined ) {
                options.body = JSON.stringify(body)
            }
            let resp = await fetch(url, options)

            if (resp.status / 100 === 2) {
                return resp.json()
            }

            router.route('/')
        }

        return $
    }
</script>
