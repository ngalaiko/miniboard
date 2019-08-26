export class LocalStorage {
  constructor() {
    this.storage = window.localStorage
  }

  set(key, value) {
    this.storage.setItem(key, value)
  }
  
  get(key) {
    return this.storage.getItem(key)
  }
}
