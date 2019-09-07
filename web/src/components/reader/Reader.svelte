<script>
    export let articles

    export let name
</script>

<div>
    {#await articles.get(name)}
        loading...
    {:then article}
        {@html atob(article.content)}
    {/await}
</div>
