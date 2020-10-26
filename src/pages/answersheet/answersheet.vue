<template>
  <view class="answersheet" :style="{height: screenHeight+'px'}">
    <view class="index" :style="{height: scrollBarHeight+'px'}">

      <!-- 问题列表 -->
      <AtList>
      <view v-for="(item, index) in questionnaire.questions" :key="index">
        <!-- 单选题 -->
        <AtCard
          :title='((index+"").length == 1 ? "0":"") + (index+1) + "、" + item.title'
          v-if="item.type == 1"
          class="at-card"
        >
          <AtRadio
            :index='index'
            :options="handleFormatOptions(item.options)"
            :value='answers[index].value'
            :onClick='handleSetChoice'
          />
        </AtCard>

        <!-- 多选题 -->
        <AtCard
          :title='((index+"").length == 1 ? "0":"") + (index+1) + "、" + item.title'
          v-else-if="item.type == 2"
          class="at-card"
        >
          <AtCheckbox
            :index='index'
            :options='handleFormatOptions(item.options)'
            :selectedList="answers[index].value"
            :onChange='handleSetChoice'
          />
        </AtCard>

        <!-- 填空题 -->
        <AtCard
          :title='((index+"").length == 1 ? "0":"") + (index+1) + "、" + item.title'
          v-else-if="item.type == 3"
          class="at-card"
        >
          <input class="at-input__input" type="text" :key="index" v-model="answers[index].value" placeholder="请输入..."/>
        </AtCard>

        <!-- 问答题 -->
        <AtCard
          :title='((index+"").length == 1 ? "0":"") + (index+1) + "、" + item.title'
          v-else-if="item.type == 4"
          class="at-card"
        >
          <textarea
            class="at-textarea__textarea"
            v-model="answers[index].value"
            :maxlength="200"
            placeholder="请输入..."
          />
          <view class="at-textarea__counter">
            {{ answers[index].value.length }}/200
          </view>
        </AtCard>

        <!-- 分隔符 -->
        <AtDivider v-if="index < answers.length - 1" />
      </view>
      </AtList>

      <!-- 导航栏 -->
      <AtTabBar
        :tabList='tarbar.list'
        :fixed='tarbar.fixed'
        :onClick='handleClickTabBar'
      />

      <!-- 消息通知 -->
      <AtMessage />
    </view>
  </view>
</template>

<script>
// 样式文件
import './answersheet.scss'

// 单选按钮: 自定义组件
import AtRadio from "../../components/radio"

// 多选按钮: 自定义组件
import AtCheckbox from "../../components/checkbox"

// Taro
import Taro from '@tarojs/taro'
import { getCurrentInstance } from '@tarojs/taro'

// 卡片
import { AtCard } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/card.scss"

// 分隔符
import { AtDivider } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/divider.scss"

// 列表
import { AtList } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/list.scss"

// 消息通知
import { AtMessage } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/message.scss"

// 底部导航栏
import { AtTabBar } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/tab-bar.scss"
import "taro-ui-vue/dist/style/components/badge.scss"

// 图标
import "taro-ui-vue/dist/style/components/icon.scss"

// 输入框
import "taro-ui-vue/dist/style/components/input.scss"

// 文本框
import "taro-ui-vue/dist/style/components/textarea.scss"

export default {
  components: {
    AtTabBar,
    AtList,
    AtCard,
    AtRadio,
    AtCheckbox,
    AtDivider,
    AtMessage,
  },
  data () {
    return {
      // 屏幕高度
      screenHeight: 600,
      // 滚动条高度
      scrollBarHeight: 540,
      // 导航栏
      tarbar: {
        list: [
          { title: '提交', iconType: 'check' },
        ],
        // 固定到底部
        fixed: true,
      },
      // 问卷
      questionnaire: {
        id: 0,
        title: "",
        questions: []
      },
      // 问卷ID
      questionnaireID: null,
      // 答案
      answers: [],
    }
  },
  created() {
    // 获取屏幕高度
    this.screenHeight = Taro.getSystemInfoSync().windowHeight;
    // 获取查询参数
    this.questionnaireID = parseInt(getCurrentInstance().router.params.id);

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
  mounted() {
    // ES6: 箭头函数
    // 获取底部导航栏高度
    setTimeout(() => {
      Taro.createSelectorQuery().select('.at-tab-bar').boundingClientRect(rect=>{
        this.scrollBarHeight = this.screenHeight - Math.floor(rect.height) - 1;
      }).exec()
    }, 300)

    // 获取问卷
    Taro.request({
      url: API_GATEWAY + "/questionnaires/" + this.questionnaireID,
      header: {
        'content-type': 'application/json'
      },
      methods: "GET",
      success: (res) => {
        this.questionnaire = res.data;
        // 初始化答案
        this.questionnaire.questions.map(element => {
          if (element.type == 2) {
            this.answers.push({
              value: [],
              question_id: element.id
            });
          } else {
            this.answers.push({
              value: "",
              question_id: element.id
            });
          }
        });
      }
    });
  },
  methods: {
    // 导航栏: 提交问卷
    handleClickTabBar(value){
      Taro.request({
        url: API_GATEWAY + '/answersheets',
        method: 'POST',
        data: {
          "user_id": Taro.getStorageSync("user_id"),
          "answers": this.answers,
          "questionnaire_id": this.questionnaireID,
        },
        success: (res) => {
          // 消息通知
          Taro.atMessage({
            'message': '提交成功',
            'type': 'success',
          });
        },
        fail: (res) => {
          // 消息通知
          Taro.atMessage({
            'message': '提交失败',
            'type': 'error',
          });
        }
      });
    },

    // 格式化选项
    handleFormatOptions(options) {
      return options.map(itme => {
        return {
          "label": itme.value,
          "value": itme.key,
        }
      })
    },

    // 设置选择题选项
    handleSetChoice(value, index) {
      this.answers[index].value = value
    },
  },
}
</script>
