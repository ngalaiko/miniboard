const apiUrl = window.location.hostname !== 'localhost'
    ? 'https://api.miniboard.app'
    : 'http://localhost:80';

const loginButton = document.getElementById('login-button')
const inputUsername = document.getElementById('username')
const inputPassword = document.getElementById('password')

const handleButtonClick = async (e) => {
    e.preventDefault()

    const response = await fetch(apiUrl + '/v1/authorizations', {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
            username: inputUsername.value,
            password: inputPassword.value,
        }),
        headers: new Headers({
            "Content-Type": "application/json",
        })
    })

    const body = await response.json()
    if (response.status !== 200) {
        alert(`error: ${body.message}`)
        return
    }

    document.location = '/users'
}

loginButton.addEventListener('click', handleButtonClick)
