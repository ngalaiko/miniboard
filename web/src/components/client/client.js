import Api from './api/api'
import Authorizations from './authorizations/authorizations'
import Users from './users/users'

export default function client() {
    let $ = {}

    $.api = new Api()
    $.authorizations = new Authorizations($.api)
    $.users = new Users($.api)

    return $
}
