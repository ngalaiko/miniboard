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
    <form>
        <input type="text" bind:value={url} placeholder="add" required="" />
        {#if (error != "")}
            <div class="alert">{error}</div>
        {/if}
        <button on:click|preventDefault={onAdd} />
    </form>
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
        border: 1px solid;
        width: 100%;
        font-size: 1.1em;
        padding: 5px;
        padding-left: 7px;
    }

    input:focus{
        outline-width: 0
    }

    form {
        display: flex;
        flex-direction: row;
        margin: 0px;
        margin-bottom: 20px;
    }

    button {
        width: 0;
        height: 0;
        padding-left: 0; padding-right: 0;
        border-left-width: 0; border-right-width: 0;
        white-space: nowrap;
        overflow: hidden;
    }
</style>
