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

            if (resp.status != 401) {
                return resp.json()
            }

            router.route('/')
        }

        return $
    }
</script>
