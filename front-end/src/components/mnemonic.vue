<template>
    <div class="signup">
        <div class="text-center" >
            <div class="signupHead">
                    <span v-if="account!==undefined">
                    <b>Welcome</b>
                    </span>
                <span v-else>
                        <b>Import wallet</b>
                    </span>
            </div>
            <img class="signupImg" src="@/assets/images/smile.png">
        </div>
        <div>
            <form class="login-form" v-if="account===undefined">
                <div class="login-form-item">
                    <div class="login-form-label">Mnemonic:</div><input class="login-form-input" v-model="mnemonic" type="text">
                </div>
                <div class="login-form-button">
                    <button type="button" class="alpha-button alpha-button-primary" @click="importWallet">Confirm</button>
                </div>
            </form>
            <div class="mnemonic-info" v-else>
                <div>
                    <b>Your address:</b> {{account.address}}
                </div>
                <div class="login-form-button mnemonic-info-button">
                    <button type="button" class="alpha-button alpha-button-primary" @click="goAhead">Go ahead</button>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
    import {createWallet} from "../assets/js/wallet";
    export default {
        data() {
            return {
                mnemonic: '',
                account: undefined
            }
        },
        methods: {
            importWallet: function() {
                // console.log(this.mnemonic)
                let account = createWallet(this.mnemonic);
                this.account = account;
                this.$store.commit('setAccount',account);
                this.$cookies.set('account',account);
            },
            goAhead: function() {
                this.$router.replace('/')
            },
        },
        mounted: function() {
          // if(this.$cookies.isKey("account")) {
          //   this.account = this.$cookies.get("account");
          //   this.$store.commit('setAccount', this.account);
          //   console.log(this.account);
          // }
        },
    }
</script>
<style scoped>
    .mnemonic-info {
        margin-top: 50px;
        display: flex;
        display: -webkit-flex;
        flex-direction: column;
        align-items: center;
    }

    .mnemonic-info-button {
        margin-top: 50px;
    }

    .signup {
        margin-top: 40px;
    }
</style>}
