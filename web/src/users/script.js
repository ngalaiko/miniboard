let username = undefined

const init = async () => {
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

init()
