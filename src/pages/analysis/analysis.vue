<template>
  <view class="analysis" :style="{height: screenHeight+'px'}">
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
            <view v-if="statistics.length != 0 && JSON.stringify(statistics[index].statistics) != '{}'">
              <view v-for="(value, key) in statistics[index]['statistics']" :key="key">
                {{ key }}: {{ questionnaire.questions[index].options[letter[key]].value }}
                <AtProgress :percent='value' />
              </view>

              <!-- 饼状图 -->
              <view class="chart" style="height: 200px;">
                <ec-canvas :ec="ecs[index]"></ec-canvas>
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
            <view v-if="statistics.length != 0 && JSON.stringify(statistics[index].statistics) != '{}'">
              <view v-for="(value, key) in statistics[index]['statistics']" :key="key">
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
      letter: {
        "A": 0,
        "B": 1,
        "C": 2,
        "D": 3,
        "E": 4,
        "F": 5,
        "G": 6,
        "H": 7,
        "I": 8,
        "J": 9,
        "K": 10,
        "L": 11,
        "M": 12,
        "N": 13,
        "O": 14,
        "P": 15,
        "Q": 16,
        "R": 17,
        "S": 18,
        "T": 19,
        "U": 20,
        "V": 21,
        "W": 22,
        "X": 23,
        "Y": 24,
        "Z": 25
      }
    }
  },
  created(){
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
                  for (let key in this.statistics[index].statistics) {
                    data.push({
                      "value": this.statistics[index].statistics[key],
                      "name": key,
                    })
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
  mounted() {
    setTimeout(() => {
      Taro.createSelectorQuery().select('.index').boundingClientRect(rect=>{
        console.log(Math.floor(rect.height));
      }).exec()
    }, 300)
  },
  methods: {
  },
}
</script>
