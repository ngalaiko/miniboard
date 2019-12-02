<script context='module'>
    import proto from './tokens_service_grpc_web_pb.js'

    export const Tokens = (apiUrl) => {
        let $ = {}

        const client = new proto.TokensServicePromiseClient(apiUrl)

        $.exchangeCode = async (code) => {
            const request = new proto.CreateTokenRequest()
                .setCode(code)
            const response = await client.createToken(request)
            return response.toObject()
        }

        return $
    }
</script>
