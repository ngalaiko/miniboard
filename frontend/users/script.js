const apiUrl = window.location.hostname == 'localhost'
    ? 'http://localhost:80'
    : 'https://api.miniboard.app'

fetch(apiUrl + '/v1/feeds', {
    credentials: 'include',
})
