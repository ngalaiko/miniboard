declare module 'svelte-select' {
    export default interface Select {
        items: any[]
        selectedValue: any
        isClearable: boolean
        isSearchable: boolean
        listAutoWidth: boolean
    }
}
