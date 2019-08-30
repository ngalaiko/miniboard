<script>
    import ArticlesForm from '../../components/articlesform/ArticlesForm.svelte'

    export let api
    export let user
</script>

<div>
    <ArticlesForm api={api} user={user} />
</div>
