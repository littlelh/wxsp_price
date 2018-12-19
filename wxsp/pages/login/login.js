//logs.js
const app = getApp()
const util = require('../../utils/util.js')

Page({
  data: {
    userInfo: {},
    hasUserInfo: false,
    canIUse: wx.canIUse('button.open-type.getUserInfo')
  },
  onLoad: function () {
    console.log('login ing...')
  },
  getUserInfo: function (e) {
    if (e.detail.userInfo) {
      console.log(e)
      app.globalData.userInfo = e.detail.userInfo,
      app.globalData.isAuth = true
      this.setData({
        userInfo: e.detail.userInfo,
        hasUserInfo: true
      })
      wx.reLaunch({
        url: '/pages/index/index',
      })
    } else {
      app.globalData.isAuth = false
    }
  }
})
