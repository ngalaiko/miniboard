export class Api {
  constructor() {
    this.authentication = null
  }

  postJSON(url, data) {
    return this.post(url, JSON.stringify(data), 'application/json')
  }

  post(url, data, contentType) {
    return fetch(url, {
        method: 'POST',
        body: data,
        headers: {
            'Content-Type': contentType,
        },
    })
    .then(response => response.json())
  }

  authenticate(auth) {
    this.authentication = auth
  }
}
