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

const feedsContainer = document.getElementById('feeds-container')

const loadFeeds = async () => {
    let response = await fetch(`/api/v1/${username}/feeds`)
    // todo: handle error
    let body = (await response.json())
    let nextPageToken = body.nextPageToken
    let feeds = body.feeds

    if (nextPageToken === '' && feeds.length == 0) {
        feedsContainer.innerText = 'Empty'
    }
}

const articlesContainer = document.getElementById('articles-container')

const loadArticles = async () => {
    let response = await fetch(`/api/v1/${username}/articles`)
    // todo: handle error
    let body = (await response.json())
    let nextPageToken = body.nextPageToken
    let articles = body.articles

    if (nextPageToken === '' && articles.length == 0) {
        articlesContainer.innerText = 'Empty'
    }
}

const init = async () => {
    await setCurrentUser()
    loadFeeds()
    loadArticles()
}

init()
