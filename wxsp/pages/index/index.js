//index.js
//获取应用实例
const app = getApp()

Page({
  data: {
    motto: 'Hello World',
    inputValue: '',
    product_name: ''
  },
  //事件处理函数
  // bindViewTap: function() {
  //   wx.navigateTo({
  //     url: '../logs/logs'
  //   })
  // },
  onShareAppMessage: function (res) {
    var that = this;
    return {
      title: '你不懂的',
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
    
  }
})
