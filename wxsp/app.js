//app.js

App({
  onLaunch: function () {
    // 展示本地存储能力
    // var logs = wx.getStorageSync('logs') || []
    // logs.unshift(Date.now())
    // wx.setStorageSync('logs', logs)

    let login_flag = wx.getStorageSync('skey')
    if (login_flag) {
      wx.checkSession({
        success: res => {
          // session key 未过期
        },
        fail: res => {
          this.UserLogin()
        }
      })
    } else {
      // 无 skey，首次登陆
      this.UserLogin()
    }
    
    // 获取用户信息
    wx.getSetting({
      success: res => {
        if (res.authSetting['scope.userInfo']) {
          this.globalData.isAuth = true
          // 已经授权，可以直接调用 getUserInfo 获取头像昵称，不会弹框
          wx.getUserInfo({
            success: res => {
              // 可以将 res 发送给后台解码出 unionId
              this.globalData.userInfo = res.userInfo

              // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
              // 所以此处加入 callback 以防止这种情况
              if (this.userInfoReadyCallback) {
                this.userInfoReadyCallback(res)
              }
            }
          })
        } else {
          wx.reLaunch({
            url: '/pages/login/login',
          })
        }
      }
    })
  },
  UserLogin: function() {
    // 登录
    wx.login({
      success: res => {
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
        if (res.code) {
          console.log("code:" + res.code)
          this.globalData.user_code = res.code
        } else {
          console.log('登录失败' + res.errMsg)
        }
      }
    })
  },
  globalData: {
    userInfo: null,
    isAuth: false,
    user_code: null
  }
})