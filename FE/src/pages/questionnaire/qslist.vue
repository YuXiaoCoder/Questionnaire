<template>
  <view class="qs-list at-row" >
  <view class="delect-icon at-col at-col-1 at-col--auto"
  @tap="handleClick"
  @click="handleClick"
  >
  </view>
  <view class="qs-input at-col at-col-10">
    <slot></slot>
  </view>
  </view>
</template>

<script>
import "taro-ui-vue/dist/style/components/flex.scss";

export default {
  name: 'QSList',
  props:{
    onClick: {
      type: Function,
      default: function () {
        return () => {}
      },
    },
    index: {
      type: Number,
      default:0
    },
    qsItem: {
      type: Object,
      default:() => ({}),
    }
  },
  methods: {
    handleClick(event) {
      const { index, qsItem } = this;
      this.onClick && this.onClick(index, qsItem, event)
    }
  }
}

</script>

<style lang='scss'>
.qs-list{
  margin-bottom:6rpx;
  //border-bottom: 2rpx solid rgba(135, 170, 184, 0.6);
  padding: 11rpx 10rpx 15rpx 0;
}
.delect-icon {
  font-family: iconfont;
  color: red;
  font-size: 48rpx;
  margin-left: 10rpx;
  text-align: left;
  line-height: 66rpx;
  &::before{
    content: "\e961";
  }
}
.qs-input {
  padding: 10rpx 10rpx;
}
</style>
