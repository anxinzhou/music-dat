<template>

  <body id="page-top">
  <nav class="navbar navbar-expand-lg navbar-dark fixed-top navbar-shrink" id="mainNav">
    <div class="container">
      <a class="navbar-brand  js-scroll-trigger" href="#" v-scroll-to="'#page-top'">Music Hotpot</a>
      <p id="user-addr" class="js-scroll-trigger user-addr" href="#" v-scroll-to="'#page-top'"></p>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarResponsive"
              aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation">
        Menu
        <i class="fas fa-bars"></i>
      </button>
      <div class="collapse navbar-collapse" id="navbarResponsive">
        <ul class="navbar-nav  ml-auto">
          <li class="nav-item text-uppercase">
            <router-link to="/"><a class="nav-link">Home</a></router-link>
          </li>
          <li class="nav-item text-uppercase">
            <img style="width: 50px;height: 50px;" :src="avatarUrl"/>
          </li>
          <li class="nav-item text-uppercase">
            <el-button round style="margin-left:20px;margin-top:2px;height: 45px; width: 75px;" @click="logout"><span style="margin-left:-10px;">Logout</span></el-button>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <header class="masthead">
    <div class="container">
      <div class="intro-text">
        <div class="intro-lead-in"></div>
      </div>
    </div>
  </header>
  <section class="bg-white">
    <el-row style="margin-bottom: 100px;">
      <!--          <img src="../assets/images/avatar.jpg">-->
      <el-col :span="12" :offset="6">
        <el-row>
          <el-col :span="6">
            <b>Avatar: </b>
          </el-col>
          <div v-if="!avatarEditing">
            <el-col :span="6" class="text-center">
              <img style="margin-top:-20px;width: 80px; height: 80px;" :src="avatarUrl">
            </el-col>
            <el-col :span="6" class="text-center">
              <button class="alpha-button" @click="editAvatar">Edit</button>
            </el-col>
          </div>
          <div v-else>
            <el-col :span="6" class="text-center">
              <el-upload
                class="avatar-uploader"
                :show-file-list="false"
                :auto-upload="false"
                action=""
                :on-change="previewAvatar">
                <img v-if="imageUrl" :src="imageUrl" class="avatar" style="width: 80px;height: 80px; margin-top:-20px;">
                <el-row v-else class="el-upload" type="flex" align="middle">
                  <el-col style="height: 100%">
                    <span class="el-icon-plus avatar-uploader-icon" style="margin-top: 42%"></span>
                  </el-col>
                </el-row>
              </el-upload>
            </el-col>
            <el-col :span="3" class="text-center">
              <button class="alpha-button" @click="saveAvatar">Save</button>
            </el-col>
            <el-col :span="2" class="text-center">
              <button class="alpha-button" @click="cancelEditAvatar">Cancel</button>
            </el-col>
          </div>
        </el-row>
      </el-col>
    </el-row>
    <el-row style="margin-bottom: 100px;">
      <!--          <img src="../assets/images/avatar.jpg">-->
      <el-col :span="12" :offset="6">
        <el-row>
          <el-col :span="6">
            <b>Nickname: </b>
          </el-col>
          <div v-if="!nicknameEditing">
            <el-col :span="6" class="text-center">
              <span>{{nickname}}</span>
            </el-col>
            <el-col :span="6" class="text-center">
              <button class="alpha-button" @click="editNickname">Edit</button>
            </el-col>
          </div>
          <div v-else>
            <el-col :span="6" class="text-center">
              <el-input v-model="newNickname" placeholder="nickname"></el-input>
            </el-col>
            <el-col :span="3" class="text-center">
              <button class="alpha-button" @click="saveNickname">Save</button>
            </el-col>
            <el-col :span="2" class="text-center">
              <button class="alpha-button" @click="cancelEditNickname">Cancel</button>
            </el-col>
          </div>
        </el-row>
      </el-col>
    </el-row>
    <el-row style="margin-bottom: 100px;">
      <el-col :span="12" :offset="6">
        <el-row>
          <el-col :span="6">
            <b>Wallet Address: </b>
          </el-col>
          <div v-if="!walletEditing">
            <el-col :span="12" class="text-center">
              {{address}}
            </el-col>
            <el-col :span="6" class="text-center">
              <button class="alpha-button" @click="editWallet">Edit</button>
            </el-col>
          </div>
          <div v-else>
            <el-col :span="6" class="text-center">
              <el-input v-model="newAddress" placeholder="address"></el-input>
            </el-col>
            <el-col :span="3" class="text-center">
              <button class="alpha-button" @click="saveWallet">Save</button>
            </el-col>
            <el-col :span="2" class="text-center">
              <button class="alpha-button" @click="cancelEditWallet">Cancel</button>
            </el-col>
          </div>
        </el-row>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="12" :offset="6">
        <el-row>
          <el-col :span="6">
            <b>Intro: </b>
          </el-col>
          <div v-if="!introEditing">
          <el-col :span="6" class="text-center">
            {{intro}}
          </el-col>
          <el-col :span="6" class="text-center">
            <button class="alpha-button" @click="editIntro">Edit</button>
          </el-col>
          </div>
          <div v-else>
            <el-col :span="6" class="text-center">
            <el-input
              type="textarea"
              :rows="2"
              placeholder="Intro"
              v-model="newIntro">
            </el-input>
            </el-col>
            <el-col :span="3" class="text-center">
              <button class="alpha-button" @click="saveIntro">Save</button>
            </el-col>
            <el-col :span="2" class="text-center">
              <button class="alpha-button" @click="cancelEditIntro">Cancel</button>
            </el-col>
          </div>
        </el-row>
      </el-col>
    </el-row>
    <!--    <el-row>-->
    <!--      <el-col :span="8" :offset="6">-->
    <!--        <el-row>-->
    <!--          <el-col :span="12">-->
    <!--&lt;!&ndash;            <button :>Edit: </button>&ndash;&gt;-->
    <!--          </el-col>-->
    <!--          <el-col :span="8" class="text-center" :offset="2">-->
    <!--            3dsadsadasdjfdkslfjsdklfjsdkf-->
    <!--            fdsfsdfdsadasdas-->
    <!--          </el-col>-->
    <!--        </el-row>-->
    <!--      </el-col>-->
    <!--    </el-row>-->
  </section>
  </body>
