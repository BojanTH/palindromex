FROM postgres:10.1

COPY wrapper.sh /
RUN chmod +x /wrapper.sh
ENTRYPOINT ["/wrapper.sh"]
CMD ["postgres", "-c", "logging_collector=on", "-c", "log_directory=/logs", "-c", "log_filename=postgresql.log", "-c", "log_statement=all"]
