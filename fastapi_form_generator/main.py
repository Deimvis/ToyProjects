from fastapi import Depends, FastAPI, Request
from fastapi.templating import Jinja2Templates
from models import Movie

app = FastAPI()
templates = Jinja2Templates(directory='templates')


@app.get('/movie')
async def get_movie(request: Request):
    return templates.TemplateResponse('movie.html', {'request': request, 'Movie': Movie})


@app.post('/movie')
async def post_movie(movie: Movie = Depends(Movie.form)):
    return movie.model_dump()