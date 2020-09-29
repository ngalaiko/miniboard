const feedList = document.querySelector('#feed-list')

const articlesContainer = document.querySelector('#articles-container')
const articlesBackButton = document.querySelector('#articles-back-button')
const articleList = document.querySelector('#article-list')

const articleReader = document.querySelector('#article-reader')

const showArticles = () => {
    feedList.classList.add('hidden')
    articlesContainer.classList.remove('hidden')
}

const showFeeds = () => {
    feedList.classList.remove('hidden')
    articlesContainer.classList.add('hidden')
}

const loadState = () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    const feedId = urlParams.get('feed')
    const articleId = urlParams.get('article')

    articleList.setAttribute('feed', feedId)
    articleReader.setAttribute('article', articleId)

    if (feedId) {
        showArticles()
    } else {
        showFeeds()
    }
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

const deleteState = (key) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    urlParams.delete(key)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)
}

articleList.addEventListener('ArticleSelected', (e) => {
    articleReader.setAttribute('article', e.detail.id)
    storeState('article', e.detail.id)
})

feedList.addEventListener('FeedSelected', (e) => {
    articleList.setAttribute('feed', e.detail.id)
    storeState('feed', e.detail.id)

    showArticles()
})

articlesBackButton.addEventListener('click', (e) => {
    articleList.removeAttribute('feed')
    deleteState('feed')
    showFeeds()
})

window.addEventListener('popstate', loadState)

loadState()
