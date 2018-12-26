#coding:utf-8

import sys
import os
import re
import requests
from lxml import etree
import time
import pymysql
#import json

sys.path.append("/usr/local/lib/python3.6")

img_path = '/data/todd/wxsp_image/'
price_url = 'https://p.3.cn/prices/mgets?pduid=1280795126&skuIds=J_'
favourable_url = 'https://cd.jd.com/promotion/v2?skuId=%s&area=1_72_2799_0&shopId=%s&venderId=%s&cat=%s'

db = pymysql.connect(host="localhost", user="todd", password="temppwd", db="wxsp_price", port=5049)

g_shop_name = {}
# 1:京东
# 其他：未定

  # 京东另一种获取 skuid 的方式
  #  html = etree.HTML(r.text)
  #  datas = html.xpath('//script[contains(@charset,"gbk")]')
  #  for data in datas:
  #      head_data = data.xpath('text()')
  #      skuid_list = re.findall(r"skuid: (.+?),", head_data[0])
  #      print(skuid_list[0])

  #  datas = html1.xpath('//div[contains(@class,"summary-price-wrap")]')

# 查找id对应的店名
def QueryShopName():
    cursor = db.cursor()
    cursor.execute("select * from t_shop_type_name")
    # data = cursor.fetchone()
    datas = cursor.fetchall()
    for data in datas:
        g_shop_name[data[0]] = data[1]
    # print(g_shop_name)
    cursor.close()

def CrawlJingDong(url, pid):
    price = ''  # 价格
    coupon = '0'  # 优惠券
    discount = ''  # 折扣

    head = {'scheme': 'https', 'user-agent': 'Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36', 'x-requested-with': 'XMLHttpRequest', 'Cookie':'qrsc=3; pinId=RAGa4xMoVrs; xtest=1210.cf6b6759; ipLocation=%u5E7F%u4E1C; _jrda=5; TrackID=1aUdbc9HHS2MdEzabuYEyED1iDJaLWwBAfGBfyIHJZCLWKfWaB_KHKIMX9Vj9_2wUakxuSLAO9AFtB2U0SsAD-mXIh5rIfuDiSHSNhZcsJvg; shshshfpa=17943c91-d534-104f-a035-6e1719740bb6-1525571955; shshshfpb=2f200f7c5265e4af999b95b20d90e6618559f7251020a80ea1aee61500; cn=0; 3AB9D23F7A4B3C9B=QFOFIDQSIC7TZDQ7U4RPNYNFQN7S26SFCQQGTC3YU5UZQJZUBNPEXMX7O3R7SIRBTTJ72AXC4S3IJ46ESBLTNHD37U; ipLoc-djd=19-1607-3638-3638.608841570; __jdu=930036140; user-key=31a7628c-a9b2-44b0-8147-f10a9e597d6f; areaId=19; __jdv=122270672|direct|-|none|-|1529893590075; PCSYCityID=25; mt_xid=V2_52007VwsQU1xaVVoaSClUA2YLEAdbWk5YSk9MQAA0BBZOVQ0ADwNLGlUAZwQXVQpaAlkvShhcDHsCFU5eXENaGkIZWg5nAyJQbVhiWR9BGlUNZwoWYl1dVF0%3D; __jdc=122270672; shshshfp=72ec41b59960ea9a26956307465948f6; rkv=V0700; __jda=122270672.930036140.-.1529979524.1529984840.85; __jdb=122270672.1.930036140|85.1529984840; shshshsID=f797fbad20f4e576e9c30d1c381ecbb1_1_1529984840145'
    }

    response = requests.get(url, headers=head)
    response.encoding='gbk'

    # 抓取图片
    save_path = img_path + str(pid) + '.jpg'
    if not os.path.exists(save_path):
        image_url = '' # 图片url
        data_url = []
        html = etree.HTML(response.text)
        datas = html.xpath('//div[contains(@id, "preview")]')
        for data in datas:
            data_url = data.xpath('div[@class="jqzoom main-img"]/img[@id="spec-img"]/@data-origin')
        image_url = 'https:' + data_url[0]

        img_response = requests.get(image_url)
        img = img_response.content

        with open(save_path, 'wb') as f:
            f.write(img)

    # 抓取价格
    skuid = re.findall(r'https://item.jd.com/(.*?).html', response.url)[0]
    price_response = requests.get(price_url + str(skuid), headers=head)
    price_response.encoding='gbk'

    #price = price_response.json()[0]['p']
    price_list = re.findall(r'"p":"(.*?)"', price_response.text)
    if len(price_list) == 0:
        print('Error, get price list is null')
        return
    else:
        price = price_list[0]
        if len(price) == 0:
            print('Error, get price is null')
            return

    # 抓取优惠券信息
    vender_id = ''
    shop_id = ''
    ids = re.findall(r"venderId:(.*?),\s.*?shopId:'(.*?)'", response.text)
    if not ids:
        ids = re.findall(r"venderId:(.*?),\s.*?shopId:(.*?),", response.text)
    vender_id = ids[0][0]
    if not vender_id:
        print('Error, get vender_id failed')
    shop_id = ids[0][1]
    # print(shop_id)
    if not shop_id:
        print('Error, get shop_id failed')

    category = ''
    cats = re.findall(r"cat: \[(.*?)\]", response.text)
    for cat in cats:
        category = category + cat + ','
    category = category.rstrip(',')
    coupon_url = favourable_url % (skuid, shop_id, vender_id, category.replace(',', '%2c'))
    coupon_response = requests.get(coupon_url, headers=head)
    coupon_response.encoding='gbk'
    fav_data = coupon_response.json()
    if fav_data['skuCoupon']:
        for item in fav_data['skuCoupon']:
            if float(price) >= item['quota']:
                if float(coupon) < item['discount']:
                    coupon = str(item['discount'])
            #desc1.append(u'有效期%s至%s,满%s减%s' % (start_time, end_time, fav_price, fav_count))
        #print(coupon)
    
    # 促销活动，不好判断，作为提示展示
    if fav_data['prom'] and fav_data['prom']['pickOneTag']:
        for item in fav_data['prom']['pickOneTag']:
            # desc2.append(item['content']) 
            discount = discount + item['name'] + ','
        #productsItem['favourableDesc1'] = ';'.join(desc2)
        discount = discount.rstrip(',')
        #print(discount)
    
    return price, coupon, discount


