import Api from './api/api';

export default function client() {
    let $ = {}

    $.api = new Api()

    return $
}
