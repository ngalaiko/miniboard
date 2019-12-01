<script context='module'>
    import { Api } from './api/Api.svelte'
    import { Articles } from './articles/Articles.svelte'
    import { Codes } from './codes/Codes.svelte'
    import { Users } from './users/Users.svelte'

    export const Client = async (router) => {
        let $ = {}

        $.api = Api(router)
        $.articles = await Articles($.api)
        $.codes = Codes()
        $.users = Users($.api)

        return $
    }
</script>
