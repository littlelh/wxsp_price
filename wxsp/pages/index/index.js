//index.js
//获取应用实例
const app = getApp()

Page({
  data: {
    inputValue: '',
    product_name: '',
    // goods_list: [{Price: 1999, Coupon: 300, Discount: '折扣', Name: 'test', Describe: '4GB+64GB', ImgUrl:'https://www.jiantong.xyz/loadimage?index=1'}]
    goods_list: [],
    scroll_height: 1000,
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
  onLoad: function () {
    // TODO: 后续加入openId的逻辑
    if (app.globalData.openId) {
    } else {
      app.openIdReadyCallback = res => {
        console.log(app.globalData.openId)
      }
    }

    var that = this
    wx.request({
      url: 'https://www.jiantong.xyz/goodsinfo?first=1&last=5',
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
  },
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
  loadMore: function () {
    var that = this
    var first = this.data.goods_list.length + 1
    var last = first + 4
    wx.request({
      url: 'https://www.jiantong.xyz/goodsinfo?first=' + String(first) + '&last=' + String(last),
      method: 'GET',
      header: {
        //设置参数内容类型为json
        'content-type': 'application/json'
      },
      success: function (res) {
        for (var index in res.data.data) {
          res.data.data[index].ImgUrl = 'https://www.jiantong.xyz/loadimage?index=' + res.data.data[index].Id
        }
        console.log(res.data.data.length)
        if (res.data.data.length == 0) {
          wx.showToast({ //如果全部加载完成了也弹一个框
            title: '我也是有底线的',
            icon: 'success',
            duration: 1000
          });
        } else {
          wx.showLoading({ //期间为了显示效果可以添加一个过度的弹出框提示“加载中”  
            title: '加载中...',
            icon: 'loading',
            mask: true
          });
          setTimeout(() => {
            wx.hideLoading();
          }, 1000)
        }

        var result = that.data.goods_list.concat(res.data.data);
        that.setData({
          goods_list: result
        })
      }
    })
  },
  previewImage: function (e) {
    var current = e.target.dataset.src
    var img_list = [current]
    wx.previewImage({
      current: current, // 当前显示图片的https链接
      urls: img_list // 需要预览的图片https链接列表
    })
  },
  viewPriceInfo: function(e) {
    var that = this
    wx.navigateTo({
      url: '../chart/chart?goods_id=' + e.currentTarget.id,
      success: function (res) {
        console.log("navigateTo price info page success.")
      },
      fail: function () {
        console.log("navigateTo price info page fail.")
      },
      complete: function () {
        // 无论成功或者失败都会调用
        console.log("navigateTo price info page complete.")
      }
    })
  }
})
