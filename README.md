
# book_boy

**WIP**
-local data.json file for now 
-data.json for format reference

**ToDo**:
- Scripts:
    - restart server
- Reduce image size (multi stage build)
- API
    - fix TODOs in code
    - call to public library api to find exact version of book/ebook 
    -better conversion method than percent (have not tested to see how accurate yet)
    - add tests for progress
- PYTHON CLI:
- .gitignore
- fix the readme


NEED: Docker >= 20.x for docker compose support

**To run**
> prod image: ./scripts/book_boy
> dev image: MODE=dev ./scripts/book_boy
