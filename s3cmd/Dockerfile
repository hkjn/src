FROM hkjn/arch

ENV UNPRIVILEGED_UID=500

RUN pacman -Syyy && \
    pacman --noconfirm -S s3cmd && \
    useradd -m user \
            -u ${UNPRIVILEGED_UID} \
            -s /bin/bash 
COPY ["s3cmd_entrypoint", "/usr/local/bin/"]

ENTRYPOINT ["bash", "s3cmd_entrypoint"]
CMD [""]
