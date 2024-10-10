import authHeader from "./auth-header.js";
import axios from "axios";
const API_URL = "http://localhost:8080/api/";

class UserService {
  getPublicContent() {
    return axios.get(API_URL);
  }

  getUserBoard() {
    return axios.get(API_URL + 'user', { headers: authHeader() });
  }
}

export default new UserService();