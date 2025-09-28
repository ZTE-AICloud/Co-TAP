FROM zxnp-pict-release-docker.artsh.zte.com.cn/os/alpine:3.20.6 as openresty
LABEL MAINTAINER="cos"

ARG OPENRESTY_VERSION="1.25.3.2"
ARG TONGSUO_VERSION="8.3.3"
ARG PCRE_VERSION="8.45"
ARG LUA_REST_EVENTS_VERSION="0.3.1"
# ARG LUA_RESTY_LMDB_VERSION="1.4.4"
ARG LUA_KONG_NGINX_MODULE_VERSION="0.17.0"
ARG ATC_ROUTER_VERSION="1.7.1"
ARG LUAROCKS_VERSION="3.12.2"


COPY build/packages/openresty-${OPENRESTY_VERSION}.tar.gz /tmp
COPY build/packages/Tongsuo-${TONGSUO_VERSION}.tar.gz /tmp
COPY build/packages/pcre-${PCRE_VERSION}.tar.gz /tmp
COPY build/packages/lua-resty-events-${LUA_REST_EVENTS_VERSION}.tar.gz /tmp
# COPY packages/lua-resty-lmdb-${LUA_RESTY_LMDB_VERSION}.tar.gz /tmp
COPY build/packages/lua-kong-nginx-module-${LUA_KONG_NGINX_MODULE_VERSION}.tar.gz /tmp
COPY build/packages/luarocks-${LUAROCKS_VERSION}.tar.gz /tmp
COPY build/packages/atc-router-${ATC_ROUTER_VERSION}.tar.gz /tmp

