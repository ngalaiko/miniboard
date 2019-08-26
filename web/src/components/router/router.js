import navid from '../navaid/navaid'

export class Router {
  constructor() {
    this.router = navid()
    this.currentComponent = []

  }

  listen() {
    this.router.listen()
  }

  register(path, component, props) {
    this.router.on(path, params => {
        this.currentComponent = [{
          component: component,
          props: { ...props, ...params }
      }]
    })
  }

  current() {
    return this.currentComponent
  }
}
