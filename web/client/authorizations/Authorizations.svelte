<script context='module'>
    export const Authorizations = (api) => {
        let $ = {}

        $.getAuthorization = async (username, password) => {
            return await api.post(`/api/v1/authorizations`, {
                username: username,
                password: password,
                grant_type: 'password',
            })
        }

        return $
    }
</script>
