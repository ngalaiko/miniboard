import { LocalStorage } from '../localstorage/local-storage'

export class Api {
  constructor() {
    this.localStorage = new LocalStorage()
  }

  get(url) {
    let options = {
        method: 'GET',
        headers: {
          'Authorization': this.authorization()
      }
    }
    return fetch(url, options)
      .then(response => response.json())
      // todo: handle errors
  }

  post(url, data) {
    let options = {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-Type':  'application/json',
            'Authorization': this.authorization()
        },
    }
    return fetch(url, options)
      .then(response => response.json())
      // todo: handle errors
  }

  patch(url, data) {
    let options = {
        method: 'PATCH',
        body: JSON.stringify(data),
        headers: {
            'Content-Type':  'application/json',
            'Authorization': this.authorization()
        },
    }
    return fetch(url, options)
      .then(response => response.json())
      // todo: handle errors
  }

  delete(url) {
    let options = {
        method: 'DELETE',
        headers: {
          'Authorization': this.authorization()
      }
    }
    return fetch(url, options)
      .then(response => response.json())
      // todo: handle errors
  }

  authenticate(auth) {
    this.localStorage.set('authentication.access_token', auth.access_token)
    this.localStorage.set('authentication.token_type', auth.token_type)
  }

  authorization() {
    let access_token = this.localStorage.get('authentication.access_token')
    let token_type = this.localStorage.get('authentication.token_type')
    if (access_token !== null && token_type !== null) {
      return `${token_type} ${access_token} `
    }
    return ''
  }

  authorized() {
    return this.authorization() != ''
  }

  logout() {
    this.localStorage.remove('authentication.access_token')
    this.localStorage.remove('authentication.token_type')
  }

  subject() {
    let token = this.localStorage.get('authentication.access_token')
    if (token == null) {
      return ''
    }
    let base64Url = token.split('.')[1]
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''));
    return JSON.parse(jsonPayload).sub
  }
}
