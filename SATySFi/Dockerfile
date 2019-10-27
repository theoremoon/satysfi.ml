FROM ocaml/opam2:ubuntu-18.04-ocaml-4.07 as builder

USER root
RUN apt update &&\
    apt -y --no-install-recommends install build-essential autoconf git m4 unzip wget ca-cacert ca-certificates ruby

USER opam
WORKDIR /home/opam
RUN git clone https://github.com/gfngfn/SATySFi

WORKDIR /home/opam/opam-repository
RUN git pull && eval `opam env` && opam repository add satysfi-external https://github.com/gfngfn/satysfi-external-repo.git && opam update

WORKDIR /home/opam/SATySFi
RUN (opam pin add -y satysfi . || true) && opam install -y satysfi
RUN sh download-fonts.sh

FROM ubuntu:18.04

WORKDIR /satysfi
COPY --from=builder /home/opam/SATySFi/lib-satysfi ./lib-satysfi
COPY --from=builder /home/opam/SATySFi/install-libs.sh .

RUN sh install-libs.sh
COPY --from=builder /home/opam/.opam/4.07/bin/satysfi /usr/bin/satysfi

WORKDIR /mount
