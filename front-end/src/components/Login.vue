<script src="../router.js"></script>
<template>
  <div class="signup">
    <div class="text-center">
      <div class="signupHead">

                        <b>Login</b>
      </div>
      <img class="signupImg" src="@/assets/images/smile.png">
    </div>
    <!--<div class="createAccount text-center">-->
    <!--Your account: {{account.address}}-->
    <!--<br>-->
    <!--<br>-->
    <!--<button class="alpha-button alpha-button-primary" type="button" @click="login">Go Ahead !</button>-->
    <!--</div>-->
    <div>
      <form class="login-form">
        <div class="login-form-item">
          <div class="login-form-label">Username:</div>
          <input class="login-form-input" v-model="username" type="text">
        </div>
        <div class="login-form-item">
          <div class="login-form-label">Password:</div>
          <input class="login-form-input" v-model="password" type="password">
        </div>
        <div class="login-form-button">
          <button type="button" class="alpha-button alpha-button-primary" @click="signin">Sign in</button>
        </div>
      </form>
    </div>
  </div>
</template>
<script>
  export default {
    data() {
      return {
        username: undefined,
        password: undefined,
      }
    },
    methods: {
      signin: function () {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.post(`${httpPath}/admin`,{
          username: this.username,
          password: this.password,
        }).then(res => {
          this.$cookies.set("uuid",res.data.uuid);
          this.$cookies.set("access-token", res.data.accessToken);
          this.$router.replace('/');
        }).catch(err=>{
          if (err.response.status === 401) {
            this.$store.state.notifyError("wrong username or password");
          } else {
            this.$store.state.notifyError("server internal error");
          }
        });
      }
    },
    mounted: function () {
    },
  }
</script>
<style>
  .signin {
    margin-top: 20px;
  }

  .formcenter {
    padding-top: 100px;
  }

  .formpos {
    margin-left: -50px;
  }

  .signup {
    margin-top: 40px;
  }
</style>}
