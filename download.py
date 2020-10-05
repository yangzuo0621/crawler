import os
import requests
import logging
from lxml import etree

class Crawler:
    def __init__(self, rootFolder, category):
        self.folder = rootFolder
        self.category = category
        self.imageUrls = []
        self.imageFolder = rootFolder + category + '/'
        os.makedirs(self.imageFolder, exist_ok=True)

        logging.basicConfig(level = logging.INFO,format = '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        self.log = logging.getLogger(__name__)

    def __getImageUrls(self, startUrl):
        while True:
            response = requests.get(startUrl)
            response.encoding = response.apparent_encoding

            html = etree.HTML(response.text)

            currentPageIndex = str(html.xpath('//div[@id="pages"]/span/text()')[0])
            self.log.info("current page: " + currentPageIndex)

            imgs = html.xpath('//img[@class="tupian_img"]')
            for item in imgs:
                self.imageUrls.append(item.values()[0])
                self.log.info(item.values()[0])

            result = html.xpath('//div[@id="pages"]//a[text()="下一页"]')
            nextPage = str(result[0].values()[1])
            startUrl = nextPage
            self.log.info("next page url: " + nextPage)

            s = nextPage.split("/")[-1]
            nextPageIndex = s.split(".")[0]
            self.log.info("next page: " + nextPageIndex)

            if currentPageIndex == nextPageIndex:
                break

    def __downloadImages(self):
        for link in self.imageUrls:
            name = link.split("/")[-1]
            response = requests.get(link)

            path = self.imageFolder + name
            with open(path, 'wb') as f:
                f.write(response.content)
            self.log.info("download: " + path)

    def crawlerImages(self, startUrl):
        self.__getImageUrls(startUrl)
        self.__downloadImages()


if __name__ == "__main__":
    crawler = Crawler("", "")
    crawler.crawlerImages("")


