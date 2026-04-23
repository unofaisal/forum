# ROUTES

## AUTH




* Post /register





* Post /login





* Post /logout





### Posts





* Get /posts





* Get /posts?id=1





* Post /posts





### Comments





* Post /comments





* Get /comments?post_id=1





### Likes





* Post /like





### Filters





* Get /posts?category=tech





* Get /posts?user=me





* Get /posts?liked=true





# Flow





### Login 





 Sends email + password\


 Check the DB


 if valid  


* Create session_token


* Store in DB


* Send cookie





### Posting





User creates post


* Request comes in


* Read cookie


* Validate session


* Get user_id


* insert post into DB





### Liking





User Likes post


* Check session


* Insert/Update like table