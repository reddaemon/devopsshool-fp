import TokenService from './TokenService';
import api from './api'

const API_URL = process.env.VUE_APP_ROOT_API;

class AuthService {
    async login({email, password }) {
      return api
        .post(API_URL + 'login', {
            email,
            password
        })
        .then(response => {
            if (response.data.access_token) {
                TokenService.setUser(response.data);
            }
            return response.data;
        });
    }

    logout() {
      return api
        .post(API_URL + 'logout')  
        .then(response => {
            console.log(response)
            TokenService.removeUser();
        });

    }

    register(email, password) {
        return api
          .post(API_URL + 'register', {
              email,
              password
          });
        }
    }

export default new AuthService();