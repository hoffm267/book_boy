
# book_boy

**Current Project**  
API functional, still developing python cli and additional features.  
Base URL: http://localhost:8080  
Endpoints: /books /audiobooks /users /progress  
All endpoints support basic crud, check controller files for additional functionality  
Test data is added on db creation for easier testing, will be removed later

**To run**  
> ./script/book_boy  

from root dir --- runs db and backend only  

> make clean  

from root dir --- removes db image but not db volume

## ToDo

- Scripts
  - restart server
  - update scripts after adding backend  
- API
  - rename audiobook_time to timestamp and book_page to page
  - fix TODOs in code
  - call to public library api to find exact version of book/ebook
    - give option to enter isbn for book (optional field). If they do, find info through api
  - complete unfinished tests
  - add searching for book/audiobook already added by psql fuzzy score
- PYTHON CLI
  - add api calls
  - finish command parsing
- Docker
  - Reduce image size (multi stage build)
  - finish builds for frontend
- Bruno
  - finish CRUD for all endpoints

NEED: Docker >= 20.x for docker compose support
