import api from './api'

const API_URL = 'http://localhost:8080/v1/'

class DataService {
    getRate() {
        return api
          .get(API_URL + 'rate');
    }
}

export default new DataService();