COPY build/patches/* /tmp/patches/


# Docker Build Arguments
ARG OPENRESTY_INSTALL="/usr/local/openresty"
ARG RESTY_PCRE_BUILD_OPTIONS="--enable-jit"
ARG RESTY_J="2"
ARG RESTY_CONFIG_OPTIONS="\
    --prefix=${OPENRESTY_INSTALL} \
    --with-file-aio \
    --with-threads \
    --with-http_realip_module \
    --with-http_ssl_module \
    --with-http_sub_module \
    --with-http_stub_status_module \
    --with-http_v2_module \
    --with-http_secure_link_module \
    --with-stream_realip_module \
    --with-stream_ssl_preread_module \
    --without-luajit-lua52 \
    --without-http_fastcgi_module \
    --without-http_uwsgi_module \
    --without-http_scgi_module \
    --without-http_autoindex_module \
    --add-module=/tmp/lua-kong-nginx-module-${LUA_KONG_NGINX_MODULE_VERSION} \
    --add-module=/tmp/lua-kong-nginx-module-${LUA_KONG_NGINX_MODULE_VERSION}/stream \
    --add-module=/tmp/lua-resty-events-${LUA_REST_EVENTS_VERSION} \
    "
ARG RESTY_CONFIG_OPTIONS_MORE=""
ARG RESTY_LUAJIT_OPTIONS="--with-luajit-xcflags='-DLUAJIT_NUMMODE=2 -DLUAJIT_ENABLE_LUA52COMPAT'"
ARG RESTY_PCRE_OPTIONS="--with-pcre-jit"

ARG _RESTY_CONFIG_DEPS="--with-pcre \
    --with-cc-opt='-I${OPENRESTY_INSTALL}/pcre/include -I${OPENRESTY_INSTALL}/openssl/include' \
    --with-ld-opt='-L${OPENRESTY_INSTALL}/pcre/lib -L${OPENRESTY_INSTALL}/openssl/lib -Wl,-rpath,${OPENRESTY_INSTALL}/pcre/lib:${OPENRESTY_INSTALL}/openssl/lib' \
    "

# hadolint ignore=SC2086,SC2231,DL3003,DL4006
RUN echo "==> Installing dependencies ..." \
 && umask 027 \
 && addgroup -g 3000 -S nginx \
 && adduser nginx -u 3000 -H -D -s /sbin/nologin -G nginx \
 && sed -i "s/dl-cdn.alpinelinux.org/mirrors.zte.com.cn/g" /etc/apk/repositories \
 && apk add --no-cache --no-scripts \
    dumb-init \
    tzdata \
    libgcc \
 && apk add --no-cache --no-scripts --virtual .build-deps \
    build-base \
    linux-headers \
    perl-dev \
    readline-dev \
    zlib-dev \
    gd-dev \
    geoip-dev \
    libxslt-dev \
    curl \
 && echo "# Build OpenSSL" \
 && cd /tmp \
 && tar -xzf Tongsuo-${TONGSUO_VERSION}.tar.gz \
 && cd Tongsuo-${TONGSUO_VERSION} \
 && ./config \
    -g \
    -O3 \
    shared \
    -DPURIFY \
    no-threads \
    no-tests \
    enable-ntls \
    --prefix=${OPENRESTY_INSTALL}/openssl \
    --libdir=lib \
    -Wl,-rpath,${OPENRESTY_INSTALL}/openssl/lib \
 && make -j${RESTY_J} \
 && make -j${RESTY_J} install_sw \
 && echo "# Build PCRE" \
 && cd /tmp \
 && tar -xzf pcre-${PCRE_VERSION}.tar.gz \
 && cd /tmp/pcre-${PCRE_VERSION} \
 && ./configure \
     --prefix=${OPENRESTY_INSTALL}/pcre \
     --disable-cpp \
     --enable-utf \
     --enable-unicode-properties \
     ${RESTY_PCRE_BUILD_OPTIONS} \
 && make -j${RESTY_J} \
 && make -j${RESTY_J} install \
 && echo "# Build lua-kong-nginx-module" \
 && cd /tmp \
 && tar -xzf lua-kong-nginx-module-${LUA_KONG_NGINX_MODULE_VERSION}.tar.gz \
 && cd /tmp/lua-kong-nginx-module-${LUA_KONG_NGINX_MODULE_VERSION} \
 && make install LUA_LIB_DIR=${OPENRESTY_INSTALL}/lualib \
 && echo "Build lua-resty-event" \
 && cd /tmp \
 && tar -xzf lua-resty-events-${LUA_REST_EVENTS_VERSION}.tar.gz \
 && cd /tmp/lua-resty-events-${LUA_REST_EVENTS_VERSION} \
 && make install LUA_LIB_DIR=${OPENRESTY_INSTALL}/lualib \
#  && echo "Build lua-resty-lmdb" \
#  && cd /tmp \
#  && tar -xzf lua-resty-lmdb-${LUA_RESTY_LMDB_VERSION}.tar.gz \
#  && cd /tmp/lua-resty-lmdb-${LUA_RESTY_LMDB_VERSION} \
#  && make install LUA_LIB_DIR=$OPENRESTY_INSTALL/lualib \
 && echo "Build atc-router" \
 && cd /tmp \
 && tar -zxf atc-router-${ATC_ROUTER_VERSION}.tar.gz \
 && cd /tmp/atc-router-${ATC_ROUTER_VERSION} \
 && make install-lualib LUA_LIB_DIR=$OPENRESTY_INSTALL/lualib \
 && echo "==> Build OpenResty" \
 && cd /tmp \
 && tar -xzf openresty-${OPENRESTY_VERSION}.tar.gz \
 && cd /tmp/openresty-${OPENRESTY_VERSION} \
 && for i in /tmp/patches/*.patch; do patch -p1 < $i; done \
 && eval ./configure -j${RESTY_J} ${_RESTY_CONFIG_DEPS} ${RESTY_CONFIG_OPTIONS} ${RESTY_CONFIG_OPTIONS_MORE} ${RESTY_LUAJIT_OPTIONS} ${RESTY_PCRE_OPTIONS} \
 && make -j${RESTY_J} \
 && make -j${RESTY_J} install \
 && echo "==> Installing LuaRocks ..." \
 && cd /tmp \
 && tar -xzf luarocks-${LUAROCKS_VERSION}.tar.gz \
 && cd luarocks-${LUAROCKS_VERSION} \
 && ./configure \
      --prefix=${OPENRESTY_INSTALL}/luajit \
      --lua-suffix=jit \
      --with-lua=${OPENRESTY_INSTALL}/luajit \
      --with-lua-include=${OPENRESTY_INSTALL}/luajit/include/luajit-2.1 \
      --lua-version=5.1 \
 && make build && make install \
 && echo "==> Finishing..." \
 && echo "Strip binaries" \
 && strip -s ${OPENRESTY_INSTALL}/nginx/sbin/nginx ${OPENRESTY_INSTALL}/openssl/bin/openssl \
 && strip -s ${OPENRESTY_INSTALL}/luajit/bin/luajit-* ${OPENRESTY_INSTALL}/luajit/lib/*so* \
 && strip -s ${OPENRESTY_INSTALL}/pcre/lib/*so* ${OPENRESTY_INSTALL}/openssl/lib/*so* \
 && echo "Prepare env" \
 && chmod 750 -R ${OPENRESTY_INSTALL}/luajit/bin ${OPENRESTY_INSTALL}/bin ${OPENRESTY_INSTALL}/nginx/sbin \
 && chown 3000:3000 -R ${OPENRESTY_INSTALL} \
 && ln -s ${OPENRESTY_INSTALL}/openssl/bin/openssl /usr/local/bin/openssl \
 && echo "Remove useless files" \
 && rm -rf ${OPENRESTY_INSTALL}/openssl/lib/*.a ${OPENRESTY_INSTALL}/pcre/lib/*.a ${OPENRESTY_INSTALL}/pcre/share \
 && rm -rf ${OPENRESTY_INSTALL}/pod/*  ${OPENRESTY_INSTALL}/nginx/conf/*.default ${OPENRESTY_INSTALL}/nginx/html/* \
 && rm -rf ${OPENRESTY_INSTALL}/nginx/conf/*cgi* ${OPENRESTY_INSTALL}/nginx/conf/*params \
 && rm -rf ${OPENRESTY_INSTALL}/nginx/conf/koi-* ${OPENRESTY_INSTALL}/nginx/conf/win-utf ${OPENRESTY_INSTALL}/nginx/conf/nginx.conf \
 && rm -rf ${OPENRESTY_INSTALL}/nginx/conf/mime.types ${OPENRESTY_INSTALL}/COPYRIGHT ${OPENRESTY_INSTALL}/bin/restydoc ${OPENRESTY_INSTALL}/resty.index \
 && rm -rf ${OPENRESTY_INSTALL}/pod ${OPENRESTY_INSTALL}/site/pod \
 && rm -rf ${OPENRESTY_INSTALL}/luajit/share/man \
 && rm -rf /etc/ssl/misc/*.pl ${OPENRESTY_INSTALL}/bin/*.pl \
 && rm -rf /var/cache/apk/* /tmp/* \
 && apk del --no-scripts .build-deps


FROM zxnp-pict-release-docker.artsh.zte.com.cn/os/alpine:3.20.6

COPY --from=openresty /usr/local/openresty /usr/local/openresty
COPY build/dockerfiles/entrypoint.sh /usr/local/openresty/entrypoint.sh
COPY bin/* /usr/local/openresty/bin/
COPY lib/* /usr/local/lib/
COPY cos.conf.default /etc/cos/cos.conf
COPY logrotate.conf /etc/cos/
COPY cos/ kong/ cos-*.rockspec /tmp/
# COPY kong/ /tmp/
# COPY cert/ /tmp/
# COPY cos-*.rockspec /tmp/
# COPY cos.conf.default /tmp/
# COPY logrotate.conf /tmp/

ENV PATH $PATH:${OPENRESTY_INSTALL}/nginx/sbin:${OPENRESTY_INSTALL}/luajit/bin:${OPENRESTY_INSTALL}/bin

# hadolint ignore=SC2016,DL3003,DL4006
RUN echo "==> Add group and user ..." \
 && umask 027 \
 && addgroup -g 3000 -S nginx \
 && adduser nginx -u 3000 -H -D -s /sbin/nologin -G nginx \
 && echo "==> Installing dependencies ..." \
 && sed -i "s/dl-cdn.alpinelinux.org/mirrors.zte.com.cn/g" /etc/apk/repositories \
 && sed -i "s/umask 022/umask 027/g" /etc/profile \
 && echo 'export PATH=$PATH:/usr/local/openresty/nginx/sbin:/usr/local/openresty/luajit/bin' >> /etc/profile \
 && apk add --no-cache --no-scripts \
    zip \
    logrotate \
 && apk add --no-cache --no-scripts --virtual .build-deps \
    bsd-compat-headers \
    linux-headers \
    musl-dev \
    expat-dev \
    make \
    gcc \
    m4 \
    curl \
 && echo "Installing Cos ..." \
 && export http_proxy="http://proxyxa.zte.com.cn:80" \
 && export https_proxy="https://proxyxa.zte.com.cn:80" \
 && cd /tmp && /usr/local/openresty/luajit/bin/luarocks make cos-*.rockspec OPENSSL_DIR=/usr/local/openresty/openssl CRYPTO_DIR=/usr/local/openresty/openssl \
 && mkdir -p /etc/cos \
 && cp /tmp/cos.conf.default /etc/cos/cos.conf \
 && cp /tmp/logrotate.conf /etc/cos/logrotate.conf \
 && cp /tmp/cert /etc/cos/ \
 && rm -rf /var/cache/apk/* /tmp/* \
 && apk del --no-scripts .build-deps

WORKDIR /usr/local/openresty

ENTRYPOINT ["/usr/local/openresty/entrypoint.sh"]

EXPOSE 8000 8443 8001 8444

STOPSIGNAL SIGQUIT

HEALTHCHECK --interval=60s --timeout=10s --retries=10 CMD kong-health

CMD ["kong", "docker-start"]