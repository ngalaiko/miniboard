import Api from './api/api'
import Authorizations from './authorizations/authorizations'
import Articles from './articles/articles'
import Labels from './labels/labels'
import Users from './users/users'

export default function client() {
    let $ = {}

    $.api = new Api()
    $.authorizations = new Authorizations($.api)
    $.articles = new Articles($.api)
    $.labels = new Labels($.api)
    $.users = new Users($.api)

    return $
}
