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
      {@html decode(article)}
    {/await}
  {/if}
</div>

<style>
  .reader {
    max-height: 100%;
    max-width: 100%;
    overflow-y: scroll;
    padding: 5px;
  }
</style>
