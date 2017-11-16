FROM hkjn/arch

ENV UNPRIVILEGED_UID=500

RUN pacman -Syyy && \
    pacman --noconfirm -S s3cmd && \
    useradd -m user \
            -u ${UNPRIVILEGED_UID} \
            -s /bin/bash 

ENTRYPOINT ["bash", "s3cmd_entrypoint"]
CMD [""]
