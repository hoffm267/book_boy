
# book_boy

**Current Project**  
Minimals endpoints, but actively working to add more.  
Base URL: localhost:8080  
Endpoints: check main.go for currently available

**To run**  
> make backend-dev

## ToDo

- Scripts
  - restart server
  - update scripts after adding backend  
- API
  - rename audiobook_time to timestamp and book_page to page
  - fix TODOs in code
  - make search function that dynamically searched for params in request string
  - call to public library api to find exact version of book/ebook
    - give option to enter isbn for book (optional field). If they do, find info through api
  - better conversion method than percent (have not tested to see how accurate yet)
  - add tests for audio book time in progress
  - switch make backend-dev do go run, not build image
- PYTHON CLI
  - add api calls
  - finish command parsing
- Reduce image size (multi stage build)
- GET RID OF BACKEND-DEV IMAGE RUN GO RUN INSTEAD (check scripts)
- .gitignore
- fix the readme

NEED: Docker >= 20.x for docker compose support
