<script>
    export let articles

    export let name
</script>

<div class='reader'>
    {#await articles.get(name)}
        loading...
    {:then article}
        {#if article.content === undefined}
            no saved content, redirecting to <a href={article.url}>{article.url}</a>
            <div hidden>
                {location.replace(article.url)}
            </div>
        {:else}
            {@html atob(article.content)}
        {/if}
    {/await}
</div>

<style>
    .reader{
        overflow: auto;
    }
</style>
