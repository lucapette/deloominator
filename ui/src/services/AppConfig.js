class AppConfig {
  static port() {
    return this.root().getAttribute('data-port');
  }

  static root() {
    return document.getElementById('root');
  }
}

export default AppConfig;
