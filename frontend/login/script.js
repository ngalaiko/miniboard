const loginButton = document.getElementById('login-button')
const inputUsername = document.getElementById('username')
const inputPassword = document.getElementById('password')

const apiUrl = window.location.hostname == 'localhost'
    ? 'http://localhost:80'
    : 'https://api.miniboard.app'

const handleButtonClick = async (e) => {
    e.preventDefault()

    const response = await fetch(apiUrl + '/v1/authorizations', {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify({
            username: inputUsername.value,
            password: inputPassword.value,
        })
    })

    if (response.status === 200) {
        document.location = '/users'
        return
    }

    const body = await response.json()
    if (body.message) {
        alert(`error: ${body.message}`)
    } else {
        alert(`error: something went very wrong`)
    }
}

loginButton.addEventListener('click', handleButtonClick)
