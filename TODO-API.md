# TODO API ADMIN

## API for the Movies

### POST /api/admin/login (/)
REQ
{
    username
    password
}
RESP
    202 success
    401 failure

### GET /api/admin/movies (/)
Req {
    last_seen_date
    limit
}
RESP
200 [] {
    ID
    name
    duration
}

### GET /api/admin/movies/:id (/)
RESP
200 {
    ID
    state: downloading | error | finished 
    filetype : mp4 | webm  
}
404

### DELETE /api/admin/movies/:id (/)
RESP
204
404


### POST /api/admin/movies/link (/)
REQ 
{

}
RESP
200 {
    id
}
401
500

### PUT /api/admin/movies/:id (/)
REQ
{
    name
}
RESP
200 {
    ID
    name
}
404

### PUT /api/admin/movies/:id/select (/)
RESP
204
404

### DELETE /api/admin/movies/selected (/)
RESP
204
404

### GET /api/admin/movies/:id/subtitles
RESP
200 []{
    id
    movie_id
    srclang
}
404


### POST /api/admin/movies/:id/subtitles
REQ
Form upload file and srclang
RESP
204 success
404
422 not correct file
422 not correct srclang

### DELETE /api/admin/movies/:id/subtitles/:id
RESP
204
404 
500


### GET /api/cinema/movie
RESP 
200 - with the movie
404 - there is not any movie selected yet


### GET /api/cinema/movie/info
RESP 
200 - with the movie
404 - there is not any movie selected yet

## Websocket API
### GET /api/admin/ws
- Send (to admin) Progress of movie 'download_progress_movie'
- Send (to admin) Movie has been selected 'new_selected_movie'
- Receive (from admin) start of movie 'start_movie'
- Receive (from admin) stop the movie 'stop_movie'
- Receive (from admin) pause the movie 'pause_movie'
- Receive (from admin) change the time 'change_time_movie'

### GET /api/ws
- Send (to user) a movie has been selected 'new_selected_movie'
- Send (to user) start of movie 'start_movie'
- Send (to user) stop the movie 'stop_movie'
- Send (to user) pause the movie 'pause_movie'
- Send (to user) change the time 'change_time_movie'

## API for configuration
### GET /api/config
RESP
200 []{
    authentication: public | common_password | user_mode
}





## Curl Commands
- curl -H "Content-Type: application/json" \
 -d '{"username":"admin","password":"admin"}' \
 -X POST  http://localhost:8080/api/admin/login

- curl -v -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTQxNzk2Mzg0fQ.V_sJ6M3sY4brDGL27he2D7eyCCvBJ9UD4af_SA9s-rA" \
    -X GET http://localhost:8080/api/admin/isAdmin