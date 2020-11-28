<template>
  <view class="questionnaire" :style="{height: screenHeight+'px'}">
    <view class="index" :style="{height: scrollBarHeight+'px'}">
      <!-- 问卷标题 -->
      <AtInput
        required
        name='title'
        title='标题'
        type='text'
        placeholder='请输入标题'
        :value='title'
        :cursor='title.length'
        :onBlur='handleChangeTitle'
      />

      <AtDivider />

      <!-- 问题列表 -->
      <AtList>
      <view v-for="(item, index) in questions" :key="index">
        <QSList
          :key="index"
          :index="index"
          :qsItem="item"
          :onClick="handleDeleteQuestion"
        >
          <view @tap="editQuestion(item)">{{item.title}}</view>
        </QSList>
      </view>
      </AtList>

      <!-- 导航栏 -->
      <AtTabBar
        :tabList='tarbar.list'
        :fixed='tarbar.fixed'
        :onClick='handleClickTabBar'
      />

      <!-- 动作面板 -->
      <AtActionSheet :isOpened='actionPanelFlag' cancelText='取消' :onCancel='handleCloseActionPanel' :onClose='handleCloseActionPanel'>
        <AtActionSheetItem :onClick='handleAddSingleChoice'>
          单选
        </AtActionSheetItem>
        <AtActionSheetItem :onClick='handleAddMultipleChoice'>
          多选
        </AtActionSheetItem>
        <AtActionSheetItem :onClick='handleAddFillBlank' v-if="questionnaireType == 2">
          填空
        </AtActionSheetItem>
        <AtActionSheetItem :onClick='handleAddQuestionAndAnswer' v-if="questionnaireType == 2">
          问答
        </AtActionSheetItem>
      </AtActionSheet>

      <!-- 单选题弹窗 -->
      <SelectModal
        v-if="addSingleChoiceFlag"
        modalTitle="单选"
        :isShow.sync='addSingleChoiceFlag'
        :questionItem.sync='questionItem'
        :questionType=1
        :questionnaireID='questionnaireID'
        :questionnaireTitle='title'
        :questionnaireType='questionnaireType'
        @updataQuestions='handleUpdateQuestions'
      />

      <!-- 多选题弹窗 -->
      <SelectModal
        v-if="addMultipleChoiceFlag"
        modalTitle="多选"
        :isShow.sync='addMultipleChoiceFlag'
        :questionItem.sync='questionItem'
        :questionType=2
        :questionnaireID='questionnaireID'
        :questionnaireTitle='title'
        :questionnaireType='questionnaireType'
        @updataQuestions='handleUpdateQuestions'
      />

      <!-- 填空题弹窗 -->
      <TitleModal
        v-if="addFillBlankFlag"
        modalTitle="填空"
        :isShow.sync='addFillBlankFlag'
        :questionItem.sync='questionItem'
        :questionType=3
        :questionnaireID='questionnaireID'
        :questionnaireTitle='title'
        :questionnaireType='questionnaireType'
        @updataQuestions='handleUpdateQuestions'
      />

      <!-- 问答题弹窗 -->
      <TitleModal
        v-if="addQuestionAndAnswerFlag"
        modalTitle="问答"
        :isShow.sync='addQuestionAndAnswerFlag'
        :questionItem.sync='questionItem'
        :questionType=4
        :questionnaireID='questionnaireID'
        :questionnaireTitle='title'
        :questionnaireType='questionnaireType'
        @updataQuestions='handleUpdateQuestions'
      />

      <!-- 消息通知 -->
      <AtMessage />
    </view>
  </view>
</template>

<script>
// Taro
import Taro from '@tarojs/taro'
import { getCurrentInstance } from '@tarojs/taro'

// 布局
import "taro-ui-vue/dist/style/components/flex.scss"

// 分隔符
import { AtDivider } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/divider.scss"

// 列表
import { AtList } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/list.scss"

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

// 自定义组件
import QSList from './qslist'
import SelectModal from './selectModal'
import TitleModal from './titleModel'

// 输入框
import { AtInput } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/input.scss"

// 自定义样式
import './questionnaire.scss'

