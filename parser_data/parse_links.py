import requests
import re
from lxml import html
import MySQLdb
import MySQLdb.cursors
import json
import time


with open('../rest_api/config/credentials.json') as json_file:  
	config = json.load(json_file)
	DATABASE = config['mysql']['database']
	SQL_USER = config['mysql']['login']
	SQL_PASSWD = config['mysql']['password']
	IP, PORT = config['mysql']['ip'], config['mysql']['port']

GAMES_LINKS = []
BASE = 'https://www.metacritic.com'
HEADERS = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36'}


def connectEntities():
	db = MySQLdb.connect(
		cursorclass=MySQLdb.cursors.DictCursor,
		host=IP,
		port=int(PORT),
		db=DATABASE,
		user=SQL_USER,
		passwd=SQL_PASSWD,
		use_unicode=True,
		charset='utf8'
	)
	return db


def get_data(page_num=0):
	url = 'https://www.metacritic.com/browse/games/release-date/available/ps4/metascore?page={0!s}'.format(page_num)

	r = requests.get(url, headers=HEADERS)

	links = [i.split('"')[1] for i in re.findall(r'<a href="\/game\/playstation-4\/.{5,25}">', r.text)]

	return [i for i in links if i.count('/') == 3]

for i in range(10):
	GAMES_LINKS.extend(get_data(i))
	time.sleep(1)


def get_info(link):
	r = requests.get(BASE + link, headers=HEADERS)
	tree = html.fromstring(r.content)

	data = tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[2]/div[2]/ul')[0]
	result = {
		'genres': [],
		'rating': 0,
		'developer': '',
		'of_players': 0,
		'name': link.split('/')[-1]
	}

	for i in data.getchildren():
		val = i.text_content().replace('\n', '').replace(' ', '')

		if val.find('Developer:') != -1:
			result['developer'] = val.replace('Developer:', '')
		elif val.find('Genre(s):') != -1:
			result['genres'] = val.replace('Genre(s):', '').split(',')
		elif val.find('#ofplayers:') != -1:
			result['of_players'] = val.replace('#ofplayers:', '').replace('Upto', '').replace('NoOnlineMultiplayer', '0').replace('Player', '').replace('morethan')
		elif val.find('Rating:') != -1:
			result['rating'] = val.replace('Rating:', '')[0]

	result['img_link'] = re.findall(r'<img class="product_image large_image" src=".{20,300}>', r.text)[0].split('"')[3]
	result['summary'] = tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[2]/div[1]/ul/li/span[2]/span/span[2]')[0].text_content().replace('\'', '')
	result['metascore'] = int(tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[1]/div[1]/div/div/a/div/span')[0].text_content())
	result['user_score'] = float(tree.xpath('//*[@id="main"]/div/div[3]/div/div/div[2]/div[1]/div[2]/div[1]/div/a/div')[0].text_content())

	return result

def insert_sql(data):
	print(data)
	query = '''
		INSERT INTO games 
			(`genres`,`rating`,`developer`,`ofplayers`,`name`,`img_link`,`summary`,`metascore`,`users_score`)
		VALUES
			("{0!s}", "{1!s}", "{2!s}", {3!s}, "{4!s}", "{5!s}", '{6!s}', {7!s}, {8!s})
	'''.format(
		','.join(data['genres']), 
		data['rating'], 
		data['developer'], 
		data['of_players'],
		data['name'],
		data['img_link'],
		data['summary'],
		data['metascore'],
		data['user_score']
	)
	print(query)

	db = connectEntities()
	cursor = db.cursor()
	cursor.execute(query)
	db.commit()


for i in range(1999):
	print('%s/1999' % i)
	try:
		data = get_info(GAMES_LINKS[i])
		time.sleep(0.5)
	except:
		continue
	insert_sql(data)
	print('\n'*3)





