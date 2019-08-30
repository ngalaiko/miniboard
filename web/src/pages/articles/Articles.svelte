<script>
    import Article from '../article/Article.svelte'
    import { ArticlesService } from './articles-service.js'
    import ArticlesForm from '../../components/articlesform/ArticlesForm.svelte'

    export let api
    export let user

    let articlesService = new ArticlesService(api, user)
    let articles = articlesService.next()
</script>

<div>
    <ArticlesForm api={api} user={user} />
    {#await articles}
        loading articles...
    {:then articles}
        {#each articles as article}
            <Article {...article} />
        {/each}
    {/await}
</div>
