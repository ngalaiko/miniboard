<script lang="ts">
  import SvelteInfiniteScroll from 'svelte-infinite-scroll'
  // @ts-ignore
  import { ArticlesClient, Articles, Article } from '../../../clients/articles.ts'
  import ArticleView from './article/Article.svelte'

  export let username: string = ''
  export let articlesClient: ArticlesClient

  let pageToken = ''
  let hasMore = true
  let articles: Articles
  let articlesList: Article[] = []

  const loadMore = () => {
    if (!hasMore) return
    articlesClient.list(username, 25, pageToken).then((articles: Articles) => {
      articlesList = [
        ...articlesList,
        ...articles.articles,
      ]
      hasMore = articles.nextPageToken !== ''
      pageToken = articles.nextPageToken
    })
  }

  loadMore()
</script>

<div class="list">
  {#each articlesList as article}
    <ArticleView article={article} />
  {/each}
  <SvelteInfiniteScroll 
    threshold={100} 
    on:loadMore={loadMore} 
  />
</div>

<style>
  .list{
    display: flex;
    flex-direction: column;
    max-height: 100%;
    overflow-y: scroll;
  }
</style>
