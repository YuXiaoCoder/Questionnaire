// 多选按钮: 添加索引
import { AtCheckbox } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/checkbox.scss"

export default {
  name: 'AtCheckbox',
  extends: AtCheckbox,
  props: {
    index: {
      type: Number,
      default: 0,
    },
  },
  methods: {
    handleClick(idx) {
      const { selectedList, options, index } = this
      const option = options[idx]
      const { disabled, value } = option
      if (disabled) return

      const selectedSet = new Set(selectedList)
      if (!selectedSet.has(value)) {
        selectedSet.add(value)
      } else {
        selectedSet.delete(value)
      }
      this.onChange([...selectedSet], index)
    },
  },
}
