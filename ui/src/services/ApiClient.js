import AppConfig from './AppConfig';

class ApiClient {
  static post(url, data) {
    return fetch(`http://localhost:${AppConfig.port()}/${url}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
  }
}

export default ApiClient;
