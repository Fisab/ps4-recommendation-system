import requests
import re

GAMES_LINKS = []
BASE = 'https://www.metacritic.com/'

def get_data(page_num=0):
	headers = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36'}
	url = 'https://www.metacritic.com/browse/games/release-date/available/ps4/metascore?page={0!s}'.format(page_num)

	r = requests.get(url, headers=headers)
	data = r.content

	links = [i.split('"')[1] for i in re.findall(r'<a href="\/game\/playstation-4\/.{5,25}">', r.text)]

	return [i for i in links if i.count('/') == 3]

for i in range(5):
	GAMES_LINKS.extend(get_data(i))

print(GAMES_LINKS)