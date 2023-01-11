Reader

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