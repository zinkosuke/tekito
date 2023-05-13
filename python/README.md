## python

```console
$ docker-compose build [--no-cache]
$ docker-compose run --rm dev bash

# fmt & lint
root@0bd9939b967d:/opt/app# make -k
poetry run isort .
poetry run black .
reformatted /opt/app/main.py

All done! ✨ 🍰 ✨
1 file reformatted.
poetry run isort --check-only .
poetry run pflake8 .
0
poetry run mypy --install-types --non-interactive
Success: no issues found in 1 source file
poetry run mypy .
Success: no issues found in 1 source file

# main.py
root@0bd9939b967d:/opt/app# poetry run python main.py -h
usage: main.py [-h] {main} ...

Sample command.

positional arguments:
  {main}
    main

options:
  -h, --help  show this help message and exit
```
