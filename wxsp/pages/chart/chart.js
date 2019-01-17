//chart.js
const wx_chart = require('../../utils/wxcharts.js')
const app = getApp()
var line_chart = null

Page({
  data: {
    currentData: 1,
    goods_id: null,
    query_type: '0',
    price_infos: [],
    time_infos: []
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
  //获取当前滑块的index
  bindchange: function (e) {
    var that = this;
    that.setData({
      currentData: e.detail.current
    })
  },
  //点击切换，滑块index赋值
  checkCurrent: function (e) {
    if (this.data.currentData === e.target.dataset.current) {
      return false;
    } else {
      var type = '0'
      if (e.target.dataset.current == 0) {
        type = '0'
      }
      else if (e.target.dataset.current == 1) {
        type = '1'
      }
      else if (e.target.dataset.current == 2) {
        type = '2'
      }
      else if (e.target.dataset.current == 3) {
        type = '3'
      }

      this.setData({
        currentData: e.target.dataset.current,
        query_type: type
      })

      this.getPriceInfo()
    }
  },
  touchHandler: function (e) {
    line_chart.scrollStart(e);
  },
  moveHandler: function (e) {
    line_chart.scroll(e);
  },
  touchEndHandler: function (e) {
    line_chart.scrollEnd(e);
    line_chart.showToolTip(e, {
      format: function (item, category) {
        return category + ' ' + item.name + ':' + item.data
      }
    });
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
        // that.setData({
        //   price_info: [],
        //   time_infos: []
        // })

        that.setData({
          price_infos: res.data.price_data,
          time_infos: res.data.time_data
          // topline_price: parseFloat(res.data.price_data[0]) + parseFloat(res.data.price_data[0] / 5),
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
    line_chart = new wx_chart({
      canvasId: 'hourPriceCanvas',
      type: 'line',
      enableScroll: true,
      categories: this.data.time_infos,
      // animation: true,
      // background: '#f5f5f5',
      series: [{
        name: '京东自营价格',
        data: this.data.price_infos,
        format: function (val, name) {
          return val.toFixed(2);
        }
      }],
      xAxis: {
        disableGrid: true
      },
      yAxis: {
        title: '商品价格（￥）',
        format: function (val) {
          return val.toFixed(2);
        },
        max: this.data.price_infos[0] + this.data.price_infos[0] / 5,
        min: this.data.price_infos[0] - this.data.price_infos[0] / 5
      },
      width: windowWidth,
      height: 250,
      // dataPointShape: true,
      extra: {
        lineStyle: 'curve'
      }
    });
  }
})
