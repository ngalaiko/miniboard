<script lang="ts">
  // @ts-ignore
  import { Article } from '../../../../clients/articles.ts'
  import TimeAgo from './TimeAgo.svelte'

  export let article: Article
  export let isSelected = false
</script>

<div class="article {isSelected ? 'selected' : ''}" on:click>
  <div class="header">
    <div class="meta">{article.siteName !== '' ? article.siteName : new URL(article.url).hostname}</div>
    <TimeAgo date={new Date(article.createTime)} />
  </div>
  <div class="title">{article.title}</div>
</div>

<style>
  .header {
    display: flex;
    justify-content: space-between;
    opacity: 70%;
    margin-bottom: 2px;
    font-size: 0.8em;
  }

  .meta {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .article {
    display: flex;
    flex-direction: column;
    border-bottom: 1px solid;
    padding: 5px;
    line-height: 1.1em;
    cursor: pointer;
  }

  .selected {
    background: gainsboro;
  }
</style>
