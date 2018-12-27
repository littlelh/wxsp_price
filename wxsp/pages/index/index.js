//index.js
//获取应用实例
const app = getApp()

Page({
  data: {
    inputValue: '',
    product_name: '',
    // goods_list: [{Price: 1999, Coupon: 300, Discount: '折扣', Name: 'test', Describe: '4GB+64GB', ImgUrl:'https://www.jiantong.xyz/loadimage?index=1'}]
    goods_list: [],
    shop_type: {
      1: '京东自营',
      2: '不知名商家',
    }
  },
  //事件处理函数
  // bindViewTap: function() {
  //   wx.navigateTo({
  //     url: '../logs/logs'
  //   })
  // },
  onShareAppMessage: function (res) {
    var that = this
    return {
      title: '最真实的价格',
      path: '/pages/index/index?id=' + that.data.scratchId,
      success: function (res) {
        // 转发成功
        that.shareClick();
      },
      fail: function (res) {
        // 转发失败
      }
    }
  },
  onPullDownRefresh: function () { // 下拉加载
    setTimeout(function () {
      wx.stopPullDownRefresh();
      console.log(1);
    }, 1000)
  },
  SearchInput: function (e) {
    this.setData({
      product_name: e.detail.value,
      inputValue: '',
      motto: e.detail.value // test
    })
    console.log("input product name: ", e.detail.value)
  },
  onLoad: function () {
    let that = this
    wx.request({
      url: 'https://www.jiantong.xyz/allgoodsinfo',
      method: 'GET',
      header: {
        //设置参数内容类型为json
        'content-type': 'application/json'
      },
      success: function (res) {
        for (var index in res.data.data) {
          // console.log(res.data.data[index])
          res.data.data[index].ImgUrl = 'https://www.jiantong.xyz/loadimage?index=' + res.data.data[index].Id
        }
        that.setData({
          goods_list: res.data.data
        })
        // console.log(that.data.goods_list)
      }
    })
  }
})
