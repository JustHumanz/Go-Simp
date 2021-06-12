import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import bootstrap from 'bootstrap'
import 'bootswatch/dist/darkly/bootstrap.min.css';
import { library } from '@fortawesome/fontawesome-svg-core'
import {faTwitter,faYoutube,faTwitch} from '@fortawesome/free-brands-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import Config from './config.json'

library.add(faTwitter,faYoutube,faTwitch)
console.log(Config.REST_API,Config.ENV)
createApp(App).use(router).use(bootstrap).component('font-awesome-icon', FontAwesomeIcon).mount('#app')