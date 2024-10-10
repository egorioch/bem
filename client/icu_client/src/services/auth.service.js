import axios from "axios";

const API_URL = "http://localhost:8080/api/auth/";

class AuthService {
    // constructor() {
    //     this.setupInterceptors();
    // }

    // ловится баг по error.response.status
    // setupInterceptors = () => {
    //     if (localStorage.getItem('user')) {
    //         axios.interceptors.request.use(
    //             async (config) => {
    //                 const token = store.getters['auth/accessToken'];
    //                 if (token) {
    //                     config.headers['Authorization'] = 'Bearer ' + token;
    //                 }
    //                 return config;
    //             },
    //             (error) => {
    //                 return Promise.reject(error);
    //             }
    //         );

    //         axios.interceptors.response.use(
    //             (response) => response,
    //             async (error) => {
    //                 console.log("error config: " + JSON.stringify(error))
    //                 const originalRequest = error.config;
    //                 if (error.response.status === 401 && !originalRequest._retry) {
    //                     originalRequest._retry = true;
    //                     try {
    //                         const newAccessToken = await store.dispatch('auth/refreshToken');
    //                         axios.defaults.headers.common['Authorization'] = 'Bearer ' + newAccessToken;
    //                         originalRequest.headers['Authorization'] = 'Bearer ' + newAccessToken;
    //                         return axios(originalRequest);
    //                     } catch (e) {
    //                         store.dispatch('auth/logout');
    //                         this.$router.push('/login');
    //                         return Promise.reject(e);
    //                     }
    //                 }
    //                 return Promise.reject(error);
    //             }
    //         );
    //     } else {
    //         return Promise.resolve()
    //     }
    // };

    login(user) {
        let token
        return axios
            .post(API_URL + "sign_in", {
                email: user.email,
                password: user.password
            })
            .then(response => {
                if (response.data.access_token) {
                    localStorage.setItem('user', JSON.stringify(response.data));
                    axios.defaults.headers.common['Authorization'] = 'Bearer' + response.data.access_token
                }
                token = JSON.stringify(response.data)
                console.log("TOKEN: " + JSON.stringify(response.data, null, 2))

                console.log("RESPONSE DATA: " + response.data)
                return response.data;
            });
    }

    logout() {
        localStorage.removeItem('user');
    }

    register(user) {
        return axios.post(API_URL + "sign_up", {
            email: user.email,
            username: user.username,
            password: user.password
        });
    }

    refreshToken(email, refreshToken) {
        return axios.post(API_URL + 'refresh', {
            email: email,
            refresh_token: refreshToken,
        }).then(response => {
            console.log("response new access: " + response.data)
            if (response.data.access_token) {
                let user = JSON.parse(localStorage.getItem('user'));
                user.access_token = response.data.access_token;
                localStorage.setItem('user', JSON.stringify(user));
            }
            return response.data;
        });
    }
}

export default new AuthService();
