FROM debian
LABEL MAINTAINER="Thorsten Hersam <thorsten.hersam@neutryno.de>"
COPY ./dist/app /app
ENTRYPOINT /app