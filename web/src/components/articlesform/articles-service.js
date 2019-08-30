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

  ifError(error) {
    if (error.code !== undefined) {
      throw error.message
    }
    console.error(error)
    throw "something went wrong, try again"
  }
}
