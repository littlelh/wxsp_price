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
          this.globalData.openId = login_flag
          if (this.openIdReadyCallback) {
            this.openIdReadyCallback(res)
          }
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
          this.GetUserData()
        } else {
          console.log('登录失败' + res.errMsg)
        }
      }
    })
  },
  GetUserData: function() {
    wx.request({
      url: 'https://www.jiantong.xyz/login',
      method: 'GET',
      // data: 'pageSize=1&pageNum=10',
      data: {
        appid: 'wx682e545eba2a1a0a',
        secret: '7537280e6207537aab605e1d97cf5f39',
        js_code: this.globalData.user_code,
        grant_type: 'authorization_code'
      },
      header: {
        //设置参数内容类型为json
        'content-type': 'application/json'
      },
      success: res => {
        console.log(res.data)
        wx.setStorageSync('skey', res.data.openid)
        this.globalData.openId = res.data.openid
        if (this.openIdReadyCallback) {
          this.openIdReadyCallback(res)
        }
      }
    })
  },
  globalData: {
    userInfo: null,
    isAuth: false,
    user_code: null,
    openId: null
  }
})