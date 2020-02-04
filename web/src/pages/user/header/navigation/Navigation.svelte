<script context='module'>
    import { writable } from 'svelte/store'
    const selectedPaneStore = writable(2)
</script>

<script>
    import { quintOut } from 'svelte/easing'
    import { Link } from 'svelte-routing'

    import BookOpen from '../../../../icons/BookOpen.svelte'
    import List from '../../../../icons/List.svelte'
    import Star from '../../../../icons/Star.svelte'

    export let username

    const titles = {}
    titles[`/${username}/starred`] = "Starred - Miniboard"
    titles[`/${username}/unread`] = "Unread - Miniboard"
    titles[`/${username}/all`] = "All - Miniboard"

    const getProps = (defaultStyle) => {
        return ({ location, href, isPartiallyCurrent, isCurrent }) => {
            if (isCurrent) {
                document.title = titles[href]
                return { class: defaultStyle + ' navigation-selected' }
            }
            return { class: defaultStyle + ' navigation-not-selected'}
        }
    }
</script>

<span class='container'>
    <Link to="/{username}/starred" getProps={getProps('navigation-border navigation-border-left')}>
        <Star size=20 />
    </Link>
    <Link to="/{username}/unread" getProps={getProps('navigation-border navigation-border-middle')}>
        <BookOpen size=20 />
    </Link>
    <Link to="/{username}/all" getProps={getProps('navigation-border navigation-border-right')}>
        <List size=20 />
    </Link>
</span>

<style>
    .container {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    :global(.navigation-selected) {
        background: black;
        color: white;
    }

    :global(.navigation-not-selected) {
        background: white;
        color: black;
    }

    :global(.navigation-border) {
        padding-top: 0.2em;
        padding-bottom: 0.2em;
        padding-left: 0.5em;
        padding-right: 0.5em;
    }

    :global(.navigation-border-left) {
        border-top: 1px solid black;
        border-bottom: 1px solid black;
        border-left: 1px solid black;
    }

    :global(.navigation-border-middle) {
        border: 1px solid black;
    }

    :global(.navigation-border-right) {
        border-top: 1px solid black;
        border-bottom: 1px solid black;
        border-right: 1px solid black;
    }
</style>
