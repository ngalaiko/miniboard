import FeedService from './services/FeedService.js'

const feedList = document.querySelector('#feed-list')

const controllContainer = document.querySelector('#controll__container')
const articlesContainer = document.querySelector('#articles-container')
const articlesBackButton = document.querySelector('#article-list-button')
const articleList = document.querySelector('#article-list')
const articleListTitle = document.querySelector('#article-list-title')

const articleReader = document.querySelector('#article-reader')

const showArticles = async (feedId) => {
    feedList.classList.add('hidden')
    articlesContainer.classList.remove('hidden')

    const feed = await FeedService.get(feedId)
    articleListTitle.innerText = feed.title
}

const showFeeds = () => {
    feedList.classList.remove('hidden')
    articlesContainer.classList.add('hidden')
}

const showReader = () => {
    controllContainer.classList.add('mobile-hidden')
    articleReader.classList.remove('hidden')
}

const showControll = () => {
    controllContainer.classList.remove('mobile-hidden')
    articleReader.classList.add('hidden')
}

const loadState = () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    const feedId = urlParams.get('feed')
    const articleId = urlParams.get('article')

    articleList.setAttribute('feed', feedId)
    articleReader.setAttribute('article', articleId)

    if (articleId) {
        showReader()
    } else {
        showControll()
    }

    if (feedId) {
        showArticles(feedId)
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

    showReader()
})

feedList.addEventListener('FeedSelected', (e) => {
    articleList.setAttribute('feed', e.detail.id)
    storeState('feed', e.detail.id)

    showArticles(e.detail.id)
})

articlesBackButton.addEventListener('click', (e) => {
    articleList.removeAttribute('feed')
    deleteState('feed')

    showFeeds()
})

articleReader.addEventListener('CloseClicked', (e) => {
    articleList.removeAttribute('article')
    deleteState('article')

    showControll()
})

window.addEventListener('popstate', loadState)

loadState()
