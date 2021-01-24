const apiUrl = window.location.hostname == 'localhost'
    ? 'http://localhost:80'
    : 'https://api.miniboard.app'

const loadFeeds = async () => {
    const response = await fetch(apiUrl + '/v1/feeds', {
        credentials: 'include',
    })

    const body = await response.json()

    if (response.status !== 200) {
        console.error(`failed to fetch feeds: ${body.message}`)
        return
    }

    const list = document.querySelector('#feeds-list')
    body.feeds.map(feed => {
        const item = list.appendChild(
            document.createElement('li').appendChild(
                document.createElement('x-feed')
            )
        )

        item.setAttribute('title', feed.title)
        if (item.icon_url !== undefined) {
            item.setAttribute('icon', feed.icon_url)
        }
    })
}

document.querySelector("#add-button").addEventListener('click', (e) => {
    import('./components/modal.js')

    const addModal = document.createElement('add-modal')
    document.body.appendChild(addModal)

    addModal.addEventListener('Closed', (e) => {
        document.body.removeChild(addModal)
    })

    addModal.addEventListener('FileAdded', (e) => {
        new Promise((resolve, reject) => {
            const reader = new FileReader()
            reader.readAsBinaryString(e.detail.file)
            reader.onload = e => {
                resolve(e.target.result)
            }
            reader.onerror = () => reject(new Error('failed to read file'))
        })
        .then(console.log)
    })

    addModal.addEventListener('UrlAdded', (e) => console.log(e.detail.url))
})

loadFeeds()
