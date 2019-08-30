export class ArticlesService {
  constructor(api, user) {
    this.api = api
    this.user = user
  }

  // create returns a Promise with the article.
  add(url) {
    return this.api.post(`/api/v1/${this.user.name}/articles`, {
      url: url
    })
    .then(this.ifError)
  }

  ifError(resp) {
    if (resp.error === undefined) {
      return resp
    }
    if (resp.code !== undefined) {
      throw resp.message
    }
    console.error(resp)
    throw "something went wrong, try again"
  }
}
