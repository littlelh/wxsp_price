//chart.js
const wx_chart = require('../../utils/wxcharts.js')
const app = getApp()
var hourline_chart = null
var dayline_chart = null
var monthline_chart = null

Page({
  data: {
    goods_id: null,
    query_type: '0',
    price_infos: [],
    coupon_infos: [],
    time_infos: [],
    goods_maxprice: 0
  },
  onLoad: function (options) {
    this.setData({
      goods_id: options.goods_id
    })

    this.getPriceInfo()
    // this.getMothElectro()
    // wx.navigateBack({
    // })
  },
  getPriceInfo: function () {
    var that = this
    wx.request({
      url: 'https://www.jiantong.xyz/goodsprice?goods_id=' + that.data.goods_id + '&query_type=' + that.data.query_type,
      method: 'GET',
      header: {
        //设置参数内容类型为json
        'content-type': 'application/json'
      },
      success: function (res) {
        that.setData({
          price_info: [],
          time_infos: []
        })

        var result_price = []
        var result_time = []
        var result_coupon = []
        var max_price = 0
        for (var index in res.data.data) {
          result_price.push(res.data.data[index].Price)
          result_time.push(res.data.data[index].Time)
          result_coupon.push(res.data.data[index].Price - res.data.data[index].Coupon)
          if (parseFloat(res.data.data[index].Price) > max_price) {
            max_price = parseFloat(res.data.data[index].Price)
          }
        }
        that.setData({
          price_infos: result_price,
          coupon_infos: result_coupon,
          time_infos: result_time,
          goods_maxprice: max_price
        })
        that.getMothElectro()
      }
    })
  },
  getMothElectro: function () {
    var windowWidth = 320;
    try {
      var res = wx.getSystemInfoSync();
      windowWidth = res.windowWidth;
    } catch (e) {
      console.error('getSystemInfoSync failed!');
    }
    hourline_chart = new wx_chart({
      canvasId: 'hourPriceCanvas',
      type: 'line',
      categories: this.data.time_infos,
      animation: true,
      // background: '#f5f5f5',
      series: [{
        name: '原价格',
        data: this.data.price_infos,
        format: function (val, name) {
          return val.toFixed(2) + '￥';
        }
      },
      {
        name: '加优惠券价格',
        data: this.data.coupon_infos,
        format: function (val, name) {
          return val.toFixed(2) + '￥';
        }
      }],
      xAxis: {
        disableGrid: true
      },
      yAxis: {
        // title: '价格（￥）',
        format: function (val) {
          return val.toFixed(2);
        },
        max: this.data.goods_maxprice + this.data.goods_maxprice/5,
        min: 0
      },
      width: windowWidth,
      height: 250,
      dataLabel: false,
      dataPointShape: true,
      extra: {
        lineStyle: 'curve'
      }
    });
  }
})
