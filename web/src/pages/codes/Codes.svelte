<script lang="ts">
  import { navigate } from "svelte-routing"
  // @ts-ignore
  import { TokensClient, Token } from '../../clients/tokens.ts'

  export let code: string = ''

  export let tokensClient: TokensClient

  const parseJwt = (token: string) => {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
  }

  tokensClient.exchangeCode(code).then((response: Token) => {
    return response.getToken()
  }).then((token: string) => {
    return parseJwt(token).sub
  }).then((subject: string) => {
    navigate(`/${subject}/unread`)
  }).catch((e: Error) => {
    navigate(`/?error=${e.message}`)
  })
</script>
