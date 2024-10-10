import axios from "axios";

const API_URL = "http://localhost:8080/api/auth/";

class AuthService {
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
        console.log("user fields: " + JSON.stringify(user))
        return axios.post(API_URL + "sign_up", {
            email: user.email,
            username: user.username,
            password: user.password,
            admin_token: user.admin_token
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
