<script context='module'>
    import proto from './codes_service_grpc_web_pb.js'

    export const Codes = (apiUrl) => {
        let $ = {}

        const client = new proto.CodesServicePromiseClient(apiUrl)

        $.sendCode = async (email) => {
            const request = new proto.CreateCodeRequest()
                .setEmail(email)
            return await client.createCode(request)
        }

        return $
    }
</script>