export default {
  components: {
    AtTabBar,
    AtDivider,
    AtInput,
    AtActionSheet,
    AtActionSheetItem,
    AtList,
    QSList,
    SelectModal,
    TitleModal,
    AtMessage,
  },
  data () {
    return {
      // 屏幕高度
      screenHeight: 600,
      // 滚动条高度
      scrollBarHeight: 540,
      // 问卷标题
      title: '新建问卷',
      // 导航栏
      tarbar: {
        list: [
          { title: '添加', iconType: 'add' },
        ],
        // 固定到底部
        fixed: true,
      },
      // 动作面板
      actionPanelFlag: false,
      // 添加单选题
      addSingleChoiceFlag: false,
      // 添加多选题
      addMultipleChoiceFlag: false,
      // 添加填空题
      addFillBlankFlag: false,
      // 添加问答题
      addQuestionAndAnswerFlag: false,
      // 问题列表
      questions: [],
      // 问卷ID
      questionnaireID: null,
      // 问卷类型: 投票/问卷
      questionnaireType: 1,
      // 问题对象
      questionItem: null,
      // 用户ID
      userID: 0,
    }
  },
  methods: {
    // 问卷标题
    handleChangeTitle(value) {
      console.log(value);
      console.log(this.questionnaireID);
      if(!this.questionnaireID){
        Taro.request({
          url: API_GATEWAY + '/questionnaires',
          method: 'POST',
          data: {
            "title": value,
            "user_id": this.userID,
            "type": this.questionnaireType,
          },
          success: (res) => {
            this.questionnaireID = res.data.id;
            // 消息通知
            Taro.atMessage({
              'message': '保存成功',
              'type': 'success',
            });
          },
          fail: (res) => {
            // 消息通知
            Taro.atMessage({
              'message': '保存失败',
              'type': 'error',
            });
          }
        });
      } else {
        Taro.request({
          url: API_GATEWAY + '/questionnaires/' + this.questionnaireID,
          method: 'PUT',
          data: {
            "title": value,
          },
          success: (res) => {
            // 消息通知
            Taro.atMessage({
              'message': '更新成功',
              'type': 'success',
            });
          },
          fail: (res) => {
            // 消息通知
            Taro.atMessage({
              'message': '更新失败',
              'type': 'error',
            });
          }
        });
      }
    },

    // 动作面板
    handleClickTabBar() {
      this.actionPanelFlag = true
    },
    handleCloseActionPanel() {
      this.actionPanelFlag = false
    },

    //编辑题目
    editQuestion(item){
      // 填充问题
      this.questionItem = JSON.parse(JSON.stringify(item));
      // 显示模态框
      if(item.type == 1) {
        this.addSingleChoiceFlag = true;
      } else if(item.type == 2) {
        this.addMultipleChoiceFlag = true;
      } else if(item.type == 3) {
        this.addFillBlankFlag = true;
      } else if(item.type == 4) {
        this.addQuestionAndAnswerFlag = true;
      }
    },

    // 添加单选
    handleAddSingleChoice() {
      this.actionPanelFlag = false;
      this.addSingleChoiceFlag = true;
    },

    // 添加多选
    handleAddMultipleChoice() {
      this.actionPanelFlag = false;
      this.addMultipleChoiceFlag = true;
    },

    // 添加填空
    handleAddFillBlank() {
      this.actionPanelFlag = false;
      this.addFillBlankFlag = true;
    },

    // 添加问答
    handleAddQuestionAndAnswer() {
      this.actionPanelFlag = false;
      this.addQuestionAndAnswerFlag = true;
    },

    // 删除问题
    handleDeleteQuestion(index){
      Taro.request({
        url: API_GATEWAY + "/questions/" + this.questions[index].id,
        method: 'DELETE',
        success: (res) => {
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
      // 动态渲染
      this.questions.splice(index, 1);
    },

    // 更新问题列表
    handleUpdateQuestions(id){
      Taro.request({
        url: API_GATEWAY + '/questionnaires/' + id,
        method: 'GET',
        header: {
          'content-type': 'application/json'
        },
        success: (res) => {
          this.title = res.data.title;
          this.questionnaireID = res.data.id;
          this.questions = res.data.questions;
        }
      });
    },
  },
  created(){
    // 获取屏幕高度
    this.screenHeight = Taro.getSystemInfoSync().windowHeight;
    // 获取查询参数
    this.questionnaireID = getCurrentInstance().router.params.id;
    this.questionnaireType = parseInt(getCurrentInstance().router.params.type);

    // 用户ID
    if (this.userID == 0) {
      // 登录
      Taro.login({
        success: res => {
          Taro.request({
            url: API_GATEWAY + "/login",
            data: {
              code: res.code
            },
            success: res => {
              // 存储用户ID
              this.userID = res.data['user_id'];
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
  },
  onReady(){
    // ES6: 箭头函数
    // 获取底部导航栏高度
    setTimeout(() => {
      Taro.createSelectorQuery().select('.at-tab-bar').boundingClientRect(rect=>{
        this.scrollBarHeight = this.screenHeight - Math.floor(rect.height) - 10;
      }).exec()
    }, 300)
  },
  mounted(){
    if (this.questionnaireType == 1) {
      this.title = "新建投票"
    } else {
      this.title = "新建问卷"
    }

    // 编辑问卷
    if(this.questionnaireID){
      Taro.request({
        url: API_GATEWAY + '/questionnaires/' + this.questionnaireID,
        method: 'GET',
        header: {
          'content-type': 'application/json'
        },
        success: (res) => {
          this.title = res.data.title;
          this.questions = res.data.questions;
        }
      });
      if (this.questionnaireType == 1) {
        Taro.setNavigationBarTitle({
          title: '编辑投票'
        });
      } else {
        Taro.setNavigationBarTitle({
          title: '编辑问卷'
        });
      }
    } else {
      if (this.questionnaireType == 1) {
        Taro.setNavigationBarTitle({
          title: '添加投票'
        });
      } else {
        Taro.setNavigationBarTitle({
          title: '添加问卷'
        });
      }
    }
  }
}
</script>
