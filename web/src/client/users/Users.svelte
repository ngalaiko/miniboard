<script context='module'>
    import proto from './users_service_grpc_web_pb.js'

        const client = new proto.UsersServiceClient('http://localhost:8080')
        const request = new proto.GetUserRequest()
        request.setName('name')

        const call = client.getUser(request, {
            'grpc-encoding': 'gzip',
        }, (err, response) => {
            console.log(response);
            console.log(err);
        })


    export const Users = (api) => {
        let $ = {}

        $.get = async (name) => {
            return await api.get(`/api/v1/${name}`)
        }

        return $
    }
</script>
