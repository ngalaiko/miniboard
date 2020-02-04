<script>
    import { navigate } from "svelte-routing"

    export let code
    export let tokensClient

    const parseJwt = (token) => {
        var base64Url = token.split('.')[1];
        var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        return JSON.parse(jsonPayload);
    }

    tokensClient.exchangeCode(code).then((response) => {
        return response.getToken()
    }).then((token) => {
        return parseJwt(token).sub
    }).then((subject) => {
        navigate(`/${subject}/unread`)
    }).catch(e => {
        navigate(`/?error=${e.message}`)
    })
</script>
