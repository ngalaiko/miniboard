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

<div id='article'>
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

    /* modified Safari reader styles below */

    @media screen {
        body {
            margin: 0;
            padding: 0;
            -webkit-user-select: none;
            overflow-x: hidden;
            -webkit-text-size-adjust: none;
        }
        body.mac {
            background-color: transparent;
        }

        #article {
            pointer-events: auto;
            -webkit-user-select: auto;
            overflow: visible;
        }

        #article:focus {
            outline: none;
        }

        .page-number {
            display: block;
        }

        #article :nth-child(1 of .page):nth-last-child(1 of .page) .page-number {
            display: none;
        }

        .page {
            margin-left: auto;
            margin-right: auto;
            padding-top: 35px;
            padding-bottom: 35px;
            position: relative;
        }
        body.watch .page {
            padding: 0 4px;
        }

        .page:last-of-type {
            padding-bottom: 45px;
        }

        .page video {
            height: auto;
            position: relative;
        }

        .page div.scrollable {
            -webkit-overflow-scrolling: touch;
        }
    }

    @media screen and (-webkit-min-device-pixel-ratio:2) {
        hr {
            height: 0.5px;
        }
    }

    #article .extendsBeyondTextColumn {
        max-width: none;
    }

    .iframe-wrapper {
        background-color: black;
        max-width: none;
        text-align: center;
    }

    iframe {
        border: 0;
    }

    @media screen and (min-width: 161px) {
        /* Apple Watch 40mm */
        body.watch .page {
            padding: 0 8.5px;
        }
    }

    @media screen and (min-width: 183px) {
        /* Apple Watch 44mm */
        body.watch .page {
            padding: 0 9.5px;
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
        h2 {
            font-size: 1.2777em;
        }
        h3 {
            font-size: 1.15em;
        }
        .subhead {
            font-size: 1.2222em;
        }
        .metadata {
            font-size: 0.9em;
            line-height: 1.6em;
        }
        .title + .metadata {
            margin-top: -0.65em;
        }
    }

    @media screen and (min-width: 704px) {
        /* iPad in landscape with the sidebar open */
        .page { padding-left: 42px; padding-right: 42px; }
    }

    @media only screen and (min-width: 780px) {
        #article {
            max-width: 800px;
            margin: 0 auto;
        }

        /* Readable margins. */
        body.system #article { max-width: 83.2ex; }
        body.athelas #article { max-width: 104ex; }
        body.charter #article { max-width: 86ex; }
        body.georgia #article { max-width: 94ex; }
        body.iowanoldstyle #article { max-width: 90ex; }
        body.palatino #article { max-width: 97ex; }
        body.seravek #article { max-width: 87ex; }
        body.timesnewroman #article { max-width: 97ex; }
        body.applesystemuiserif #article { max-width: 93ex; }

        :matches(body.pingfangsc, body.pingfangtc) #article { max-width: 87.6ex; }
        :matches(body.heitisc, body.heititc) #article { max-width: 74.8ex; }
        :matches(body.songtisc, body.songtitc) #article { max-width: 102ex; }
        :matches(body.kaitisc, body.kaititc) #article { max-width: 102ex; }
        :matches(body.yuantisc, body.yuantitc) #article { max-width: 86.2ex; }
        :matches(body.libiansc, body.libiantc) #article { max-width: 95ex; }
        :matches(body.weibeisc, body.weibeitc) #article { max-width: 99ex; }
        :matches(body.yuppysc, body.yuppytc) #article { max-width: 87.6ex; }

        body.hiraginosansw3 #article { max-width: 75.7ex; }
        body.hiraginokakugothicpron #article { max-width: 76.4ex; }
        body.hiraginominchopron #article { max-width: 77.5ex; }
        body.hiraginomarugothicpron #article { max-width: 75.1ex; }

        body.applesdgothicneo #article { max-width: 82ex; }
        body.nanumgothic #article { max-width: 88.6ex; }
        body.nanummyeongjo #article { max-width: 94.1ex; }

        .page {
            /* We don't want the lines seperating pages to extend beyond the primary text column. */
            padding-left: 0px;
            padding-right: 0px;
            margin-left: 70px;
            margin-right: 70px;
        }
    }

    #article {
        -webkit-font-smoothing: subpixel-antialiased;
    }

    /* Reader's paper appearance. */
    html.paper {
        height: 100%;
    }

    html.paper body {
        height: calc(100% - 44px);
    }

    html.paper body:after {
        content: "";
        height: 22px;
        display: block;
    }

    /* Clearfix */
    html.paper .page::after {
        content: "";
        display: table;
        clear: both;
    }

    html.paper #article {
        min-height: 100%;
        margin: 22px auto 0 auto;
    }

    html.paper #article :nth-child(1 of .page), html.paper #article :nth-child(1 of .page):nth-last-child(1 of .page) {
        padding-top: 53px;
    }

    html.paper #article html.paper #article {
        /* Stop lining this text up with .page's right margin. */
        top: 14px;
        right: 0px;
    }
</style>
