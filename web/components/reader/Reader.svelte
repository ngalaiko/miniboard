<script>
    export let articles

    export let name
    export let title


    // https://stackoverflow.com/questions/30106476/using-javascripts-atob-to-decode-base64-doesnt-properly-decode-utf-8-strings
    const b64DecodeUnicode = (article) => {
        document.title = `Miniboard - ${article.title}`
        return decodeURIComponent(atob(article.content).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
    }
</script>

<div class='article'>
    {#await articles.get(name)}
        loading...
    {:then article}
        {#if article.content === undefined}
            no saved content, redirecting to <a href={article.url} target='_blank'>{article.url}</a>
            <div hidden>
                {window.open(article.url, '_blank')}
            </div>
        {:else}
            <div class='page'>
                <h1>{article.title}</h1>
                {@html b64DecodeUnicode(article)}
            </div>
        {/if}
    {:catch e}
        failed to fetch, are you online?
    {/await}
</div>

<style>
    .article {
        min-height: 100%;
        margin: 22px auto 0 auto;
        overflow: visible;
        font-size: 18px;
        line-height: 29px;
    }

    .page {
        text-align: start;
        word-wrap: break-word;
        padding-top: 53px;
        margin-top: 0;
        border-top: none;
        padding-left: 0px;
        padding-right: 0px;
        margin-left: 70px;
        margin-right: 70px;
        padding-bottom: 45px;
        max-width: 83.2ex;
    }

    h1 {
        font-weight: bold;
        font-size: 1.95552em;
        line-height: 1.2141em;
        margin-top: 0;
        margin-bottom: 0.5em;
    }

    :global(pre) {
        overflow: auto;
    }

    :global(img) {
        width: 100%;
        height: auto;
    }

    :global(blockquote) {
        font-style: italic;
        border-left: 3px solid #ccc;
        margin-left: 2px;
        margin-right: 6px;
        padding-left: 16px;
    }
</style>
