<script context='module'>
    import proto from './codes_service_grpc_web_pb.js'

    export const Codes = () => {
        let $ = {}

        const client = new proto.CodesServiceClient('http://localhost:8080')
        $.sendCode = async (email) => {
            const request = new proto.CreateCodeRequest()
            request.setEmail(email)

            let error
            await client.createCode(request, {}, (err, response) => {
                error = err
            })
            if (error) throw error
            return
        }

        return $
    }
</script>
