<template>
  <div id="app">
    <Navbar />

    <div class="container">
      <router-view />
    </div>

    <footer class="footer">
      <img src="./assets/rect_footer.svg" class="footer-svg">
    </footer>
  </div>
</template>

<script>
import Navbar from './components/Navbar.vue';
export default {
  components: {
    Navbar
  },
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
   
  },
  methods: {
    logOut() {
      this.$store.dispatch("auth/logout");
      this.$router.push("/login")
    },
  }
}
</script>


<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@300;500;700&display=swap');

body,
html {
  width: 100vw;
  height: 100vh;
}

#app {
  background-color: #232323;
  width: 100vw;
  height: 100vh;
  font-family: 'Montserrat', sans-serif;
}

.container {
  display: flex;
  justify-content: center;
  margin-top: 3em;
  color: rgb(217, 215, 215);
}

.footer {
  position: fixed;
  bottom: 0;
  width: 100%;
  text-align: center;
}

.footer:after {
  content: "";
  display: block;
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  background: transparent;
  background-image: linear-gradient(to top, #11de9304, #232323);
  opacity: 1;
}

</style>