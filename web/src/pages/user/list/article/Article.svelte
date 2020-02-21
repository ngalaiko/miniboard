<script lang="ts">
  // @ts-ignore
  import { Article } from '../../../../clients/articles.ts'

  export let article: Article
</script>

<div class="article" on:click>
  {article.title}
</div>

<style>
  .article {
    border-top: 1px solid;
    border-bottom: 1px solid;
    padding: 3px;
  }
</style>
