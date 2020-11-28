<template>
  <view class="select-modal">
    <!-- 模态框 -->
    <AtModal
    :isOpened="isShow"
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

        <!-- 分隔符 -->
        <AtDivider />

        <AtList class="selectModel">
          <view v-for="(item, index) in question.options" :key="index" >
            <QSList :onClick="handleDeleteQuestionOption" :index="index" :key="index">
              <view>
                <input type="text" :key="item.value" v-model="item.value" placeholder="请输入选项"/>
              </view>
            </QSList>
          </view>
        </AtList>

        <AtButton
          type="secondary"
          :circle="true"
          :onClick="handleAddQuestionOption"
        >
          添加选项
        </AtButton>
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
// 自定义组件
import QSList from "./qslist";

// Taro
import Taro from '@tarojs/taro'

// 模态框
import { AtModal, AtModalHeader, AtModalContent, AtModalAction} from "taro-ui-vue";
import "taro-ui-vue/dist/style/components/modal.scss"

// 分隔符
import { AtDivider } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/divider.scss"

// 输入框
import { AtInput } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/input.scss"

// 按钮
import { AtButton } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/button.scss"
import "taro-ui-vue/dist/style/components/loading.scss"

// 列表
import { AtList } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/list.scss"

// 消息通知
import { AtMessage } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/message.scss"

export default {
  name: "SelectModal",
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
      default: '新建问卷'
    },
    questionType: {
      type: Number,
      default: 0
    },
    questionnaireType: {
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
        options: [],
        type: this.questionType
      },
      // 用户ID
      userID: 0,
    }
  },
  components: {
    AtButton,
    AtInput,
    AtDivider,
    AtModal,
    AtModalHeader,
    AtModalContent,
    AtModalAction,
    QSList,
    AtList,
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

    // 添加选项
    handleAddQuestionOption() {
      this.question.options.push({
        value: "选项" + (this.question.options.length + 1).toString(),
        label: "选项" + (this.question.options.length + 1).toString(),
      });
    },

    // 删除选项
    handleDeleteQuestionOption(index) {
      // 删除选项
      this.question.options.splice(index, 1)
      // 更新选项
      this.question.options = this.question.options.map((item) => {
        return {
          label: item.value,
          value: item.value,
        }
      })
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
      // 更新选项
      this.question.options = this.question.options.map((item) => {
        return {
          label: item.value,
          value: item.value,
        }
      })

      if(!this.questionnaireID){
        Taro.request({
          url: API_GATEWAY + '/questionnaires',
          method: 'POST',
          data:{
            "title": this.questionnaireTitle,
            "user_id": this.userID,
            "type": this.questionnaireType,
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
        this.question['questionnaire_id'] = parseInt(this.questionnaireID)
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

<style lang="scss">
.selectModel{
    max-height: 400rpx;
    overflow-x: hidden;
}
</style>
