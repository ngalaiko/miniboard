<script context='module'>
    import proto from './users_service_grpc_web_pb.js'

    export const Users = (apiUrl) => {
        let $ = {}

        const client = new proto.UsersServicePromiseClient(apiUrl)

        $.get = async (name) => {
            const request = new proto.GetUserRequest()
                .setName(name)
            let response = await client.getUser(request)
            return response.toObject()
        }

        $.me = async () => {
            const request = new proto.GetMeRequest()
            let response = await client.getMe()
            return response.toObject()
        }

        return $
    }
</script>
