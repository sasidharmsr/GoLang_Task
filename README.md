# GoLang_Task

## About the project

Building a RESTful API that can get/create/update/delete user data from a persistence database and I Used Golang and For DB I Used MongoDB.


### Design

The template follows project convention doc.

MY User Model

// SAMPLE USER
```JSON
{
                "id": "63a73dfbff2b1a88f5c92408",
                "name": "geetha",
                "dob": "2001-03-01T00:00:00Z",
                "address": {
                    "longitude": 80.41185411766716,
                    "latitude": 16.30411413387836,
                    "address": "guntur"
                },
                "description": "My Looks Will Kill Your Eyes",
                "created_at": "2022-12-24T23:29:23.357+05:30",
                "following": null,
                "followers": null
  }
  ```
  
 ### Functionality's

    The API should follow typical RESTful API design pattern.

    The data should be saved in the DB.

    Proper unit tests are Written.
    
    The address of the user includes a geographic coordinates
    
    If You hit this API `http://localhost:8090/user/getnearusers/{userid}` you will return all users Near to Location of this User.
    
    Responce Be Like :-     Here the Responce is Sorted According to the Distance
    ```JSON
            {
                "Distance": 335.72366925836286,
                "UserName": "geetha",
                "UserId": "63a73dfbff2b1a88f5c92408",
                "Address": "guntur"
            },
            {
                "Distance": 356.7349194624298,
                "UserName": "sita",
                "UserId": "63a73e45ff2b1a88f5c9240a",
                "Address": "sattenapalli"
            },
            {
                "Distance": 435.22877554306353,
                "UserName": "srijan",
                "UserId": "63a73e97ff2b1a88f5c9240e",
                "Address": "karimnagar"
            } 
   ```
   
    
