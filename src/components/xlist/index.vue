<template>
  <view class="xlist" :style="{height: scrollBarHeight+'px'}">
    <!-- 问卷列表 -->
    <AtList>
      <AtListItem
        :title='itme.title'
        :iconInfo='questionnaireIcon'
        v-for="(itme, index) in questionnaires"
        :key="index"
        :onClick="handleClickQSItme"
        :id='itme.id'
      />
    </AtList>

    <!-- 导航栏 -->
    <AtTabBar
      :tabList='tarbar.list'
      :fixed='tarbar.fixed'
      :current="tarbar.current"
      :onClick="handleClickTabBar"
    />

    <!-- 动作面板 -->
    <AtActionSheet :isOpened='actionPanelFlag' cancelText='取消' :onCancel='handleCloseActionPanel' :onClose='handleCloseActionPanel'>
      <AtActionSheetItem
      :onClick="handleEditQuestionnaire"
      >
        编辑
      </AtActionSheetItem>
      <AtActionSheetItem
      :onClick="handleShareQuestionnaire"
      >
        分享
      </AtActionSheetItem>
      <AtActionSheetItem
      :onClick="handleAnalyzeQuestionnaire"
      >
        统计
      </AtActionSheetItem>
      <AtActionSheetItem
      :onClick="handleDeleteQuestionnaire"
      >
        删除
      </AtActionSheetItem>
    </AtActionSheet>

    <!-- 消息通知 -->
    <AtMessage />
  </view>
</template>

<script>
// Taro
import Taro from '@tarojs/taro'

// 列表
import { AtList } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/list.scss"
import AtListItem from "./item"

// 消息通知
import { AtMessage } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/message.scss"

// 动作面板
import { AtActionSheet, AtActionSheetItem } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/action-sheet.scss"

// 底部导航栏
import { AtTabBar } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/tab-bar.scss"
import "taro-ui-vue/dist/style/components/badge.scss"

// 图标
import "taro-ui-vue/dist/style/components/icon.scss"

// 自定义样式
import './index.scss'

export default {
  name: 'Questionnaire',
  props: {
    xlistType: {
      type: Number,
      default: 1,
    },
    xlistTitle: {
      type: String,
      default: "列表页面",
    }
  },
  components: {
    AtTabBar,
    AtList,
    AtListItem,
    AtActionSheet,
    AtActionSheetItem,
    AtMessage,
  },
  data () {
    return {
      // 屏幕高度
      screenHeight: 600,
      // 滚动条高度
      scrollBarHeight: 540,
      // 问卷图标
      questionnaireIcon: {
        size: '16',
        color: '#33A6B8',
        value: 'file-generic',
      },
      // 导航栏
      tarbar: {
        list: [
          { title: '我的', iconType: 'user' },
          { title: '添加', iconType: 'add' }
        ],
        // 固定到底部
        fixed: true,
        current: 0,
      },
      // 问卷列表
      questionnaires: [],
      // 问卷ID
      questionnaireID: null,
      // 动作面板
      actionPanelFlag: false,
    }
  },
  methods: {
    // 导航栏
    handleClickTabBar(value){
      if (this.tarbar.current == value) {
        return
      } else {
        this.tarbar.current = value
      }

      if (this.tarbar.current === 1) {
        Taro.redirectTo({
          url: '/pages/questionnaire/questionnaire?type=' + this.xlistType
        })
      }
    },

    // 动作面板
    handleCloseActionPanel() {
      this.actionPanelFlag = false
    },

    // 点击指定的问卷
    handleClickQSItme(id){
      this.questionnaireID = id;
      this.actionPanelFlag = true
    },

    // 编辑问卷
    handleEditQuestionnaire() {
      Taro.redirectTo({
        url: '/pages/questionnaire/questionnaire?id=' + this.questionnaireID + '&type=' + this.xlistType,
      });
    },

    // 分享问卷
    handleShareQuestionnaire() {
      Taro.setClipboardData({
        data: this.questionnaireID.toString(),
        success: (res) => {}
      });
      // 关闭动作面板
      this.actionPanelFlag = false
    },

    // 分析问卷
    handleAnalyzeQuestionnaire() {
      Taro.redirectTo({
        url: '/pages/analysis/analysis?id=' + this.questionnaireID,
      });
    },

    // 删除问卷
    handleDeleteQuestionnaire() {
      Taro.request({
        url: API_GATEWAY + '/questionnaires/' + this.questionnaireID,
        method: 'DELETE',
        success: (res) => {
          // 获取问卷列表
          Taro.request({
            url: API_GATEWAY + "/questionnaires?type=" + this.xlistType,
            header: {
              'content-type': 'application/json'
            },
            methods: "GET",
            success: (res) => {
              this.questionnaires = res.data.questionnaires;
            }
          });
          // 消息通知
          Taro.atMessage({
            'message': '删除成功',
            'type': 'success',
          });
        },
        fail: (res) => {
          // 消息通知
          Taro.atMessage({
            'message': '删除失败',
            'type': 'error',
          });
        }
      });
      // 关闭动作面板
      this.actionPanelFlag = false
    }
  },
  created(){
    // 获取屏幕高度
    this.screenHeight = Taro.getSystemInfoSync().windowHeight;

    // 获取Cookie
    let cookie = Taro.getStorageSync("cookie");

    if (cookie == undefined || cookie == null || cookie == "") {
      // 清除所有缓存
      Taro.clearStorageSync();
      // 跳转到首页
      Taro.redirectTo({
        url: "/pages/index/index",
      });
    }
  },
  mounted(){
    // ES6: 箭头函数
    // 获取底部导航栏高度
    setTimeout(() => {
      Taro.createSelectorQuery().select('.at-tab-bar').boundingClientRect(rect=>{
        this.scrollBarHeight = this.screenHeight - Math.floor(rect.height) - 1;
      }).exec()
    }, 300)

    // 设置页面标题
    Taro.setNavigationBarTitle({
      title: this.xlistTitle
    });

    // 获取问卷列表
    Taro.request({
      url: API_GATEWAY + "/questionnaires?type=" + this.xlistType + "&user_id=" + Taro.getStorageSync("user_id"),
      header: {
        'content-type': 'application/json',
        "Cookie": Taro.getStorageSync("cookie")
      },
      methods: "GET",
      success: (res) => {
        this.questionnaires = res.data.questionnaires;
      }
    });
  }
}
</script>
