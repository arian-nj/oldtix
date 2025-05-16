 2025-02-14 13:08

Status: 

Tags: [[card-game]] [[authentication]]

# Endpoints
/auth/guest
	check uid and send to server
	server create or send back token
/auth/register
/auth/

# Authentication
how i should connect multiple auth methods to one user
**person** table is main one

| id  | name  | username  |
| --- | ----- | --------- |
| 1   | arian | king kong |

username password (delete in production)

| id  | person_id | password   |
| --- | --------- | ---------- |
| 3   | 1         | snhknhnhad |

phone number connects to it 

| id  | person_id | number     |
| --- | --------- | ---------- |
| 2   | 1         | 0912112112 |
telegram 
can request phone number
or in app give a link when user clicks goes to a unique /start in telegram bot of game

| id  | person_id | telegram_user_id |
| --- | --------- | ---------------- |
| 2   | 1         | 19912134231      |


**Authentication**
- During the authentication flow, we log users in by sending their username and password to the authentication server. To make this experience seamless, we can generate a user ID and a user password directly on the player's machine and save them invisibly. This way, players don't have to do anything. We call this "credentialless login".
- The server sends back an auth token to identify the user and keep them logged in while the Godot game is running.


# References

https://docs.w4.gd

