export class ArticlesService {
  constructor(api, user) {
    this.api  = api
    this.pageSize = 50
    this.user = user
    this.from = ''
    this.items = []
  }

  // next returns a Promise with articles.
  next() {
    return this.api.get(`/api/v1/${this.user.name}/articles?page_size=${this.pageSize}&from=${this.from}`)
    .then(resp => { 
      this.from = resp.page_token 
      return resp.articles
    })
  }

  // create returns a Promise with the article.
  add(url) {
    return this.api.post(`/api/v1/${this.user.name}/articles`, {
      url: url
    })
    .then(this.ifError)
  }

  // delete deletes the article, returns nothing.
  delete(article) {
      return this.api.delete(`/api/v1/${article.name}`)
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
