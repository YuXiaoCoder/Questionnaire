// 按钮
import { AtButton } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/button.scss"
import "taro-ui-vue/dist/style/components/loading.scss"

export default {
  name: 'AtButton',
  extends: AtButton,
  props: {
    index: {
      type: Number,
      default: 0,
    },
  },
  methods: {
    handelOnClick(event): void {
      const { index } = this
      if (!this.disabled) {
        this.onClick && this.onClick(event, index)
      }
    },
  },
}
