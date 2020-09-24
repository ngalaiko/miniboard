let readerController = document.querySelector('#reader-controller')
let articleController = document.querySelector('#article-controller')
let feedController = document.querySelector('#feed-controller')

articleController.addEventListener('ArticleSelected', (e) => {
    readerController.articleId = e.detail.id
})

feedController.addEventListener('FeedSelected', (e) => {
    articleController.feedId = e.detail.id
})
