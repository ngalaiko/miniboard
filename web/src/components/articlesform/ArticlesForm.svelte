<script>
    import { ArticlesService } from './articles-service.js'

    export let api
    export let user

    let articlesService = new ArticlesService(api, user)

    let url = ""
    let error = ""

    function handleClick() {
        articlesService.add(url)
            .catch(err => { error = err })
    }
</script>

<div>
    <input type="text" bind:value={url} placeholder="add" required="" />
    {#if (error != "")}
    <div class="alert">{error}</div>
    {/if}
    <button on:click|preventDefault={handleClick}>Add</button>
</div>
