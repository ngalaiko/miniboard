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
    function onAdd() {
        error = ""

        let rnd = Math.random()
        articlesList = [{
          'url': url,
          'title': 'Loading...',
          'random': rnd
        }].concat(articlesList)

        articlesService.add(url)
            .catch(err => { error = err })
            .then(article => {
              articlesList = [article].concat(articlesList.filter(article => article.random != rnd ))
            })
            .then(() => { url = "" })
    }

    function onDeleted(name) {
        articlesList = articlesList.filter(article => article.name != name)
    }
</script>

<div>
    <div>
        <input type="text" bind:value={url} placeholder="add" required="" />
        {#if (error != "")}
            <div class="alert">{error}</div>
        {/if}
        <button on:click|preventDefault={onAdd}>+</button>
    </div>
    <div>
        {#each articlesList as article}
            <Article on:deleted={(e) => onDeleted(e.detail)} api={api} {...article} } />
        {/each}
    </div>
</div>
