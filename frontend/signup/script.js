import UsersService from '/users/services/users.js'

const signupButton = document.getElementById('signup-button')
const inputUsername = document.getElementById('username')
const inputPassword = document.getElementById('password')

const apiUrl = window.location.hostname == 'localhost'
    ? 'http://localhost:80'
    : 'https://api.miniboard.app'

const handleButtonClick = async (e) => {
    e.preventDefault()

    UsersService.create({
        username: inputUsername.value,
        password: inputPassword.value,
    }).then(() => {
        alert('You are now signed up')
        document.location = '/login'
    }).catch((error) => {
        alert(`error: ${error}`)
    })
}

signupButton.addEventListener('click', handleButtonClick)
