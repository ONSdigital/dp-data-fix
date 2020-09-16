# dp-data-fix


# Content mover

Script for moving published content. Copies the published content into the new collection at the desired location. 
The script will then find a fix any broken links in `.json` pages in the published content directory. Any fixed content
 will also added to the collection.
 
 _Note:_ content can only be moved if that content is not already in another collection. 

## Setting up from scratch

1. SSH into the environment

2. Run a golang container with a volume that maps the content directory on the publishing box:
```
sudo docker run -i -t --name dp-data-fix \
   --userns=host \
   -v <CONTENT_DIR>:<VOLUME_NAME>:rw \
   golang /bin/bash
```

3. Update the container and install vim (always useful)
```
apt-get update && apt-get install vim
```

4. Make the appropriate dir structure:
```
cd src && mkdir -p github.com/ONSdigital && cd github.com/ONSdigital
```

5. Clone the content move script:
```
git clone -b feature/content-mover https://github.com/ONSdigital/dp-zebedee-utils.git
```

6. Go get the script dependencies:
```
go get github.com/satori/go.uuid
go get github.com/ONSdigital/log.go/log
```

7. Move into the moves dir
```
cd dp-zebedee-utils/moves
```

## Running the script.

### Config

| Flag       | Description                                    |
|------------|:-----------------------------------------------|
| zeb_root   | The zebedee root directory                     |
| collection | The name of the collection to use for the move |
| src        | The uri of the published content to be moved   |
| dest       | The uri to move the content to                 |
| create     | Should a new collection be created?            |

### Example

Compile:
```
go build -o moves
```
Run:
```
./moves -zeb_root="/zebedee_root" \
            -create=true \
            -collection="testCollection" \
            -src="/aaa/bbb/ccc" \
            -dest="/aaa/bbb/ccc/ddd"
```