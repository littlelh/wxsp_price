//chart.js
const wx_chart = require('../../utils/wxcharts.js')
const app = getApp()
var hourline_chart = null
var dayline_chart = null
var monthline_chart = null

Page({
  data: {
    goods_id: null
  },
  onLoad: function (options) {
    this.setData({
      goods_id: options.goods_id
    })
    this.getMothElectro()
    // wx.navigateBack({
    // })
  },
  getMothElectro: function () {
    var windowWidth = 320;
    try {
      var res = wx.getSystemInfoSync();
      windowWidth = res.windowWidth;
    } catch (e) {
      console.error('getSystemInfoSync failed!');
    }
    hourline_chart = new wx_chart({ //当月用电折线图配置
      canvasId: 'hourPriceCanvas',
      type: 'line',
      categories: ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12'], //categories X轴
      animation: true,
      background: '#f5f5f5',
      series: [{
        name: '总用电量',
        //data: yuesimulationData.data,
        data: [1, 6, 9, 1, 0, 2, 9, 2, 8, 6, 0, 9],
        format: function (val, name) {
          return val.toFixed(2) + 'kWh';
        }
      }],
      xAxis: {
        disableGrid: true
      },
      yAxis: {
        title: '当月用电(kWh)',
        format: function (val) {
          return val.toFixed(2);
        },
        max: 20,
        min: 0
      },
      width: windowWidth,
      height: 200,
      dataLabel: false,
      dataPointShape: true,
      extra: {
        lineStyle: 'curve'
      }
    });
  }
})
