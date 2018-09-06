####################################################
# GOLANG BUILDER
####################################################
FROM golang:1.11 as go_builder

COPY . /go/src/github.com/malice-plugins/fprot
WORKDIR /go/src/github.com/malice-plugins/fprot
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan

####################################################
# PLUGIN BUILDER
####################################################
FROM ubuntu:bionic

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/fprot.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware

RUN buildDeps='ca-certificates wget' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 --no-install-recommends \
  && set -x \
  && echo "===> Install F-PROT..." \
  && wget https://github.com/maliceio/malice-av/raw/master/fprot/fp-Linux.x86.32-ws.tar.gz \
  -O /tmp/fp-Linux.x86.32-ws.tar.gz \
  && tar -C /opt -zxvf /tmp/fp-Linux.x86.32-ws.tar.gz \
  && ln -fs /opt/f-prot/fpscan /usr/local/bin/fpscan \
  && ln -fs /opt/f-prot/fpscand /usr/local/sbin/fpscand \
  && ln -fs /opt/f-prot/fpmon /usr/local/sbin/fpmon \
  && cp /opt/f-prot/f-prot.conf.default /opt/f-prot/f-prot.conf \
  && ln -fs /opt/f-prot/f-prot.conf /etc/f-prot.conf \
  && chmod a+x /opt/f-prot/fpscan \
  && chmod u+x /opt/f-prot/fpupdate \
  && ln -fs /opt/f-prot/man_pages/scan-mail.pl.8 /usr/share/man/man8/ \
  && echo "===> Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/*

# Ensure ca-certificates is installed for elasticsearch to use https
RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=go_builder /bin/avscan /bin/avscan

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Update F-PROT Definitions
RUN mkdir -p /opt/malice && /opt/f-prot/fpupdate

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]
CMD ["--help"]

####################################################
# http://files.f-prot.com/files/unix-trial/fp-Linux.x86.64-fs.tar.gz
# http://files.f-prot.com/files/unix-trial/fp-Linux.x86.32-ws.tar.gz
# http://files.f-prot.com/files/unix-trial/fp-Linux.x86.64-ws.tar.gz