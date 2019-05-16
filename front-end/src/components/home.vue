<template>

  <body id="page-top">
  <nav class="navbar navbar-expand-lg navbar-dark fixed-top" id="mainNav">
    <div class="container">
      <a class="navbar-brand  js-scroll-trigger" href="#" v-scroll-to="'#page-top'">Alphaslot</a>
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

      window.addEventListener('scroll', scrollHandle)
    },
    created: function () {
      // this.account = this.$store.state.account
    },
    beforeCreate: function () {
    }
  }
</script>
