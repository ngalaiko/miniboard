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

    let url = ''
    let error = ''
    function onAdd() {
        error = ''

        let rnd = Math.random()
        articlesList = [{
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'random': rnd
        }].concat(articlesList)

        articlesService.add(url)
            .catch(err => { error = err })
            .then(article => {
              articlesList = [article].concat(articlesList.filter(article => article.random != rnd ))
            })

        url = ''
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
    <div class='list'>
        {#each articlesList as article}
            <Article on:deleted={(e) => onDeleted(e.detail)} api={api} {...article} } />
        {/each}
    </div>
</div>

<style>
    .list {
      display: flex;
      flex-direction: column;
      justify-content: center;
    }

    input {
        border: 1px solid #ccc;
        padding: 3px;
        line-height: 20px;
        width: 250px;
        font-size: 99%;
        margin-bottom: 10px;
        margin-top: 5px;
        -webkit-appearance: none
    }

    input:focus{
        outline-width: 0
    }
</style>
