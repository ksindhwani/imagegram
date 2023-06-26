# imagegram

1. This Project was my bandlabs assingment to create a mini instagram called imagegram

### What to build:
* “Imagegram” - a system that allows you to upload images and comment on them
* no frontend/UI is required
  
### User stories (where the user is an API consumer):
* As a user, I should be able to create posts with images (1 post - 1 image)
* As a user, I should be able to set a text caption when I create a post
* As a user, I should be able to comment on a post
* As a user, I should be able to delete a comment (created by me) from a post
* As a user, I should be able to get the list of all posts along with the last 2 comments to each post


### Functional requirements:
* RESTful Web API (JSON)
* Maximum image size - 100MB
* Allowed image formats: .png, .jpg, .bmp.
* Save uploaded images in the original format
* Convert uploaded images to .jpg format and resize to 600x600
* Serve images only in .jpg format
* Posts should be sorted by the number of comments (desc)
* Retrieve posts via a cursor-based pagination
  
### Non-functional requirements:
* Maximum response time for any API call except uploading image files - 50 ms
* Minimum throughput handled by the system - 100 RPS
*  Users have a slow and unstable internet connection

## Setup

* This project uses docker for its development environment. Hence docker is needed to run the project.
* After Docker installation, run the following steps
* Git clone the project. `git clone <project url>`
* go to project directory . `cd <project-directory>`
* rename `.env.template` to `.env` as we use environment variables from the env file
* This project mounts a host volume with docker volumne to save upload images hence create a seperate directory in your system.
* Copy the absolute path of that directory and write it in the `HOST_IMAGE_DIRECTORY` environment variable in the .env File
* run `docker-compose up` in the project directory



### Endpoint and applications to satisfy use cases

`POST /posts` with form-data parameters - Create new Posts
#### Example

```
curl --location '0.0.0.0:8001/posts' \
--form 'image=@"/Users/kunalsindhwani/Desktop/Screenshot 2023-06-26 at 1.24.53 PM.png"' \
--form 'caption="Test Second post with docker"' \
--form 'userId="2"'

```


`POST /posts/{postId}/comments` - Comment on a post

#### Example

```
curl --location '0.0.0.0:8001/posts/2/comments' \
--header 'Content-Type: application/json' \
--data '{
    "userId" : 1,
    "content" : "third Comment on post 2"
}'

```

`DELETE /posts/{postId}/comments/{commentId}` Delete a comment

#### Example

```
curl --location --request DELETE '0.0.0.0:8001/posts/1/comments/7'
```

`GET /posts?cursor={cursorValue}&pageSize={pageSize}` - Get  the list of all posts along with the last 2 comments to each post
#### Example

```
curl --location '0.0.0.0:800/posts?cursor=11&pageSize=10'
```
