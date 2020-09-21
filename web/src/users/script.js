import './components/UserFeedController/UserFeedController.js'
import './components/UserArticleController/UserArticleController.js'

let userFeedController = document.querySelector('#user-feed-controller')
let userArticleController = document.querySelector('#user-article-controller')

userFeedController.addEventListener('FeedSelected', (e) => {
    console.log("feed", e.detail.id, "selected")
    userArticleController.feedId = e.detail.id
})

userArticleController.addEventListener('ArticleSelected', (e) => {
    console.log("article", e.detail.id, "selected")
})
