<template>

  <body id="page-top">
  <nav class="navbar navbar-expand-lg navbar-dark fixed-top" id="mainNav">
    <div class="container">
      <a class="navbar-brand  js-scroll-trigger" href="#" v-scroll-to="'#page-top'">Alphaslot</a>
      <p id="user-addr" class="js-scroll-trigger user-addr" href="#" v-scroll-to="'#page-top'">This is addr</p>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarResponsive"
              aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation">
        Menu
        <i class="fas fa-bars"></i>
      </button>
      <div class="collapse navbar-collapse" id="navbarResponsive">
        <ul class="navbar-nav  ml-auto">
          <li class="nav-item text-uppercase">
            <a class="nav-link" href="#" v-scroll-to="'#dat'">Dat</a>
          </li>
          <li class="nav-item text-uppercase">
            <a class="nav-link" href="#" v-scroll-to="'#avatar'">Avatar</a>
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
    <a id="dat">
      <vue-dat>

      </vue-dat>
    </a>
    <a id="avatar">
      <vue-avatar>

      </vue-avatar>
    </a>
  </div>
  </body>
</template>
<script>
  export default {
    components: {
      'vue-dat': () => import('@/components/dat.vue'),
      'vue-avatar': () => import('@/components/avatar.vue')
    },
    data () {
      return {
      }
    },
    methods: {
    },
    mounted: function () {
      var el = document.getElementById('mainNav')

      function scrollHandle () {
        if (document.documentElement.scrollTop > 100) {
          el.classList.add('navbar-shrink')
        } else {
          el.classList.remove('navbar-shrink')
        }
      }

      // set user address from cookie
      let userAddr = document.getElementById('user-addr')
      if (this.$cookies.isKey('account')) {
        userAddr.innerText = 'Address: '  + this.$cookies.get('account').address
      } else {
        userAddr.innerText = 'Not Logged in'
      }

      window.addEventListener('scroll', scrollHandle)
    },
    created: function () {
      // this.account = this.$store.state.account
    },
    beforeCreate: function () {
    }
  }
</script>
