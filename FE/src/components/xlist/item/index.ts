// 列表
import { AtListItem } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/list.scss"

export default {
  name: 'AtListItem',
  extends: AtListItem,
  props: {
    id: {
      type: Number,
      default: 0,
    },
  },
  methods: {
    handleClick(event) {
      const { id } = this
      if (typeof this.onClick === 'function' && !this.disabled) {
        this.onClick(id,event)
      }
    },
  },
}
