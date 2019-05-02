
## Rest api endpoints  
  
### /register [POST]  
Required params:  
- `login` 
- `password`
- `email`(optional)  
  
Return:  
 
 - Bad login/pass(len is 0)
	 - `{"status_code": 400, "status_msg": ...}` 
 - Login occupied
	 - `{"status_code": 409, "status_msg": ...}`
- Done
	- `{"status_code": 202, "status_msg": ...}`
---
### /auth [POST]
Required params:  
- `login`  
- `password`  

Return:
 - Bad login/pass(len is 0)
	- `{"status_code": 400, "status_msg": ...}` 
- Something went wrong when writing session to mysql
	- `{"status_code": 500, "status_msg": ...}` 
- Wrong pass/login
	- `{"status_code": 403, "status_msg": ...}` 
- Done
	- `{"status_code": 200, "status_msg": ...}` 
	- At header will attach cookie
---
### /getFavoriteGenres [GET]
Required params:  
- 🤨(cookie at header with key: "session_key")

Return:
- Nothing for key "session_key"
	- `{"status_code": 400, "status_msg": ...}` 
- Wrong cookie
	- `{"status_code": 403, "status_msg": ...}` 
- Everything goes right
	- `{"status_code": 200, "result": ["someGenre", ...], "status_msg": ...}` 
---
### /getTopGames [GET]
Required params:  
- 🤨(cookie at header with key: "session_key")
- `games_limit`(`int`) at `URL` `<optional>`
- `game_genre`(`string`) at `URL` `<optional>`

Return:
- Nothing for key "session_key"
	- `{"status_code": 400, "status_msg": ...}` 
- Wrong cookie
	- `{"status_code": 403, "status_msg": ...}` 
- Everything goes right
	- `{"status_code": 200, "result": [{struct_game}, ...], "status_msg": ...}` 
---
### /searchGames [GET]
Required params:  
- 🤨(cookie at header with key: "session_key")
- `game_name`(`string`) at `URL`

Return:
- Nothing for key "session_key"
	- `{"status_code": 400, "status_msg": ...}` 
- Wrong cookie
	- `{"status_code": 403, "status_msg": ...}` 
- Everything goes right
	- `{"status_code": 200, "result": [{struct_game}, ...], "status_msg": ...}` 
---
### /getGameById [GET]
Required params:  
- 🤨(cookie at header with key: "session_key")
- `game_id`(`int`) at `URL`

Return:
- Nothing for key "session_key"
	- `{"status_code": 400, "status_msg": ...}` 
- Wrong cookie
	- `{"status_code": 403, "status_msg": ...}` 
- Everything goes right
	- `{"status_code": 200, "result": [{struct_game}, ...], "status_msg": ...}` 
---
### /getGenres [GET]
Required params:  
- 🤨(cookie at header with key: "session_key")

### Here is struct of game:
Return:
- Nothing for key "session_key"
	- `{"status_code": 400, "status_msg": ...}` 
- Wrong cookie
	- `{"status_code": 403, "status_msg": ...}` 
- Everything goes right
	- `{"status_code": 200, "result": ["genre_name", ...], "status_msg": ...}` 
---

```	
GameId        int
Genres        string
Rating        string
Developer     string
OfPlayers     int
Name          string
ImgLink       string
Summary       string
Metascore     int   
UsersScore    float3
ProcessedName string
 ```
- Example of `struct_game`:
```json
{
   "GameId": 264,
   "Genres": "Action,Shooter,First-Person,Arcade",
   "Rating": "T",
   "Developer": "Bungie",
   "OfPlayers": 0,
   "Name": "destiny-house-of-wolves",
   "ImgLink": "https://static.metacritic.com/images/products/games/8/6bbec7228b25ee4c8b5f8dfbbe9147e2-98.jpg",
   "Summary": "Expand your Destiny adventure with a myriad of weapons, armor, and gear to earn in new story missions, 3 new competitive multiplayer maps, and a new cooperative Strike. Expansion II introduces a new competitive elimination mode in the Crucible and an all-new arena activity – The Prison of Elders. The Reef is open. Join the Awoken and hunt down the Fallen rising against us.",
   "Metascore": 72,
   "UsersScore": 4.2,
   "ProcessedName": "destiny house of wolves"
}
```