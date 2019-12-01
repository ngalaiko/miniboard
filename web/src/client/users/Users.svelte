<script context='module'>
    import proto from './users_service_grpc_web_pb.js'

    export const Users = () => {
        let $ = {}

        const client = new proto.UsersServicePromiseClient('http://localhost:8080')

        $.get = async (name) => {
            const request = new proto.GetUserRequest()
                .setName(name)
            let response = await client.getUser(request)
            return response.toObject()
        }

        return $
    }
</script>
