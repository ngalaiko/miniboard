<script>
    export let articles

    export let name
</script>

<div class='reader'>
    {#await articles.get(name)}
        loading...
    {:then article}
        {@html atob(article.content)}
    {/await}
</div>

<style>
    .reader{
        overflow: auto;
    }
</style>
