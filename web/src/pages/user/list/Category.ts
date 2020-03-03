import { ListParams } from '../../../clients/articles'

export class Category {
  value: string
  label: string
  listParams: ListParams

  constructor (value: string, label: string, listParams: ListParams) {
    this.value = value
    this.label = label
    this.listParams = listParams
  }
}

export const Categories = {
    All: new Category('all', 'All', new ListParams()),
    Unread: new Category('unread', 'Unread', new ListParams().withRead(false)),
    Favorite: new Category('favorite', 'Favorite', new ListParams().withFavorite(true)),
}
