<script>
    import Article from '../article/Article.svelte'
    import { ArticlesService } from './articles-service.js'
    import ArticlesForm from '../../components/articlesform/ArticlesForm.svelte'

    export let api
    export let user

    let articlesService = new ArticlesService(api, user)

    let articlesList = []
    articlesService.next()
        .then(list => { articlesList = articlesList.concat(list) })

    let url = ""
    let error = ""
    function handleAdd() {
        error = ""
        articlesService.add(url)
            .catch(err => { error = err })
            .then(article => { articlesList = [article].concat(articlesList) } )
            .then(() => { url = "" })
    }
</script>

<div>
    <div>
        <input type="text" bind:value={url} placeholder="add" required="" />
        {#if (error != "")}
            <div class="alert">{error}</div>
        {/if}
        <button on:click|preventDefault={handleAdd}>Add</button>
    </div>
    <div>
        {#each articlesList as article}
            <Article {...article} />
        {/each}
    </div>
</div>
