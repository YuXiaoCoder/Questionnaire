<template>
  <view class="title-modal">
    <AtModal
    :isOpened="isShow"
    :onConfirm="handleSaveQuestion"
    >
      <AtModalHeader>{{modalTitle}}</AtModalHeader>
      <AtModalContent>
        <AtInput
          required
          autoFocus
          name="title"
          title="标题"
          type="text"
          placeholder="请输入标题"
          :value="question.title"
          :cursor='question.title.length'
          :onChange="handleChangeTitle"
        />
      </AtModalContent>

      <AtModalAction>
        <AtButton
          type="secondary"
          :full="true"
          :onClick="handleCloseAddSingleChoice"
        >
          取消
        </AtButton>
        <AtButton
          type="primary"
          :full="true"
          :onClick='handleSaveQuestion'
        >
          确定
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

// 模态框
import { AtModal, AtModalHeader, AtModalContent, AtModalAction} from "taro-ui-vue";
import "taro-ui-vue/dist/style/components/modal.scss"

// 输入框
import { AtInput } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/input.scss"

// 按钮
import { AtButton } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/button.scss"
import "taro-ui-vue/dist/style/components/loading.scss"

// 消息通知
import { AtMessage } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/message.scss"

export default {
  name: "TitleModal",
  props: {
    isShow: {
      type: Boolean,
      default: false,
    },
    modalTitle: {
      type: String,
      default: '选择题'
    },
    questionItem: null,
    questionnaireID: null,
    questionnaireTitle: {
      type: String,
      default:'新建问卷'
    },
    questionType: {
      type: Number,
      default: 0
    },
  },
  data () {
    return {
      // 问题
      question: {
        id: 0,
        title: "新建问题",
        type: this.questionType
      },
      // 用户ID
      userID: 0,
    }
  },
  components: {
    AtButton,
    AtInput,
    AtModal,
    AtModalHeader,
    AtModalContent,
    AtModalAction,
    AtMessage,
  },
  mounted(){
    // 编辑问题: 转换问题对象
    if(this.questionItem){
      this.question = JSON.parse(JSON.stringify(this.questionItem))
    }

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
  methods: {
    // 修改问题标题
    handleChangeTitle(value) {
      this.question.title = value;
    },

    // 取消按钮
    handleCloseAddSingleChoice() {
      // 清理缓存
      this.$emit('update:questionItem', null);

      // 关闭模态框
      this.$emit('update:isShow', false);
    },

    // 保存按钮
    handleSaveQuestion(){
      if(!this.questionnaireID){
        Taro.request({
          url: API_GATEWAY + '/questionnaires',
          method: 'POST',
          data:{
            "title": this.questionnaireTitle,
            "user_id": this.userID,
            "questions": [
              this.question
            ]
          },
          success: (res) => {
            // 消息通知
            Taro.atMessage({
              'message': '保存成功',
              'type': 'success',
            });

            // 更新问题列表
            this.$emit('updataQuestions', res.data.id);

            // 关闭模态框
            this.$emit('update:isShow', false);
          },
          fail: (res) => {
            // 消息通知
            Taro.atMessage({
              'message': '保存失败',
              'type': 'error',
            });

            // 关闭模态框
            this.$emit('update:isShow', false);
          }
        });
      }
      else{
        // 添加问卷ID
        this.question['questionnaire_id'] = parseInt(this.questionnaireID);
        if(this.question.id != 0) {
          Taro.request({
            url: API_GATEWAY + '/questions/' + this.question.id,
            method: 'PUT',
            data: this.question,
            success: (res) => {
              // 消息通知
              Taro.atMessage({
                'message': '更新成功',
                'type': 'success',
              });

              // 更新问题列表
              this.$emit('updataQuestions', parseInt(this.questionnaireID));

              // 关闭模态框
              this.$emit('update:isShow', false);
            },
            fail: (res) => {
              // 消息通知
              Taro.atMessage({
                'message': '更新失败',
                'type': 'error',
              });

              // 关闭模态框
              this.$emit('update:isShow', false);
            }
          });
        } else {
          Taro.request({
            url: API_GATEWAY + '/questions',
            method: 'POST',
            data: this.question,
            success: (res) => {
              // 消息通知
              Taro.atMessage({
                'message': '保存成功',
                'type': 'success',
              });

              // 更新问题列表
              this.$emit('updataQuestions', parseInt(this.questionnaireID));

              // 关闭模态框
              this.$emit('update:isShow', false);
            },
            fail: (res) => {
              // 消息通知
              Taro.atMessage({
                'message': '保存失败',
                'type': 'error',
              });

              // 关闭模态框
              this.$emit('update:isShow', false);
            }
          });
        }
      }
    },
  }
};
</script>
