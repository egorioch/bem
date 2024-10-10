import axios from 'axios';
import store from '../store'

const setupInterceptors = () => {
  if (store.getters['auth/loggedIn']) {
    axios.interceptors.request.use(
      async (config) => {
        const token = store.getters['auth/accessToken'];
        if (token) {
          config.headers['Authorization'] = 'Bearer ' + token;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // ловится баг по error.response.status
    axios.interceptors.response.use(
      (response) => response,
      async (error) => {
        console.log("error config: " + JSON.stringify(error))
        const originalRequest = error.config;
        if (error.response.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          try {
            const newAccessToken = await store.dispatch('auth/refreshToken');
            axios.defaults.headers.common['Authorization'] = 'Bearer ' + newAccessToken;
            originalRequest.headers['Authorization'] = 'Bearer ' + newAccessToken;
            return axios(originalRequest);
          } catch (e) {
            store.dispatch('auth/logout');
            this.$router.push('/login');
            return Promise.reject(e);
          }
        }
        return Promise.reject(error);
      }
    );
  } else {
    return Promise.resolve()
  }
};

export default setupInterceptors;