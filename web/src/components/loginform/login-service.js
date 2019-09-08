export default function LoginService(authorizations, users) {
    let $ = {}

    $.login = async (username, password) => {
        let resp = await authorizations.getAuthorization(username, password)

        if (resp.error == undefined) {
            // no error, user exists
            return resp
        }

        if (resp.code === 3) { // 3 - bad request, show the error message
            throw resp.message
        } else if (resp.code != 5) { // 5 means that user is not found, create one
            // unknown error
            console.error(resp)
            throw 'something went wrong, try again?'
        }

        // user doesn't exist
        await users.create(username, password)

        return await $.login(username, password)
    }

    return $
}
