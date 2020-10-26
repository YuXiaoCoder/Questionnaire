import Vue from 'vue'
import './app.scss'
import Taro from '@tarojs/taro'

const App = new Vue({
  onShow(options) {
    // 清除Cookie
    setInterval(() => {
      Taro.removeStorageSync("cookie")
      Taro.removeStorageSync("user_id")
    }, 3600000)
  },
  render(h) {
    // this.$slots.default 是将要会渲染的页面
    return h('block', this.$slots.default)
  }
})

export default App
