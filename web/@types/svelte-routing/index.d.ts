declare module 'svelte-routing' {
    interface State {
        replace?: boolean
        state?: any
    }

    export function navigate(to: string, state?: State): void

    export interface Router {
        basepath?: string
        url?: string
    }

    export interface Route {
        path?: string
        component?: any
    }

    export interface Link {
        to: string
        replace?: boolean
        state?: State
        getProps?: (any) => any
    }
}
