import axios from 'axios';
import TokenService from './TokenService';

const API_URL = 'http://localhost:8080/v1/';

class AuthService {
    async login({email, password }) {
      return axios
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

    async logout() {
      return axios
        .post(API_URL + 'logout')  
        .then(response => {
            console.log(response)
            TokenService.removeUser();
        });

    }

    register(email, password) {
        return axios
          .post(API_URL + 'signup', {
              email,
              password
          });
        }
    }

export default new AuthService();