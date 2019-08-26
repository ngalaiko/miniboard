import { LocalStorage } from '../localstorage/local-storage'

export class Api {
  constructor() {
    this.localStorage = new LocalStorage();
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
}
