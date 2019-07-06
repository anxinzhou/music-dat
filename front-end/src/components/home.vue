<script src="../store.js"></script>
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
            <a class="nav-link" href="#" v-scroll-to="'#page-top'">Home</a>
          </li>
          <li class="nav-item text-uppercase">
            <a class="nav-link" href="#" v-scroll-to="'#upload'">Market</a>
          </li>
          <li class="nav-item text-uppercase">
            <router-link :to="{name:'Profile'}"><a><img style="width: 50px;height: 50px;" :src="avatarUrl"/></a></router-link>
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
  <div>
    <a id="upload">
      <vue-upload>

      </vue-upload>
    </a>
  </div>
  </body>
</template>
<script>
  export default {
    components: {
      'vue-upload': () => import('@/components/upload.vue'),
    },
    data() {
      return {
        avatarUrl: undefined,
      }
    },
    methods: {
      logout: function () {
        console.log("logout")
        this.$cookies.remove("access-token");
        this.$cookies.remove("uuid");
        this.$router.replace('/login');
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
      let uuid = this.$cookies.get("uuid");
      this.getAvatarUrl(uuid);
    },
    beforeCreate: function () {
    }
  }
</script>
