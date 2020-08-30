let username = undefined

const setCurrentUser = async () => {
    let response = await fetch('/api/v1/users/me')
    switch (response.status / 100) {
    case 2:
        username = (await response.json()).name
        break
    default:
        alert(`Error: ${(await response.json()).error}`)
        document.location = '/'
        break
    }
}

//

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
    let response = await fetch(`/api/v1/${username}/sources`, {
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

    let response = await fetch(`/api/v1/${username}/sources`, {
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
const feedsListPlaceholder = document.getElementById('feeds-list-placeholder')

let feedsPageToken = undefined

const loadFeeds = async () => {
    if (feedsPageToken === "") return
    if (feedsPageToken === undefined) feedsPageToken = ""

    let response = await fetch(`/api/v1/${username}/feeds?page_size=10&page_token=${feedsPageToken}`)
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

    if (urlParams.get('source') === e.target.id) return

    urlParams.set('source', e.target.id)

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
    li.id = `${feed.name}`
    li.innerHTML = `${feed.title}`

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

const articlesList = document.getElementById('articles-list')
const articlesListPlaceholder = document.getElementById('articles-list-placeholder')

let articlesPageToken = undefined

const loadArticles = async () => {
    if (articlesPageToken === "") return
    if (articlesPageToken === undefined) articlesPageToken = ""

    const urlParams = new URLSearchParams(window.location.search.slice(1))

    let articlesUrl = `/api/v1/${username}/articles?`
        + `page_size=10`
        + `&page_token=${articlesPageToken}`

    if (urlParams.get('source')) {
        articlesUrl += `&source_name=${urlParams.get('source')}`
    }

    let response = await fetch(articlesUrl)
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
        return
    }

    let body = (await response.json())
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
    li.id = `${article.name}-container`
    li.classList.add('article')

    let div = document.createElement('div')
    div.id = article.name
    div.innerHTML = `
        ${article.title}
    `
    div.addEventListener('click', handleSelectArticle)

    li.appendChild(div)

    let child = articlesList.firstChild
    if (child === null) {
        articlesList.appendChild(li)
        return
    }

    while (child && child.firstElementChild && li.id < child.firstElementChild.id) {
        child = child.nextSibling
    }

    articlesList.insertBefore(li, child)
}

articlesList.addEventListener('scroll', handleScroll(loadArticles))

//

const loadReader = async () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    const selectedArticleName = urlParams.get('article')

    if (!selectedArticleName) return

    await displayArticle(selectedArticleName)
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

const displayArticle = async (articleName) => {
    if (readerContent.name === articleName) return

    let response = await fetch(`/api/v1/${articleName}?view=ARTICLE_VIEW_FULL`)
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
    readerContent.name = articleName
}

//

const updateVisibility = async () => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get('source')) {
        setVisible(articlesList, true)
        setVisible(feedsList, false)
    } else {
        setVisible(articlesList, false)
        setVisible(feedsList, true)
    }

    setVisible(readerContainer, urlParams.get('article'))
}

window.addEventListener('popstate', updateVisibility)
window.addEventListener('popstate', loadReader)

//

const init = async () => {
    await setCurrentUser()
    loadFeeds()
    loadArticles()
    loadReader()
    updateVisibility()
}

init()
