
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

Return:
- Nothing for key "session_key"
	- `{"status_code": 400, "status_msg": ...}` 
- Wrong cookie
	- `{"status_code": 403, "status_msg": ...}` 
- Everything goes right
	- `{"status_code": 200, "result": [{...}, ...], "status_msg": ...}` 

Here is struct of game:

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