export class LabelsService {
  constructor(api, user) {
    this.api = api
    this.user= user
  }

  create(title) {
    return this.api.post(`/api/v1/${this.user.name}/labels`, {
      title: title,
    })
  }

  get(labelName) {
    return this.api.get(`/api/v1/${labelName}`)
  }
}
