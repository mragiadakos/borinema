# Borinema

## Description:
Borinema is a private online cinema that allows you and your friends watch the same movie together online.
It synchronizes the movie between the browsers.


## List of TODOs
- Admin installs the software in a VPS with specific configurations
..* HTTP Port
..* HTTPS Certificate file
..* Admin password and username
..* Folder that will contain the movies
..* Authorization options: public; common password for all visitors; enable user mode
....* If common password, then what is the password
....* If user mode is enabled, then enable them after approval or or new users are allowed to enter
..* On user authorization, who press the buttons for the videos: only admin, operators, all users
- Login page for administrator panel
- Administrator's panel provide specific pages
..* Movie browser where the admin can organize movies into folders and uploads subs for movies.
..* Playlist where it will contain list of movies, breaks between them and in the middle of movies 
....* Also enable autoplay, after pressing play
..* Users list, unless the configuration option is public or common password 
....* Users can be approved to enter
....* Deleted
....* Blocked
....* Become operators
- Login page for the main page 
..* if it is public then it will not show up. 
..* If it is common password, then it will request only the 
- Main page for watching the movie
..* Movie player
..* Chat (which will ask a username, if users are not enabled) 