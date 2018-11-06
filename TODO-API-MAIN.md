# TODO API MAIN

## Login
Not necessary if it is public

### POST /api/main/login
REQ if common password
{
    password
}
REQ if user mode
{
    email
    password
}
RESP
202
401

### POST /api/main/signup
REQ
{
    email
    password
}
RESP
204
422 empty email or not correct format
422 password is lower than 8 characters