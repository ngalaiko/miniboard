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

<div class='reader'>
  <button class="navigation-bar" on:click={() => {
    history.pushState(null, "", `#`)
    dispatch('close')
  }}>
    <ChevronLeftIcon size="24" />
  </button>
  {#if articleName}
    {#await articlesClient.get(articleName)}
      loading...
    {:then article}
      <div class='page'>
        <h1 href={article.url} target="_blank">{article.title}</h1>
        {@html decode(article)}
      </div>
    {/await}
  {/if}
</div>

<style>
  .reader {
    display: flex;
    flex-direction: column;
    max-height: 100%;
    max-width: 100%;
    overflow-y: scroll;
    align-items: center;
  }

  h1 {
    margin: 5px;
  }

  a {
    color: inherit;
    text-decoration: none;
  }

  .header:hover {
    background: gainsboro;
    cursor: pointer;
  }

  .navigation-bar {
    display: none;
  }

  @media screen and (max-width: 414px) {
    .navigation-bar {
      display: flex;
      width: 100%;
      padding: 0;
      font: inherit;
      border: 0;
      background: inherit;
    }
  }

  :global(pre) {
      overflow: auto;
  }

  :global(.page) {
    max-width: 660px;
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

  .reader {
    -webkit-font-smoothing: subpixel-antialiased;
    padding: 5px;
  }
</style>
