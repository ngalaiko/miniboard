<script context='module'>
    import proto from './users_service_grpc_web_pb.js'


    export const Users = (api) => {
        let $ = {}

        $.get = async (name) => {
            return await api.get(`/api/v1/${name}`)
        }

        return $
    }
</script>
