FROM hkjn/alpine:0.2.0

ARG tf_arch
ARG tf_version

ENV TF_VERSION=$tf_version \
    USER=tfuser \
    TF_ARCH=$tf_arch

COPY hashicorp.asc .

RUN echo "fetching https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_${TF_ARCH}.zip"
RUN apk add --no-cache ca-certificates curl gnupg && \
    gpg --import hashicorp.asc && \
    curl -fsSLO https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_${TF_ARCH}.zip && \
    curl -fsSLO https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS && \
    curl -fsSLO https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS.sig && \
    echo "$(grep terraform_${TF_VERSION}_${TF_ARCH}.zip terraform_${TF_VERSION}_SHA256SUMS)" > SHA256SUMS && \
    gpg --verify terraform_${TF_VERSION}_SHA256SUMS.sig terraform_${TF_VERSION}_SHA256SUMS && \
    sha256sum -c SHA256SUMS && \
    unzip terraform_${TF_VERSION}_${TF_ARCH}.zip && \
    mv terraform /usr/bin/ && \
    rm terraform_${TF_VERSION}_${TF_ARCH}.zip terraform_${TF_VERSION}_SHA256SUMS terraform_${TF_VERSION}_SHA256SUMS.sig SHA256SUMS && \
    adduser -D $USER

USER $USER
WORKDIR /home/$USER

ENTRYPOINT ["terraform"]
