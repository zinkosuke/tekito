FROM python:3.11-slim

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV POETRY_CACHE_DIR=/usr/share/pypoetry_cache

WORKDIR /opt/app

COPY poetry_requirements.txt /tmp/

RUN set -eux \
 && apt-get update -yqq \
 && apt-get install -y --no-install-recommends \
        curl \
        make \
 # poetry \
 && pip install -r /tmp/poetry_requirements.txt \
 # clean \
 && apt-get autoremove -yqq --purge \
 && apt-get -y clean \
 && rm -rf \
        /tmp/* \
        /usr/share/doc \
        /usr/share/doc-base \
        /usr/share/man \
        /var/lib/apt/lists/* \
        /var/tmp/* \
 && true

COPY app/ .

RUN set -eux \
 && poetry install --no-root \
 && true
