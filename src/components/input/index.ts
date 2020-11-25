import { AtInput } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/input.scss"

export default {
  name: 'AtInput',
  extends: AtInput,
  props: {
    index: {
      type: Number,
      default: 0,
    },
  },
  methods: {
    handleBlur(event): void {
      const { index } = this

      if (typeof this.onBlur === 'function') {
        this.onBlur(event.detail.value, index, event)
      }
      if (event.type === 'blur' && !this.inputClearing) {
        // fix # 583 AtInput 不触发 onChange 的问题
        this.onChange(event.detail.value, index)
      }
      // 还原状态
      this.inputClearing = false
    },
  },
}
