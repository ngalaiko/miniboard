<script lang="ts">
  import SvelteInfiniteScroll from 'svelte-infinite-scroll'
  // @ts-ignore
  import { ArticlesClient, Articles, Article, ListParams } from '../../../clients/articles.ts'
  import ArticleView from './Article.svelte'
  import { onMount } from 'svelte'

  export let username: string = ''
  export let articlesClient: ArticlesClient
  export let listParams: ListParams = new ListParams(false, false)

  let pageToken = ''
  let hasMore = true
  let articlesList: Article[] = []

  const loadMore = async () => {
    const articles = await articlesClient.list(username, 25, pageToken, listParams)

    articlesList = [
      ...articlesList,
      ...articles.articles,
    ]

    pageToken = articles.nextPageToken
    hasMore = pageToken !== ''
  }

  onMount(loadMore)
</script>

<div class="list">
  {#each articlesList as article}
    <ArticleView 
      article={article} 
      on:selected
    />
  {/each}
  <SvelteInfiniteScroll 
    threshold={100} 
    hasMore={hasMore}
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
