<script context='module'>
    import proto from './users_service_grpc_web_pb.js'

    export const Users = (apiUrl) => {
        let $ = {}

        const client = new proto.UsersServicePromiseClient(apiUrl)

        $.me = async () => {
            const request = new proto.GetMeRequest()
            return await client.getMe(request)
        }

        return $
    }
</script>
