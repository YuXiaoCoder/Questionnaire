<template>
  <view class="index" :style="{height: screenHeight+'px'}">
    <!-- 菜单 -->
    <AtGrid v-if="openMenu" :data="items" :hasBorder="false" mode='rect' :columnNum	="1" :onClick="handleClickMenu" />

    <!-- 投票/问卷 -->
    <XList v-if="openXList" :xlistType="xlistType" :xlistTitle="xlistTitle"/>

    <!-- 跳转 -->
    <AtModal
    :isOpened="isShow"
    >
      <AtModalHeader>回答</AtModalHeader>
      <AtModalContent>
        <AtInput
          autoFocus
          required
          name="questionnaireID"
          title="题号"
          type="number"
          placeholder="请输入数字"
          :value='questionnaireID'
          :onBlur='handleChangeQuestionnaireID'
        />
      </AtModalContent>

      <AtModalAction>
        <AtButton
          type="secondary"
          :full="true"
          :onClick="handleCloseModal"
        >
          取消
        </AtButton>
        <AtButton
          type="primary"
          :full="true"
          :onClick='handleJumpPage'
        >
          跳转
        </AtButton>
      </AtModalAction>
    </AtModal>

    <!-- 消息通知 -->
    <AtMessage />
  </view>
</template>

<script>
// Taro
import Taro from '@tarojs/taro'

// 自定义组件
import XList from "../../components/xlist"

// 模态框
import { AtModal, AtModalHeader, AtModalContent, AtModalAction} from "taro-ui-vue";
import "taro-ui-vue/dist/style/components/modal.scss"

// 输入框
import { AtInput } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/input.scss"

// 按钮
import { AtButton } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/button.scss"

// 消息通知
import { AtMessage } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/message.scss"

// 栅格布局
import { AtGrid } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/grid.scss"

// 图标
import "taro-ui-vue/dist/style/components/icon.scss"

// 自定义样式
import './index.scss'

export default {
  components: {
    AtGrid,
    AtMessage,
    XList,
    AtInput,
    AtButton,
    AtModal,
    AtModalHeader,
    AtModalContent,
    AtModalAction,
  },
  data () {
    return {
      // 显示菜单
      openMenu: true,
      // 显示列表
      openXList: false,
      // 投票: 1, 问卷: 2
      xlistType: 1,
      // 标题
      xlistTitle: "我的问卷",
      // 屏幕高度
      screenHeight: 600,
      // 导航项
      items: [
        {
          image: 'https://s1.ax1x.com/2020/10/24/BVUtSO.png',
          value: '投票'
        },
        {
          image: 'https://s1.ax1x.com/2020/10/24/BVUNlD.png',
          value: '问卷'
        },
        {
          image: 'https://s1.ax1x.com/2020/10/25/BeFMCV.png',
          value: '回答'
        }
      ],
      // 模态框
      isShow: false,
      // 问卷ID
      questionnaireID: null,
    }
  },
  methods: {
    // 点击菜单
    handleClickMenu(item, index) {
      if (index == 0) {
        this.openMenu = false
        this.xlistType = 1
        this.xlistTitle = "我的投票"
        this.openXList = true
      } else if (index == 1) {
        this.openMenu = false
        this.xlistType = 2
        this.xlistTitle = "我的问卷"
        this.openXList = true
      } else if (index == 2) {
        this.isShow = true
      } else {
        // 消息通知
        Taro.atMessage({
          'message': '暂不支持',
          'type': 'error',
        });
      }
    },

    // 取消按钮
    handleCloseModal() {
      // 清理缓存
      this.questionnaireID = null
      // 关闭模态框
      this.isShow = false
    },

    // 跳转页面
    handleJumpPage() {
      // 通过正则表达式判断是否为数字
      let re = /^\d+(\.\d+)?$/
      if (re.test(this.questionnaireID)) {
        Taro.request({
          url: API_GATEWAY + '/questionnaires/' + this.questionnaireID,
          method: 'GET',
          header: {
            'content-type': 'application/json'
          },
          success: (res) => {
            if (res.statusCode == 200) {
              Taro.redirectTo({
                url: '/pages/answersheet/answersheet?id=' + this.questionnaireID,
              });
            } else {
              // 消息通知
              Taro.atMessage({
                'message': '请输入正确的题目编号',
                'type': 'error',
              });
            }
          },
        });
      } else {
       // 消息通知
        Taro.atMessage({
          'message': '请输入数字',
          'type': 'error',
        });
      }
    },

    // 修改问卷ID
    handleChangeQuestionnaireID(value) {
      this.questionnaireID = value;
    }
  },
  created() {
    // 获取屏幕高度
    this.screenHeight = Taro.getSystemInfoSync().windowHeight;
  },
  mounted() {
    // 获取Cookie
    let cookie = Taro.getStorageSync("cookie");

    if (cookie == undefined || cookie == null || cookie == "") {
      // 登录
      Taro.login({
        success: res => {
          Taro.request({
            url: API_GATEWAY + "/login",
            data: {
              code: res.code
            },
            success: res => {
              // 存储Cookie
              Taro.setStorageSync("cookie", res.header['Set-Cookie'])
              Taro.setStorageSync("user_id", res.data['user_id'])
            },
            fail: res => {
              // 消息通知
              Taro.atMessage({
                'message': '请检查后端服务',
                'type': 'error',
              });
            }
          })
        },
        fail: res => {
          // 消息通知
          Taro.atMessage({
            'message': '请检查网络',
            'type': 'error',
          });
        }
      });
    }
  }
}
</script>
