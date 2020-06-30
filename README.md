# twitter-clone using GraphQl and Go
To use this clone. First Setup a local Mysql database and make changes in .env file for connection string.
Create a local database according to the name provided in .env file.
# go run server.go to run the server
All dependencies would be automatically downloaded

/graphql endpoint cannot be accesed by default. It is prootected.
Generate a token first to access the /graphql endpoint.

To generate a token vistit /signup to create new user. This will generate a token for you as well no need to login manually after /signup until server is running.

If account already exist then use /login endpoint.

Request to login or signup can be made using postman by sending the given below sample request as json body.

{   
        
	"email": "diwakarsingh@gmail.com",
	
    "password":"dev"
 
 }
 
# After authentication no need to pass the token from header. It is done automatically.

# Query
query allUsers{
  allUsers {
    email
  }
}


query myPost{
  myPost {
    email
    text
    time
  }
}

query FollowedPost{
  followedPost{
		email
    text
    time
    
    
}
}

query FollowedUser{
  followedUser{
		email
  }
}

# Mutation

mutation createPost{
  createPost(input:{text:"Post 2"}) {
  text
  time
 
    
}
}

mutation FollowUser{
  followUser(input:{Email:"diwakarsingh"}) {
			MyEmail
 			followedEmail
    
}
}

