<script context='module'>
    import proto from './users_service_grpc_web_pb.js'

    export const Users = (api) => {
        let $ = {}

        const client = new proto.UsersServiceClient('http://localhost:8080')

        $.get = async (name) => {
            const request = new proto.GetUserRequest()
            request.setName(name)

            await client.getUser(request, {}, (err, response) => {
                console.log(response);
                console.log(err);
            })

            return await api.get(`/api/v1/${name}`)
        }

        return $
    }
</script>
