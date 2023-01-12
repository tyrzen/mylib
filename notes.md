## Auth endpoints 

>`POST /user {userObject}` -> To create user resource, i.e, Signup.
> 
>`POST /user/:userName {password: ${user_password}}` -> To identify the user by the username and authenticate using the password, i.e, Login.
> 
>`GET /user` -> To list all users from your table/collection
> 
>`GET /user/:id` -> To get details of a particular user (:id can be replaced with :username, or :email or the primary key you have setup for the user.)
> 
>`PUT /user/:id` -> To Update the user object.
> 
>`DELETE /user/:id` -> To delete the user object.

## Client ID
The client_id is a public identifier for apps. Even though it’s public, it’s best that it isn’t guessable by third parties, so many implementations use something like a 32-character hex string.

## Client Secret
The client_secret is a secret known only to the application and the authorization server. It is essential the application’s own password. It must be sufficiently random to not be guessable, which means you should avoid using common UUID libraries which often take into account the timestamp or MAC address of the server generating it. A great way to generate a secure secret is to use a cryptographically-secure library to generate a 256-bit value and then convert it to a hexadecimal representation.

```
{
    "access_token": "qwErtY8zyW1abcdefGHI",
    "token_timeout": "3600",
    "user_name": "john.doe@example.com",
    "token_type": "Bearer",
    "refresh_token": "zxcvbn1JKLMNOPQRSTU",
    "refresh_token_timeout": "5184000"
}

{
    "access_token": "cdf01657-110d-4155-99a7-f986b2ff13a0:int",
    "token_type": "bearer",
    "expires_in": 3599,
    "scope": "apis@acmeinc.com"
}
```

## OAuth2
[read me](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2)