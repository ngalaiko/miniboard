const storeState = (key, value) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))

    if (urlParams.get(key) === value) return

    urlParams.set(key, value)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)
}

const getState = (key) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    const value = urlParams.get(key)
    return value ? value : undefined
}

const deleteState = (key) => {
    const urlParams = new URLSearchParams(window.location.search.slice(1))
    urlParams.delete(key)

    let refresh = window.location.protocol +
        "//" + window.location.host + window.location.pathname +
        `?${urlParams.toString()}`
    window.history.pushState({ path: refresh }, '', refresh)
}

const onSubscriptionSelected = (subscriptionId) => {
    deleteState('tag')
    storeState('subscription', subscriptionId)

    sendMessage("items:load", {
        subscriptionId: subscriptionId,
    })
}

const onItemSelected = (itemId) => {
    sendMessage("items:select", {id: itemId})
    storeState('item', itemId)
}

const onTagSelected = (tagId) => {
    deleteState('subscription')
    storeState('tag', tagId)

    sendMessage("items:load", {
        tagId: tagId,
    })
}

document.querySelector('#items-list').addEventListener('scroll', (e) => {
    const { scrollTop, scrollHeight, clientHeight } = e.target
    const needMore = scrollTop + clientHeight >= scrollHeight - 50
    if (!needMore) return

    const subscriptionId = getState('subscription')
    const tagId = getState('tag')
    const createdLt = document.querySelector('#items-list').lastElementChild.getAttribute('created')

    if (!createdLt) return

    document.querySelector('#items-list').insertAdjacentHTML('beforeend', '<div class="page-separator"></div')

    sendMessage("items:loadmore", {
        tagId: tagId,
        subscriptionId: subscriptionId,
        createdLt: createdLt,
    })
})

window.addEventListener('keydown', (e) => {
    if (isModalClosed()) return

    switch (e.key) {
    case 'Enter':
        closeModal()

        const url = document.querySelector("#input-url").value
        if (url === '') return

        document.querySelector("#input-url").value = ''

        sendMessage("subscriptions:create", {
            url: url,
        })
        break
    case 'Escape':
        closeModal()
        break
    }
})

document.querySelector("#input-file").addEventListener('change', async (e) => {
    const files = e.target.files

    closeModal()

    if (files.length === 0) return

    const promise = new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsBinaryString(files[0])
        reader.onload = e => {
            resolve(e.target.result)
        }
        reader.onerror = () => reject(new Error('failed to read file'))
    })
    .then((raw) => sendMessage("subscriptions:import", { file: raw }))
})

const isModalClosed = () => {
    return document.querySelector('#modal').hidden
}

const closeModal = () => {
    document.querySelector('#modal').hidden = true
}

const showModal = () => {
    document.querySelector('#modal').hidden = false
}

document.querySelector("#background").addEventListener('click', () => {
    if (!isModalClosed()) closeModal()
})

document.querySelector('#add-button').addEventListener('click', () => {
    showModal()
})
