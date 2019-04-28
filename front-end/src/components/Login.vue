<template>
        <div class="signup">
            <div class="text-center" >
                <div class="signupHead">
                    <span v-if="account!==undefined">
                    <b>Welcome</b>
                    </span>
                    <span v-else>
                        <b>Login</b>
                    </span>
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
                      <div class="login-form-label">Email:</div><input class="login-form-input" v-model="email" type="text">
                    </div>
                    <div class="login-form-item">
                      <div class="login-form-label">Password:</div><input class="login-form-input" v-model="password" type="password">
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
            email: undefined,
            password: undefined,
            account: undefined
        }
    },
    methods: {
        signin: function () {
           this.$cookies.set("access-token",'todo');
           this.$router.replace('/mnemonic')
        }
    },
    mounted: function() {
        if(this.$cookies.isKey("account")) {
            this.account = this.$cookies.get("account");
            this.$store.commit('setAccount', this.account);
            console.log(this.account);
        }
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