</template>
<script>
  export default {
    data() {
      return {
        avatarUrl: '',
        nickname: '',
        balance: '',
        nicknameEditing: false,
        avatarEditing: false,
        introEditing: false,
        walletEditing: false,
        newAddress: undefined,
        newNickname: '',
        newIntro: '',
        newAvatar: undefined,
        imageUrl: undefined,
        intro: '',
        address: undefined,
      }
    },
    methods: {
      logout: function () {
        console.log("logout")
        this.$cookies.remove("access-token");
        this.$cookies.remove("uuid");
        this.$router.replace('/login');
      },
      editNickname: function () {
        console.log("edit nickname");
        this.newNickname = this.nickname
        this.nicknameEditing = true;
      },
      saveNickname: function() {
        if(this.newNickname === this.nickname) {
          this.$store.state.notifySuccess("set nickname success")
          this.nicknameEditing = false;
          return
        }
        console.log("save nickname");
        let uuid = this.$cookies.get('uuid');
        let httpPath = this.$store.state.config.httpPath;
        this.axios.post(`${httpPath}/profile/${uuid}/nickname`,{
          nickname: this.newNickname,
        }).then(res=>{
          this.nickname = this.newNickname;
          this.nicknameEditing = false;
          this.$store.state.notifySuccess("set nickname success")
        }).catch(err=>{
          console.log(err.response.data.reason);
          this.$store.state.notifyError(err.response.data.reason)
        });
      },
      cancelEditNickname: function() {
        this.nicknameEditing = false;
      },
      saveIntro: function() {
        if(this.newIntro === this.intro) {
          this.introEditing = false;
          this.$store.state.notifySuccess("set intro success")
          return
        }
        console.log("save intro");
        let uuid = this.$cookies.get('uuid');
        let httpPath = this.$store.state.config.httpPath;
        this.axios.post(`${httpPath}/profile/${uuid}/intro`,{
          intro: this.newIntro,
        }).then(res=>{
          this.intro = this.newIntro;
          this.introEditing = false;
          this.$store.state.notifySuccess("set intro success")
        }).catch(err=>{
            console.log(err.response.data.reason);
          this.$store.state.notifyError(err.response.data.reason)
        });
      },
      cancelEditIntro: function() {
        this.introEditing = false;
      },
      editIntro: function() {
          this.newIntro = this.intro;
          this.introEditing = true;
      },
      saveWallet: function() {
        if(this.newAddress === this.address) {
          this.walletEditing = false;
          this.$store.state.notifySuccess("set wallet success")
          return
        }
        console.log("save wallet");
        this.newAddress = this.newAddress.replace(' ','');
        let httpPath = this.$store.state.config.httpPath;
        let uuid = this.$cookies.get('uuid');
        this.axios.post(`${httpPath}/profile/${uuid}/wallet`,{
          wallet: this.newAddress,
        }).then(res=>{
          this.address = this.newAddress;
          this.walletEditing = false;
          this.$store.state.notifySuccess("set wallet success")
        }).catch(err=>{
          console.log(err.response.data.reason);
          this.$store.state.notifyError(err.response.data.reason)
        });
      },
      editWallet: function() {
        this.newAddress = this.address;
          this.walletEditing = true;
      },
      cancelEditWallet: function() {
          this.walletEditing = false;
      },
      saveAvatar: function() {
        if(this.newAvatar===undefined) {
          this.$store.state.notifyError("Please choose a picture")
          return
        }
        console.log("save avatar");
        let httpPath = this.$store.state.config.httpPath;
        let nickname = this.nickname;
        let formData = new FormData();
        formData.append("avatar",this.newAvatar);
        let uuid = this.$cookies.get('uuid');
        this.axios.post(`${httpPath}/profile/${uuid}/avatar`,formData,{
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }).then(res=>{
            this.avatarEditing = false;
            let url = res.data.avatarUrl;
            this.avatarUrl = undefined;
            this.avatarUrl  = url;
            this.newAvatar = undefined;
            this.imageUrl = undefined;
          this.$store.state.notifySuccess("set avatar success")
        }).catch(err=>{
          this.$store.state.notifyError(err.response.data.reason)
          console.log(err.response.data.reason);
        });
      },
      cancelEditAvatar: function() {
        this.avatarEditing = false;
      },
      editAvatar: function() {
          this.avatarEditing  = true;
      },
      previewAvatar: function (file, fileList) {
        if (file !== undefined) {
          this.imageUrl = URL.createObjectURL(file.raw);
          this.newAvatar = file.raw
        } else {
          this.imageUrl = undefined;
          this.newAvatar = undefined
        }
        console.log(file.raw);
      },
      getIntro: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/intro`).then(res=>{
          this.intro = res.data.intro;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getNickname: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/nickname`).then(res=>{
          this.nickname = res.data.nickname;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getWallet: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/wallet`).then(res=>{
          this.address = res.data.wallet;
        }).catch(err=>{
          console.log(err.response.data.reason)
        });
      },
      getAvatarUrl: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/avatar`).then(res=>{
          this.avatarUrl = res.data.avatarUrl;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
    },
    mounted: function () {
      // var el = document.getElementById('mainNav')
      // el.classList.add('navbar-shrink');
      // function scrollHandle () {
      //   if (document.documentElement.scrollTop > 100) {
      //     el.classList.add('navbar-shrink')
      //   } else {
      //     el.classList.remove('navbar-shrink')
      //   }
      // }

      // window.addEventListener('scroll', scrollHandle)
    },
    created: function () {
      let uuid = this.$cookies.get('uuid');
      this.getIntro(uuid);
      this.getNickname(uuid);
      this.getWallet(uuid);
      this.getAvatarUrl(uuid);
    },
    beforeCreate: function () {
    }
  }
</script>
<style scoped>
  .avatar-uploader .el-upload {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    width: 150px;
    height: 150px;
    justify-content: center;
  }

  .avatar-uploader .el-upload:hover {
    border-color: #409EFF;
  }
</style>
