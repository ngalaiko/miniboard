<script context='module'>
    import proto from './codes_service_grpc_web_pb.js'

    export const Codes = () => {
        let $ = {}

        const client = new proto.CodesServicePromiseClient('http://localhost:8080')

        $.sendCode = async (email) => {
            const request = new proto.CreateCodeRequest()
                .setEmail(email)
            return await client.createCode(request)
        }

        return $
    }
</script>
