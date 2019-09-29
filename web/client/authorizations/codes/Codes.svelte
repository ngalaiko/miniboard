<script context='module'>
    export const Codes = (api) => {
        let $ = {}

        $.sendCode = async (email) => {
            return await api.post(`/api/v1/authorizations/codes`, {
                email: email,
            })
        }

        return $
    }
</script>
