import api from './api'

const API_URL = process.env.VUE_APP_ROOT_API

class DataService {
    getRate() {
        return api
          .get(API_URL + 'rate');
    }
}

export default new DataService();