<template>
  <div class="signup">
    <div class="text-center">
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
          <div class="login-form-label">Address:</div>
          <input class="login-form-input" v-model="mnemonic" type="text">
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
      importWallet: function () {
        // console.log(this.mnemonic)
        let httpPath =  this.$store.state.config.httpPath;
        let address = this.mnemonic.replace(' ','');
        // this.mnemonic = this.mnemonic.replace(/\s+/g,' ').replace(/^\s+|\s+$/,'');
        // let account = createWallet(this.mnemonic);
        let nickname = this.$cookies.get('nickname');
        console.log("address",address,nickname);
        this.axios.post(`${httpPath}/profile/${nickname}/wallet`,{
          address: address,
        }).then(res=>{
          this.$cookies.set('address',address);
          this.$router.replace('/');
        }).catch(err=>{
          this.$store.state.notifyError(err.response.data.reason);
          console.log(err.response.data.reason);
        });
      },
      goAhead: function () {
        this.$router.replace('/')
      },
    },
    mounted: function () {
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
