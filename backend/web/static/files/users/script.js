const itemUrl = (id) => {
    const path = window.location.pathname;
    if (path.endsWith('/items/')) return path + `${id}/`;
    return path.split('/items')[0] + '/items/' + `${id}/`;
}

const subscriptionUrl = (id) => {
    const path = window.location.pathname;
    if (path.endsWith('/items/')) return `/users/subscriptions/${id}/items/`
    const split = path.split('/');
    const lastId = split[split.length-2];
    if (path.includes('/items/')) return `/users/subscriptions/${id}/items/${lastId}/`;
    return `/users/subscriptions/${id}/items/` + path.replace('/users/', '');
}

const tagUrl = (id) => {
    const path = window.location.pathname;
    if (path.endsWith('/items/')) return `/users/tags/${id}/items/`
    const split = path.split('/');
    const lastId = split[split.length-2];
    if (path.includes('/items/')) return `/users/tags/${id}/items/${lastId}/`;
    return `/users/tags/${id}/items/` + path.replace('/users/', '');
}

const nav = (to) => {
    let refresh = window.location.protocol + "//" + window.location.host + to;
    window.history.pushState({ path: refresh }, '', refresh)
}

const getState = (key) => {
    const path = window.location.pathname;
    const split = path.split('/');
    const keyIndex = split.indexOf(key);
    return keyIndex > 0 ? split[keyIndex + 1] : undefined;
}

const onSubscriptionSelected = (subscriptionId) => {
    nav(subscriptionUrl(subscriptionId))
    showItems()

    sendMessage("items:load", {
        subscriptionId: subscriptionId,
    })
}

const onItemSelected = (itemId) => {
    sendMessage("items:select", {id: itemId})
    nav(itemUrl(itemId))
}

const onTagSelected = (tagId) => {
    nav(tagUrl(tagId))
    showItems()

    sendMessage("items:load", {
        tagId: tagId,
    })
}

const hideItems = () => {
    showTagsSubscriptions()
    const itemId = getState('items')
    if (itemId) {
        nav(`/users/${itemId}/`)
    } else {
        nav("/users/")
    }
}

const showItems = () => {
    document.querySelector('#tags-subscriptions').classList.add('d-hide')
    document.querySelector('#items').classList.remove('d-hide')
}

const showTagsSubscriptions = () => {
    document.querySelector('#tags-subscriptions').classList.remove('d-hide')
    document.querySelector('#items').classList.add('d-hide')
}

document.querySelector('#items-list').addEventListener('scroll', (e) => {
    const { scrollTop, scrollHeight, clientHeight } = e.target
    const needMore = scrollTop + clientHeight >= scrollHeight - 50
    if (!needMore) return

    const subscriptionId = getState('subscriptions')
    const tagId = getState('tags')
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
    return !document.querySelector('#modal').classList.contains('active')
}

const closeModal = () => {
    document.querySelector('#modal').classList.remove('active')
}

const showModal = () => {
    document.querySelector('#modal').classList.add('active')
}

document.querySelector("#background").addEventListener('click', () => {
    if (!isModalClosed()) closeModal()
})

document.querySelector('#add-button').addEventListener('click', () => {
    showModal()
})
