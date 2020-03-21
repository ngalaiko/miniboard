<script lang="ts">
  // @ts-ignore
  import { ArticlesClient, Article } from '../../../clients/articles.ts'
  import { ChevronLeftIcon } from 'svelte-feather-icons'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let articleName: string|null
  export let articlesClient: ArticlesClient

  const decode = (article: Article) => {
    return decodeURIComponent(atob(article.content).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''))
  }
</script>

<button class="navigation-bar" on:click={() => {
  dispatch('close')
}}>
  <ChevronLeftIcon size="30" />
</button>
<div class='reader'>
  {#if /users\/.+\/articles\/.+/.test(articleName)}
    {#await articlesClient.get(articleName)}
      loading...
    {:then article}
      <div class='page'>
        <a href={article.url} target="_blank"><h1>{article.title}</h1></a>
        {@html decode(article)}
      </div>
    {/await}
  {/if}
</div>

<style>
  .reader {
    max-height: 100%;
    overflow-y: scroll;
    padding: 0.5em;
    word-wrap: break-word;
    text-align: start;
    line-height: 1.4em;
  }

  h1 {
    margin: 0;
    line-height: 1.2em;
  }

  a {
    display: flex;
    color: inherit;
    text-decoration: none;
  }

  a:hover {
    cursor: pointer;
    background: gainsboro;
  }

  .navigation-bar {
    display: none;
  }

  .page {
    max-width: 660px;
  }

  :global(pre) {
      overflow: auto;
  }

  :global(img) {
    max-width: 100%;
    height: auto;
  }

  :global(iframe) {
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
    .page {
      margin-left: auto;
      margin-right: auto;
      padding-top: 35px;
      padding-bottom: 35px;
    }

    .page video {
      height: auto;
    }
  }

  @media screen and (max-width: 414px) {
    .navigation-bar {
      display: flex;
      width: 100%;
      font: inherit;
      border: 0;
      background: inherit;
      padding: 0;
      margin: 0;
    }
    .reader {
      padding: 1em;
    }
    .page {
      padding-top: 0;
    }
  }

  @media screen and (max-width: 569px) {
    h1.title {
      font-size: 1.5558em;
    }
    h1 {
      font-size: 1.4em;
    }
  }

  @media screen and (min-width: 704px) {
    .page {
      padding-left: 42px;
      padding-right: 42px;
    }
  }

  @media only screen and (min-width: 780px) {
    .reader {
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
</style>
