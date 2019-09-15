<script context='module'>
    import { Api } from './api/Api.svelte'
    import { Authorizations } from './authorizations/Authorizations.svelte'
    import { Articles } from './articles/Articles.svelte'
    import Labels from './labels/labels'
    import { Users } from './users/Users.svelte'

    export const Client = async () => {
        let $ = {}

        $.api = Api()
        $.authorizations = Authorizations($.api)
        $.articles = await Articles($.api)
        $.labels = Labels($.api)
        $.users = Users($.api)

        return $
    }
</script>
