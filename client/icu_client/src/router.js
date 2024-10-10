import { createRouter, createWebHistory } from 'vue-router'
import Home from "./components/Home.vue";
import Login from "./components/Login.vue";
import Register from "./components/Register.vue";
import store from './store';
import { jwtDecode } from "jwt-decode";
// lazy-loaded
const Profile = () => import("./components/Profile.vue")
const BoardUser = () => import("./components/BoardUser.vue")

const routes = [
  {
    path: "/",
    name: "home",
    component: Home,
    meta: { requiresAuth: false },

  },
  {
    path: "/home",
    component: Home,
  },
  {
    path: "/login",
    component: Login,
    meta: { requiresAuth: false },
  },
  {
    path: "/register",
    component: Register,
    meta: { requiresAuth: false },
  },
  {
    path: "/profile",
    name: "profile",
    // lazy-loaded
    component: Profile,
    meta: { requiresAuth: true },
  },
  {
    path: "/user",
    name: "user",
    // lazy-loaded
    component: BoardUser,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// const publicPages = ['/login', '/register', '/home', '/'];

router.beforeEach(async (to, from, next) => {
  const isAuthenticated = store.getters['auth/loggedIn'];
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!isAuthenticated) {
      next('/login');
    } else {
      const token = store.getters['auth/accessToken'];

      if (isTokenExpired(token)) {
        try {
          await store.dispatch('auth/refreshToken');
          next();
        } catch (error) {
          await store.dispatch('auth/logout')
          next('/login');
        }
      } else {
        next();
      }
    }
  } else {
    next();
  }
});

const isTokenExpired = (token) => {
  try {
    const decodedToken = jwtDecode(token);
    console.log("decodedTme" + decodedToken.exp)
    const currentTime = Date.now() / 1000; // Текущее время в секундах
    console.log("currentTime: " + currentTime);
    return decodedToken.exp < currentTime;
  } catch (error) {
    return true; // Если токен не декодируется, считаем его протухшим
  }
};

export default router