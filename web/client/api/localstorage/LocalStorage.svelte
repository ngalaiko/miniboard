<script context='module'>
    export const LocalStorage = () => {
        let $ = {}

        const storage = window.localStorage

        $.set = (key, value) => {
            storage.setItem(key, value)
        }

        $.get = (key) => {
            return storage.getItem(key)
        }

        $.remove = (key) => {
            return storage.removeItem(key)
        }

        return $
    }
</script>
