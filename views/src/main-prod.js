import Vue from 'vue'

import Cookies from 'js-cookie'

import 'normalize.css/normalize.css' // a modern alternative to CSS resets

import '@/styles/index.scss' // global css

// 复制粘贴
import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

import './utils/flexible'
// excel导出
import JsonExcel from 'vue-json-excel'
Vue.component('downloadExcel', JsonExcel)

import App from './App'
import store from './store'
import router from './router'

import './permission' // permission control
import './utils/error-log' // error log
import './icons' // error log

// 全局注册组件
import Editor from 'vue2-ace-editor'
Vue.component('editor', Editor)

import * as filters from './filters' // global filters

Vue.use(ELEMENT, {
  size: Cookies.get('size') || 'medium' // set element-ui default size
})

// 引入form-create 表单生成器
import formCreate from '@form-create/element-ui'

Vue.use(formCreate)

// 注册全局过滤器
Object.keys(filters).forEach(key => {
  Vue.filter(key, filters[key])
})

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
