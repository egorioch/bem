<template>
  <div class="common">
    <h3>
      Profile
      <strong>{{ currentUser }}</strong>
    </h3>

    <p>
      <strong>Token:</strong>
      {{ accessToken }}
    </p>
    <p>
      <strong>Username:</strong>
      {{ currentUser.username }}
    </p>
    <p>
      <strong>Email:</strong>
      {{ currentUser.email }}
    </p>
    <p>
      <strong>Authorities:</strong>
      {{ currentUser.role }}
    </p>
  </div>
</template>

<script>
export default {
  name: 'Profile',
  computed: {
    currentUser() {
      console.log("USER: " + JSON.stringify(this.$store.getters["auth/user"]))
      return this.$store.getters["auth/user"]
    },
    accessToken() {
  
      return this.$store.getters["auth/access_token"]
    }
  },
  mounted() {
    if (!this.currentUser) {
      this.$router.push('/login');
    }
  }
};
</script>

<style lang="scss" scoped>


</style>