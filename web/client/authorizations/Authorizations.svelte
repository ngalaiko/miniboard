<script context='module'>
    export const Authorizations = (api) => {
        let $ = {}

        $.getAuthorization = async (authorization_code) => {
            return await api.post(`/api/v1/authorizations`, {
                authorization_code: authorization_code,
                grant_type: 'authorization_code',
            })
        }

        return $
    }
</script>
