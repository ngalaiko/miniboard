<script lang="ts">
  // @ts-ignore
  import { ArticlesClient, Article } from '../../../clients/articles.ts'

  export let articleName: string|null
  export let articlesClient: ArticlesClient

  const decode = (article: Article) => {
    return decodeURIComponent(atob(article.content).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''))
  }
</script>

<div class='reader'>
  {#if articleName}
    {#await articlesClient.get(articleName)}
      loading...
    {:then article}
      <a class="header" href={article.url} target="_blank">
        <h1>{article.title}</h1>
      </a>
      {@html decode(article)}
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
    overflow-x: hidden;
    padding-left: 40px;
    padding-right: 40px;
    align-items: center;
  }

  .header {
    padding: 20px;
    padding-right: 20px;
    margin-top: 20px;
    width: 100%;
    display: flex;
    flex-direction: column;
    max-width: 660px;
  }

  h1 {
    margin: 0;
  }

  a {
    color: inherit;
    text-decoration: none;
  }

  .header:hover {
    background: gainsboro;
    cursor: pointer;
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
</style>
