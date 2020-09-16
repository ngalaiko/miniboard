const isVisible = (e) => {
    return e.style.display !== 'none'
}

const setVisible = (e, visible) => {
    e.style.display = visible ? '' : 'none'
}

//

const addFormButton = document.getElementById('add-form-button')
const addFormInput = document.getElementById('add-form-input')

const handleAddFormButton = async (e) => {
    e.preventDefault()

    await handleAdd(addFormInput.value)

    addFormInput.value = ''
}

const handleAdd = async (url) => {
    let response = await fetch(`/api/v1/sources`, {
        method: 'POST',
        body: JSON.stringify({
            url: url
        })
    })
    switch (response.status) {
        case 200:
            let resourceResponse = await fetch(`/api/v1/${((await response.json()).name)}`)
            let resource = await resourceResponse.json()
            switch (resource.name.split('/')[2]) {
                case "articles":
                    addArticle(resource)
                    break
                case "feeds":
                    addFeed(resource)
                    break
            }
            break
        case 409:
            break
        default:
            alert(`Error: ${(await response.json()).message}`)
            break
    }
}

addFormButton.addEventListener('click', handleAddFormButton)

//

const addFormFile = document.getElementById('add-form-file')

const handleAddFormFile = async (e) => {
    if (e.target.files.length === 0) return

    const file = e.target.files[0]

    const content = await new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsBinaryString(file)
        reader.onload = e => {
            if (e.target === null) return
            resolve(e.target.result)
        }
        reader.onerror = () => reject(new Error('failed to read file'))
    })

    let response = await fetch(`/api/v1/sources`, {
        method: 'POST',
        body: JSON.stringify({
            raw: btoa(content)
        })
    })
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
    }
}

addFormFile.addEventListener('change', handleAddFormFile)

//

const handleScroll = (loadMore) => {
    return async (e) => {
        const { scrollTop, scrollHeight, clientHeight } = e.target

        if (scrollTop + clientHeight >= scrollHeight - 25) {
            if (isVisible(e.target.lastElementChild)) return

            setVisible(e.target.lastElementChild, true)

            await loadMore()

            setVisible(e.target.lastElementChild, false)
        }
    }
}

//

const feedsList = document.getElementById('feeds-list')
const feedsContainer = document.getElementById('feeds-container')
const feedsListPlaceholder = document.getElementById('feeds-list-placeholder')

let feedsPageToken = undefined

const loadFeeds = async () => {
    if (feedsPageToken === "") return
    if (feedsPageToken === undefined) feedsPageToken = ""

    let response = await fetch(`/api/v1/feeds?page_size=10&page_token=${feedsPageToken}`)
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
        return
    }

    let body = (await response.json())
    feedsPageToken = body.nextPageToken
    let feeds = body.feeds

    feedsList.lastElementChild.style.visibility = 'hidden'
    feeds.forEach(addFeed)

    const { scrollTop, scrollHeight, clientHeight } = feedsList
    let isListFull = (scrollTop + clientHeight >= scrollHeight)

    if (isListFull && feedsPageToken !== '') await loadFeeds()
}

const handleSelectFeed = async (e) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get('feed') === e.target.id) return

    urlParams.set('feed', e.target.id)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)

    updateVisibility()
    await reloadArticles()
}

const addFeed = (feed) => {
    feedsListPlaceholder.hidden = true

    let li = document.createElement('li')
    li.id = feed.id
    li.innerText = feed.title
    li.classList.add('feed')
    li.title = feed.title

    li.addEventListener('click', handleSelectFeed)

    let child = feedsList.firstChild
    if (child === null) {
        feedsList.appendChild(li)
        return
    }

    while (child && child && li.id < child.id) {
        child = child.nextSibling
    }

    feedsList.insertBefore(li, child)
}

feedsList.addEventListener('scroll', handleScroll(loadFeeds))

//

const articlesContainer = document.getElementById('articles-container')
const articlesTitle = document.getElementById('articles-title')
const articlesList = document.getElementById('articles-list')
const articlesListPlaceholder = document.getElementById('articles-list-placeholder')
const articlesBackButton = document.getElementById('articles-back-button')

