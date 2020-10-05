import requests
from lxml import etree

import os

category=""
rootPath = '/home/yangzuo/images/'+category+'/'
os.makedirs(rootPath, exist_ok=True)
startUrl = ""
imgUrl = []

while True:
    response = requests.get(startUrl)
    response.encoding=response.apparent_encoding

    html = etree.HTML(response.text)
    imgs = html.xpath('//img[@class="tupian_img"]')

    for item in imgs:
        imgUrl.append(item.values()[0])
        print(item.values()[0])

    currentPageIndex = str(html.xpath('//div[@id="pages"]/span/text()')[0])
    print("current page: " + currentPageIndex)

    result = html.xpath('//div[@id="pages"]//a[text()="下一页"]')
    nextPage = str(result[0].values()[1])
    startUrl = nextPage
    print("next page: " + nextPage)

    s = nextPage.split("/")[-1]
    nextPageIndex = s.split(".")[0]
    print(nextPageIndex)

    if currentPageIndex == nextPageIndex:
        break

print(imgUrl)

for link in imgUrl:
    name = link.split("/")[-1]
    print(name)
    response = requests.get(link)

    with open(rootPath+name, 'wb') as f:
        f.write(response.content)

