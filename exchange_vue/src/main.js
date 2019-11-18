import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store/'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-default/index.css'

Vue.config.productionTip = false;

Vue.use(ElementUI);
const API = require('../config/api_list');
window.API = API;

new Vue({
	el: '#app',
	router,
	store,
	template: '<App/>',
	components: { App }
})
