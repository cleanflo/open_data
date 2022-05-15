# ----- All Data ----- #
FROM mcr.microsoft.com/mssql/server:2019-latest AS base
USER root

ENV DEFAULT_MSSQL_SA_PASSWORD=Dfo54E*SF
ENV ACCEPT_EULA=Y

ADD data.tar.gz /var/opt/mssql/

COPY restore-db.sh entrypoint.sh /opt/mssql/bin/
RUN chmod +x /opt/mssql/bin/restore-db.sh /opt/mssql/bin/entrypoint.sh
RUN chown -R mssql:root /var/opt/mssql/data && \
                chmod 0755 /var/opt/mssql/data && \
                chmod -R 0650 /var/opt/mssql/data/*

USER mssql
RUN /opt/mssql/bin/restore-db.sh

CMD [ "/opt/mssql/bin/sqlservr" ]
ENTRYPOINT [ "/opt/mssql/bin/entrypoint.sh" ]

FROM base as all-data

# ----- Alberta ----- #
FROM base AS alberta
RUN mv /var/opt/mssql/data/alberta* ./ && \
        rm -rf /var/opt/mssql/data/* &&  \
        mv ./alberta* /var/opt/mssql/data/
USER mssql

# ----- British Columbia ----- #
FROM base AS british-columbia
RUN mv /var/opt/mssql/data/british-columbia* ./ && \
        rm -rf /var/opt/mssql/data/* && \
        mv ./british-columbia* /var/opt/mssql/data/
USER mssql

# ----- Nova Scotia ----- #
FROM base AS nova-scotia
RUN mv /var/opt/mssql/data/nova-scotia* ./ && \
        rm -rf /var/opt/mssql/data/* && \
        mv ./nova-scotia* /var/opt/mssql/data/
USER mssql

# ----- Ontario ----- #
FROM base AS ontario
RUN mv /var/opt/mssql/data/ontario* ./ && \
        rm -rf /var/opt/mssql/data/* && \
        mv ./ontario* /var/opt/mssql/data/
USER mssql

# ----- Saskatchewan ----- #
FROM base AS saskatchewan
RUN mv /var/opt/mssql/data/saskatchewan* ./ && \
        rm -rf /var/opt/mssql/data/* && \
        mv ./saskatchewan* /var/opt/mssql/data/ 
USER mssql

# ----- OPEN DATA ----- #
# create images for postalcodes, censusdata, etc. in mssql, mysql, etc.
FROM base AS open-data
RUN mv /var/opt/mssql/data/open-data* ./ && \
        rm -rf /var/opt/mssql/data/* && \
        mv ./open-data* /var/opt/mssql/data/ 
USER mssql
