FROM amazoncorretto:11

ARG DIGDAG_VERSION=0.10.5
ARG OWNER=digdag

RUN set -eux \
 && yum install -y \
    shadow-utils \
    unzip \
 # XXX UID/GID \
 && useradd -s /bin/false -m -U -u 1000 ${OWNER} \
 # digdag \
 && curl -fLsS -o /usr/local/bin/digdag \
    "https://dl.digdag.io/digdag-${DIGDAG_VERSION}" \
 && chmod +x /usr/local/bin/digdag \
 && mkdir -p /var/log/digdag \
 && chown -R ${OWNER} /var/log/digdag \
 # awscli \
 && curl -fLsS -o awscliv2.zip \
    "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" \
 && unzip awscliv2.zip \
 && ./aws/install \
 && rm -rf awscliv2.zip aws \
 # clean \
 && yum remove -y \
    shadow-utils \
 && rm -rf /var/cache/yum/* \
 && yum clean all \
 && true

WORKDIR /home/digdag

USER ${OWNER}

COPY config/ ./config/
COPY projects/ ./projects/
COPY scripts/ ./scripts/
COPY shared/ ./shared/

CMD ["java", "-jar", "/usr/local/bin/digdag"]
