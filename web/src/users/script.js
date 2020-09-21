let userReaderController = document.querySelector('#user-reader-controller')
let userArticleController = document.querySelector('#user-article-controller')
let userFeedController = document.querySelector('#user-feed-controller')

userArticleController.addEventListener('ArticleSelected', (e) => {
    userReaderController.articleId = e.detail.id
})

userFeedController.addEventListener('FeedSelected', (e) => {
    userArticleController.feedId = e.detail.id
})
