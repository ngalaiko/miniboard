<script context='module'>
    import proto from './codes_service_grpc_web_pb.js'

    export const Codes = (api) => {
        let $ = {}

        const client = new proto.CodesServiceClient('http://localhost:8080')
        const request = new proto.CreateCodeRequest()
        request.setEmail('asdasd')

        const call = client.createCode(request, {
            'grpc-encoding': 'gzip',
        }, (err, response) => {
            console.log(response);
        })

        $.sendCode = async (email) => {
            return await api.post(`/api/v1/codes`, {
                email: email,
            })
        }

        return $
    }
</script>
