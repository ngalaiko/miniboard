import navid from '../navaid/navaid'

import NotFound from '../../pages/notfound/NotFound.svelte'

export class Router {
  constructor() {
    let router = navid()
    this.routes = {}

    this.notFound = [{
      component: NotFound
    }]

    router.listen()
  }

  register(path, component, props) {
    this.routes[path] = [{
      component: component,
      props: props
    }]
  }

  route(path) {
    if (path in this.routes) {
      return this.routes[path]
    }
    return this.notFound
  }
}
