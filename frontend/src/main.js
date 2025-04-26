import { createApp } from 'vue'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';
import { createPinia } from 'pinia';
import piniaPersist from 'pinia-plugin-persist';
import mitt from "mitt";

import App from './App.vue';
import router from './router';
import useStore from '@/store/index';
import paas from 'paas-component-library';

const pinia = createPinia()

pinia.use(piniaPersist)


const app = createApp(App)
app.use(Antd);
app.use(router)
app.use(pinia)
app.use(paas)
app.config.globalProperties.$store = useStore();
app.config.globalProperties.$bus = mitt();

app.mount('#app')
