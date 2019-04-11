import requests
import re
from lxml import html

GAMES_LINKS = []
BASE = 'https://www.metacritic.com'
HEADERS = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36'}


def get_data(page_num=0):
	url = 'https://www.metacritic.com/browse/games/release-date/available/ps4/metascore?page={0!s}'.format(page_num)

	r = requests.get(url, headers=HEADERS)

	links = [i.split('"')[1] for i in re.findall(r'<a href="\/game\/playstation-4\/.{5,25}">', r.text)]

	return [i for i in links if i.count('/') == 3]

for i in range(1):
	GAMES_LINKS.extend(get_data(i))

def get_info(link):
	r = requests.get(BASE + link, headers=HEADERS)
	tree = html.fromstring(r.content)

	data = tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[2]/div[2]/ul')[0]
	result = {
		'genres': [],
		'rating': None,
		'developer': None,
		'of_players': None,
		'name': link.split('/')[-1]
	}

	for i in data.getchildren():
		val = i.text_content().replace('\n', '').replace(' ', '')

		if val.find('Developer:') != -1:
			result['developer'] = val.replace('Developer:', '')
		elif val.find('Genre(s):') != -1:
			result['genres'] = val.replace('Genre(s):', '').split(',')
		elif val.find('#ofplayers:') != -1:
			result['of_players'] = val.replace('#ofplayers:', '')
		elif val.find('Rating:') != -1:
			result['rating'] = val.replace('Rating:', '')	

	result['img_link'] = re.findall(r'<img class="product_image large_image" src=".{20,300}>', r.text)[0].split('"')[3]
	result['summary'] = tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[2]/div[1]/ul/li/span[2]/span/span[2]')[0].text_content()
	result['metascore'] = int(tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[1]/div[1]/div/div/a/div/span')[0].text_content())
	result['user_score'] = float(tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[1]/div[2]/div[1]/div/a/div')[0].text_content())

	return result

for i in range(3):
	data = get_info(GAMES_LINKS[i])
	print(data)






