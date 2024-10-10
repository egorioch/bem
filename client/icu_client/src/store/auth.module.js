import AuthService from '../services/auth.service';
import createPersistedState from "vuex-persistedstate";

const localStorageUser = JSON.parse(localStorage.getItem('user'))
export const auth = {
  namespaced: true,
  plugins: [createPersistedState()],
  state: {
    user: localStorageUser,
    loggedIn: localStorageUser ? true : false,
    access_token: "",
    refresh_token: ""
  },
  // actions представляют асинхронные операции для обращения к серверу
  actions: {
    login({ commit }, user) {
      return AuthService.login(user).then(
        user => {
          commit('loginSuccess', user);
          commit('setAccessToken', user.access_token)
          commit('setRefreshToken', user.refresh_token)
          return Promise.resolve(user);
        },
        error => {
          commit('loginFailure');
          return Promise.reject(error);
        }
      );
    },
    logout({ commit }) {
      AuthService.logout();
      commit('logout');
    },
    register({ commit }, user) {
      return AuthService.register(user).then(
        response => {
          commit('registerSuccess');
          return Promise.resolve(response.data);
        },
        error => {
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    refreshToken({commit, state}) { 
      console.log("REFRESH TOKEN IN STATE: " + state.user.refresh_token + "\n, email: " + state.user.user_data.email)
      try {
        response_data = AuthService.refreshToken(state.user.user_data.email, state.user.refresh_token);
        newAccessToken = response_data.access_token;
        const updatedUser = {
          ...state.user,
          access_token: response_data.access_token,
        }
        commit("loginSuccess", updatedUser)
        return updatedUser.access_token;
      } catch(error) {
        commit("logout");
        throw error;
      }
    }
  },
  // для изменения состояния в actions(всегда синхронны)
  mutations: {
    loginSuccess(state, user) {
      state.loggedIn = true;
      state.user = user;
    },
    loginFailure(state) {
      state.loggedIn = false;
      state.user = null;
    },
    logout(state) {
      state.loggedIn = false;
      state.access_token = "";
      state.refresh_token = "";
      state.user = null;
    },
    registerSuccess(state) {
      state.loggedIn = true;
    },
    registerFailure(state) {
      state.loggedIn = false;
    },
    setRefreshToken(state, refreshToken) {
      state.refresh_token = refreshToken
    },
    setAccessToken(state, accessToken) {
      state.access_token = accessToken
    }
  },
  getters: {
    user: (state) => state.user ? state.user.user_data : null,
    accessToken: (state) => state.user.access_token,
    loggedIn: (state) => state.loggedIn,
  }
};