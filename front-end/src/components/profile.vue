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
<!--            <el-col :span="6" class="text-center">-->
<!--              <button class="alpha-button" @click="editNickname">Edit</button>-->
<!--            </el-col>-->
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
        address: '',
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
      }
    },
    methods: {
      logout: function () {
        console.log("logout")
        this.$cookies.remove("avatarUrl");
        this.$cookies.remove("nickname");
        this.$cookies.remove("access-token");
        this.$cookies.remove("address");
        this.$cookies.remove("account");
        this.$router.replace('/login');
      },
      editNickname: function () {
        console.log("edit nickname");
        this.nicknameEditing = true;
      },
      saveNickname: function() {
        console.log("save nickname");
        this.nicknameEditing = false;
      },
      cancelEditNickname: function() {
        this.nicknameEditing = false;
      },
      saveIntro: function() {
        console.log("save intro");
        let httpPath = this.$store.state.config.httpPath;
        let nickname = this.nickname;
        this.axios.post(`${httpPath}/profile/${nickname}/intro`,{
          intro: this.newIntro,
        }).then(res=>{
          this.intro = this.newIntro;
          this.introEditing = false;
        }).catch(err=>{
            this.$store.state.notifyError(err.response.data.reason);
            console.log(err.response.data.reason);
        });
      },
      cancelEditIntro: function() {
        this.introEditing = false;
      },
      editIntro: function() {
          this.introEditing = true;
      },
      saveWallet: function() {
        console.log("save wallet");
        this.newAddress = this.newAddress.replace(' ','');
        let httpPath = this.$store.state.config.httpPath;
        let nickname = this.nickname;
        this.axios.post(`${httpPath}/profile/${nickname}/wallet`,{
          address: this.newAddress,
        }).then(res=>{
          this.address = this.newAddress;
          this.$cookies.set('address',this.address);
          this.walletEditing = false;
        }).catch(err=>{
          this.$store.state.notifyError(err.response.data.reason);
          console.log(err.response.data.reason);
        });
      },
      editWallet: function() {
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
        this.axios.post(`${httpPath}/profile/${nickname}/avatar`,formData,{
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }).then(res=>{
            this.avatarEditing = false;
            let url = res.data.avatarUrl;
            this.$cookies.set('avatarUrl',url);
            this.avatarUrl = undefined;
            this.avatarUrl  = url;
            this.newAvatar = undefined;
            this.imageUrl = undefined;
        }).catch(err=>{
          this.$store.state.notifyError(err.response.data.reason);
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
      this.avatarUrl = this.$cookies.get('avatarUrl');
      this.nickname = this.$cookies.get('nickname');
      this.address = this.$cookies.get('address');
      let httpPath = this.$store.state.config.httpPath;
      let nickname = this.nickname;
      this.axios.get(`${httpPath}/profile/${nickname}/intro`).then(res=>{
        this.intro = res.data.intro;
      }).catch(console.log);
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
