<template>
  <body>
  <vue-snotify></vue-snotify>
  <router-view/>
  </body>
</template>
<script>

    export default {
        name: 'App',
        data() {
            return {};
        },
        methods: {

        },
        beforeCreate: function() {
            let config = require('./etc/config.json');
            this.$store.commit('setConfig', config);

            this.$store.commit('setNotifyInfo', (message)=>{
                this.$snotify.info(message,{
                    position: "rightTop",
                    timeout: 1500
                })
            });
            this.$store.commit('setNotifyError', (message)=>{
                this.$snotify.error(message, {
                    position: "rightTop",
                    timeout: 1500
                });
            });
            this.$store.commit('setNotifySuccess',(message)=>{
                this.$snotify.success(message, {
                    position: "rightTop",
                    timeout: 1500
                });
            });
        }
    }
</script>
<style>
  /*@import './vendor/css/all.min.css';*/
  @import 'https://fonts.googleapis.com/css?family=Montserrat:400,700';
  @import 'https://fonts.googleapis.com/css?family=Kaushan+Script';
  @import 'https://fonts.googleapis.com/css?family=Droid+Serif:400,700,400italic,700italic';
  @import 'https://fonts.googleapis.com/css?family=Roboto+Slab:400,100,300,700';
  @import './assets/css/my.css';
</style>