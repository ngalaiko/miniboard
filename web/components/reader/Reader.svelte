<script>
    export let articles

    export let name
    export let title


    // https://stackoverflow.com/questions/30106476/using-javascripts-atob-to-decode-base64-doesnt-properly-decode-utf-8-strings
    function b64DecodeUnicode(article) {
        document.title = `Miniboard - ${article.title}`
        return decodeURIComponent(atob(article.content).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
    }
</script>

<div class='reader'>
    {#await articles.get(name)}
        loading...
    {:then article}
        {#if article.content === undefined}
            no saved content, redirecting to <a href={article.url} target='_blank'>{article.url}</a>
            <div hidden>
                {window.open(article.url, '_blank')}
            </div>
        {:else}
            {@html b64DecodeUnicode(article)}
        {/if}
    {:catch e}
        failed to fetch, are you online?
    {/await}
</div>

<style>
    .reader{
        overflow: auto;
    }
</style>
