<script lang="ts">
  // @ts-ignore
  import { Article } from '../../../../clients/articles.ts'
  import TimeAgo from './TimeAgo.svelte'

  export let article: Article
</script>

<div class="article" class:read={article.isRead} on:click>
  <div class="header">
    <img class="header-icon" src="{article.iconUrl}" alt=""/>
    <div class="header-title">{article.siteName !== '' ? article.siteName : new URL(article.url).hostname}</div>
    <div style="margin-left: auto;">
      <TimeAgo date={new Date(article.createTime)} />
    </div>
  </div>
  <div class="title">{article.title}</div>
</div>

<style>
  .header {
    display: flex;
    align-items: center;
    opacity: 70%;
    margin-bottom: 0.2em;
    font-size: 0.8em;
  }

  .header-icon {
    height: 1em;
    width: auto;
    padding-right: 0.3em;
  }

  .header-title {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .article {
    display: flex;
    flex-direction: column;
    line-height: 1.1em;
    cursor: pointer;
    padding: 0.3em;
  }

  .read {
    opacity: 0.5;
  }
</style>
