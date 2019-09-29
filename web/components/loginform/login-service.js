export default function LoginService(authorizations, codes, users) {
    let $ = {}

    $.sendCode = async (email) => {
        let resp = codes.sendCode(email)
    }

    $.login = async (code) => {
        let resp = await authorizations.getAuthorization(code)

        if (resp.error == undefined) {
            return resp
        }

        if (resp.code === 3) {
            throw resp.message
        }

        // unknown error
        console.error(resp)
        throw 'something went wrong, try again?'
    }

    return $
}