# 根据不同的 shop_type 去对应网站抓取数据
def CrawlShopTypeProduct(url, pid, shop_type):
    swicher = {  #定义一个map，相当于定义case：func()
        1: CrawlJingDong,
        2: lambda :print('invalid shop_type')
    }
    func = swicher.get(shop_type, 2)  #从map中取出方法
    return func(url, pid)


def StartCrawlObj():
    cursor = db.cursor()
    cursor.execute("select * from t_spider_obj")
    # data = cursor.fetchone()
    datas = cursor.fetchall()
    for data in datas:
        pid = data[0]
        pname = data[1]
        pdescribe = data[2]
        url = data[3]
        shop_type = data[4]
        info_table = 't_product_info_' + str(pid)

        price, coupon, discount = CrawlShopTypeProduct(url, pid, shop_type)
        cursor.execute('create table if not exists ' + info_table + '(price varchar(10) NOT NULL, '
            + 'coupon varchar(10) DEFAULT NULL, discount varchar(60) DEFAULT NULL, shop_type smallint(5) unsigned NOT NULL, '
            + 'time_stamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)')

        insert_sql = 'insert into ' + info_table + ' (price, coupon, discount, shop_type) values ('
        if (len(price) > 0):
            insert_sql = insert_sql + '\'' + price + '\''
        else:
            insert_sql = insert_sql + '\'\''

        if (len(coupon) > 0):
            insert_sql = insert_sql + ', \'' + coupon+ '\''
        else:
            insert_sql = insert_sql + ', \'\''

        if (len(discount) > 0):
            insert_sql = insert_sql + ', \'' + discount + '\''
        else:
            insert_sql = insert_sql + ', \'\''

        insert_sql = insert_sql + ', \'' + str(shop_type) + '\')'
        
        try:
            cursor.execute(insert_sql)
            db.commit()
        except Exception as e:
            print('Error, exec sql:', insert_sql, " failed!")
            db.rollback()
        #time.sleep(1) # 休眠1秒
    cursor.close()
        

if __name__=='__main__':
    QueryShopName()
    
    try:
        StartCrawlObj()
    except Exception as e:
        print(e)

    db.close()
