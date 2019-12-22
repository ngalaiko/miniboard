<script>
    export let articles

    export let name

    const decoder = new TextDecoder()

    const decode = (article) => {
        document.title = `${article.getTitle()} - Miniboard`
        return decoder.decode(article.getContent())
    }
</script>

<div id='article'>
    {#await articles.get(name)}
        loading...
    {:then article}
        {#if article.getContent() === ""}
            no saved content, redirecting to <a href={article.url} target='_blank'>{article.getUrl()}</a>
            <div hidden>
                {window.open(article.getUrl(), '_blank')}
            </div>
        {:else}
            <div class='page'>
                <h1>{article.getTitle()}</h1>
                {@html decode(article)}
            </div>
        {/if}
    {:catch e}
        failed to fetch, are you online?
    {/await}
</div>

<style>
    :global(pre) {
        overflow: auto;
    }

    :global(img) {
        max-width: 100%;
        height: auto;
    }

    :global(blockquote) {
        font-style: italic;
        border-left: 3px solid #ccc;
        margin-left: 2px;
        margin-right: 6px;
        padding-left: 16px;
    }

    @media screen {
        body {
            margin: 0;
            padding: 0;
            -webkit-user-select: none;
            overflow-x: hidden;
            -webkit-text-size-adjust: none;
        }

        #article {
            pointer-events: auto;
            -webkit-user-select: auto;
            overflow: visible;
        }

        #article:focus {
            outline: none;
        }

        .page {
            margin-left: auto;
            margin-right: auto;
            padding-top: 35px;
            padding-bottom: 35px;
            position: relative;
        }

        .page video {
            height: auto;
            position: relative;
        }
    }

    @media screen and (max-width: 569px) {
        /* iPhone 5 in landscape (568px) and smaller, including all iPhones in portrait */
        h1.title {
            font-size: 1.5558em;
        }
        h1 {
            font-size: 1.4em;
        }
    }

    @media screen and (min-width: 704px) {
        .page { padding-left: 42px; padding-right: 42px; }
    }

    @media only screen and (min-width: 780px) {
        #article {
            max-width: 800px;
            margin: 0 auto;
        }

        .page {
            padding-left: 0px;
            padding-right: 0px;
            margin-left: 70px;
            margin-right: 70px;
        }
    }

    #article {
        -webkit-font-smoothing: subpixel-antialiased;
        padding: 5px;
    }
</style>
