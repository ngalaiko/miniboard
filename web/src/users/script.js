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

const feedsList = document.getElementById('feeds-list')

const loadFeeds = async () => {
    let response = await fetch(`/api/v1/${username}/feeds?page_size=10`)
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
        return
    }

    let body = (await response.json())
    let nextPageToken = body.nextPageToken
    let feeds = body.feeds

    if (nextPageToken === '' && feeds.length == 0) {
        feedsList.innerText = 'Empty'
        return
    }

    feedsList.innerText = ''
    feeds.forEach(addFeed)
}

const addFeed = (feed) => {
    let li = document.createElement('li')
    li.id = `$feedarticle.name}-container`
    li.innerHTML = `
    <div id="${feed.name}">
        ${feed.title}
    </div>`

    let child = feedsList.firstChild
    if (child === null) {
        feedsList.appendChild(li)
        return
    }

    while (child && li.id < child.firstElementChild.id) {
        child = child.nextSibling
    }

    feedsList.insertBefore(li, child)
}

//

const articlesList = document.getElementById('articles-list')

const loadArticles = async () => {
    let response = await fetch(`/api/v1/${username}/articles?page_size=10`)
    if (response.status !== 200) {
        alert(`Error: ${(await response.json()).message}`)
        return
    }

    let body = (await response.json())
    let nextPageToken = body.nextPageToken
    let articles = body.articles

    if (nextPageToken === '' && articles.length == 0) {
        articlesList.innerText = 'Empty'
        return
    }

    articlesList.innerText = ''
    articles.forEach(addArticle)
}

const addArticle = (article) => {
    let li = document.createElement('li')
    li.id = `${article.name}-container`
    li.innerHTML = `
    <div id="${article.name}">
        ${article.title}
    </div>`

    let child = articlesList.firstChild
    if (child === null) {
        articlesList.appendChild(li)
        return
    }

    while (child && li.id < child.firstElementChild.id) {
        child = child.nextSibling
    }

    articlesList.insertBefore(li, child)
}

//

const init = async () => {
    await setCurrentUser()
    loadFeeds()
    loadArticles()
}

init()
