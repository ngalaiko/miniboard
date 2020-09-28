let articleReader = document.querySelector('#article-reader')
let articleList = document.querySelector('#article-list')
let feedList = document.querySelector('#feed-list')

const loadState = () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    articleList.setAttribute('feed', urlParams.get('feed'))
    articleReader.setAttribute('article', urlParams.get('article'))
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
    articleReader.setAttribute('article', e.detail.id)
    storeState('article', e.detail.id)
})

feedList.addEventListener('FeedSelected', (e) => {
    articleList.setAttribute('feed', e.detail.id)
    storeState('feed', e.detail.id)
})

window.addEventListener('popstate', loadState)

loadState()
