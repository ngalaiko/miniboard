let readerController = document.querySelector('#reader-controller')
let articleList = document.querySelector('#article-list')
let feedController = document.querySelector('#feed-controller')

const loadState = () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    articleList.setAttribute('feed', urlParams.get('feed'))
}

const storeState = (key, value) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get(key) === value) return 

    urlParams.set(key, value)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)
}

articleList.addEventListener('ArticleSelected', (e) => {
    readerController.articleId = e.detail.id
})

feedController.addEventListener('FeedSelected', (e) => {
    articleList.setAttribute('feed', e.detail.id)
    storeState('feed', e.detail.id)
})

window.addEventListener('popstate', loadState)

loadState()
