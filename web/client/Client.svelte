<script context='module'>
    import { Api } from './api/Api.svelte'
    import { Authorizations } from './authorizations/Authorizations.svelte'
    import { Articles } from './articles/Articles.svelte'
    import { Codes } from './authorizations/codes/Codes.svelte'
    import { Users } from './users/Users.svelte'

    export const Client = async () => {
        let $ = {}

        $.api = Api()
        $.authorizations = Authorizations($.api)
        $.articles = await Articles($.api)
        $.codes = Codes($.api)
        $.users = Users($.api)

        return $
    }
</script>
