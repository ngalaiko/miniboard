<script context='module'>
    import { Api } from './api/Api.svelte'
    import { Authorizations } from './authorizations/Authorizations.svelte'
    import { Articles } from './articles/Articles.svelte'
    import Labels from './labels/labels'
    import { Users } from './users/Users.svelte'

    export const Client = () => {
        let $ = {}

        $.api = new Api()
        $.authorizations = new Authorizations($.api)
        $.articles = new Articles($.api)
        $.labels = new Labels($.api)
        $.users = new Users($.api)

        return $
    }
</script>
