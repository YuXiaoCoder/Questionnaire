// 单选按钮: 添加索引
import { AtRadio } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/radio.scss"

export default {
  name: 'AtRadio',
  extends: AtRadio,
  props: {
    index: {
      type: Number,
      default: 0,
    },
  },
  methods: {
    handleClick(option) {
      const { index } = this
      if (option.disabled) return
      this.onClick(option.value, index)
    }
  },
}
