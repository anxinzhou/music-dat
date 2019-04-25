import Vue from 'vue'
import App from './App'
import axios from 'axios'
import VueAxios from 'vue-axios'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import VueCookies from 'vue-cookies'
import Snotify from 'vue-snotify'
import 'vue-snotify/styles/material.css'
import router from "./router";
import store from "./store";
import {Upload,
  Button} from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

Vue.use(Upload);
Vue.use(Button);
Vue.config.productionTip = false;
Vue.use(VueCookies);
Vue.use(Snotify);
var VueScrollTo = require('vue-scrollto');
Vue.use(VueScrollTo);
Vue.use(VueAxios,axios);


new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
