export default function users(api) {
    let $ = {}

    $.create = async (username, password) => {
        return await api.post(`/api/v1/users`, {
            username: username,
            password: password ,
        })
    }

    $.get = async (name) => {
        return await api.get(`/api/v1/${name}`)
    }

    return $
}
