# TODO API ADMIN

## API for the Movies

### POST /api/admin/login
REQ
{
    username
    password
}
RESP
    202 success
    401 failure

### GET /api/admin/movies
RESP
200 [] {
    ID
    name
    duration
}

### GET /api/admin/movies/:id
RESP
200 {
    ID
    name
}
404

### DELETE /api/admin/movies/:id
RESP
204
404


### POST /api/admin/movies
REQ 
Form to upload the movie
RESP
200 {
    ID
}
401
500

### PUT /api/admin/movies/:id
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

## API for configuration
### GET /api/config
RESP
200 []{
    authentication: public | common_password | user_mode
}

## API for Playlist

### GET /api/admin/playlist
RESP
200 [] {
    movie_id : optional id
    order
    break : optional seconds
}

### POST /api/admin/playlist
REQ
{
    movie_id : optional id
    order
    break : optional seconds
}
RESP
204 Success
422 movie_id does not exists
422 break is 0 or lower
422 taken place of the order
422 order can not be less than 0

### DELETE /api/admin/playlist/:order
RESP
204 Success
404 does not exists

## API for users
these are enabled if in the configuration the "authentication" equals "user_mode"


### GET /api/admin/users
RESP
200 [] {
    id
    email
    is_operator
    is_blocked
}

### PUT /api/admin/users/:id
REQ
{
    id
    is_operator
    is_blocked
    is_approved
}
RESP
204
404

### DELETE /api/admin/users/:id
RESP
204
404




## Curl Commands
- curl -H "Content-Type: application/json" \
 -d '{"username":"admin","password":"admin"}' \
 -X POST  http://localhost:8080/api/admin/login

- curl -v -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTQxNzk2Mzg0fQ.V_sJ6M3sY4brDGL27he2D7eyCCvBJ9UD4af_SA9s-rA" \
    -X GET http://localhost:8080/api/admin/isAdmin