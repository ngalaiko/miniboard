export class ArticlesService {
  constructor(api, user) {
    this.api  = api
    this.pageSize = 50
    this.user = user
    this.from = ''
    this.items = []
  }

  next() {
    return this.api.get(`/api/v1/${this.user.name}/articles?page_size=${this.pageSize}&from=${this.from}`)
    .then(resp => { 
      this.from = resp.page_token 
      return resp.articles
    })
  }
}
