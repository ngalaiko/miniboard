<script context='module'>
    export const Users = (api) => {
        let $ = {}

        $.get = async (name) => {
            return await api.get(`/api/v1/${name}`)
        }

        return $
    }
</script>
