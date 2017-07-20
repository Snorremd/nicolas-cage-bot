FROM scratch
ADD cacert.pem /etc/ssl/certs/
ADD main /
CMD ["/main"]