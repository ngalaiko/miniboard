const urlParams = new URLSearchParams(window.location.search)
const codeParam = urlParams.get('code')

fetch('/api/v1/tokens', {
    method: 'POST',
    body: JSON.stringify({
        code: codeParam
    })
}).then(async (response) => {
    switch (response.status) {
    case 200:
        document.location = '/users'
        break
    case 400:
        alert(`Error: ${(await response.json()).error}`)
        document.location = '/'
        break
    default:
        alert('Error: Something went wrong\nPlease try again')
        document.location = '/'
        break
    }
})
