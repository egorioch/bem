<template>
  <nav class="nav">
    <div class="nav-common">
      <div class="nav-item home">
        <router-link to="/"><img src="../assets/С.svg" class="icu-image"></router-link>
      </div>

      <div class="nav-auth">
        <div class="nav-item"><router-link v-if="!this.$store.getters['auth/loggedIn']" to="/login">Логин</router-link></div>
        <div class="nav-item"><router-link v-if="!this.$store.getters['auth/loggedIn']" to="/register">Регистрация</router-link></div>
        <div class="nav-item"><router-link v-if="this.$store.getters['auth/loggedIn']" to="/profile">Профиль</router-link></div>
        <div><button v-if="this.$store.getters['auth/loggedIn']" @click="logout">Выход</button></div>
      </div>

    </div>
  </nav>
</template>

<script>
export default {
  data() {
    return {
      loggedIn: localStorage.getItem('user')
    }
  },
  props: {
    userApp: {

    } 
  },
  methods: {
    logout() {
      localStorage.removeItem('user')
      this.loggedIn = false
      this.$router.push('/login')
      this.$store.dispatch("auth/logout")
    }
  }
}
</script>

<style scoped lang="scss">
.nav {
  width: 100%;
  background-color: #2F3133;
  /* Черный цвет для фона */
  padding: 0.6rem;
  display: flex;
}

.nav ul {
  list-style: none;
  display: flex;
  width: 100%;

}

.nav-common {
  display: flex;
  width: 100%;
  align-items: center;

  .nav-auth {
    display: flex;
    margin-left: auto;
    margin-right: 2em;
  }
}

.nav-item a {
  text-decoration: none;
  color: #d061ffd0;

  font-size: 1.43em;
  padding: 0 6px 0 6px;

}

.home {
  padding-left: 0.8em;
  color: #DDE0E7;


  .icu-image {
    width: 1.7em;
  }
}



.nav-item a:hover {
  transition: 0.5s;
  color: #9260FF;

  /* Темно-фиолетовый цвет при наведении */
}

.container {
  margin-top: 2rem;
  padding: 1rem;
}
</style>