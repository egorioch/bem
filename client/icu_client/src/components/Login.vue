<template>
  <div class="common">
    <div class="icon-eyes">
      <div class="eye left">
        <div class="eye_white-group">
          <img class="eye_white" src="../assets/eyewhite.svg" alt="">
          <img class="pupil" src="../assets/pupil.svg" alt="">
          <div class="int-eye"></div>
          <!-- <img class="upper_eyelid" src="../assets/closed_eye.svg"> -->
          <!-- <img class="upper_eyelid closed_eye" src="../assets/closed_eye.svg"> -->
        </div>
      </div>

      <div class="eye right">
        <div class="eye_white-group">
          <img class="eye_white" src="../assets/eyewhite.svg" alt="">
          <img class="pupil" src="../assets/pupil.svg" alt="">
        </div>
      </div>

    </div>

    <div class="login">
      <div class="login-header">Авторизация</div>
      <div class="login-fields">
        <Form @submit="handleLogin" :validation-schema="schema" class="form">
          <div class="login-fields-row">
            <label for="username" class="email">Email</label>
            <Field name="email" type="email" class="input-form" />
            <ErrorMessage name="email" class="error-feedback" />
          </div>
          <div class="login-fields-row">
            <label for="password" class="password">Password</label>
            <Field name="password" type="password" class="input-form" />
            <ErrorMessage name="password" class="error-feedback" />
          </div>

          <div class="button-form">
            <button class="button">Login</button>
          </div>

        </Form>
      </div>
    </div>
  </div>
  <div>
    <div v-if="message" class="alert alert-danger" role="alert">
      {{ message }}
    </div>
  </div>
  
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";
import * as yup from "yup";

export default {
  name: "Login",
  components: {
    Form,
    Field,
    ErrorMessage,
  },
  data() {
    const schema = yup.object().shape({
      email: yup.string().required("Email is required!"),
      password: yup.string().required("Password is required!"),
    });

    return {
      loading: false,
      message: "",
      schema,
      closedEyes: false,
    };
  },
  computed: {
    loggedIn() {
      return this.$store.getters.loggedIn;
    },
  },
  created() {
    if (this.loggedIn) {
      this.$router.push("/profile");
    }
  },
  methods: {
    handleLogin(user) {
      this.loading = true;

      this.$store.dispatch("auth/login", user).then(
        () => {
          this.$router.push("/profile");
        },
        (error) => {
          this.loading = false;
          this.message =
            (error.response &&
              error.response.data &&
              error.response.data.message) ||
            error.message ||
            error.toString();
        }
      );
    },
    closeEyes() {

    }
  },
};
</script>

<style scoped lang="scss">
.common {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.icon-eyes {
  height: 7em;
  display: flex;
  flex-direction: row;
  justify-content: center;
  width: 200px;


  .eye {
    width: 100%;
    padding: 0 1em 0 1em;

    &.left {
      .int-eye {
        background-color: green;
        animation: move-int-eye linear infinite 1s;
      }
    }


    .eye_white-group {
      display: flex;
      justify-content: center;
      align-items: center;
      overflow: hidden;

      .pupil {
        position: absolute;
        width: inherit;
      }

      .upper_eyelid {
        position: absolute;
        margin-bottom: 2.3em;
        // rotate: 180deg;
      }

      .closed_eye {
        rotate: 180deg;
      }
    }
  }
}

.icon-eyes:hover .pupil {
  animation: rotate180 3s;
}

@keyframes rotate180 {
  0% {
    transform: rotate(0deg);
  }

  50% {
    transform: rotate(-180deg);
  }

  100% {
    transform: rotate(0deg);
  }
}

.login {
  border: 3px solid #9260ff4f;
  border-radius: 10px;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  height: 20em;
  width: 30em;

  .login-header {
    // margin-top: 5px;
    font-size: 2em;
    display: flex;
    justify-content: center;
    margin-bottom: 0.8em;
  }

  .form {
    width: 100%;
  }

  .login-fields {
    display: flex;
    width: 95%;
    // padding-left: 10px;

    .login-fields-row {
      display: flex;
      width: 100%;
      // margin-left: 10px;
      flex-direction: column;
      margin-bottom: 1.5em;
    }

    .login-fields-row:focus .pupil {
      animation: rotate180 1s;
    }

    .input-form {
      width: 100%;
    }
  }

}



.button-form {
  display: flex;
  justify-content: center;

  .button {
    width: 20%;
    display: flex;
    justify-content: center;
    border: solid 1.5px #ffffffa4;
    border-radius: 10px;
    background-color: #11de9304;
    color: #ffffffa4;

  }

  .button:hover {
    transition-duration: 0.5s;
    color: #ebebebc0;
    background-color: #9260ff4f;
  }
}
</style>