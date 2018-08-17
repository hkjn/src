#
# OpenVPN server image.
#
FROM hkjn/alpine

# Install openvpn.
RUN apk add --no-cache bash libressl openvpn

COPY run /usr/local/sbin/
CMD ["run"]

