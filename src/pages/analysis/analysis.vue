<template>
  <view class="analysis">
    <view class="index">
      <!-- 问题列表 -->
      <AtList>
        <view v-for="(item, index) in questionnaire.questions" :key="index">
          <AtCard
            :title='((index+"").length == 1 ? "0":"") + (index+1) + "、" + item.title'
            v-if="item.type == 1 || item.type == 2"
            class="at-card"
          >
            <!-- 详细占比 -->
            <view v-if="statistics.length != 0 && JSON.stringify(statistics[index].result) != '{}'">
              <view v-for="(value, key, optionIndex) in statistics[index].result" :key="key">
                {{ letter[optionIndex] }}: {{ key }}
                <AtProgress :percent='value' />
              </view>

              <!-- 饼状图 -->
              <view class="chart">
                <ec-canvas :ec="ecs[index]" :id="`chart-${index}`" :canvas-id="`chart-${index}`"></ec-canvas>
              </view>
            </view>
            <view v-else>
              暂无数据
            </view>
          </AtCard>

          <AtCard
            :title='((index+"").length == 1 ? "0":"") + (index+1) + "、" + item.title'
            v-else-if="item.type == 3 || item.type == 4"
            class="at-card"
          >
            <AtList>
            <view v-if="statistics.length != 0 && JSON.stringify(statistics[index].result) != '{}'">
              <view v-for="(value, key) in statistics[index].result" :key="key">
                <AtListItem :title='key' :hasBorder="false" />
              </view>
            </view>
            </AtList>
          </AtCard>

          <!-- 分隔符 -->
          <AtDivider v-if="index < questionnaire.questions.length - 1" />
        </view>
      </AtList>

      <!-- 消息通知 -->
      <AtMessage />
    </view>
  </view>
</template>

<script>
// Echarts
import * as echarts from '../../components/ec-canvas/echarts';

// Taro
import Taro from '@tarojs/taro'
import { getCurrentInstance } from '@tarojs/taro'

// 卡片
import { AtCard } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/card.scss"

// 分隔符
import { AtDivider } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/divider.scss"

// 进度条
import { AtProgress } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/progress.scss"

// 列表
import { AtList, AtListItem } from "taro-ui-vue"
import "taro-ui-vue/dist/style/components/list.scss"

// 消息通知
import { AtMessage } from 'taro-ui-vue'
import "taro-ui-vue/dist/style/components/message.scss"

// 图标
import "taro-ui-vue/dist/style/components/icon.scss"

// 自定义样式
import './analysis.scss';

export default {
  components: {
    AtCard,
    AtList,
    AtListItem,
    AtProgress,
    AtDivider,
    AtMessage,
  },
  data () {
    return {
      // 屏幕高度
      screenHeight: 600,
      // 滚动条高度
      scrollBarHeight: 540,
      // 问卷ID
      questionnaireID: 0,
      // 问卷
      questionnaire: {
        id: 0,
        title: "",
        questions: []
      },
      // 图表元数据
      ecs: [],
      // 统计分析
      statistics: [],
      // 字母
      letter: [
        "A",
        "B",
        "C",
        "D",
        "E",
        "F",
        "G",
        "H",
        "I",
        "J",
        "K",
        "L",
        "M",
        "N",
        "O",
        "P",
        "Q",
        "R",
        "S",
        "T",
        "U",
        "V",
        "W",
        "X",
        "Y",
        "Z"
      ]
    }
  },
  created(){
    // 获取屏幕高度
    this.screenHeight = Taro.getSystemInfoSync().windowHeight;
    // 获取查询参数
    this.questionnaireID = parseInt(getCurrentInstance().router.params.id);

    // 获取问卷
    Taro.request({
      url: API_GATEWAY + '/questionnaires/' + this.questionnaireID,
      method: 'GET',
      header: {
        'content-type': 'application/json'
      },
      success: (res) => {
        this.questionnaire = res.data

        // 统计分析
        Taro.request({
          url: API_GATEWAY + '/questionnaires/' + this.questionnaireID + '/analysis',
          method: 'GET',
          header: {
            'content-type': 'application/json'
          },
          success: (res) => {
            this.statistics = res.data
            // 初始化图表数据
            this.questionnaire.questions.map((element, index) => {
              this.ecs.push({
                onInit: (canvas, width, height, dpr) => {
                  const chart = echarts.init(canvas, null, {
                    width: width,
                    height: height,
                    devicePixelRatio: dpr
                  });
                  canvas.setChart(chart);

                  let option = {
                    backgroundColor: "#ffffff",
                    color: ["#37A2DA", "#32C5E9", "#67E0E3", "#91F2DE", "#FFDB5C", "#FF9F7F"],
                    series: [{
                      label: {
                        normal: {
                          fontSize: 14
                        }
                      },
                      type: 'pie',
                      center: ['50%', '50%'],
                      radius: ['40%', '60%'],
                    }]
                  };

                  // 数据
                  let data = []

                  // 选项编号
                  let i = 0
                  let result = this.statistics[index].result;

                  // 添加结果
                  for (let v in result) {
                    data.push({
                      "value": result[v],
                      "name": this.letter[i]
                    })
                    i++
                  }
                  option.series[0].data = data

                  chart.setOption(option);
                  return chart;
                }
              })
            });
          },
          fail: (res) => {
            // 消息通知
            Taro.atMessage({
              'message': '请求失败',
              'type': 'error',
            });
          }
        });
      },
      fail: (res) => {
        // 消息通知
        Taro.atMessage({
          'message': '请求失败',
          'type': 'error',
        });
      }
    });
  },
}
</script>

<style lang="scss">
.chart {
  width: 100%;
  height: 300px;
}
</style>
