from bs4 import BeautifulSoup
import requests

URL = "http://10.201.10.133/0配布用サーバ/_AdobeStock"
page = requests.get(URL)
soup = BeautifulSoup(page.content, 'html.parser')
print(soup.prettify())