let articlesPageToken = undefined

const handleArticlesBackButton = async (e) => {
    e.preventDefault()

    const urlParams = new URLSearchParams(window.location.search.slice(1))

    urlParams.delete('feed')

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`

    window.history.pushState({ path: refresh }, '', refresh)
    updateVisibility()
}

articlesBackButton.addEventListener('click', handleArticlesBackButton)

const updateArticlesTitle = async (feedId) => {
    let feedElement = document.getElementById(feedId)
    if (feedId) {
        articlesTitle.innerText = feedElement.title
        return
    }

    let response = await fetch(`/api/v1/feeds/${feedId}`)
    if (response.status !== 200)  {
        alert(`Error: ${await response.json().message}`)
        return
    }

    let body = await response.json()
    articlesTitle.innerText = body.title
}

const loadArticles = async () => {
    if (articlesPageToken === "") return
    if (articlesPageToken === undefined) articlesPageToken = ""

    const urlParams = new URLSearchParams(window.location.search.slice(1))

    let articlesUrl = `/api/v1/articles?`
        + `page_size=10`
        + `&page_token=${articlesPageToken}`

    if (urlParams.get('feed')) {
        let feedId = urlParams.get('feed')
        articlesUrl += `&feed_id_eq=${feedId}`
        updateArticlesTitle(feedId)
    }

    let response = await fetch(articlesUrl)
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
        return
    }

    let body = await response.json()
    articlesPageToken = body.nextPageToken
    let articles = body.articles

    articlesList.lastElementChild.hidden = true
    articles.forEach(addArticle)

    const { scrollTop, scrollHeight, clientHeight } = articlesList

    let isListFull = (scrollTop + clientHeight >= scrollHeight)

    if (isListFull && articlesPageToken !== '') await loadArticles()
}

const reloadArticles = async () => {
    document.querySelectorAll('.article').forEach(n => n.remove())
    setVisible(articlesList.lastElementChild, true)

    articlesPageToken = undefined
    await loadArticles()

    setVisible(articlesList.lastElementChild, false)
}

const addArticle = (article) => {
    articlesListPlaceholder.hidden = true

    let li = document.createElement('li')
    li.id = article.id
    li.classList.add('article')
    li.innerText = article.title
    li.addEventListener('click', handleSelectArticle)

    let child = articlesList.firstChild
    if (child === null) {
        articlesList.appendChild(li)
        return
    }

    while (child && li.id < child.id) {
        child = child.nextSibling
    }

    articlesList.insertBefore(li, child)
}

articlesList.addEventListener('scroll', handleScroll(loadArticles))

//

const loadReader = async () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    const selectedArticleId = urlParams.get('article')

    if (!selectedArticleId) return

    await displayArticle(selectedArticleId)
}

const handleSelectArticle = async (e) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get('article') === e.target.id) return

    urlParams.set('article', e.target.id)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`

    window.history.pushState({ path: refresh }, '', refresh)
    updateVisibility()

    await displayArticle(e.target.id)
}

const readerContainer = document.getElementById('reader-container')
const readerContent = document.getElementById('reader-content')
const readerLink = document.getElementById('reader-link')
const readerTitle = document.getElementById('reader-title')

const displayArticle = async (articleId) => {
    if (readerContent.id === articleId) return

    let response = await fetch(`/api/v1/articles/${articleId}?view=ARTICLE_VIEW_FULL`)
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
        return
    }

    let article = await response.json()

    let decodedContent = decodeURIComponent(atob(article.content).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))

    readerTitle.innerText = article.title
    readerLink.href = article.url
    readerContent.innerHTML = decodedContent
    readerContent.id = articleId
}

//

const updateVisibility = async () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get('feed')) {
        setVisible(articlesContainer, true)
        setVisible(feedsContainer, false)
    } else {
        setVisible(articlesContainer, false)
        setVisible(feedsContainer, true)
    }

    setVisible(readerContainer, urlParams.get('article'))
}

window.addEventListener('popstate', updateVisibility)
window.addEventListener('popstate', loadReader)

//

const init = async () => {
    loadFeeds()
    loadArticles()
    loadReader()
    updateVisibility()
}

init()
