export default function users(api) {
    let $ = {}

    $.create = async (username, password) => {
        return await api.post(`/api/v1/users`, {
            username: username,
            password: password ,
        })
    }

    return $
}
